package testdata

import (
	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/core-go"
	cmock "github.com/cmd-stream/core-go/testdata/mock"
)

type Want struct {
	Cmd     cmock.Cmd
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
