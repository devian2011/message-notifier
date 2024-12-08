package transport

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"notifier/internal/entity"

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

	type responseData struct {
		Error error     `json:"error"`
		Code  int       `json:"code"`
		Id    uuid.UUID `json:"id"`
	}

	rtr.POST("/api/v1/notify/message/template", func(ctx *fasthttp.RequestCtx) {
		command := &entity.TemplateCommand{}
		unmarshalErr := json.Unmarshal(ctx.Request.Body(), command)
		if unmarshalErr != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			rData, _ := json.Marshal(responseData{
				Error: unmarshalErr,
				Code:  http.StatusBadRequest,
			})
			ctx.Response.SetBody(rData)

			return
		}

		id, err := handler.HandleTemplateMessage(*command)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			rData, _ := json.Marshal(responseData{
				Error: err,
				Code:  http.StatusBadRequest,
			})
			ctx.Response.SetBody(rData)

			return
		}
		ctx.SetStatusCode(http.StatusOK)
		rData, _ := json.Marshal(responseData{
			Error: nil,
			Code:  http.StatusOK,
			Id:    id,
		})
		ctx.Response.SetBody(rData)
	})

	rtr.POST("/api/v1/notify/message/custom", func(ctx *fasthttp.RequestCtx) {
		command := &entity.CustomCommand{}
		unmarshalErr := json.Unmarshal(ctx.Request.Body(), command)
		if unmarshalErr != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			rData, _ := json.Marshal(responseData{
				Error: unmarshalErr,
				Code:  http.StatusBadRequest,
			})
			ctx.Response.SetBody(rData)

			return
		}

		id, err := handler.HandleCustomMessage(*command)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			rData, _ := json.Marshal(responseData{
				Error: err,
				Code:  http.StatusBadRequest,
			})
			ctx.Response.SetBody(rData)

			return
		}
		ctx.SetStatusCode(http.StatusOK)
		rData, _ := json.Marshal(responseData{
			Error: nil,
			Code:  http.StatusOK,
			Id:    id,
		})
		ctx.Response.SetBody(rData)
	})

	return &HttpTransport{
		cfg:    cfg,
		srv:    srv,
		router: rtr,
	}
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
