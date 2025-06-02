package hooks

// CircuitBreaker defines the interface for the Circuit Breaker pattern.
type CircuitBreaker interface {
	Open() bool
	Fail()
	Success()
}
