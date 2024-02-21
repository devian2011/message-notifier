package transport

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type HttpConfig struct {
	Addr string `json:"addr" yaml:"addr"`
	Tls  struct {
		CertFile string `json:"cert_file" yaml:"certFile"`
		KeyFile  string `json:"key_file" yaml:"keyFile"`
	} `yaml:"tls" json:"tls"`
}

type HttpTransport struct {
	cfg     *HttpConfig
	srv     *fasthttp.Server
	router  *router.Router
	handler Handler
}

func NewHttpTransport(cfg *HttpConfig, handler Handler) *HttpTransport {
	rtr := router.New()
	srv := &fasthttp.Server{
		Handler: rtr.Handler,
	}

	rtr.POST("/api/notify/message/template", func(ctx *fasthttp.RequestCtx) {

	})

	rtr.POST("/api/notify/message/custom", func(ctx *fasthttp.RequestCtx) {

	})

	return &HttpTransport{
		cfg:    cfg,
		srv:    srv,
		router: rtr,
	}
}

func (s *HttpTransport) Register() {

}

func (s *HttpTransport) Run(errCh chan error) {
	go func() {
		if s.cfg.Tls.CertFile != "" && s.cfg.Tls.KeyFile != "" {
			errCh <- s.srv.ListenAndServeTLS(s.cfg.Addr, s.cfg.Tls.CertFile, s.cfg.Tls.KeyFile)
		} else {
			errCh <- s.srv.ListenAndServe(s.cfg.Addr)
		}
	}()
}

func (s *HttpTransport) Shutdown() error {
	return s.srv.Shutdown()
}
