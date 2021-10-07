package httputil

import (
	"context"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

type Server struct {
	handler    http.Handler
	network    string
	addr       string
	httpServer *http.Server
}

func (s *Server) Start(ctx context.Context) error {
	lc := net.ListenConfig{
		Control: func(network, address string, conn syscall.RawConn) error {
			var operr error
			if err := conn.Control(func(fd uintptr) {
				operr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
				if operr != nil {
					return
				}
				operr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			}); err != nil {
				return err
			}
			return operr
		},
	}

	l, err := lc.Listen(ctx, s.network, s.addr)
	if err != nil {
		return err
	}

	s.httpServer = &http.Server{
		Handler: s.handler,
	}

	errs := make(chan error, 1)
	go func() {
		errs <- s.httpServer.Serve(l)
	}()

	select {
	case <-ctx.Done():
		go s.Stop(context.Background())
		return ctx.Err()
	case err := <-errs:
		return err
	case <-time.After(2 * time.Second):
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func NewServer(network, addr string, handler http.Handler) *Server {
	// REUSEADDR and REUSEPORT are not defined on unix socket.
	// need to remove the addrs first
	if network == "unix" {
		err := os.Remove(addr)
		if err != nil && !os.IsNotExist(err) {
			panic(err)
		}
	}

	return &Server{
		handler: handler,
		network: network,
		addr:    addr,
	}
}
