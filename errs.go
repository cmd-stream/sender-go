package sender

import "errors"

// ErrTimeout is returned when a command is sent but no result is received
// within the expected time.
//
// Deprecated: use github.com/cmd-stream/cmd-stream-go instead.

var ErrTimeout = errors.New("timeout")
