package hooks

import (
	"context"

	"github.com/cmd-stream/core-go"
)

// HooksFactory provides a way to create new Hooks instances.
type HooksFactory[T any] interface {
	New() Hooks[T]
}

// Hooks defines an interface for customizing behavior during the command
// sending process. Implementations can provide hooks for events such as
// BeforeSend, OnError, OnResult, and OnTimeout.
type Hooks[T any] interface {
	BeforeSend(ctx context.Context, cmd core.Cmd[T]) (context.Context, error)
	OnError(ctx context.Context, sentCmd SentCmd[T], err error)
	OnResult(ctx context.Context, sentCmd SentCmd[T], recvResult ReceivedResult,
		err error)
	OnTimeout(ctx context.Context, sentCmd SentCmd[T], err error)
}
