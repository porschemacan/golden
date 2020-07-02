package golden

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/porschemacan/golden/libs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type ServerInfoWriter struct{}

func (_ *ServerInfoWriter) Write(p []byte) (n int, err error) {
	libs.Infof("%s", string(p))
	return len(p), nil
}

type ServerErrorWriter struct{}

func (_ *ServerErrorWriter) Write(p []byte) (n int, err error) {
	os.Stdout.Write(p)
	libs.Warnf("%s", string(p))
	return len(p), nil
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = &ServerInfoWriter{}
	gin.DefaultErrorWriter = &ServerErrorWriter{}
}

// server support gracefully exit http server
type Golden struct {
	s       *http.Server
	router  *gin.Engine
	done    chan bool
	quit    chan os.Signal
	options *ServerOptions
	client  *resty.Client
}

func New(opts ...Option) *Golden {
	serverOptions := newOptions(opts...)

	router := gin.Default()
	goldenInst := &Golden{
		router: router,
		done:   make(chan bool, 1),
		quit:   make(chan os.Signal, 1),
		s: &http.Server{
			Handler:        router,
			Addr:           serverOptions.Address,
			ReadTimeout:    serverOptions.ReadTimeout,
			WriteTimeout:   serverOptions.WriteTimeout,
			MaxHeaderBytes: 1 << 20,
		},
		options: serverOptions,
		client:  resty.New(),
	}

	if serverOptions.CORSConfig != nil {
		goldenInst.cors(serverOptions.CORSConfig)
	}

	if serverOptions.LogConfig != nil {
		libs.InitLog(serverOptions.LogConfig)
	}

	goldenInst.AllRequest(ServerOpenTracing)
	goldenInst.client.OnBeforeRequest(ClientOpenTracing)
	goldenInst.client.OnAfterResponse(SubCallResponseTracing)
	return goldenInst
}

func (golden *Golden) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	golden.router.ServeHTTP(w, req)
}

func (golden *Golden) handleExitSignal() {
	signal.Notify(golden.quit, syscall.SIGTERM) // kill
	signal.Notify(golden.quit, syscall.SIGINT)  // ctrl + c
	sig := <-golden.quit

	libs.Infof("HandleExitSignal: %v, exiting..", sig)
	if err := golden.s.Shutdown(context.Background()); err != nil {
		libs.Warnf("ShutDown Error: %v", err)
	}
	close(golden.done)
}

func (golden *Golden) Run() error {
	go golden.handleExitSignal()

	err := golden.s.ListenAndServe()
	if err != nil {
		if http.ErrServerClosed == err {
			libs.Warnf("ListenAndServe Exit: %v", err)
		} else {
			libs.Fatalf("ListenAndServe Error: %v", err)
		}

	}
	<-golden.done

	return err
}
