package sender_test

import (
	"context"
	"errors"
	"testing"
	"time"

	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/core-go"
	cmock "github.com/cmd-stream/core-go/test/mock"
	sndr "github.com/cmd-stream/sender-go"
	"github.com/cmd-stream/sender-go/test/helpers"
	mock "github.com/cmd-stream/sender-go/test/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestSender(t *testing.T) {
	t.Run("Send", func(t *testing.T) {
		t.Run("Should succeed", func(t *testing.T) {
			var (
				want = helpers.Want{
					Cmd: cmock.NewCmd(),
					Results: []helpers.WantResult{
						{
							Seq:       core.Seq(1),
							BytesRead: 20,
							Result:    cmock.NewResult(),
							Err:       nil,
						},
					},

					CmdSeq:     core.Seq(1),
					ClientID:   grp.ClientID(2),
					CmdSize:    10,
					CmdSendErr: nil,

					Err: nil,
				}
				group = mock.NewClientGroup().RegisterSend(
					func(cmd core.Cmd[any], results chan<- core.AsyncResult) (
						seq core.Seq, clientID grp.ClientID, n int, err error,
					) {
						asserterror.EqualDeep(t, cmd, want.Cmd)
						results <- core.AsyncResult{
							Seq:       want.Results[0].Seq,
							BytesRead: want.Results[0].BytesRead,
							Result:    want.Results[0].Result,
							Error:     want.Results[0].Err,
						}
						return want.CmdSeq, want.ClientID, want.CmdSize, want.CmdSendErr
					},
				)
			)
			helpers.TestSuccess(group, want, helpers.Test, t)
		})

		t.Run("If hooks.BeforeSend fails with an error, Send should return it", func(t *testing.T) {
			helpers.TestFailedHooksBeforeSend(helpers.Test, t)
		})

		t.Run("If ClientGroup.Send fails with an error, Send should return it",
			func(t *testing.T) {
				var (
					want = helpers.Want{
						Cmd: cmock.NewCmd(),

						CmdSeq:   core.Seq(1),
						ClientID: 1,
						CmdSize:  10,

						Err: errors.New("ClientGroup.Send error"),
					}
					group = mock.NewClientGroup().RegisterSend(
						func(cmd core.Cmd[any], results chan<- core.AsyncResult) (seq core.Seq,
							clientID grp.ClientID, n int, err error,
						) {
							seq = want.CmdSeq
							clientID = want.ClientID
							n = want.CmdSize
							err = want.Err
							return
						},
					)
				)
				helpers.TestFailedSend(group, want, helpers.Test, t)
			})

		t.Run("Should be able to timeout", func(t *testing.T) {
			var (
				want = helpers.Want{
					Cmd: cmock.NewCmd(),

					CmdSeq:   core.Seq(1),
					ClientID: grp.ClientID(1),
					CmdSize:  10,

					Err: sndr.ErrTimeout,
				}
				group = mock.NewClientGroup().RegisterSend(
					func(cmd core.Cmd[any], results chan<- core.AsyncResult) (
						seq core.Seq, clientID grp.ClientID, n int, err error,
					) {
						seq = want.CmdSeq
						clientID = want.ClientID
						n = want.CmdSize
						return
					},
				)
			)
			helpers.TestTimeout(group, want, helpers.Test, t)
		})
	})

	t.Run("SendWithDeadline", func(t *testing.T) {
		t.Run("Should succeed", func(t *testing.T) {
			var (
				want = helpers.Want{
					Cmd: cmock.NewCmd(),
					Results: []helpers.WantResult{
						{
							Seq:       core.Seq(1),
							BytesRead: 20,
							Result:    cmock.NewResult(),
							Err:       nil,
						},
					},

					CmdSeq:     core.Seq(1),
					ClientID:   grp.ClientID(2),
					CmdSize:    10,
					CmdSendErr: nil,

					Err: nil,
				}
				wantDeadline = time.Now()
				group        = mock.NewClientGroup().RegisterSendWithDeadline(
					func(cmd core.Cmd[any], results chan<- core.AsyncResult,
						deadline time.Time,
					) (seq core.Seq, clientID grp.ClientID, n int, err error) {
						asserterror.EqualDeep(t, cmd, want.Cmd)
						asserterror.Equal(t, deadline, wantDeadline)
						results <- core.AsyncResult{
							Seq:       want.Results[0].Seq,
							BytesRead: want.Results[0].BytesRead,
							Result:    want.Results[0].Result,
							Error:     want.Results[0].Err,
						}
						return want.CmdSeq, want.ClientID, want.CmdSize, want.CmdSendErr
					},
				)
				fn = helpers.WrapTestDeadline(wantDeadline)
			)
			helpers.TestSuccess(group, want, fn, t)
		})

		t.Run("If hooks.BeforeSend fails with an error, Send should return it", func(t *testing.T) {
			var (
				wantDeadline = time.Now()
				fn           = helpers.WrapTestDeadline(wantDeadline)
			)
			helpers.TestFailedHooksBeforeSend(fn, t)
		})

		t.Run("If ClientGroup.Send fails with an error, SendWithDeadline should return it",
			func(t *testing.T) {
				var (
					want = helpers.Want{
						Cmd: cmock.NewCmd(),

						CmdSeq:   core.Seq(1),
						ClientID: 1,
						CmdSize:  10,

						Err: errors.New("ClientGroup.Send error"),
					}
					deadline = time.Now()
					group    = mock.NewClientGroup().RegisterSendWithDeadline(
						func(cmd core.Cmd[any], results chan<- core.AsyncResult, deadline time.Time) (seq core.Seq,
							clientID grp.ClientID, n int, err error,
						) {
							seq = want.CmdSeq
							clientID = want.ClientID
							n = want.CmdSize
							err = want.Err
							return
						},
					)
					fn = helpers.WrapTestDeadline(deadline)
				)
				helpers.TestFailedSend(group, want, fn, t)
			})

		t.Run("Should be able to timeout", func(t *testing.T) {
			var (
				want = helpers.Want{
					Cmd: cmock.NewCmd(),

					CmdSeq:   core.Seq(1),
					ClientID: grp.ClientID(1),
					CmdSize:  10,

					Err: sndr.ErrTimeout,
				}
				group = mock.NewClientGroup().RegisterSend(
					func(cmd core.Cmd[any], results chan<- core.AsyncResult) (
						seq core.Seq, clientID grp.ClientID, n int, err error,
					) {
						seq = want.CmdSeq
						clientID = want.ClientID
						n = want.CmdSize
						return
					},
				)
			)
			helpers.TestTimeout(group, want, helpers.Test, t)
		})
	})

	t.Run("SendMulti", func(t *testing.T) {
		t.Run("Should succeed", func(t *testing.T) {
			var (
				want = helpers.Want{
					Cmd: cmock.NewCmd(),
					Results: []helpers.WantResult{
						{
							Seq: core.Seq(1),
							Result: cmock.NewResult().RegisterLastOne(
								func() (lastOne bool) { return false },
							),
							BytesRead: 10,
							Err:       nil,
						},
						{
							Seq: core.Seq(2),
							Result: cmock.NewResult().RegisterLastOne(
								func() (lastOne bool) { return true },
							),
							BytesRead: 20,
							Err:       nil,
						},
					},

					CmdSeq:     core.Seq(1),
					ClientID:   grp.ClientID(2),
					CmdSize:    10,
					CmdSendErr: nil,

					Err: nil,
				}
				group = mock.NewClientGroup().RegisterSend(
					func(cmd core.Cmd[any], results chan<- core.AsyncResult) (
						seq core.Seq, clientID grp.ClientID, n int, err error,
					) {
						asserterror.EqualDeep(t, cmd, want.Cmd)
						for i := range want.Results {
							results <- core.AsyncResult{
								Seq:       want.Results[i].Seq,
								BytesRead: want.Results[i].BytesRead,
								Result:    want.Results[i].Result,
								Error:     want.Results[i].Err,
							}
						}
						return want.CmdSeq, want.ClientID, want.CmdSize, want.CmdSendErr
					},
				)
				handler = mock.NewResultHandler()
			)
			for i := range want.Results {
				handler.RegisterHandle(
					func(result core.Result, err error) error {
						asserterror.EqualDeep(t, result, want.Results[i].Result)
						asserterror.EqualError(t, err, want.Results[i].Err)
						return nil
					},
				)
			}
			helpers.TestMultiSuccess(group, handler, want, helpers.TestMulti, t)
		})

		t.Run("If hooks.BeforeSend fails with an error, Send should return it", func(t *testing.T) {
			helpers.TestMultiFailedHooksBeforeSend(helpers.TestMulti, t)
		})

		t.Run("If ClientGroup.Send fails with an error, SendMulti should return it",
			func(t *testing.T) {
				var (
					want = helpers.Want{
						Cmd: cmock.NewCmd(),

						CmdSeq:   core.Seq(1),
						ClientID: 1,
						CmdSize:  10,

						Err: errors.New("ClientGroup.Send error"),
					}
					group = mock.NewClientGroup().RegisterSend(
						func(cmd core.Cmd[any], results chan<- core.AsyncResult) (seq core.Seq,
							clientID grp.ClientID, n int, err error,
						) {
							seq = want.CmdSeq
							clientID = want.ClientID
							n = want.CmdSize
							err = want.Err
							return
						},
					)
				)
				helpers.TestMultiFailedSend(group, want, helpers.TestMulti, t)
			})

		t.Run("Should be able to timeout", func(t *testing.T) {
			var (
				wantCtx, cancel = context.WithCancel(context.Background())
				want            = helpers.Want{
					Cmd: cmock.NewCmd(),
					Results: []helpers.WantResult{
						{
							Seq: core.Seq(1),
							Result: cmock.NewResult().RegisterLastOne(
								func() (lastOne bool) { return false },
							),
							BytesRead: 10,
							Err:       nil,
						},
					},

					CmdSeq:     core.Seq(1),
					ClientID:   grp.ClientID(2),
					CmdSize:    10,
					CmdSendErr: nil,

					Err: nil,
				}
				group = mock.NewClientGroup().RegisterSend(
					func(cmd core.Cmd[any], results chan<- core.AsyncResult) (
						seq core.Seq, clientID grp.ClientID, n int, err error,
					) {
						asserterror.EqualDeep(t, cmd, want.Cmd)
						for i := range want.Results {
							results <- core.AsyncResult{
								Seq:       want.Results[i].Seq,
								BytesRead: want.Results[i].BytesRead,
								Result:    want.Results[i].Result,
								Error:     want.Results[i].Err,
							}
						}
						return want.CmdSeq, want.ClientID, want.CmdSize, want.CmdSendErr
					},
				)
				handler = mock.NewResultHandler()
			)
			defer cancel()
			for i := range want.Results {
				handler.RegisterHandle(
					func(result core.Result, err error) error {
						asserterror.EqualDeep(t, result, want.Results[i].Result)
						asserterror.EqualError(t, err, want.Results[i].Err)
						cancel()
						return nil
					},
				)
			}
			handler.RegisterHandle(
				func(result core.Result, err error) error {
					asserterror.EqualError(t, err, sndr.ErrTimeout)
					return nil
				},
			)
			helpers.TestMultiTimeout(wantCtx, group, handler, want, helpers.TestMulti, t)
		})
	})

	t.Run("SendMultiWithDeadline", func(t *testing.T) {
		t.Run("Should succeed", func(t *testing.T) {
			var (
				want = helpers.Want{
					Cmd: cmock.NewCmd(),
					Results: []helpers.WantResult{
						{
							Seq: core.Seq(1),
							Result: cmock.NewResult().RegisterLastOne(
								func() (lastOne bool) { return false },
							),
							BytesRead: 10,
							Err:       nil,
						},
						{
							Seq: core.Seq(2),
							Result: cmock.NewResult().RegisterLastOne(
								func() (lastOne bool) { return true },
							),
							BytesRead: 20,
							Err:       nil,
						},
					},

					CmdSeq:     core.Seq(1),
					ClientID:   grp.ClientID(2),
					CmdSize:    10,
					CmdSendErr: nil,

					Err: nil,
				}
				wantDeadline = time.Now()
				group        = mock.NewClientGroup().RegisterSendWithDeadline(
					func(cmd core.Cmd[any], results chan<- core.AsyncResult,
						deadline time.Time,
					) (seq core.Seq, clientID grp.ClientID, n int, err error) {
						asserterror.EqualDeep(t, cmd, want.Cmd)
						asserterror.Equal(t, deadline, wantDeadline)
						for i := range want.Results {
							results <- core.AsyncResult{
								Seq:       want.Results[i].Seq,
								BytesRead: want.Results[i].BytesRead,
								Result:    want.Results[i].Result,
								Error:     want.Results[i].Err,
							}
						}
						return want.CmdSeq, want.ClientID, want.CmdSize, want.CmdSendErr
					},
				)
				handler = mock.NewResultHandler()
				fn      = helpers.WrapTestMultiDeadline(wantDeadline)
			)
			for i := range want.Results {
				handler.RegisterHandle(
					func(result core.Result, err error) error {
						asserterror.EqualDeep(t, result, want.Results[i].Result)
						asserterror.EqualError(t, err, want.Results[i].Err)
						return nil
					},
				)
			}
			helpers.TestMultiSuccess(group, handler, want, fn, t)
		})

		t.Run("If hooks.BeforeSend fails with an error, Send should return it", func(t *testing.T) {
			var (
				wantDeadline = time.Now()
				fn           = helpers.WrapTestMultiDeadline(wantDeadline)
			)
			helpers.TestMultiFailedHooksBeforeSend(fn, t)
		})

		t.Run("If ClientGroup.Send fails with an error, SendMultiWithDeadline should return it",
			func(t *testing.T) {
				var (
					want = helpers.Want{
						Cmd: cmock.NewCmd(),

						CmdSeq:   core.Seq(1),
						ClientID: 1,
						CmdSize:  10,

						Err: errors.New("ClientGroup.Send error"),
					}
					wantDeadline = time.Now()
					group        = mock.NewClientGroup().RegisterSendWithDeadline(
						func(cmd core.Cmd[any], results chan<- core.AsyncResult,
							deadline time.Time,
						) (seq core.Seq, clientID grp.ClientID, n int, err error) {
							seq = want.CmdSeq
							clientID = want.ClientID
							n = want.CmdSize
							err = want.Err
							return
						},
					)
					fn = helpers.WrapTestMultiDeadline(wantDeadline)
				)
				helpers.TestMultiFailedSend(group, want, fn, t)
			})

		t.Run("Should be able to timeout", func(t *testing.T) {
			var (
				wantCtx, cancel = context.WithCancel(context.Background())
				want            = helpers.Want{
					Cmd: cmock.NewCmd(),
					Results: []helpers.WantResult{
						{
							Seq: core.Seq(1),
							Result: cmock.NewResult().RegisterLastOne(
								func() (lastOne bool) { return false },
							),
							BytesRead: 10,
							Err:       nil,
						},
					},

					CmdSeq:     core.Seq(1),
					ClientID:   grp.ClientID(2),
					CmdSize:    10,
					CmdSendErr: nil,

					Err: nil,
				}
				wantDeadline = time.Now()
				group        = mock.NewClientGroup().RegisterSendWithDeadline(
					func(cmd core.Cmd[any], results chan<- core.AsyncResult,
						deadline time.Time,
					) (seq core.Seq, clientID grp.ClientID, n int, err error) {
						asserterror.EqualDeep(t, cmd, want.Cmd)
						for i := range want.Results {
							results <- core.AsyncResult{
								Seq:       want.Results[i].Seq,
								BytesRead: want.Results[i].BytesRead,
								Result:    want.Results[i].Result,
								Error:     want.Results[i].Err,
							}
						}
						return want.CmdSeq, want.ClientID, want.CmdSize, want.CmdSendErr
					},
				)

				handler = mock.NewResultHandler()
				fn      = helpers.WrapTestMultiDeadline(wantDeadline)
			)
			defer cancel()
			for i := range want.Results {
				handler.RegisterHandle(
					func(result core.Result, err error) error {
						asserterror.EqualDeep(t, result, want.Results[i].Result)
						asserterror.EqualError(t, err, want.Results[i].Err)
						cancel()
						return nil
					},
				)
			}
			handler.RegisterHandle(
				func(result core.Result, err error) error {
					asserterror.EqualError(t, err, sndr.ErrTimeout)
					return nil
				},
			)
			helpers.TestMultiTimeout(wantCtx, group, handler, want, fn, t)
		})
	})
}
