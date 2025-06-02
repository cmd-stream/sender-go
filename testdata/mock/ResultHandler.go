package mock

import (
	"github.com/cmd-stream/core-go"
	"github.com/ymz-ncnk/mok"
)

type HandleFn func(result core.Result, err error) error

func NewResultHandler() ResultHandler {
	return ResultHandler{
		Mock: mok.New("ResultHandler"),
	}
}

type ResultHandler struct {
	*mok.Mock
}

func (h ResultHandler) RegisterHandle(fn HandleFn) ResultHandler {
	h.Register("Handle", fn)
	return h
}

func (h ResultHandler) Handle(result core.Result, err error) error {
	var mokResult []interface{}
	mokResult, err = h.Call("Handle", mok.SafeVal[core.Result](result), mok.SafeVal[error](err))
	if err != nil {
		panic(err)
	}
	aerr, _ := mokResult[0].(error)
	return aerr
}
