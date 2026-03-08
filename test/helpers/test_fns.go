package helpers

import (
	"context"
	"errors"
	"testing"
	"time"

	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/core-go"
	cmock "github.com/cmd-stream/core-go/test/mock"
	hks "github.com/cmd-stream/sender-go/hooks"
	mock "github.com/cmd-stream/sender-go/test/mock"

	sndr "github.com/cmd-stream/sender-go"
	"github.com/ymz-ncnk/mok"

	asserterror "github.com/ymz-ncnk/assert/error"
)

type TestFn func(hooks mock.Hooks[any], factory mock.HooksFactory[any],
	group mock.ClientGroup,
	cmd cmock.Cmd,
	wantResult core.Result,
	wantErr error,
	t *testing.T,
)

func TestSuccess(group mock.ClientGroup, w Want, fn TestFn, t *testing.T) {
	var (
		hooks = mock.NewHooks[any]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
				asserterror.EqualDeep(t, cmd, w.Cmd)

				actx := context.WithoutCancel(ctx)
				return actx, nil
			},
		).RegisterOnResult(
			func(ctx context.Context, sentCmd hks.SentCmd[any],
				recvResult hks.ReceivedResult, err error,
			) {
				asserterror.EqualDeep(t, sentCmd, hks.SentCmd[any]{
					Seq:  w.CmdSeq,
					Size: w.CmdSize,
					Cmd:  w.Cmd,
				})
				asserterror.EqualDeep(t, recvResult, hks.ReceivedResult{
					Seq:    w.Results[0].Seq,
					Size:   w.Results[0].BytesRead,
					Result: w.Results[0].Result,
				})
				asserterror.EqualError(t, err, w.Results[0].Err)
			},
		)
		factory = mock.NewHooksFactory[any]().RegisterNew(
			func() hks.Hooks[any] {
				return hooks
			},
		)
	)
	fn(hooks, factory, group, w.Cmd, w.Results[0].Result, w.Err, t)
}

func TestFailedHooksBeforeSend(fn TestFn, t *testing.T) {
	var (
		wantResult core.Result = nil
		wantErr                = errors.New("HooksFactory.BeforeSend error")

		hooks = mock.NewHooks[any]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
				return nil, wantErr
			},
		)
		factory = mock.NewHooksFactory[any]().RegisterNew(
			func() hks.Hooks[any] {
				return hooks
			},
		)
	)
	fn(hooks, factory, mock.NewClientGroup(), cmock.NewCmd(), wantResult, wantErr, t)
}

func TestTimeout(group mock.ClientGroup, w Want, fn TestFn, t *testing.T) {
	var (
		wantCtx, cancel = context.WithCancel(context.Background())
		hooks           = mock.NewHooks[any]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
				return wantCtx, nil
			},
		).RegisterOnTimeout(
			func(ctx context.Context, sentCmd hks.SentCmd[any], err error) {
				asserterror.Equal(t, ctx, wantCtx)
				asserterror.EqualDeep(t, sentCmd, hks.SentCmd[any]{
					Seq:  w.CmdSeq,
					Size: w.CmdSize,
					Cmd:  w.Cmd,
				})
				asserterror.EqualError(t, err, w.Err)
			},
		)
		factory = mock.NewHooksFactory[any]().RegisterNew(
			func() hks.Hooks[any] {
				return hooks
			},
		)
	)
	group.RegisterForget(
		func(seq core.Seq, clientID grp.ClientID) {
			asserterror.Equal(t, seq, w.CmdSeq)
			asserterror.Equal(t, clientID, w.ClientID)
		},
	)
	cancel()
	fn(hooks, factory, group, w.Cmd, nil, w.Err, t)
}

func TestFailedSend(group mock.ClientGroup, w Want, fn TestFn, t *testing.T) {
	var (
		wantCtx = context.WithoutCancel(context.Background())
		hooks   = mock.NewHooks[any]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
				return wantCtx, nil
			},
		).RegisterOnError(
			func(ctx context.Context, sentCmd hks.SentCmd[any], err error) {
				asserterror.Equal(t, ctx, wantCtx)
				asserterror.EqualError(t, err, w.Err)
			},
		)
		factory = mock.NewHooksFactory[any]().RegisterNew(
			func() hks.Hooks[any] {
				return hooks
			},
		)
	)
	fn(hooks, factory, group, w.Cmd, nil, w.Err, t)
}

func Test(hooks mock.Hooks[any], factory mock.HooksFactory[any],
	group mock.ClientGroup,
	cmd cmock.Cmd,
	Result core.Result,
	wantErr error,
	t *testing.T,
) {
	var (
		sender = sndr.New(group, sndr.WithHooksFactory(factory))
		mocks  = []*mok.Mock{hooks.Mock, factory.Mock, group.Mock, cmd.Mock}
	)
	result, err := sender.Send(context.Background(), cmd)
	asserterror.EqualError(t, err, wantErr)
	asserterror.EqualDeep(t, result, Result)

	asserterror.EqualDeep(t, mok.CheckCalls(mocks), mok.EmptyInfomap)
}

func WrapTestDeadline(deadline time.Time) TestFn {
	return func(hooks mock.Hooks[any], factory mock.HooksFactory[any],
		group mock.ClientGroup,
		cmd cmock.Cmd,
		Result core.Result,
		wantErr error,
		t *testing.T,
	) {
		TestDeadline(hooks, factory, group, deadline, cmd, Result, wantErr, t)
	}
}

func TestDeadline(hooks mock.Hooks[any], factory mock.HooksFactory[any],
	group mock.ClientGroup,
	deadline time.Time,
	cmd cmock.Cmd,
	Result core.Result,
	wantErr error,
	t *testing.T,
) {
	var (
		sender = sndr.New(group, sndr.WithHooksFactory(factory))
		mocks  = []*mok.Mock{hooks.Mock, factory.Mock, group.Mock, cmd.Mock}
	)
	result, err := sender.SendWithDeadline(context.Background(), cmd, deadline)
	asserterror.EqualError(t, err, wantErr)
	asserterror.EqualDeep(t, result, Result)

	asserterror.EqualDeep(t, mok.CheckCalls(mocks), mok.EmptyInfomap)
}
