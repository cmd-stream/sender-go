package sender

import (
	"crypto/tls"

	grp "github.com/cmd-stream/cmd-stream-go/group"
)

type MakeOptions[T any] struct {
	Group        []grp.SetOption[T]
	Sender       []SetOption[T]
	TLSConfig    *tls.Config
	ClientsCount int
}

type SetMakeOption[T any] func(o *MakeOptions[T])

// WithGroup sets options for the client group.
func WithGroup[T any](ops ...grp.SetOption[T]) SetMakeOption[T] {
	return func(o *MakeOptions[T]) { o.Group = ops }
}

// WithSender sets options for the sender.
func WithSender[T any](ops ...SetOption[T]) SetMakeOption[T] {
	return func(o *MakeOptions[T]) { o.Sender = ops }
}

// WithTLSConfig sets the TLS configuration for the client group.
func WithTLSConfig[T any](conf *tls.Config) SetMakeOption[T] {
	return func(o *MakeOptions[T]) { o.TLSConfig = conf }
}

// WithClientsCount sets the number of clients in the client group.
func WithClientsCount[T any](count int) SetMakeOption[T] {
	return func(o *MakeOptions[T]) { o.ClientsCount = count }
}

func ApplyMakeOptitions[T any](ops []SetMakeOption[T], o *MakeOptions[T]) {
	for i := range ops {
		if ops[i] != nil {
			ops[i](o)
		}
	}
}
