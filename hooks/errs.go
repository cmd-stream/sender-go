package hooks

import "errors"

// ErrNotAllowed indicates that sending the Command is not allowed at this
// time.
//
// Deprecated: use github.com/cmd-stream/cmd-stream-go instead.

var ErrNotAllowed = errors.New("not allowed")
