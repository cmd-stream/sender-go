package hooks

import (
	"context"

	"github.com/cmd-stream/core-go"
)

// NewCircuitBreakerHooksFactory creates a new CircuitBreakerHooksFactory.
func NewCircuitBreakerHooksFactory[T any](cb CircuitBreaker,
	factory HooksFactory[T]) CircuitBreakerHooksFactory[T] {
	return CircuitBreakerHooksFactory[T]{cb, factory}
}

// CircuitBreakerHooksFactory can be used to create hooks that incorporate
// circuit breaker logic during the command sending process.
type CircuitBreakerHooksFactory[T any] struct {
	cb      CircuitBreaker
	factory HooksFactory[T]
}

func (f CircuitBreakerHooksFactory[T]) New() Hooks[T] {
	return NewCircuitBreakerHooks(f.cb, f.factory.New())
}

// NewCircuitBreakerHooks creates a new CircuitBreakerHooks.
func NewCircuitBreakerHooks[T any](cb CircuitBreaker,
	hooks Hooks[T]) CircuitBreakerHooks[T] {
	return CircuitBreakerHooks[T]{cb, hooks}
}

// CircuitBreakerHooks checks whether the circuit breaker is open before sending.
// If so, it returns ErrCircuitOpen, otherwise the corresponding method of the
// inner Hooks is called.
type CircuitBreakerHooks[T any] struct {
	cb    CircuitBreaker
	hooks Hooks[T]
}

func (h CircuitBreakerHooks[T]) BeforeSend(ctx context.Context, cmd core.Cmd[T]) (
	context.Context, error) {
	if h.cb.Open() {
		return ctx, ErrCircuitOpen
	}
	return h.hooks.BeforeSend(ctx, cmd)
}

func (h CircuitBreakerHooks[T]) OnError(ctx context.Context, sentCmd SentCmd[T],
	err error) {
	h.cb.Fail()
	h.hooks.OnError(ctx, sentCmd, err)
}

func (h CircuitBreakerHooks[T]) OnResult(ctx context.Context, sentCmd SentCmd[T],
	recvResult ReceivedResult, err error) {
	h.cb.Success()
	h.hooks.OnResult(ctx, sentCmd, recvResult, err)
}

func (h CircuitBreakerHooks[T]) OnTimeout(ctx context.Context, sentCmd SentCmd[T],
	err error) {
	h.cb.Fail()
	h.hooks.OnTimeout(ctx, sentCmd, err)
}
