package grpc

type Options struct {
	host string
	port int
}

type Option func(*Options)

func WithHost(host string) Option {
	return func(so *Options) {
		so.host = host
	}
}

func WithPort(port int) Option {
	return func(so *Options) {
		so.port = port
	}
}
