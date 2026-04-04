package hooks

// CircuitBreaker defines the interface for the Circuit Breaker Pattern.
//
// Deprecated: use github.com/cmd-stream/cmd-stream-go instead.

type CircuitBreaker interface {
	Allow() bool
	Fail()
	Success()
}
