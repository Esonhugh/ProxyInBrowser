package http_proxy

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"

	"github.com/esonhugh/proxyinbrowser/cmd/server/define"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestDump, _ := httputil.DumpRequest(r, false)
		log.Debugf("url=%s\n\n%s\n", r.URL, requestDump)
		next.ServeHTTP(w, r)
	}
}

func CreateHttpProxyServer(TargetConn *websocket.Conn, Port string, rch chan define.RelayCommandResp, stop chan struct{}) {
	var (
		caCertFile         = "cert/cert.pem"
		caKeyFile          = "cert/key.pem"
		httpServerExitDone = &sync.WaitGroup{}
	)

	certGen, err := newCertGenerator(caCertFile, caKeyFile)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	wsproxy := createHTTPProxy(TargetConn, rch)
	forwardHandler := wsproxy.ServeHTTP

	connectHandler := newInterceptHandler(certGen.Get, logRequest(forwardHandler), httpServerExitDone, stop)

	handler := logRequest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "CONNECT" {
			connectHandler.ServeHTTP(w, r)
		} else {
			forwardHandler(w, r)
			// wsproxy.ServeHTTP(w, r)
		}
	})

	srv := &http.Server{
		Addr:    ":" + Port,
		Handler: http.HandlerFunc(handler),
	}

	go func() {
		httpServerExitDone.Add(1)
		defer httpServerExitDone.Done()
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) && err != nil {
			log.Error(err)
		}
	}()
	log.Infof("HTTP Proxy server is started on port: %v", Port)
	<-stop // block if receive stop command
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Fatalf("HTTP server Shutdown: %v", err)
	}
	httpServerExitDone.Wait()
	log.Info("HTTP server stopped")
}
