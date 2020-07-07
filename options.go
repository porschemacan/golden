package golden

import (
	"github.com/porschemacan/golden/libs"
	"time"
)

type ServerOptions struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	ServiceName  string
	Address      string
	CORSConfig   *CORSConfig
	LogConfig    *libs.LogConfig
	HtmlConfig   *HtmlConfig
}

type Option func(*ServerOptions)

func newOptions(opts ...Option) *ServerOptions {
	serverOptions := &ServerOptions{
		Address:     ":7777",
		ServiceName: "golden-http-server",
	}

	for _, opt := range opts {
		opt(serverOptions)
	}

	return serverOptions
}

func Name(name string) Option {
	return func(o *ServerOptions) {
		o.ServiceName = name
	}
}

func Timeout(readTimeout time.Duration, writeTimeout time.Duration) Option {
	return func(o *ServerOptions) {
		o.ReadTimeout = readTimeout
		o.WriteTimeout = writeTimeout
	}
}

func Address(address string) Option {
	return func(o *ServerOptions) {
		o.Address = address
	}
}

func Cors(cors *CORSConfig) Option {
	return func(o *ServerOptions) {
		o.CORSConfig = cors
	}
}

func LogConfig(cfg *libs.LogConfig) Option {
	return func(o *ServerOptions) {
		o.LogConfig = cfg
	}
}

func Html(cfg *HtmlConfig) Option {
	return func(o *ServerOptions) {
		o.HtmlConfig = cfg
	}
}
