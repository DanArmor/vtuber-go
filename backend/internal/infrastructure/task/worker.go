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
func (ts *Worker) GetTaskToRun() *RunRequest {
	ts.logger.Debug("Find task to run")
	pendingTask, err := ts.db.QueueTask.Query().
		Where(
			queuetask.And(
				queuetask.Status(string(TaskStatusPending)),
			),
		).
		ForUpdate(
			sql.WithLockAction(sql.SkipLocked),
		).
		First(context.Background())
	if err != nil && ent.IsNotFound(err) {
		ts.logger.Debug("Find task to run. Done. No task")
		return nil
	}
	if err != nil && !ent.IsNotFound(err) {
		panic(err)
	}

	ts.logger.Debug("Find task to run. Done", zap.String("TaskName", pendingTask.TaskName))
	return &RunRequest{
		ID:   pendingTask.ID,
		Name: pendingTask.TaskName,
		Data: pendingTask.Data,
	}
}

// Run starts TaskWorker in infinite loop of check-run tasks.
func (ts *Worker) Run(ctx context.Context) error {
	// Vars for runtime behaviour
	interval := time.Minute
	counter := 1
	quit := make(chan int, 1)

	for {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					ts.logger.Warn("Recovered", zap.Any("error", err))
					quit <- 1
				}
			}()
			ticker := time.NewTicker(interval)

			for {
				select {
				case <-ticker.C:
					ts.logger.Debug("TaskWorker tick start")
					task := ts.GetTaskToRun()
					if task != nil {
						ts.ExecuteTask(task) // TODO check return value of error
					}
					ts.logger.Debug("TaskWorker tick end")
				}
			}
		}()
		<-quit
		counter++
		if counter > 10 {
			panic("TaskWorker restart limit exceeded")
		}
	}
}

// ExecuteTask executes tasks with logs and stuff.
func (ts *Worker) ExecuteTask(task *RunRequest) error {
	ts.logger.Info("Executing task", zap.String("Task", task.Name))
	currentTaskInfo, err := ts.db.QueueTask.Query().
		Where(
			queuetask.ID(task.ID),
		).
		First(context.Background())
	if err != nil && !ent.IsNotFound(err) {
		panic(err)
	}
	if currentTaskInfo.Status == string(TaskStatusRunning) {
		return nil
	}

	// Lock task
	err = currentTaskInfo.Update().SetStatus(string(TaskStatusRunning)).Exec(context.Background())
	if err != nil {
		panic(err)
	}
	ts.logger.Debug("Executing task. Locked task", zap.String("TaskName", task.Name))

	taskDescriptor := getTaskInfo(task.Name)
	// Execute task
	taskDescriptor.Func(task.Data)

	// Unlock task
	err = currentTaskInfo.Update().
		SetStatus(string(TaskStatusDone)).
		Exec(context.Background())
	if err != nil {
		panic(err)
	}
	ts.logger.Debug("Executing task. Unlocked task", zap.String("TaskName", task.Name))
	return nil
}
