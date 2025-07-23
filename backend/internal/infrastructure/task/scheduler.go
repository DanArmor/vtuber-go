package task

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/ent/queuescheduledtask"
	"github.com/DanArmor/vtuber-go/ent/queuetask"
	"go.uber.org/zap"
)

// TaskScheduler - global struct to manage all running tasks.
type TaskScheduler struct {
	db        *ent.Client
	logger    *zap.Logger
	dataFuncs map[string]func() TaskInput
}

func NewTaskScheduler(db *ent.Client, logger *zap.Logger) *TaskScheduler {
	return &TaskScheduler{
		db:        db,
		logger:    logger.With(zap.String("Service", "TaskScheduler")),
		dataFuncs: map[string]func() TaskInput{},
	}
}

// AddScheduleTaskInput is input for AddTaskToRunList func.
type AddScheduleTaskInput struct {
	ScheduleTaskName string // Unique name for this scheduled task
	TaskName         string // TaskName is task to run (name from task descriptor)
	Interval         *int   // Interval to run task with
	DataFn           func() TaskInput
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
		SetInterval(*input.Interval).
		OnConflict(
			sql.ConflictColumns("schedule_name"),
			sql.ResolveWith(func(us *sql.UpdateSet) {
				columns := us.Columns()
				for _, c := range columns {
					us.SetExcluded(c)
				}
				us.SetIgnore(queuescheduledtask.FieldLastRunTimestamp)
			}),
		).
		Exec(context.Background())
	if err != nil {
		panic(err)
	}
	if input.DataFn != nil {
		ts.dataFuncs[input.ScheduleTaskName] = input.DataFn
	} else {
		ts.dataFuncs[input.ScheduleTaskName] = nil
	}
}

// TODO add `for update skip locked` ?
// SchedulerGetTasksToRun returns list of tasks to run at the moment.
func (ts *TaskScheduler) SchedulerGetTasksToRun() []RunRequest {
	ts.logger.Debug("Check scheduler tasks")
	pendingTasks, err := ts.db.QueueScheduledTask.Query().All(context.Background())
	if err != nil && !ent.IsNotFound(err) {
		panic(err)
	}
	timeNow := time.Now()
	tasksToExecute := []*ent.QueueScheduledTask{}
	for _, task := range pendingTasks {
		duration := timeNow.Sub(task.LastRunTimestamp)
		if duration > time.Duration(task.Interval)*time.Millisecond {
			tasksToExecute = append(tasksToExecute, task)
		}
	}
	result := []RunRequest{}
	for _, task := range tasksToExecute {
		result = append(result, RunRequest{
			Name:   task.TaskName,
			Data:   nil,
			DataFn: ts.dataFuncs[task.ScheduleName],
		})
	}
	ts.logger.Debug("Check scheduler tasks. Done", zap.Int("Result items count", len(result)))
	return result
}

func (ts *TaskScheduler) RequestTasksRun(tasks []RunRequest) {
	for _, task := range tasks {
		runningTasksCount, err := ts.db.QueueTask.Query().
			Where(
				queuetask.StatusIn(string(TaskStatusRunning), string(TaskStatusPending)),
			).
			Count(context.Background())
		if err != nil {
			panic(err)
		}
		// If there is a running instance of this task already - do not post new task
		if runningTasksCount != 0 {
			continue
		}
		var data TaskInput
		data = nil
		if task.DataFn != nil {
			data = task.DataFn()
		}
		err = ts.db.QueueTask.Create().
			SetTaskName(task.Name).
			SetStatus(string(TaskStatusPending)).
			SetData(data).
			Exec(context.Background())
		ts.logger.Debug("Created queue task", zap.String("TaskName", task.Name))
		if err != nil {
			panic(err)
		}
		err = ts.db.QueueScheduledTask.Update().
			SetLastRunTimestamp(time.Now()).
			Where(queuescheduledtask.TaskName(task.Name)).
			Exec(context.Background())
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
