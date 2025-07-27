package task

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/ent/queuetask"
	"go.uber.org/zap"
)

// Worker - global struct to manage all running tasks.
type Worker struct {
	name    string
	db      *ent.Client
	taskMap map[string]TaskDescriptor
	logger  *zap.Logger
}

func NewWorker(name string, db *ent.Client, logger *zap.Logger) *Worker {
	return &Worker{
		name:    name,
		db:      db,
		taskMap: map[string]TaskDescriptor{},
		logger:  logger.With(zap.String("Service", "Worker"), zap.String("Name", name)),
	}
}

// GetTaskToRun returns list of tasks to run at the moment.
func (w *Worker) GetTaskToRun(ctx context.Context) *RunRequest {
	w.logger.Debug("Find task to run")
	pendingTask, err := w.db.QueueTask.Query().
		Where(
			queuetask.And(
				queuetask.Status(string(TaskStatusPending)),
			),
		).
		ForUpdate(
			sql.WithLockAction(sql.SkipLocked),
		).
		First(ctx)
	if err != nil && ent.IsNotFound(err) {
		w.logger.Debug("Find task to run. Done. No task")
		return nil
	}
	if err != nil && !ent.IsNotFound(err) {
		panic(err)
	}

	w.logger.Debug("Find task to run. Done", zap.String("TaskName", pendingTask.TaskName))
	return &RunRequest{
		ID:   pendingTask.ID,
		Name: pendingTask.TaskName,
		Data: pendingTask.Data,
	}
}

func (w *Worker) CleanStaleJobs(ctx context.Context) error {
	w.logger.Debug("Cleaning stale jobs")
	err := w.db.QueueTask.Update().
		SetStatus(string(TaskStatusInterrupted)).
		Where(
			queuetask.And(
				queuetask.Status(string(TaskStatusRunning)),
				queuetask.WorkerName(w.name),
			),
		).
		Exec(ctx)
	w.logger.Debug("Done")
	return err
}

// Run starts TaskWorker in infinite loop of check-run tasks.
func (w *Worker) Run(ctx context.Context) error {
	// clean stale jobs on startup
	err := w.CleanStaleJobs(ctx)
	if err != nil {
		panic(err)
	}
	// Vars for runtime behaviour
	interval := time.Minute
	counter := 1
	quit := make(chan int, 1)

	const CONTINUE_FUNC = 1
	const TERMINATE_FUNC = 2

	for {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					w.logger.Warn("Recovered", zap.Any("error", err))
					quit <- CONTINUE_FUNC
				}
			}()
			ticker := time.NewTicker(interval)

			for {
				select {
				case <-ticker.C:
					w.logger.Debug("TaskWorker tick start")
					task := w.GetTaskToRun(ctx)
					if task != nil {
						w.ExecuteTask(ctx, task) // TODO check return value of error
					}
					w.logger.Debug("TaskWorker tick end")
				case <-ctx.Done():
					quit <- TERMINATE_FUNC
					return
				}
			}
		}()
		status := <-quit
		if status == TERMINATE_FUNC {
			return nil
		}
		counter++
		if counter > 10 {
			panic("TaskWorker restart limit exceeded")
		}
	}
}

// ExecuteTask executes tasks with logs and stuff.
func (w *Worker) ExecuteTask(ctx context.Context, task *RunRequest) error {
	w.logger.Info("Executing task", zap.String("Task", task.Name))
	currentTaskInfo, err := w.db.QueueTask.Query().
		Where(
			queuetask.ID(task.ID),
		).
		First(ctx)
	if err != nil && !ent.IsNotFound(err) {
		panic(err)
	}
	if currentTaskInfo.Status == string(TaskStatusRunning) {
		return nil
	}

	// Lock task
	err = currentTaskInfo.Update().
		SetStatus(string(TaskStatusRunning)).
		SetWorkerName(w.name).
		Exec(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := recover(); err != nil {
			w.logger.Warn("Recovered", zap.Any("error", err))
			err = currentTaskInfo.Update().
				SetStatus(string(TaskStatusError)).
				Exec(ctx)
			if err != nil {
				panic(err)
			}
		}
	}()
	w.logger.Debug("Executing task. Locked task", zap.String("TaskName", task.Name))

	taskDescriptor := getTaskInfo(task.Name)
	taskStatus := TaskStatusDone
	// Execute task
	err = taskDescriptor.Func(ctx, task.Data)
	if err != nil {
		w.logger.Error("Error during executiong of the taks", zap.String("TaskName", task.Name), zap.Error(err))
		taskStatus = TaskStatusError
	}

	// Unlock task
	err = currentTaskInfo.Update().
		SetStatus(string(taskStatus)).
		Exec(ctx)
	if err != nil {
		panic(err)
	}
	w.logger.Debug("Executing task. Unlocked task", zap.String("TaskName", task.Name))
	return nil
}
