package mocks

import (
	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/ymz-ncnk/mok"
)

type NewFn[T any] func(clients []grp.Client[T]) (strategy grp.DispatchStrategy[grp.Client[T]])

func NewDispatchStrategyFactory[T any]() DispatchStrategyFactory[T] {
	return DispatchStrategyFactory[T]{
		Mock: mok.New("DispatchStrategyFactory"),
	}
}

type DispatchStrategyFactory[T any] struct {
	*mok.Mock
}

func (f DispatchStrategyFactory[T]) RegisterNew(fn NewFn[T]) DispatchStrategyFactory[T] {
	f.Register("New", fn)
	return f
}

func (f DispatchStrategyFactory[T]) New(clients []grp.Client[T]) (
	strategy grp.DispatchStrategy[grp.Client[T]],
) {
	result, err := f.Call("Read", clients)
	if err != nil {
		panic(err)
	}

	strategy, _ = result[0].(grp.DispatchStrategy[grp.Client[T]])
	return
}
