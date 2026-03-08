package mocks

import (
	hks "github.com/cmd-stream/sender-go/hooks"
	"github.com/ymz-ncnk/mok"
)

type HooksFactoryNewFn[T any] func() hks.Hooks[T]

func NewHooksFactory[T any]() *HooksFactory[T] {
	return &HooksFactory[T]{
		Mock: new(mok.Mock),
	}
}

type HooksFactory[T any] struct {
	*mok.Mock
}

func (m HooksFactory[T]) RegisterNew(fn HooksFactoryNewFn[T]) HooksFactory[T] {
	m.Register("New", fn)
	return m
}

func (m HooksFactory[T]) New() (hooks hks.Hooks[T]) {
	result, err := m.Call("New")
	if err != nil {
		panic(err)
	}

	hooks, _ = result[0].(Hooks[T])
	return
}
