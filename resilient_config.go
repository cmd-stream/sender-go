package sender

import (
	"time"

	cln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	dcln "github.com/cmd-stream/delegate-go/client"
	hks "github.com/cmd-stream/sender-go/hooks"
	"github.com/ymz-ncnk/circbrk-go"
)

// ResilientConfig holds configuration for creating a resilient Sender, that
// automatically handles keepalive, reconnects, and circuit breaker behavior,
// also allows optional hooks to observe or modify sender behavior.
type ResilientConfig[T any] struct {
	// KeepaliveTime is the duration after which a keepalive probe is sent
	// if no commands have been sent. Must be non-zero.
	// Recommended value: 30 * time.Second
	KeepaliveTime time.Duration

	// KeepaliveIntvl is the interval between keepalive probes after the first one.
	// Must be non-zero.
	// Recommended value: 10 * time.Second
	KeepaliveIntvl time.Duration

	// CircuitBreakerWindowSize is the number of requests in the sliding window
	// used to calculate failure rate. Must be non-zero.
	// Recommended value: 20
	CircuitBreakerWindowSize int

	// CircuitBreakerFailureRate is the failure rate threshold (0.0–1.0)
	// that triggers the circuit breaker to open. Must be non-zero.
	// Recommended value: 0.5
	CircuitBreakerFailureRate float64

	// CircuitBreakerOpenDuration is the duration the circuit breaker stays open
	// before transitioning to half-open. Must be non-zero.
	// Recommended value: 30 * time.Second
	CircuitBreakerOpenDuration time.Duration

	// CircuitBreakerSuccessThreshold is the number of successful requests in
	// half-open state required to close the circuit breaker. Must be non-zero.
	// Recommended value: 2
	CircuitBreakerSuccessThreshold int

	// CircuitBreakerChangeStateCallback is an optional callback invoked when
	// the circuit breaker changes state.
	CircuitBreakerChangeStateCallback circbrk.ChangeStateCallback

	// HooksFactory provides hooks that can observe or modify sender behavior.
	// If nil, a NoopHooksFactory is used by default.
	HooksFactory hks.HooksFactory[T]
}

// ToOptions converts the ResilientConfig into a slice of SetMakeOption[T]
// that can be passed to a sender constructor. It panics if required fields
// are not set.
func (c ResilientConfig[T]) ToOptions() []SetMakeOption[T] {
	if c.KeepaliveTime == 0 {
		panic("KeepaliveTime is required")
	}
	if c.KeepaliveIntvl == 0 {
		panic("KeepaliveIntvl is required")
	}
	if c.CircuitBreakerWindowSize == 0 {
		panic("CircuitBreakerWindowSize is required")
	}
	if c.CircuitBreakerFailureRate == 0 {
		panic("CircuitBreakerFailureRate is required")
	}
	if c.CircuitBreakerOpenDuration == 0 {
		panic("CircuitBreakerOpenDuration is required")
	}
	if c.CircuitBreakerSuccessThreshold == 0 {
		panic("CircuitBreakerSuccessThreshold is required")
	}
	if c.HooksFactory == nil {
		c.HooksFactory = hks.NoopHooksFactory[T]{}
	}
	return []SetMakeOption[T]{
		WithGroup(
			grp.WithReconnect[T](),
			grp.WithClient[T](
				cln.WithKeepalive(
					dcln.WithKeepaliveTime(c.KeepaliveTime),
					dcln.WithKeepaliveIntvl(c.KeepaliveIntvl),
				),
			),
		),
		WithSender(
			WithHooksFactory(c.hooksFactory()),
		),
	}
}

func (c ResilientConfig[T]) hooksFactory() hks.CircuitBreakerHooksFactory[T] {
	cb := circbrk.New(circbrk.WithWindowSize(c.CircuitBreakerWindowSize),
		circbrk.WithFailureRate(c.CircuitBreakerFailureRate),
		circbrk.WithOpenDuration(c.CircuitBreakerOpenDuration),
		circbrk.WithSuccessThreshold(c.CircuitBreakerSuccessThreshold),
		circbrk.WithChangeStateCallback(c.CircuitBreakerChangeStateCallback),
	)
	return hks.NewCircuitBreakerHooksFactory(cb, c.HooksFactory)
}
