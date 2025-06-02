package sender

import "github.com/cmd-stream/core-go"

// ResultHandler is an interface that wraps the Handle method. It is used to
// process multiple results received from the server.
type ResultHandler interface {
	Handle(result core.Result, err error) error
}

// ResultHandlerFn is a function type that implements the ResultHandler
// interface.
type ResultHandlerFn func(result core.Result, err error) error

func (fn ResultHandlerFn) Handle(result core.Result, err error) error {
	return fn(result, err)
}
