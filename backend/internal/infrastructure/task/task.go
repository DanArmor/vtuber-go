package task

import "context"

type TaskInput map[string]any
type TaskFn func(context.Context, TaskInput) error

// TaskDescriptor is description of task to run.
type TaskDescriptor struct {
	Name string // Unique name for the task
	Func TaskFn // Function to run in task
}

// RunRequest is structed used to request work in queue.
type RunRequest struct {
	ID int // Unique id of the task from DB.
	// ! Worker should have valid value of ID field upon the execution of the task
	Name   string           // Unique name for the task
	Data   TaskInput        // For worker. Data to run task with
	DataFn func() TaskInput // For scheduler. To populate data on task creation
}

// Run is single running task.
type Run struct {
	Task TaskDescriptor  // Task to run
	Data TaskInput       // Data to run task with
	Ctx  context.Context // Context for task
}
