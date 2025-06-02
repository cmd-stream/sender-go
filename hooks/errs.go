package hooks

import "errors"

// ErrCircuitOpen is returned when an operation is attempted while the circuit
// breaker is in the open state.
var ErrCircuitOpen = errors.New("circuit breaker is open")
