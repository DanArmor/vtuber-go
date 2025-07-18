package task

import (
	"go.uber.org/zap"
)

var glTaskMap = map[string]TaskDescriptor{}

// RegisterTask registers task in global list of available tasks
func RegisterTask(task TaskDescriptor) {
	// TODO thing about these global loggers. How to replace this?
	zap.L().Debug("Register task", zap.String("Task", task.Name))
	glTaskMap[task.Name] = task
}

// getTaskInfo returns task info by task name if exists
// panic otherwise
func getTaskInfo(name string) TaskDescriptor {
	zap.L().Debug("Get task info", zap.String("Task", name))
	task, ok := glTaskMap[name]
	if !ok {
		panic("Not found in task map " + name)
	}
	return task
}
