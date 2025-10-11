# sender-go

[![Go Reference](https://pkg.go.dev/badge/github.com/cmd-stream/sender-go.svg)](https://pkg.go.dev/github.com/cmd-stream/sender-go)
[![GoReportCard](https://goreportcard.com/badge/cmd-stream/sender-go)](https://goreportcard.com/report/github.com/cmd-stream/sender-go)
[![codecov](https://codecov.io/gh/cmd-stream/sender-go/graph/badge.svg?token=RXPJ6ZIPK7)](https://codecov.io/gh/cmd-stream/sender-go)

A lightweight Go library for sending commands using a [cmd-stream](https://github.com/cmd-stream/cmd-stream-go)
client group, with built-in support for hooks, deadlines, and multi-result
handling.

It provides a high-level abstraction over the following code:

```go
import (
  grp "github.com/cmd-stream/cmd-stream-go/group"
  core "github.com/cmd-stream/sender-go"
)

var (
  group grp.ClientGroup = ...
  results chan core.AsyncResult = ...
)
seq, clientID, n, err := group.Send(cmd, results)
if err != nil {
  return err
}
select {
case <-ctx.Done():
  group.Forget(seq, clientID)
  err = ErrTimeout
case asyncResult := <-results:
  result = asyncResult.Result
  err = asyncResult.Error
}
```

## Usage Examples

Here are just a few lines of code to show typical usage:

```go
// Send one Command, receive one Result.
result, err := sender.Send(ctx, cmd) // ctx allows canceling if the wait takes 
// too long

// Send a Command with a deadline.
result, err := sender.SendWithDeadline(ctx, cmd, deadline)

// Send one Command, receive multiple Results.
var resultHandler ResultHandlerFn = func(result core.Result, err error) error { 
  // handle each result here
  return nil
}
err := sender.SendMulti(ctx, cmd, resultsCount, resultHandler) 
// resultsCount defines the number of expected Results

// Send a Command with a deadline, receive multiple Results.
err := sender.SendMultiWithDeadline(ctx, cmd, resultsCount, resultHandler,
  deadline)
```

More detailed examples can be found at [examples-go](https://github.com/cmd-stream/cmd-stream-examples-go).

For special cases, you can implement your own sender, it’s not hard to do.

## Hooks

sender-go also supports hooks, allowing you to customize behavior during the send
process. Hooks can be used for logging, instrumentation, circuit breaker
integration, and more. They are provided through a `HooksFactory`, which creates
a fresh `Hooks` instance for each send operation. This ensures isolation and
flexibility for each request.

The `hooks` package already includes ready-to-use implementations like
`CircuitBreakerHooks` and `NoopHooks`.

## Resilient Configuration

The [ResilientConfig](resilient_config.go) type provides a convenient way to
construct a **resilient sender** that automatically handles keepalive,
reconnects, and circuit breaker behavior.

### Example

```go
package main

import (
  "time"

  sndr "github.com/cmd-stream/sender-go"
)

func main() {
  cfg := sender.ResilientConfig[MyCmd]{
    KeepaliveTime:               30 * time.Second,
    KeepaliveIntvl:              10 * time.Second,
    CircuitBreakerWindowSize:    20,
    CircuitBreakerFailureRate:   0.5,
    CircuitBreakerOpenDuration:  30 * time.Second,
    CircuitBreakerSuccessThreshold: 2,
    // Optional HooksFactory to observe or modify sender behavior.
    // HooksFactory: ...
  }

  sender, err := sndr.Make(cfg.ToOptions()...)
  ...
}
```

### Behavior

- `ResilientConfig` panics if any required field is unset (zero).
- If no `HooksFactory` is provided, a `NoopHooksFactory` is used automatically.
- The generated options include:
  - automatic reconnects (`grp.WithReconnect`);
  - keepalive management (`cln.WithKeepalive`)
  - circuit breaker hooks factory (`WithHooksFactory`).
