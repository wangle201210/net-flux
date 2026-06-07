package network

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dellinger2023/net-flux/pkg/logger"
)

const (
	defaultName = "tcp-server"
)

type TcpServer struct {
	sync.RWMutex
	nextID   uint32
	addr     string
	name     string
	handler  EventHandler
	ctx      context.Context
	cancel   context.CancelFunc
	conns    map[uint32]TCPConn
	outgoing chan TCPConn
	options  *TCPConnOptions
}

func NewTcpServer(addr string, handler EventHandler, opts *TCPConnOptions) *TcpServer {
	return &TcpServer{
		nextID:   1,
		addr:     addr,
		name:     defaultName,
		handler:  handler,
		conns:    make(map[uint32]TCPConn),
		outgoing: make(chan TCPConn),
		options:  opts,
	}
}

func (s *TcpServer) Run(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	logger.Infof("[%s] is listening on %s", s.name, s.addr)

	acceptCh := make(chan net.Conn)
	acceptErr := make(chan error, 1)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				acceptErr <- err
				return
			}
			select {
			case acceptCh <- conn:
			case <-s.ctx.Done():
				conn.Close()
				return
			}
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			logger.Infof("[%s] is shutting down...", s.name)
			listener.Close()
			<-acceptErr
			s.closeConns()
			return nil

		case conn := <-acceptCh:
			if err := s.registerConn(conn); err != nil {
				logger.Errorf("failed to create connection: %v", err)
				conn.Close()
			}

		case conn := <-s.outgoing:
			s.Lock()
			delete(s.conns, conn.ID())
			s.Unlock()
			conn.Close()

		case err := <-acceptErr:
			if s.ctx.Err() != nil {
				s.closeConns()
				return nil
			}
			if err != nil && !errors.Is(err, net.ErrClosed) {
				return err
			}
			return nil
		}
	}
}

func (s *TcpServer) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *TcpServer) closeConns() {
	s.Lock()
	defer s.Unlock()
	for _, conn := range s.conns {
		conn.Close()
	}
}

func (s *TcpServer) registerConn(conn net.Conn) error {
	logger.Debugf("new connection from %s", conn.RemoteAddr().String())

	opts := TCPConnOptions{
		HeartbeatInterval: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReconnectInterval: 10 * time.Second,
		ReconnectMax:      3,
		Handler:           s.handler,
		OutgoingCh:        s.outgoing,
	}
	if s.options != nil {
		opts = *s.options
		opts.Handler = s.handler
		opts.OutgoingCh = s.outgoing
	}
	opts.ID = atomic.AddUint32(&s.nextID, 1)

	socket, err := NewServerConn(conn, opts)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()
	s.conns[socket.ID()] = socket
	return nil
}
