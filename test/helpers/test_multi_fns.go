package helpers

import (
	"context"
	"errors"
	"testing"
	"time"

	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/core-go"
	sndr "github.com/cmd-stream/sender-go"
	hks "github.com/cmd-stream/sender-go/hooks"
	"github.com/cmd-stream/sender-go/test/mocks"
	cmocks "github.com/cmd-stream/testkit-go/mocks/core"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type TestMultiFn func(hooks mocks.Hooks[any], factory mocks.HooksFactory[any],
	group mocks.ClientGroup,
	cmd cmocks.Cmd,
	resultsCount int,
	handler mocks.ResultHandler,
	wantErr error,
	t *testing.T,
)

func TestMultiShouldWork(group mocks.ClientGroup, handler mocks.ResultHandler,
	w Want, fn TestMultiFn, t *testing.T,
) {
	var (
		hooks = mocks.NewHooks[any]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
				asserterror.EqualDeep(cmd, w.Cmd, t)

				actx := context.WithoutCancel(ctx)
				return actx, nil
			},
		)
		factory = mocks.NewHooksFactory[any]().RegisterNew(
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
				asserterror.EqualDeep(sentCmd, hks.SentCmd[any]{
					Seq:  w.CmdSeq,
					Size: w.CmdSize,
					Cmd:  w.Cmd,
				}, t)
				asserterror.EqualDeep(recvResult, hks.ReceivedResult{
					Seq:    w.Results[i].Seq,
					Size:   w.Results[i].BytesRead,
					Result: w.Results[i].Result,
				}, t)
				asserterror.EqualError(err, w.Results[i].Err, t)
			},
		)
	}
	fn(hooks, factory, group, w.Cmd, len(w.Results), handler, w.Err, t)
}

func TestMultiFailedHooksBeforeSend(fn TestMultiFn, t *testing.T) {
	var (
		wantErr = errors.New("HooksFactory.BeforeSend error")

		hooks = mocks.NewHooks[any]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
				return nil, wantErr
			},
		)
		factory = mocks.NewHooksFactory[any]().RegisterNew(
			func() hks.Hooks[any] {
				return hooks
			},
		)
	)
	fn(hooks, factory, mocks.NewClientGroup(), cmocks.NewCmd(), 0,
		mocks.NewResultHandler(), wantErr, t)
}

func TestMultiTimeout(wantCtx context.Context, group mocks.ClientGroup,
	handler mocks.ResultHandler,
	w Want,
	fn TestMultiFn,
	t *testing.T,
) {
	var (
		hooks = mocks.NewHooks[any]().RegisterBeforeSend(
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
				asserterror.Equal(ctx, wantCtx, t)
				asserterror.EqualDeep(sentCmd, hks.SentCmd[any]{
					Seq:  w.CmdSeq,
					Size: w.CmdSize,
					Cmd:  w.Cmd,
				}, t)
				asserterror.EqualError(err, sndr.ErrTimeout, t)
			},
		)
		factory = mocks.NewHooksFactory[any]().RegisterNew(
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
	fn(hooks, factory, group, w.Cmd, len(w.Results), handler, w.Err, t)
}

func TestMultiFailedSend(group mocks.ClientGroup, w Want, fn TestMultiFn, t *testing.T) {
	var (
		wantCtx = context.WithoutCancel(context.Background())
		hooks   = mocks.NewHooks[any]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
				return wantCtx, nil
			},
		).RegisterOnError(
			func(ctx context.Context, sentCmd hks.SentCmd[any], err error) {
				asserterror.Equal(ctx, wantCtx, t)
				asserterror.EqualError(err, w.Err, t)
			},
		)
		factory = mocks.NewHooksFactory[any]().RegisterNew(
			func() hks.Hooks[any] {
				return hooks
			},
		)
	)
	fn(hooks, factory, group, w.Cmd, 0, mocks.NewResultHandler(), w.Err, t)
}

func TestMulti(hooks mocks.Hooks[any], factory mocks.HooksFactory[any],
	group mocks.ClientGroup,
	cmd cmocks.Cmd,
	resultsCount int,
	handler mocks.ResultHandler,
	wantErr error,
	t *testing.T,
) {
	var (
		sender = sndr.New(group, sndr.WithHooksFactory(factory))
		mocks  = []*mok.Mock{hooks.Mock, factory.Mock, group.Mock, handler.Mock, cmd.Mock}
	)
	err := sender.SendMulti(context.Background(), cmd, resultsCount, handler)
	asserterror.EqualError(err, wantErr, t)

	asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
}

func WrapTestMultiDeadline(deadline time.Time) TestMultiFn {
	return func(hooks mocks.Hooks[any], factory mocks.HooksFactory[any],
		group mocks.ClientGroup,
		cmd cmocks.Cmd,
		resultsCount int,
		handler mocks.ResultHandler,
		wantErr error,
		t *testing.T,
	) {
		TestMultiDeadline(hooks, factory, group, deadline, cmd, resultsCount,
			handler, wantErr, t)
	}
}

func TestMultiDeadline(hooks mocks.Hooks[any], factory mocks.HooksFactory[any],
	group mocks.ClientGroup,
	deadline time.Time,
	cmd cmocks.Cmd,
	resultsCount int,
	handler mocks.ResultHandler,
	wantErr error,
	t *testing.T,
) {
	var (
		sender = sndr.New(group, sndr.WithHooksFactory(factory))
		mocks  = []*mok.Mock{hooks.Mock, factory.Mock, group.Mock, cmd.Mock}
	)
	err := sender.SendMultiWithDeadline(context.Background(), cmd, resultsCount,
		handler, deadline)
	asserterror.EqualError(err, wantErr, t)

	asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
}
