package mock

import (
	"github.com/ymz-ncnk/mok"
)

type (
	AllowFn   func() bool
	FailFn    func()
	SuccessFn func()
)

func NewCircuitBreaker() CircuitBreaker {
	return CircuitBreaker{
		Mock: mok.New("CircuitBreaker"),
	}
}

type CircuitBreaker struct {
	*mok.Mock
}

func (c CircuitBreaker) RegisterAllow(fn AllowFn) CircuitBreaker {
	c.Register("Allow", fn)
	return c
}

func (c CircuitBreaker) RegisterFail(fn FailFn) CircuitBreaker {
	c.Register("Fail", fn)
	return c
}

func (c CircuitBreaker) RegisterSuccess(fn SuccessFn) CircuitBreaker {
	c.Register("Success", fn)
	return c
}

func (c CircuitBreaker) Allow() bool {
	result, err := c.Call("Allow")
	if err != nil {
		panic(err)
	}
	return result[0].(bool)
}

func (c CircuitBreaker) Fail() {
	_, err := c.Call("Fail")
	if err != nil {
		panic(err)
	}
}

func (c CircuitBreaker) Success() {
	_, err := c.Call("Success")
	if err != nil {
		panic(err)
	}
}
