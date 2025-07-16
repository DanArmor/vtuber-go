package queue

import (
	"context"
	"log"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/ent/queuescheduledtask"
	"go.uber.org/zap"
)

// QueueManager - global struct to manage all running tasks.
type QueueManager struct {
	db      *ent.Client
	taskMap map[string]Task
	logger  *zap.Logger
}

func NewQueueManager(db *ent.Client, logger *zap.Logger) *QueueManager {
	return &QueueManager{
		db:      db,
		taskMap: map[string]Task{},
		logger:  logger.With(zap.String("Service", "QueueManager")),
	}
}

// RegisterTask registers task in global list of available tasks
func (qm *QueueManager) RegisterTask(task Task) {
	qm.logger.Debug("Register task", zap.String("Task", task.Name))
	qm.taskMap[task.Name] = task
}

// getTaskInfo returns task info by task name if exists
// panic otherwise
func (qm *QueueManager) getTaskInfo(name string) Task {
	qm.logger.Debug("Get task info", zap.String("Task", name))
	task, ok := qm.taskMap[name]
	if !ok {
		panic("Not found in task map " + name)
	}
	return task
}

// AddTaskToRunListInput is input for AddTaskToRunList func
type AddTaskToRunListInput struct {
	Name             string
	RestartPolicy    *RestartPolicy
	FirstStartPolicy *FirstStartPolicyType
	Interval         *int
}

// AddTaskToRunList adds task to execution
func (qm *QueueManager) AddTaskToRunList(input AddTaskToRunListInput) {
	qm.logger.Debug("Add task to run list", zap.String("Task", input.Name))
	task := qm.getTaskInfo(input.Name)
	// Do not change any values of 'input' after this block
	if input.RestartPolicy == nil {
		input.RestartPolicy = &task.RestartPolicy
	}
	if input.FirstStartPolicy == nil {
		input.FirstStartPolicy = &task.FirstStartPolicy
	}
	if input.Interval == nil {
		input.Interval = &task.Interval
	}
	err := qm.db.QueueScheduledTask.Create().
		SetName(input.Name).
		SetStatus(string(TaskStatusPending)).
		SetInterval(*input.Interval).
		OnConflict(sql.ConflictColumns("name"), sql.ResolveWithNewValues()).
		Exec(context.Background())
	if err != nil {
		panic(err)
	}
}

// CheckScheduler returns list of tasks to run at the moment
func (qm *QueueManager) CheckScheduler() []Task {
	qm.logger.Debug("Check scheduler tasks")
	pendingTasks, err := qm.db.QueueScheduledTask.Query().Where(
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
			tasksNamesToExecute = append(tasksNamesToExecute, task.Name)
		}
	}
	result := []Task{}
	for _, name := range tasksNamesToExecute {
		result = append(result, qm.getTaskInfo(name))
	}
	qm.logger.Debug("Check scheduler tasks. Done", zap.Int("Result items count", len(result)))
	return result
}

func (qm *QueueManager) ExecuteTask(task Task) error {
	qm.logger.Info("Executing task", zap.String("Task", task.Name))
	currentTaskInfo, err := qm.db.QueueScheduledTask.Query().
		Where(
			queuescheduledtask.And(
				queuescheduledtask.NameEQ(task.Name),
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
	qm.logger.Debug("Executing task. Locked task", zap.String("Task", task.Name))

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
	qm.logger.Debug("Executing task. Unlocked task", zap.String("Task", task.Name))
	return nil
}

// Run starts QueueManager in infinite loop of check-run tasks
func (qm *QueueManager) Run(ctx context.Context) error {
	// Vars for runtime behaviour
	interval := time.Duration(time.Minute)
	counter := 1
	quit := make(chan int, 1)

	for {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Recovered: %v", err)
					quit <- 1
				}
			}()
			ticker := time.NewTicker(interval)

			for {
				select {
				case <-ticker.C:
					qm.logger.Debug("Queue manager tick start")
					tasks := qm.CheckScheduler()
					for _, task := range tasks {
						go qm.ExecuteTask(task)
					}
					qm.logger.Debug("Queue manager tick end")
				}
			}
		}()
		<-quit
		counter += 1
		if counter > 10 {
			panic("QueueManager restart limit exceeded")
		}
	}
}
