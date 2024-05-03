package gunet

import (
	"encoding/binary"
	"errors"
	"net"
	"syscall"
)

type TcpListener struct {
	listener net.Listener
}

type TcpConn struct {
	conn    net.Conn
	tme     []byte
	tmeList []*msg
	is      bool
}

func NewTcpS(address string) (*TcpListener, error) {
	l := new(TcpListener)
	listener, err := net.Listen("tcp", address)
	l.listener = listener
	return l, err
}

func (l *TcpListener) Accept() (*TcpConn, error) {
	tlc := new(TcpConn)
	conn, err := l.listener.Accept()
	tlc.conn = conn
	tlc.is = true
	go tlc.read()
	return tlc, err
}

func (tlc *TcpConn) ok() bool { return tlc != nil }

func (tlc *TcpConn) Read() ([]byte, error) {
	if tlc.tmeList == nil {
		tlc.tmeList = make([]*msg, 0)
	}
	for {
		if tlc.ok() {
			if len(tlc.tmeList) > 0 {
				for id, m := range tlc.tmeList {
					tlc.tmeList = append(tlc.tmeList[:id], tlc.tmeList[id+1:]...)
					return m.data, nil
				}
			} else if !tlc.is {
				break
			}
		} else {
			break
		}
	}
	return nil, errors.New("tcp conn closed")
}

func (tlc *TcpConn) read() {
	if tlc.tme == nil {
		tlc.tme = make([]byte, 0)
	}
	if tlc.tmeList == nil {
		tlc.tmeList = make([]*msg, 0)
	}
	for {
		bin := make([]byte, PacketMaxLen)
		if !tlc.ok() {
			return
		}
		n, err := tlc.conn.Read(bin)
		if err != nil {
			tlc.is = false
			return
		}
		if n == 0 {
			continue
		}
		tlc.tme = append(tlc.tme, bin[:n]...)
		if len(tlc.tme) < 4 {
			continue
		}
		tlc.dwtme()
	}
}

func (tlc *TcpConn) dwtme() {
	packetLen := binary.BigEndian.Uint32(tlc.tme[0:4])
	if packetLen > uint32(len(tlc.tme)) {
		return
	}
	mt := msgDecode(tlc.tme[:packetLen])
	if mt == nil {
		return
	}
	tlc.tmeList = append(tlc.tmeList, mt)
	tlc.tme = tlc.tme[packetLen:]
	if len(tlc.tme) < 4 {
		return
	}
	packetLen2 := binary.BigEndian.Uint32(tlc.tme[0:4])
	if packetLen2 <= uint32(len(tlc.tme)) {
		tlc.dwtme() //  保证每次都可以处理完完整的数据包
	}
}

func NewTcpC(address string) (*TcpConn, error) {
	tlc := new(TcpConn)
	conn, err := net.Dial("tcp", address)
	tlc.conn = conn
	tlc.is = true
	return tlc, err
}

func (tlc *TcpConn) Write(b []byte) (int, error) {
	if !tlc.ok() {
		return 0, syscall.EINVAL
	}
	bin := msgEncode(&msg{data: b})
	return tlc.conn.Write(bin)
}

func (tlc *TcpConn) Close() error {
	if !tlc.ok() {
		return syscall.EINVAL
	}
	tlc.is = false
	return tlc.conn.Close()
}
