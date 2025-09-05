package sender

import "github.com/cmd-stream/sender-go/hooks"

type Options[T any] struct {
	HooksFactory hooks.HooksFactory[T]
}

type SetOption[T any] func(o *Options[T])

// WithHooksFactory sets a factory that creates new hooks for each send
// operation. Hooks can customize behavior during the sending process, such as
// logging or instrumentation.
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
