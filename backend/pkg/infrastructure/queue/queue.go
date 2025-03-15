package queue

import (
	"context"
	"time"
)

// RestartPolicyType for enum of restart policies
type RestartPolicyType int

// Restart Policy enum
const (
	RestartPolicyPanic RestartPolicyType = 1
	RestartPolicyError RestartPolicyType = 2
)

// RestartPolicy
// Restart Count = 0 - infinite restart
// RestartOnError - default is false. Doesn't restart on panic in any case
// OnRestartCountExceeded - reaction to end of restart count. Default is RestartPolicyError
type RestartPolicy struct {
	RestartCount           int
	RestartOnError         bool
	OnRestartCountExceeded RestartPolicyType
}

// Task is description of task to run
type Task struct {
	RestartPolicy RestartPolicy              // RestartPolicy for the task
	Name          string                     // Unique name for the task
	Func          func(map[string]any) error // Function to run in task
}

// TaskRun is single running task
type TaskRun struct {
	RestartPolicy RestartPolicy              // RestartPolicy for the task
	Func          func(map[string]any) error // Function to run in task
	Data          map[string]any             // Data to run task with
	ctx           context.Context            // Context for task
}

type SheduledTask struct {
	Task     Task
	Interval time.Duration
}

// QueueManager - global struct to manage all running tasks.
type QueueManager struct {
	sheduledTasks []SheduledTask
	queueTasks    []Task
}

// TODO check for unique name

func (qm *QueueManager) AddQueueTask(task Task) {
	qm.queueTasks = append(qm.queueTasks, task)
}

func (qm *QueueManager) AddSheduledTask(task Task, interval time.Duration) {
	qm.sheduledTasks = append(qm.sheduledTasks, SheduledTask{
		Task:     task,
		Interval: interval,
	})
}
