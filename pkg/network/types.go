package network

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	Version    = 1 // 协议版本
	Flags      = 0 // 标志
	Reserved   = 0 // 保留字段
	HeaderSize = 8 // 包头大小

	defaultReadBufferSize  = 4096
	defaultWriteBufferSize = 4096
	maxPacketSize          = 65535 // 协议 Length 字段为 uint16
)

var (
	ErrClosed  = errors.New("connection is closed")
	ErrTimeout = errors.New("timeout")
	ErrNoReady = errors.New("connection is not ready")
)

// 包头定义, 长度固定为8字节
type PacketHeader struct {
	// 包长度
	Length uint16

	// 命令
	CMD uint8

	// 子命令
	SubCMD uint8

	// 版本
	Version uint8

	// 标志
	Flags uint16

	// 保留字段
	Reserved uint8
}

// TCPConnOptions is the options of a connection.
type TCPConnOptions struct {
	// ID is the unique identifier of the connection.
	ID uint32

	// HeartbeatInterval is the interval of heartbeat.
	HeartbeatInterval time.Duration

	// ReadTimeout is the timeout of reading.
	ReadTimeout time.Duration

	// WriteTimeout is the timeout of writing.
	WriteTimeout time.Duration

	ReconnectInterval time.Duration

	// ReconnectMax is the maximum number of reconnecting.
	ReconnectMax int

	// Handler is the handler of the connection.
	Handler EventHandler

	// OutgoingCh is the channel of outgoing connections.
	OutgoingCh chan TCPConn

	// ParentCtx is cancelled when the server/client shuts down.
	ParentCtx context.Context

	// ReaderWG tracks active read loops; server shutdown waits on it.
	ReaderWG *sync.WaitGroup
}

// TCPConn is an interface of methods that are used as callbacks on a connection.
type TCPConn interface {

	// ID is the unique identifier of the connection.
	ID() uint32

	// RemoteAddr is the remote address of the connection.
	RemoteAddr() net.Addr

	// LocalAddr is the local address of the connection.
	LocalAddr() net.Addr

	// SetReadDeadline sets the read deadline of the connection.
	SetReadDeadline(t time.Time) error

	// SetWriteDeadline sets the write deadline of the connection.
	SetWriteDeadline(t time.Time) error

	// SetHeartbeatInterval sets the interval of heartbeat.
	SetHeartbeatInterval(d time.Duration) error

	// SetNoNagle sets the option of no nagle.
	SetNoNagle(enable bool) error

	// SetKeepAlive sets the option of keep alive.
	SetKeepAlive(enable bool) error

	// SetKeepAlivePeriod sets the period of keep alive.
	SetKeepAlivePeriod(d time.Duration) error

	// SetReuseAddr sets the option of reuse addr.
	SetReuseAddr(enable bool) error

	// SetReusePort sets the option of reuse port.
	SetReusePort(enable bool) error

	// returns the raw net.Conn
	RawConn() net.Conn

	// Write writes the data to the connection.
	// It returns the number of bytes written and any error encountered.
	Write(p []byte) (n int, err error)

	// Read reads data from the connection.
	// It returns the number of bytes read and any error encountered.
	Read(p []byte) (n int, err error)

	// Connects to the address.
	Connect(addr string) error

	// ConnectWithTLS opens the connection with TLS.
	ConnectWithTLS(addr string, tlsConfig *tls.Config) error

	// Close closes the connection.
	// It returns any error encountered.
	Close() error

	// IsClosed returns true if the connection is closed.
	IsClosed() bool

	// WritePacket writes a protobuf packet to the connection.
	WritePacket(cmd, subcmd uint8, pkt proto.Message) error
}

// connHandler is an interface of methods that are used as callbacks on a connection
type connHandler interface {
	// OnConnect is called when the connection was accepted,
	// It returns an error if the connection is not accepted.
	OnConnect(TCPConn) error

	// OnMessage is called when the connection receives a message.
	OnMessage(TCPConn, *PacketHeader, []byte) error

	// OnClose is called when the connection closed
	OnClose(TCPConn)
}

type EventHandler interface {
	// OnConnect is called when the connection was accepted,
	// It returns an error if the connection is not accepted.
	OnConnect(TCPConn) error
	// OnClose is called when the connection closed
	OnClose(TCPConn)
	// OnCmdSystem is called when the connection receives a system command.
	OnCmdSystem(TCPConn, proto.Message) error
	// OnCmdDiscovery is called when the connection receives a discovery command.
	OnCmdDiscovery(TCPConn, proto.Message) error
	// OnCmdDataReport is called when the connection receives a data report command.
	OnCmdDataReport(TCPConn, uint8, proto.Message) error
	// OnCmdConfig is called when the connection receives a config command.
	OnCmdConfig(TCPConn, proto.Message) error
	// OnCmdEvent is called when the connection receives a event command.
	OnCmdEvent(TCPConn, proto.Message) error
	// OnCmdControl is called when the connection receives a control command.
	OnCmdControl(TCPConn, proto.Message) error
}

type TCPReader interface {
	// ReadPacket reads a packet from the connection.
	ReadPacket() (*PacketHeader, []byte, error)
	// Close closes the reading loop.
	Close()
}

type TCPWriter interface {
	// WritePacket writes a packet to the connection.
	WritePacket(cmd, subcmd uint8, pkt proto.Message) error
}
