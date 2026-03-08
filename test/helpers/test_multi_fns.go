package helpers

import (
	"context"
	"errors"
	"testing"
	"time"

	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/core-go"
	cmock "github.com/cmd-stream/core-go/test/mock"
	sndr "github.com/cmd-stream/sender-go"
	hks "github.com/cmd-stream/sender-go/hooks"
	mock "github.com/cmd-stream/sender-go/test/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type TestMultiFn func(hooks mock.Hooks[any], factory mock.HooksFactory[any],
	group mock.ClientGroup,
	cmd cmock.Cmd,
	resultsCount int,
	handler mock.ResultHandler,
	wantErr error,
	t *testing.T,
)

func TestMultiSuccess(group mock.ClientGroup, handler mock.ResultHandler,
	w Want, fn TestMultiFn, t *testing.T,
) {
	var (
		hooks = mock.NewHooks[any]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
				asserterror.EqualDeep(t, cmd, w.Cmd)

				actx := context.WithoutCancel(ctx)
				return actx, nil
			},
		)
		factory = mock.NewHooksFactory[any]().RegisterNew(
			func() hks.Hooks[any] {
				return hooks
			},
		)
	)
	for i := range w.Results {
		hooks.RegisterOnResult(
			func(ctx context.Context, sentCmd hks.SentCmd[any],
				recvResult hks.ReceivedResult, err error,
			) {
				asserterror.EqualDeep(t, sentCmd, hks.SentCmd[any]{
					Seq:  w.CmdSeq,
					Size: w.CmdSize,
					Cmd:  w.Cmd,
				})
				asserterror.EqualDeep(t, recvResult, hks.ReceivedResult{
					Seq:    w.Results[i].Seq,
					Size:   w.Results[i].BytesRead,
					Result: w.Results[i].Result,
				})
				asserterror.EqualError(t, err, w.Results[i].Err)
			},
		)
	}
	fn(hooks, factory, group, w.Cmd, len(w.Results), handler, w.Err, t)
}

func TestMultiFailedHooksBeforeSend(fn TestMultiFn, t *testing.T) {
	var (
		wantErr = errors.New("HooksFactory.BeforeSend error")

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
	fn(hooks, factory, mock.NewClientGroup(), cmock.NewCmd(), 0,
		mock.NewResultHandler(), wantErr, t)
}

func TestMultiTimeout(wantCtx context.Context, group mock.ClientGroup,
	handler mock.ResultHandler,
	w Want,
	fn TestMultiFn,
	t *testing.T,
) {
	var (
		hooks = mock.NewHooks[any]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
				return wantCtx, nil
			},
		).RegisterOnResult(
			func(ctx context.Context, sentCmd hks.SentCmd[any],
				recvResult hks.ReceivedResult, err error,
			) {
				// nothing to do
			},
		).RegisterOnTimeout(
			func(ctx context.Context, sentCmd hks.SentCmd[any], err error) {
				asserterror.Equal(t, ctx, wantCtx)
				asserterror.EqualDeep(t, sentCmd, hks.SentCmd[any]{
					Seq:  w.CmdSeq,
					Size: w.CmdSize,
					Cmd:  w.Cmd,
				})
				asserterror.EqualError(t, err, sndr.ErrTimeout)
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
	fn(hooks, factory, group, w.Cmd, len(w.Results), handler, w.Err, t)
}

func TestMultiFailedSend(group mock.ClientGroup, w Want, fn TestMultiFn, t *testing.T) {
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
	fn(hooks, factory, group, w.Cmd, 0, mock.NewResultHandler(), w.Err, t)
}

func TestMulti(hooks mock.Hooks[any], factory mock.HooksFactory[any],
	group mock.ClientGroup,
	cmd cmock.Cmd,
	resultsCount int,
	handler mock.ResultHandler,
	wantErr error,
	t *testing.T,
) {
	var (
		sender = sndr.New(group, sndr.WithHooksFactory(factory))
		mocks  = []*mok.Mock{hooks.Mock, factory.Mock, group.Mock, handler.Mock, cmd.Mock}
	)
	err := sender.SendMulti(context.Background(), cmd, resultsCount, handler)
	asserterror.EqualError(t, err, wantErr)

	asserterror.EqualDeep(t, mok.CheckCalls(mocks), mok.EmptyInfomap)
}

func WrapTestMultiDeadline(deadline time.Time) TestMultiFn {
	return func(hooks mock.Hooks[any], factory mock.HooksFactory[any],
		group mock.ClientGroup,
		cmd cmock.Cmd,
		resultsCount int,
		handler mock.ResultHandler,
		wantErr error,
		t *testing.T,
	) {
		TestMultiDeadline(hooks, factory, group, deadline, cmd, resultsCount,
			handler, wantErr, t)
	}
}

func TestMultiDeadline(hooks mock.Hooks[any], factory mock.HooksFactory[any],
	group mock.ClientGroup,
	deadline time.Time,
	cmd cmock.Cmd,
	resultsCount int,
	handler mock.ResultHandler,
	wantErr error,
	t *testing.T,
) {
	var (
		sender = sndr.New(group, sndr.WithHooksFactory(factory))
		mocks  = []*mok.Mock{hooks.Mock, factory.Mock, group.Mock, cmd.Mock}
	)
	err := sender.SendMultiWithDeadline(context.Background(), cmd, resultsCount,
		handler, deadline)
	asserterror.EqualError(t, err, wantErr)

	asserterror.EqualDeep(t, mok.CheckCalls(mocks), mok.EmptyInfomap)
}
