package task

import (
	"context"
	"time"

	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/ent/queuescheduledtask"
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

// ExecuteTask executes tasks with logs and stuff
func (tw *TaskWorker) ExecuteTask(task TaskDescriptor) error {
	tw.logger.Info("Executing task", zap.String("Task", task.Name))
	currentTaskInfo, err := tw.db.QueueScheduledTask.Query().
		Where(
			queuescheduledtask.And(
				queuescheduledtask.TaskName(task.Name),
			),
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
	tw.logger.Debug("Executing task. Locked task", zap.String("TaskName", task.Name))

	// Execute task
	task.Func(nil)

	// Unlock task
	err = currentTaskInfo.Update().
		SetStatus(string(TaskStatusPending)).
		SetLastRunTimestamp(time.Now()).
		Exec(context.Background())
	if err != nil {
		panic(err)
	}
	tw.logger.Debug("Executing task. Unlocked task", zap.String("TaskName", task.Name))
	return nil
}
