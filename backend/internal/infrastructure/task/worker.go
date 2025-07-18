package task

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/ent/queuetask"
	"go.uber.org/zap"
)

// TaskWorker - global struct to manage all running tasks.
type TaskWorker struct {
	db      *ent.Client
	taskMap map[string]TaskDescriptor
	logger  *zap.Logger
}

func NewTaskWorker(db *ent.Client, logger *zap.Logger) *TaskWorker {
	return &TaskWorker{
		db:      db,
		taskMap: map[string]TaskDescriptor{},
		logger:  logger.With(zap.String("Service", "TaskWorker")),
	}
}

// GetTaskToRun returns list of tasks to run at the moment
func (ts *TaskWorker) GetTaskToRun() *TaskRunRequest {
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
	return &TaskRunRequest{
		ID:   pendingTask.ID,
		Name: pendingTask.TaskName,
		Data: pendingTask.Data,
	}
}

// Run starts TaskWorker in infinite loop of check-run tasks
func (ts *TaskWorker) Run(ctx context.Context) error {
	// Vars for runtime behaviour
	interval := time.Duration(time.Minute)
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
		counter += 1
		if counter > 10 {
			panic("TaskWorker restart limit exceeded")
		}
	}
}

// ExecuteTask executes tasks with logs and stuff
func (ts *TaskWorker) ExecuteTask(task *TaskRunRequest) error {
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
