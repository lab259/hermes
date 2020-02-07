package hermes

import (
	"context"

	"github.com/lab259/go-rscsrv"
	"github.com/valyala/fasthttp"
)

// ServiceConfig TODO
type ServiceConfig struct {
	Name   string
	Router Router

	Bind string
	TLS  *FasthttpServiceConfigurationTLS
}

// Service TODO
type Service interface {
	rscsrv.Service
	rscsrv.StartableWithContext
}

type service struct {
	config ServiceConfig
	server fasthttp.Server
}

func (srv *service) listenAndServe() error {
	srv.server.Handler = srv.config.Router.Handler()

	if srv.config.TLS == nil {
		return srv.server.ListenAndServe(srv.config.Bind)
	}

	return srv.server.ListenAndServeTLS(srv.config.Bind, srv.config.TLS.CertFile, srv.config.TLS.KeyFile)
}

func (srv *service) StartWithContext(ctx context.Context) (err error) {
	done := make(chan bool, 1)

	go func() {
		<-ctx.Done()
		err = srv.server.Shutdown()
		close(done)
	}()

	if err := srv.listenAndServe(); err != nil {
		return err
	}

	<-done
	return
}

func (srv *service) Name() string {
	return srv.config.Name
}

// NewService TODO
func NewService(config ServiceConfig) Service {
	return &service{
		config: config,
	}
}
