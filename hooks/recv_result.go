package hooks

import "github.com/cmd-stream/core-go"

// Deprecated: use github.com/cmd-stream/cmd-stream-go instead.
type ReceivedResult struct {

	Seq    core.Seq
	Size   int
	Result core.Result
}
