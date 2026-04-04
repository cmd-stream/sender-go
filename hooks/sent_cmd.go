package hooks

import "github.com/cmd-stream/core-go"

// Deprecated: use github.com/cmd-stream/cmd-stream-go instead.
type SentCmd[T any] struct {

	Seq  core.Seq
	Size int
	Cmd  core.Cmd[T]
}
