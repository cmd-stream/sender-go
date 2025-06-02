package hooks

import "github.com/cmd-stream/core-go"

type ReceivedResult struct {
	Seq    core.Seq
	Size   int
	Result core.Result
}
