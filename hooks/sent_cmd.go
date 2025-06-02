package hooks

import "github.com/cmd-stream/core-go"

type SentCmd[T any] struct {
	Seq  core.Seq
	Size int
	Cmd  core.Cmd[T]
}
