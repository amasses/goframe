package goframe

import (
	"bufio"
	"net"
	"sync"
)

type delimiterBasedFrameConn struct {
	delimiter byte
	c         net.Conn
	r         *bufio.Reader
	w         *bufio.Writer
	m         sync.RWMutex
}

// NewDelimiterBasedFrameConn returns a Frame conn framed with delimiter.
func NewDelimiterBasedFrameConn(delimiter byte, conn net.Conn) FrameConn {
	return &delimiterBasedFrameConn{
		delimiter: delimiter,
		c:         conn,
		r:         bufio.NewReader(conn),
		w:         bufio.NewWriter(conn),
		m:         sync.RWMutex{},
	}
}

func (fc *delimiterBasedFrameConn) ReadFrame() ([]byte, error) {
	var (
		isPrefix bool
		err      error
		line, ln []byte
	)

	fc.m.RLock()
	defer fc.m.RUnlock()

	for isPrefix && err == nil {
		line, err = fc.r.ReadBytes(fc.delimiter)
		ln = append(ln, line...)
		if err != nil {
			return ln, err
		}
	}

	return ln, nil
}

func (fc *delimiterBasedFrameConn) WriteFrame(p []byte) error {
	fc.m.Lock()
	defer fc.m.Unlock()

	_, err := fc.w.Write(p)
	if err != nil {
		return err
	}
	err = fc.w.WriteByte(fc.delimiter)
	if err != nil {
		return err
	}
	fc.w.Flush()
	return nil
}

func (fc *delimiterBasedFrameConn) Close() error {
	return fc.c.Close()
}

func (fc *delimiterBasedFrameConn) Conn() net.Conn {
	return fc.c
}
