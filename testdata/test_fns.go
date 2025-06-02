package testdata

import (
	"context"
	"errors"
	"testing"
	"time"

	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/core-go"
	cmock "github.com/cmd-stream/core-go/testdata/mock"
	hks "github.com/cmd-stream/sender-go/hooks"
	hmock "github.com/cmd-stream/sender-go/hooks/testdata/mock"
	"github.com/cmd-stream/sender-go/testdata/mock"

	sndr "github.com/cmd-stream/sender-go"
	"github.com/ymz-ncnk/mok"

	asserterror "github.com/ymz-ncnk/assert/error"
)

type TestFn func(hooks hmock.Hooks[any], factory hmock.HooksFactory[any],
	group mock.ClientGroup,
	cmd cmock.Cmd,
	wantResult core.Result,
	wantErr error,
	t *testing.T,
)

func TestShouldWork(group mock.ClientGroup, w Want, fn TestFn, t *testing.T) {
	var (
		hooks = hmock.NewHooks[any]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
				asserterror.EqualDeep(cmd, w.Cmd, t)

				actx := context.WithoutCancel(ctx)
				return actx, nil
			},
		).RegisterOnResult(
			func(ctx context.Context, sentCmd hks.SentCmd[any],
				recvResult hks.ReceivedResult, err error) {
				asserterror.EqualDeep(sentCmd, hks.SentCmd[any]{
					Seq:  w.CmdSeq,
					Size: w.CmdSize,
					Cmd:  w.Cmd,
				}, t)
				asserterror.EqualDeep(recvResult, hks.ReceivedResult{
					Seq:    w.Results[0].Seq,
					Size:   w.Results[0].BytesRead,
					Result: w.Results[0].Result,
				}, t)
				asserterror.EqualError(err, w.Results[0].Err, t)
			},
		)
		factory = hmock.NewHooksFactory[any]().RegisterNew(
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

		hooks = hmock.NewHooks[any]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
				return nil, wantErr
			},
		)
		factory = hmock.NewHooksFactory[any]().RegisterNew(
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
		hooks           = hmock.NewHooks[any]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
				return wantCtx, nil
			},
		).RegisterOnTimeout(
			func(ctx context.Context, sentCmd hks.SentCmd[any], err error) {
				asserterror.Equal(ctx, wantCtx, t)
				asserterror.EqualDeep(sentCmd, hks.SentCmd[any]{
					Seq:  w.CmdSeq,
					Size: w.CmdSize,
					Cmd:  w.Cmd,
				}, t)
				asserterror.EqualError(err, w.Err, t)
			},
		)
		factory = hmock.NewHooksFactory[any]().RegisterNew(
			func() hks.Hooks[any] {
				return hooks
			},
		)
	)
	group.RegisterForget(
		func(seq core.Seq, clientID grp.ClientID) {
			asserterror.Equal(seq, w.CmdSeq, t)
			asserterror.Equal(clientID, w.ClientID, t)
		},
	)
	cancel()
	fn(hooks, factory, group, w.Cmd, nil, w.Err, t)
}

func TestFailedSend(group mock.ClientGroup, w Want, fn TestFn, t *testing.T) {
	var (
		wantCtx = context.WithoutCancel(context.Background())
		hooks   = hmock.NewHooks[any]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
				return wantCtx, nil
			},
		).RegisterOnError(
			func(ctx context.Context, sentCmd hks.SentCmd[any], err error) {
				asserterror.Equal(ctx, wantCtx, t)
				asserterror.EqualError(err, w.Err, t)
			},
		)
		factory = hmock.NewHooksFactory[any]().RegisterNew(
			func() hks.Hooks[any] {
				return hooks
			},
		)
	)
	fn(hooks, factory, group, w.Cmd, nil, w.Err, t)
}

func Test(hooks hmock.Hooks[any], factory hmock.HooksFactory[any],
	group mock.ClientGroup,
	cmd cmock.Cmd,
	Result core.Result,
	wantErr error,
	t *testing.T,
) {
	var (
		sender = sndr.New[any](group, sndr.WithHooksFactory(factory))
		mocks  = []*mok.Mock{hooks.Mock, factory.Mock, group.Mock, cmd.Mock}
	)
	result, err := sender.Send(context.Background(), cmd)
	asserterror.EqualError(err, wantErr, t)
	asserterror.EqualDeep(result, Result, t)

	asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
}

func WrapTestDeadline(deadline time.Time) TestFn {
	return func(hooks hmock.Hooks[any], factory hmock.HooksFactory[any],
		group mock.ClientGroup,
		cmd cmock.Cmd,
		Result core.Result,
		wantErr error,
		t *testing.T,
	) {
		TestDeadline(hooks, factory, group, deadline, cmd, Result, wantErr, t)
	}
}

func TestDeadline(hooks hmock.Hooks[any], factory hmock.HooksFactory[any],
	group mock.ClientGroup,
	deadline time.Time,
	cmd cmock.Cmd,
	Result core.Result,
	wantErr error,
	t *testing.T) {
	var (
		sender = sndr.New(group, sndr.WithHooksFactory(factory))
		mocks  = []*mok.Mock{hooks.Mock, factory.Mock, group.Mock, cmd.Mock}
	)
	result, err := sender.SendWithDeadline(context.Background(), cmd, deadline)
	asserterror.EqualError(err, wantErr, t)
	asserterror.EqualDeep(result, Result, t)

	asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
}
