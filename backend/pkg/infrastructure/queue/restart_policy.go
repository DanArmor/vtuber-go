package queue

// RestartPolicyType for enum of restart policies
type RestartPolicyType string

// Restart Policy enum
const (
	RestartPolicyPanic RestartPolicyType = "panic"
	RestartPolicyError RestartPolicyType = "error"
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
