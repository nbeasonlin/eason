package common

import (
	"crypto/tls"
	"io"
	"net"

	"github.com/gorilla/websocket"
)

type Connection interface {
	Write([]byte) error
	Read() ([]byte, error)
	Close() error
}
type readWriteCloserConnectionImpl struct {
	rwc io.ReadWriteCloser
}

func NewReadWriteCloserConnection(rwc io.ReadWriteCloser) Connection {
	connection := new(readWriteCloserConnectionImpl)

	connection.rwc = rwc

	return connection
}
func intToByteArray(num int) []byte {
	b := make([]byte, 4)

	base := 256

	for i := 0; num > 0; num /= base {
		b[i] = byte(num % base)
		i++
	}

	return b
}
func byteArrToInt(b []byte) int {
	num := 0
	base := 256
	for i := len(b) - 1; i >= 0; i-- {
		num *= base
		num += int(b[i])
	}
	return num
}

func (r *readWriteCloserConnectionImpl) write(p []byte) error {
	unwrite := len(p)
	begin := 0
	for unwrite > 0 {
		num, err := r.rwc.Write(p[begin : begin+unwrite])
		unwrite -= num
		begin += num
		if err != nil && unwrite > 0 {
			return err
		}
	}
	return nil
}

func (t *readWriteCloserConnectionImpl) Write(p []byte) error {
	return t.write(append(intToByteArray(len(p)), p...))
}

// must read length of p bytes data when calling the function in one time
func (r *readWriteCloserConnectionImpl) read(p []byte) error {
	unread := len(p)
	begin := 0
	for unread > 0 {
		num, err := r.rwc.Read(p[begin : begin+unread])
		unread -= num
		begin += num
		if err != nil && unread > 0 {
			return err
		}
	}
	return nil
}

func (r *readWriteCloserConnectionImpl) Read() ([]byte, error) {
	lArr := make([]byte, 4)
	if err := r.read(lArr); err != nil {
		return nil, err
	}
	lData := byteArrToInt(lArr)
	data := make([]byte, lData)
	if err := r.read(data); err != nil {
		return nil, err
	}
	return data, nil
}
func (r *readWriteCloserConnectionImpl) Close() error {
	return r.rwc.Close()
}

func NewTCPConnection(tcp net.Conn) Connection {
	return NewReadWriteCloserConnection(tcp)
}

func NewTLSConnection(conn *tls.Conn) Connection {
	return NewReadWriteCloserConnection(conn)
}

type webSocketConnectionImpl struct {
	ws *websocket.Conn
}

func NewWebSocketConnection(ws *websocket.Conn) Connection {
	connection := new(webSocketConnectionImpl)

	connection.ws = ws

	return connection
}
func (w *webSocketConnectionImpl) Read() ([]byte, error) {
	_, data, err := w.ws.ReadMessage()
	if err == websocket.ErrCloseSent {
		err = io.EOF
	}
	return data, err
}
func (w *webSocketConnectionImpl) Write(p []byte) error {
	return w.ws.WriteMessage(websocket.BinaryMessage, p)
}
func (w *webSocketConnectionImpl) Close() error {
	return w.ws.Close()
}
