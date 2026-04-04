package sender

import (
	"time"

	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/core-go"
)

// ClientGroup represents a group of clients used to send commands and receive
// results.
//
// Deprecated: use github.com/cmd-stream/cmd-stream-go instead.

type ClientGroup[T any] interface {
	Send(cmd core.Cmd[T], results chan<- core.AsyncResult) (
		seq core.Seq, clientID grp.ClientID, n int, err error)
	SendWithDeadline(cmd core.Cmd[T], results chan<- core.AsyncResult,
		deadline time.Time,
	) (seq core.Seq, clientID grp.ClientID, n int, err error)
	Has(seq core.Seq, clientID grp.ClientID) (ok bool)
	Forget(seq core.Seq, clientID grp.ClientID)
	Done() <-chan struct{}
	Err() (err error)
	Close() (err error)
}
