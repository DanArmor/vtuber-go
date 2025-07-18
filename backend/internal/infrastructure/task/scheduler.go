package task

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/ent/queuescheduledtask"
	"go.uber.org/zap"
)

// TaskScheduler - global struct to manage all running tasks.
type TaskScheduler struct {
	db     *ent.Client
	logger *zap.Logger
}

func NewTaskScheduler(db *ent.Client, logger *zap.Logger) *TaskScheduler {
	return &TaskScheduler{
		db:     db,
		logger: logger.With(zap.String("Service", "TaskScheduler")),
	}
}

// AddScheduleTaskInput is input for AddTaskToRunList func.
type AddScheduleTaskInput struct {
	ScheduleTaskName string // Unique name for this scheduled task
	TaskName         string // TaskName is task to run (name from task descriptor)
	Interval         *int   // Interval to run task with
}

// AddScheduleTask adds task to execution.
func (ts *TaskScheduler) AddScheduleTask(input AddScheduleTaskInput) {
	ts.logger.Debug(
		"Add scheduler task",
		zap.String("ScheduleTask", input.ScheduleTaskName),
		zap.String("Task", input.TaskName),
	)

	err := ts.db.QueueScheduledTask.Create().
		SetScheduleName(input.ScheduleTaskName).
		SetTaskName(input.TaskName).
		SetStatus(string(TaskStatusPending)).
		SetInterval(*input.Interval).
		OnConflict(sql.ConflictColumns("schedule_name"), sql.ResolveWithNewValues()).
		Exec(context.Background())
	if err != nil {
		panic(err)
	}
}

// TODO add `for update skip locked` ?
// SchedulerGetTasksToRun returns list of tasks to run at the moment.
func (ts *TaskScheduler) SchedulerGetTasksToRun() []RunRequest {
	ts.logger.Debug("Check scheduler tasks")
	pendingTasks, err := ts.db.QueueScheduledTask.Query().Where(
		queuescheduledtask.And(
			queuescheduledtask.StatusEQ(string(TaskStatusPending)),
		),
	).All(context.Background())
	if err != nil && !ent.IsNotFound(err) {
		panic(err)
	}
	timeNow := time.Now()
	tasksNamesToExecute := []string{}
	for _, task := range pendingTasks {
		duration := timeNow.Sub(task.LastRunTimestamp)
		if duration > time.Duration(task.Interval)*time.Millisecond {
			tasksNamesToExecute = append(tasksNamesToExecute, task.TaskName)
		}
	}
	result := []RunRequest{}
	for _, name := range tasksNamesToExecute {
		result = append(result, RunRequest{
			Name: name,
			Data: nil,
		})
	}
	ts.logger.Debug("Check scheduler tasks. Done", zap.Int("Result items count", len(result)))
	return result
}

func (ts *TaskScheduler) RequestTasksRun(tasks []RunRequest) {
	for _, task := range tasks {
		err := ts.db.QueueTask.Create().
			SetTaskName(task.Name).
			SetData(task.Data).
			SetStatus(string(TaskStatusPending)).
			Exec(context.Background())
		ts.logger.Debug("Created queue task", zap.String("TaskName", task.Name))
		if err != nil {
			panic(err)
		}
	}
}

// Run starts TaskScheduler in infinite loop of check-run tasks.
func (ts *TaskScheduler) Run(ctx context.Context) error {
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
					ts.logger.Debug("TaskScheduler tick start")
					tasks := ts.SchedulerGetTasksToRun()
					ts.RequestTasksRun(tasks)
					ts.logger.Debug("TaskScheduler tick end")
				}
			}
		}()
		<-quit
		counter++
		if counter > 10 {
			panic("TaskScheduler restart limit exceeded")
		}
	}
}
