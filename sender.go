package sender

import (
	"context"
	"time"

	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/core-go"
	hks "github.com/cmd-stream/sender-go/hooks"
)

// New creates a new Sender with the given client group and optional hooks.
func New[T any](group ClientGroup[T], ops ...SetOption[T]) Sender[T] {
	o := Options[T]{
		HooksFactory: hks.NoopHooksFactory[T]{},
	}
	Apply(ops, &o)

	return Sender[T]{
		group:   group,
		options: o,
	}
}

// Sender provides a high-level abstraction over a client group for sending
// commands to the server.
type Sender[T any] struct {
	group   ClientGroup[T]
	options Options[T]
}

// Send sends a command to the server and waits (using the ctx) for the result.
func (s Sender[T]) Send(ctx context.Context, cmd core.Cmd[T]) (
	result core.Result, err error) {
	var (
		results = make(chan core.AsyncResult, 1)
		hooks   = s.options.HooksFactory.New()
	)
	ctx, err = hooks.BeforeSend(ctx, cmd)
	if err != nil {
		return
	}
	seq, clientID, n, err := s.group.Send(cmd, results)
	sentCmd := hks.SentCmd[T]{
		Seq:  seq,
		Size: n,
		Cmd:  cmd,
	}
	if err != nil {
		hooks.OnError(ctx, sentCmd, err)
		return
	}
	return s.receive(ctx, sentCmd, results, clientID, hooks)
}

// Send sends a command to the server with the specified deadline and waits
// (using the ctx) for the result.
func (s Sender[T]) SendWithDeadline(ctx context.Context,
	cmd core.Cmd[T], dealine time.Time) (result core.Result, err error) {
	var (
		results = make(chan core.AsyncResult, 1)
		hooks   = s.options.HooksFactory.New()
	)
	ctx, err = hooks.BeforeSend(ctx, cmd)
	if err != nil {
		return
	}
	seq, clientID, n, err := s.group.SendWithDeadline(cmd, results, dealine)
	sentCmd := hks.SentCmd[T]{
		Seq:  seq,
		Size: n,
		Cmd:  cmd,
	}
	if err != nil {
		hooks.OnError(ctx, sentCmd, err)
		return
	}
	return s.receive(ctx, sentCmd, results, clientID, hooks)
}

// Send sends a command to the server and waits (using the ctx) for multiple
// results.
func (s Sender[T]) SendMulti(ctx context.Context, cmd core.Cmd[T],
	resultsCount int, handler ResultHandler) (err error) {
	var (
		results = make(chan core.AsyncResult, resultsCount)
		hooks   = s.options.HooksFactory.New()
	)
	ctx, err = hooks.BeforeSend(ctx, cmd)
	if err != nil {
		return
	}
	seq, clientID, n, err := s.group.Send(cmd, results)
	sentCmd := hks.SentCmd[T]{
		Seq:  seq,
		Size: n,
		Cmd:  cmd,
	}
	if err != nil {
		hooks.OnError(ctx, sentCmd, err)
		return
	}
	s.receiveMulti(ctx, sentCmd, results, clientID, hooks, handler)
	return
}

// Send sends a command to the server with the specified deadline and waits
// (using the ctx) for multiple results.
func (s Sender[T]) SendMultiWithDeadline(ctx context.Context,
	cmd core.Cmd[T],
	resultsCount int,
	handler ResultHandler,
	dealine time.Time,
) (err error) {
	var (
		results = make(chan core.AsyncResult, resultsCount)
		hooks   = s.options.HooksFactory.New()
	)
	ctx, err = hooks.BeforeSend(ctx, cmd)
	if err != nil {
		return
	}
	seq, clientID, n, err := s.group.SendWithDeadline(cmd, results, dealine)
	sentCmd := hks.SentCmd[T]{
		Seq:  seq,
		Size: n,
		Cmd:  cmd,
	}
	if err != nil {
		hooks.OnError(ctx, sentCmd, err)
		return
	}
	s.receiveMulti(ctx, sentCmd, results, clientID, hooks, handler)
	return
}

// Close closes the underlying client group.
func (s Sender[T]) Close() error {
	return s.group.Close()
}

// Done returns a channel that is closed when the underlying client group is
// closed.
func (s Sender[T]) Done() <-chan struct{} {
	return s.group.Done()
}

func (s Sender[T]) receive(ctx context.Context, sentCmd hks.SentCmd[T],
	results <-chan core.AsyncResult,
	clientID grp.ClientID,
	hooks hks.Hooks[T],
) (result core.Result, err error) {
	select {
	case <-ctx.Done():
		err = ErrTimeout
		hooks.OnTimeout(ctx, sentCmd, err)
		s.group.Forget(sentCmd.Seq, clientID)
	case asyncResult := <-results:
		recvResult := hks.ReceivedResult{
			Seq:    core.Seq(1),
			Size:   asyncResult.BytesRead,
			Result: asyncResult.Result,
		}
		hooks.OnResult(ctx, sentCmd, recvResult, asyncResult.Error)
		result = asyncResult.Result
		err = asyncResult.Error
	}
	return
}

func (s Sender[T]) receiveMulti(ctx context.Context, sentCmd hks.SentCmd[T],
	results <-chan core.AsyncResult,
	clientID grp.ClientID,
	hooks hks.Hooks[T],
	handler ResultHandler,
) {
	var (
		result    core.Result
		handleErr error
		err       error
		i         = 1
	)
	for {
		select {
		case <-ctx.Done():
			err = ErrTimeout
			hooks.OnTimeout(ctx, sentCmd, err)
			s.group.Forget(sentCmd.Seq, clientID)
		case asyncResult := <-results:
			recvResult := hks.ReceivedResult{
				Seq:    core.Seq(i),
				Size:   asyncResult.BytesRead,
				Result: asyncResult.Result,
			}
			hooks.OnResult(ctx, sentCmd, recvResult, asyncResult.Error)
			result = asyncResult.Result
			err = asyncResult.Error
		}
		handleErr = handler.Handle(result, err)
		if handleErr != nil {
			err = handleErr
		}
		if err != nil || result.LastOne() {
			return
		}
		i++
	}
}
