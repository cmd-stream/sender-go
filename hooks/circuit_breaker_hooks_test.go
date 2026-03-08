package hooks_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cmd-stream/core-go"
	cmock "github.com/cmd-stream/core-go/test/mock"
	hks "github.com/cmd-stream/sender-go/hooks"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"

	mock "github.com/cmd-stream/sender-go/test/mock"
)

func TestCircuitBreakerHooks(t *testing.T) {
	t.Run("BeforeSend", func(t *testing.T) {
		t.Run("Should work", func(t *testing.T) {
			var (
				wantCtx       = context.Background()
				wantErr error = nil
				wantCmd       = cmock.NewCmd()
				cb            = mock.NewCircuitBreaker().RegisterAllow(
					func() bool { return true },
				)
				innerHooks = mock.NewHooks[any]().RegisterBeforeSend(
					func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
						asserterror.Equal(t, ctx, wantCtx)
						asserterror.EqualDeep(t, cmd, wantCmd)
						return ctx, nil
					},
				)
				hooks = hks.NewCircuitBreakerHooks(cb, innerHooks)
				mocks = []*mok.Mock{cb.Mock, innerHooks.Mock}
			)
			ctx, err := hooks.BeforeSend(wantCtx, cmock.NewCmd())
			asserterror.Equal(t, ctx, wantCtx)
			asserterror.EqualError(t, err, wantErr)

			asserterror.EqualDeep(t, mok.CheckCalls(mocks), mok.EmptyInfomap)
		})

		t.Run("Should return an error if the circuit breaker is open", func(t *testing.T) {
			var (
				wantCtx = context.Background()
				cb      = mock.NewCircuitBreaker().RegisterAllow(
					func() bool { return false },
				)
				hooks = hks.NewCircuitBreakerHooks[any](cb, nil)
				mocks = []*mok.Mock{cb.Mock}
			)
			ctx, err := hooks.BeforeSend(wantCtx, cmock.NewCmd())
			asserterror.Equal(t, ctx, wantCtx)
			asserterror.EqualError(t, err, hks.ErrNotAllowed)

			asserterror.EqualDeep(t, mok.CheckCalls(mocks), mok.EmptyInfomap)
		})
	})

	t.Run("OnError", func(t *testing.T) {
		var (
			wantCtx     = context.Background()
			wantSentCmd = hks.SentCmd[any]{}
			wantErr     = errors.New("error")
			cb          = mock.NewCircuitBreaker().RegisterFail(func() {})
			innerHooks  = mock.NewHooks[any]().RegisterOnError(
				func(ctx context.Context, sentCmd hks.SentCmd[any], err error) {
					asserterror.Equal(t, ctx, wantCtx)
					asserterror.EqualDeep(t, sentCmd, wantSentCmd)
					asserterror.EqualError(t, err, wantErr)
				},
			)
			hooks = hks.NewCircuitBreakerHooks(cb, innerHooks)
			mocks = []*mok.Mock{cb.Mock, innerHooks.Mock}
		)
		hooks.OnError(wantCtx, wantSentCmd, wantErr)

		asserterror.EqualDeep(t, mok.CheckCalls(mocks), mok.EmptyInfomap)
	})

	t.Run("OnResult", func(t *testing.T) {
		var (
			wantCtx        = context.Background()
			wantSentCmd    = hks.SentCmd[any]{}
			wantRecvResult = hks.ReceivedResult{}
			wantErr        = errors.New("error")
			cb             = mock.NewCircuitBreaker().RegisterSuccess(func() {})
			innerHooks     = mock.NewHooks[any]().RegisterOnResult(
				func(ctx context.Context, sentCmd hks.SentCmd[any],
					recvResult hks.ReceivedResult, err error,
				) {
					asserterror.Equal(t, ctx, wantCtx)
					asserterror.EqualDeep(t, sentCmd, wantSentCmd)
					asserterror.EqualDeep(t, recvResult, wantRecvResult)
					asserterror.EqualError(t, err, wantErr)
				},
			)
			hooks = hks.NewCircuitBreakerHooks(cb, innerHooks)
			mocks = []*mok.Mock{cb.Mock, innerHooks.Mock}
		)
		hooks.OnResult(wantCtx, wantSentCmd, wantRecvResult, wantErr)

		asserterror.EqualDeep(t, mok.CheckCalls(mocks), mok.EmptyInfomap)
	})

	t.Run("OnTimeout", func(t *testing.T) {
		var (
			wantCtx     = context.Background()
			wantSentCmd = hks.SentCmd[any]{}
			wantErr     = errors.New("error")
			cb          = mock.NewCircuitBreaker().RegisterFail(func() {})
			innerHooks  = mock.NewHooks[any]().RegisterOnTimeout(
				func(ctx context.Context, sentCmd hks.SentCmd[any], err error) {
					asserterror.Equal(t, ctx, wantCtx)
					asserterror.EqualDeep(t, sentCmd, wantSentCmd)
					asserterror.EqualError(t, err, wantErr)
				},
			)
			hooks = hks.NewCircuitBreakerHooks(cb, innerHooks)
			mocks = []*mok.Mock{cb.Mock, innerHooks.Mock}
		)
		hooks.OnTimeout(wantCtx, wantSentCmd, wantErr)

		asserterror.EqualDeep(t, mok.CheckCalls(mocks), mok.EmptyInfomap)
	})
}
