package queue

import "context"

// Task is description of task to run
type Task struct {
	RestartPolicy    RestartPolicy        // RestartPolicy for the task
	FirstStartPolicy FirstStartPolicyType // FirstStartPolicy for the task
	// ! RestartPolicy and FirstStartPolicy are not used at the moment.
	// ! default restart policy is fallback to panic
	// ! default first start policy is 'schedule' policy
	Interval int                        // Interval for the task (in mc)
	Name     string                     // Unique name for the task
	Func     func(map[string]any) error // Function to run in task
}

// TaskRun is single running task
type TaskRun struct {
	Task Task           // Task to run
	Data map[string]any // Data to run task with
	// ! Data  field in unused at the moment
	Ctx context.Context // Context for task
}
