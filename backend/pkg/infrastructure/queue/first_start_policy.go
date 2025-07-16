package queue

// FirstStartPolicy for enum of restart policies
type FirstStartPolicyType string

// FirstStartPolicy enum
const (
	FirstStartPolicyImmediate FirstStartPolicyType = "immediate" // Start task right after scheduling it
	FirstStartPolicyInterval  FirstStartPolicyType = "interval"  // Start task after an interval
	FirstStartPolicySchedule  FirstStartPolicyType = "schedule"  // Start task depending on schedule value
)
