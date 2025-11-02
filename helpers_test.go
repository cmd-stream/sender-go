package sender_test

import (
	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/core-go"
	cmocks "github.com/cmd-stream/testkit-go/mocks/core"
)

type Want struct {
	Cmd     cmocks.Cmd
	Results []WantResult

	CmdSeq     core.Seq
	ClientID   grp.ClientID
	CmdSize    int
	CmdSendErr error

	Err error
}

type WantResult struct {
	Seq       core.Seq
	Result    core.Result
	BytesRead int
	Err       error
}
