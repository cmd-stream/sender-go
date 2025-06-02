package mock

import (
	"time"

	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/core-go"
	"github.com/ymz-ncnk/mok"
)

type SendFn func(cmd core.Cmd[any], results chan<- core.AsyncResult) (
	seq core.Seq, clientID grp.ClientID, n int, err error)
type SendWithDeadlineFn func(cmd core.Cmd[any], results chan<- core.AsyncResult,
	deadline time.Time,
) (seq core.Seq, clientID grp.ClientID, n int, err error)
type HasFn func(seq core.Seq, clientID grp.ClientID) (ok bool)
type ForgetFn func(seq core.Seq, clientID grp.ClientID)
type DoneFn func() <-chan struct{}
type ErrFn func() error
type CloseFn func() error

func NewClientGroup() ClientGroup {
	return ClientGroup{
		Mock: mok.New("ClientGroup"),
	}
}

type ClientGroup struct {
	*mok.Mock
}

func (g ClientGroup) RegisterSend(fn SendFn) ClientGroup {
	g.Register("Send", fn)
	return g
}

func (g ClientGroup) RegisterSendWithDeadline(fn SendWithDeadlineFn) ClientGroup {
	g.Register("SendWithDeadline", fn)
	return g
}

func (g ClientGroup) RegisterHas(fn HasFn) (ok bool) {
	g.Register("Has", fn)
	return
}

func (g ClientGroup) RegisterForget(fn ForgetFn) ClientGroup {
	g.Register("Forget", fn)
	return g
}

func (g ClientGroup) RegisterDone(fn DoneFn) ClientGroup {
	g.Register("Done", fn)
	return g
}

func (g ClientGroup) RegisterErr(fn ErrFn) ClientGroup {
	g.Register("Err", fn)
	return g
}

func (g ClientGroup) RegisterClose(fn CloseFn) ClientGroup {
	g.Register("Close", fn)
	return g
}

func (g ClientGroup) Send(cmd core.Cmd[any],
	results chan<- core.AsyncResult,
) (seq core.Seq, clientID grp.ClientID, n int, err error) {
	result, err := g.Call("Send", mok.SafeVal[core.Cmd[any]](cmd), results)
	if err != nil {
		panic(err)
	}

	seq = result[0].(core.Seq)
	clientID = result[1].(grp.ClientID)
	n = result[2].(int)
	err, _ = result[3].(error)
	return
}

func (g ClientGroup) SendWithDeadline(cmd core.Cmd[any],
	results chan<- core.AsyncResult,
	deadline time.Time,
) (seq core.Seq, clientID grp.ClientID, n int, err error) {
	result, err := g.Call("SendWithDeadline", cmd, results, deadline)
	if err != nil {
		panic(err)
	}

	seq = result[0].(core.Seq)
	clientID = result[1].(grp.ClientID)
	n = result[2].(int)
	err, _ = result[3].(error)
	return
}

func (g ClientGroup) Has(seq core.Seq, clientID grp.ClientID) (ok bool) {
	result, err := g.Call("Has", seq, clientID)
	if err != nil {
		panic(err)
	}

	ok, _ = result[0].(bool)
	return
}

func (g ClientGroup) Forget(seq core.Seq, clientID grp.ClientID) {
	_, err := g.Call("Forget", seq, clientID)
	if err != nil {
		panic(err)
	}
}

func (g ClientGroup) Done() (ch <-chan struct{}) {
	result, err := g.Call("Done")
	if err != nil {
		panic(err)
	}

	ch, _ = result[0].(<-chan struct{})
	return
}

func (g ClientGroup) Err() (err error) {
	result, err := g.Call("Err")
	if err != nil {
		panic(err)
	}
	err, _ = result[0].(error)
	return
}

func (g ClientGroup) Close() (err error) {
	result, err := g.Call("Close")
	if err != nil {
		panic(err)
	}

	err, _ = result[0].(error)
	return
}
