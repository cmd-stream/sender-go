package mocks

import (
	"context"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/sender-go/hooks"
	"github.com/ymz-ncnk/mok"
)

type (
	BeforeSendFn[T any] func(ctx context.Context, cmd core.Cmd[T]) (context.Context, error)
	OnErrorFn[T any]    func(ctx context.Context, sentCmd hooks.SentCmd[T], err error)
	OnResultFn[T any]   func(ctx context.Context, sentCmd hooks.SentCmd[T],
		recvResult hooks.ReceivedResult, err error)
)
type OnTimeoutFn[T any] func(ctx context.Context, sentCmd hooks.SentCmd[T], err error)

func NewHooks[T any]() *Hooks[T] {
	return &Hooks[T]{
		Mock: mok.New("Hooks"),
	}
}

type Hooks[T any] struct {
	*mok.Mock
}

func (h Hooks[T]) RegisterBeforeSend(fn BeforeSendFn[T]) Hooks[T] {
	h.Register("BeforeSend", fn)
	return h
}

func (h Hooks[T]) RegisterOnError(fn OnErrorFn[T]) Hooks[T] {
	h.Register("OnError", fn)
	return h
}

func (h Hooks[T]) RegisterOnResult(fn OnResultFn[T]) Hooks[T] {
	h.Register("OnResult", fn)
	return h
}

func (h Hooks[T]) RegisterOnTimeout(fn OnTimeoutFn[T]) Hooks[T] {
	h.Register("OnTimeout", fn)
	return h
}

func (h Hooks[T]) BeforeSend(ctx context.Context, cmd core.Cmd[T]) (
	actx context.Context, err error,
) {
	result, err := h.Call("BeforeSend", ctx, mok.SafeVal[core.Cmd[T]](cmd))
	if err != nil {
		panic(err)
	}

	actx, _ = result[0].(context.Context)
	err, _ = result[1].(error)
	return
}

func (h Hooks[T]) OnError(ctx context.Context, sentCmd hooks.SentCmd[T],
	err error,
) {
	_, err = h.Call("OnError", ctx, sentCmd, err)
	if err != nil {
		panic(err)
	}
}

func (h Hooks[T]) OnResult(ctx context.Context, sentCmd hooks.SentCmd[T],
	recvResult hooks.ReceivedResult, err error,
) {
	_, err = h.Call("OnResult", ctx, sentCmd, recvResult, mok.SafeVal[error](err))
	if err != nil {
		panic(err)
	}
}

func (h Hooks[T]) OnTimeout(ctx context.Context, sentCmd hooks.SentCmd[T],
	err error,
) {
	_, err = h.Call("OnTimeout", ctx, sentCmd, err)
	if err != nil {
		panic(err)
	}
}
