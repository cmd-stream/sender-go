package sender

import "github.com/cmd-stream/sender-go/hooks"

// Deprecated: use github.com/cmd-stream/cmd-stream-go instead.
type Options[T any] struct {

	HooksFactory hooks.HooksFactory[T]
}

// Deprecated: use github.com/cmd-stream/cmd-stream-go instead.
type SetOption[T any] func(o *Options[T])


// logging or instrumentation.
//
// Deprecated: use github.com/cmd-stream/cmd-stream-go instead.

func WithHooksFactory[T any](factory hooks.HooksFactory[T]) SetOption[T] {
	return func(o *Options[T]) {
		o.HooksFactory = factory
	}
}

func Apply[T any](ops []SetOption[T], o *Options[T]) {
	for i := range ops {
		if ops[i] != nil {
			ops[i](o)
		}
	}
}
