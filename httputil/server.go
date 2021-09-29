package httputil

import (
	"context"
	"net"
	"net/http"
	"os"
	"syscall"
)

type Server struct {
	handler http.Handler
	network string
	addr    string
	closed  chan struct{}
}

func (s *Server) Start(ctx context.Context) error {
	errs := make(chan error, 1)

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

	server := &http.Server{
		Handler: s.handler,
	}

	go func() {
		errs <- server.Serve(l)
	}()

	select {
	case <-ctx.Done():
		return server.Shutdown(ctx)
	case err := <-errs:
		return err
	case <-s.closed:
		return server.Shutdown(ctx)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		close(s.closed)
		return nil
	}
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
		closed:  make(chan struct{}, 1),
	}
}
