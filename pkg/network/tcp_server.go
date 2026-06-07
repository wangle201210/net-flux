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
	defaultName         = "tcp-server"
	defaultShutdownWait = 5 * time.Second
	closeConnsTimeout   = 2 * time.Second
)

type TcpServer struct {
	sync.RWMutex
	nextID       uint32
	addr         string
	name         string
	handler      EventHandler
	ctx          context.Context
	cancel       context.CancelFunc
	shuttingDown atomic.Bool
	connWG       sync.WaitGroup
	conns        map[uint32]TCPConn
	outgoing     chan TCPConn
	options      *TCPConnOptions
}

func NewTcpServer(addr string, handler EventHandler, opts *TCPConnOptions) *TcpServer {
	return &TcpServer{
		nextID:   1,
		addr:     addr,
		name:     defaultName,
		handler:  handler,
		conns:    make(map[uint32]TCPConn),
		outgoing: make(chan TCPConn, 16),
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
				select {
				case acceptErr <- err:
				default:
				}
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
			return s.shutdown(listener, acceptErr)

		case conn := <-acceptCh:
			if s.shuttingDown.Load() {
				conn.Close()
				continue
			}
			if err := s.registerConn(conn); err != nil {
				logger.Errorf("failed to create connection: %v", err)
				conn.Close()
			}

		case conn := <-s.outgoing:
			s.removeConn(conn)
			if !s.shuttingDown.Load() {
				conn.Close()
			}

		case err := <-acceptErr:
			if s.ctx.Err() != nil || s.shuttingDown.Load() {
				return s.shutdown(listener, acceptErr)
			}
			if err != nil && !errors.Is(err, net.ErrClosed) {
				return err
			}
			return nil
		}
	}
}

func (s *TcpServer) shutdown(listener net.Listener, acceptErr chan error) error {
	if !s.shuttingDown.CompareAndSwap(false, true) {
		return nil
	}

	defer logger.Infof("[%s] shutdown complete", s.name)

	logger.Infof("[%s] is shutting down...", s.name)
	listener.Close()

	s.closeConns()
	s.waitConns()

	select {
	case <-acceptErr:
	default:
	}

	s.drainOutgoing()
	return nil
}

func (s *TcpServer) waitConns() {
	done := make(chan struct{})
	go func() {
		s.connWG.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(defaultShutdownWait):
		logger.Warningf("[%s] timed out waiting for connections to exit", s.name)
	}
}

func (s *TcpServer) drainOutgoing() {
	for {
		select {
		case conn := <-s.outgoing:
			s.removeConn(conn)
		default:
			return
		}
	}
}

func (s *TcpServer) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *TcpServer) removeConn(conn TCPConn) {
	s.Lock()
	delete(s.conns, conn.ID())
	s.Unlock()
}

func (s *TcpServer) closeConns() {
	s.Lock()
	conns := make([]TCPConn, 0, len(s.conns))
	for _, conn := range s.conns {
		conns = append(conns, conn)
	}
	clear(s.conns)
	s.Unlock()

	done := make(chan struct{})
	go func() {
		for _, conn := range conns {
			if tc, ok := conn.(*tcpConn); ok {
				tc.abortClose(true)
			} else {
				conn.Close()
			}
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(closeConnsTimeout):
		logger.Warningf("[%s] timed out closing connections", s.name)
	}
}

func (s *TcpServer) registerConn(conn net.Conn) error {
	logger.Debugf("new connection from %s", conn.RemoteAddr().String())

	opts := TCPConnOptions{
		HeartbeatInterval: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		Handler:           s.handler,
		OutgoingCh:        s.outgoing,
		ParentCtx:         s.ctx,
		ReaderWG:          &s.connWG,
	}
	if s.options != nil {
		opts = *s.options
		opts.Handler = s.handler
		opts.OutgoingCh = s.outgoing
		opts.ParentCtx = s.ctx
		opts.ReaderWG = &s.connWG
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
