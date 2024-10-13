package http_proxy

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

type getCertFn func(hostname string) (*tls.Config, error)

type interceptHandler struct {
	listener channelListener
	server   *http.Server
	getCert  getCertFn
}

func newInterceptHandler(getCert getCertFn, innerHandler http.HandlerFunc, wg *sync.WaitGroup, closer chan struct{}) *interceptHandler {
	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Scheme = "https"
			r.URL.Host = r.Host
			innerHandler(w, r)
		}),
	}

	listener := channelListener(make(chan net.Conn))

	go func() {
		// returns always a non-nil error if the server is not closed/shtudown
		wg.Add(1)
		defer wg.Done()
		err := server.Serve(listener)
		if err != nil {
			log.WithError(err).Error("error serving intercept")
		}
	}()

	go func() {
		<-closer
		server.Close()
		if err := server.Shutdown(context.TODO()); err != nil {
			log.Fatalf("HTTP server Shutdown: %v", err)
		}
	}()

	return &interceptHandler{
		listener: listener,
		server:   server,
		getCert:  getCert,
	}
}

func (i *interceptHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		log.Printf("split host port failed '%s': %s", r.Host, err)
		http.Error(w, http.StatusText(http.StatusBadRequest)+err.Error(), http.StatusBadRequest)
		return
	}

	tlsConfig, err := i.getCert(host)
	if err != nil {
		log.Println("failed to obtain tls config:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError)+err.Error(), http.StatusInternalServerError)
		return
	}

	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Print("hijack of connection failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError)+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	clientConn, _, err := hj.Hijack()
	if err != nil {
		log.Println("hijack failed:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	tlsConn := tls.Server(clientConn, tlsConfig)
	i.handleConnection(tlsConn)
}

func (i *interceptHandler) handleConnection(c net.Conn) {
	i.listener <- c
}

func (i *interceptHandler) Close() {
	i.server.Close()
}

// channelListener allows to send connection into a listener through a channel
type channelListener chan net.Conn

func (cl channelListener) Accept() (net.Conn, error) {
	return <-cl, nil
}

func (cl channelListener) Addr() net.Addr {
	return nil
}

func (cl channelListener) Close() error {
	close(cl)
	return nil
}
