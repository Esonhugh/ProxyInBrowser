package http_proxy

import (
	"context"
	"errors"
	"io"
	"net"
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
		log.Printf("url=%s\n%s", r.URL, requestDump)
		next.ServeHTTP(w, r)
	}
}

func tunnel(w http.ResponseWriter, r *http.Request) {
	dialer := net.Dialer{}
	serverConn, err := dialer.DialContext(r.Context(), "tcp", r.Host)
	if err != nil {
		log.Printf("failed to connect to upstream %s", r.Host)
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	defer serverConn.Close()

	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Print("hijack of connection failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	clientConn, bufClientConn, err := hj.Hijack()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	go io.Copy(serverConn, bufClientConn)
	io.Copy(bufClientConn, serverConn)
}

func forward(w http.ResponseWriter, r *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	for header, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func CreateHttpProxyServer(TargetConn *websocket.Conn, Port string, rch chan define.RelayCommandResp, stop chan struct{}) {
	var (
		caCertFile         = "cert/cert.pem"
		caKeyFile          = "cert/key.pem"
		httpServerExitDone = &sync.WaitGroup{}
	)

	certGen, err := newCertGenerator(caCertFile, caKeyFile)
	if err != nil {
		log.Print(err)
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
