package network

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"sync"

	"github.com/dellinger2023/net-flux/pkg/logger"
)

type tcpReader struct {
	conn      TCPConn
	br        *bufio.Reader
	hdrBuf    [HeaderSize]byte
	buff      []byte
	closeOnce sync.Once
	handler   connHandler
}

func newReader(conn TCPConn, handler connHandler) (TCPReader, error) {
	if conn == nil || conn.IsClosed() {
		return nil, ErrNoReady
	}
	if handler == nil {
		return nil, errors.New("handler is required")
	}
	r := &tcpReader{
		conn:    conn,
		br:      bufio.NewReaderSize(conn, defaultReadBufferSize),
		buff:    make([]byte, 0, defaultReadBufferSize),
		handler: handler,
	}
	go r.loop()
	return r, nil
}

func (r *tcpReader) loop() {
	defer func() {
		if rc := recover(); rc != nil {
			logger.Errorf("[id=%d] readLoop panic: %v", r.conn.ID(), rc)
		}
		r.conn.Close()
		logger.Infof("[id=%d] readLoop exited", r.conn.ID())
	}()
	for {
		head, payload, err := r.ReadPacket()
		if err != nil {
			if !errors.Is(err, io.EOF) && !errors.Is(err, ErrClosed) && !r.conn.IsClosed() {
				logger.Errorf("[id=%d] readLoop error: %v", r.conn.ID(), err)
			}
			return
		}
		if r.handler != nil {
			if err := r.handler.OnMessage(r.conn, head, payload); err != nil {
				logger.Errorf("[id=%d] dispatch error: %v", r.conn.ID(), err)
			}
		}
	}
}

func parseHeader(buf []byte) *PacketHeader {
	return &PacketHeader{
		Length:   binary.BigEndian.Uint16(buf[0:2]),
		CMD:      buf[2],
		SubCMD:   buf[3],
		Version:  buf[4],
		Flags:    binary.BigEndian.Uint16(buf[5:7]),
		Reserved: buf[7],
	}
}

// ReadPacket reads the next packet. The returned payload aliases the internal
// read buffer and is only valid until the next ReadPacket call.
func (r *tcpReader) ReadPacket() (*PacketHeader, []byte, error) {
	if r.conn == nil || r.conn.IsClosed() {
		return nil, nil, ErrNoReady
	}

	if _, err := io.ReadFull(r.br, r.hdrBuf[:]); err != nil {
		return nil, nil, err
	}

	header := parseHeader(r.hdrBuf[:])
	if header.Length > uint16(maxPacketSize) {
		return nil, nil, errors.New("packet too large")
	}

	payloadLen := int(header.Length)
	if payloadLen == 0 {
		return header, nil, nil
	}

	if cap(r.buff) < payloadLen {
		r.buff = make([]byte, payloadLen)
	} else {
		r.buff = r.buff[:payloadLen]
	}

	if _, err := io.ReadFull(r.br, r.buff); err != nil {
		return nil, nil, err
	}

	return header, r.buff, nil
}

func (r *tcpReader) Close() {
	r.closeOnce.Do(func() {
		r.conn.Close()
	})
}
