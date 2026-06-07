package network

import (
	"encoding/binary"
	"errors"
	"io"
	"sync"

	"google.golang.org/protobuf/proto"
)

type tcpWriter struct {
	conn TCPConn
	mu   sync.Mutex
	buff []byte
}

func newWriter(conn TCPConn) (TCPWriter, error) {
	return &tcpWriter{
		conn: conn,
		buff: make([]byte, 0, defaultWriteBufferSize),
	}, nil
}

func (w *tcpWriter) WritePacket(cmd, subcmd uint8, pkt proto.Message) error {
	if w.conn == nil || w.conn.IsClosed() {
		return ErrNoReady
	}

	data, err := proto.Marshal(pkt)
	if err != nil {
		return err
	}

	if len(data) > maxPacketSize {
		return errors.New("packet too large")
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	need := HeaderSize + len(data)

	if cap(w.buff) < need {
		w.buff = make([]byte, need)
	} else {
		w.buff = w.buff[:need]
	}

	binary.BigEndian.PutUint16(w.buff[0:2], uint16(len(data)))
	w.buff[2] = cmd
	w.buff[3] = subcmd
	w.buff[4] = Version
	binary.BigEndian.PutUint16(w.buff[5:7], Flags)
	w.buff[7] = Reserved
	copy(w.buff[HeaderSize:], data)

	_, err = writeAll(w.conn, w.buff)
	return err
}

func writeAll(w io.Writer, p []byte) (int, error) {
	total := 0
	for len(p) > 0 {
		n, err := w.Write(p)
		total += n
		if err != nil {
			return total, err
		}
		p = p[n:]
	}
	return total, nil
}
