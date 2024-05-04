package gunet

import (
	"bufio"
	"encoding/binary"
	"net"
	"syscall"
)

type TcpListener struct {
	listener net.Listener
}

type TcpConn struct {
	conn   net.Conn
	reader *bufio.Reader
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
	tlc.reader = bufio.NewReaderSize(tlc.conn, PacketMaxLen)
	return tlc, err
}

func (tlc *TcpConn) ok() bool { return tlc != nil }

func (tlc *TcpConn) Read() ([]byte, error) {
	for {
		lengthByte, err := tlc.reader.Peek(4) // 读取前4个字节的数据
		if err != nil {
			return nil, err
		}
		packetLen := binary.BigEndian.Uint32(lengthByte)
		if packetLen > PacketMaxLen { // 太大了清空
			tlc.reader.Reset(tlc.conn)
			continue
		}
		if uint32(tlc.reader.Buffered()) < packetLen { // 太小了再等等
			continue
		}
		// 读到了解包消息数据
		pack := make([]byte, packetLen)
		_, err = tlc.reader.Read(pack)
		if err != nil {
			return nil, err
		}
		mt := msgDecode(pack)
		return mt.data, nil
	}
}

func NewTcpC(address string) (*TcpConn, error) {
	tlc := new(TcpConn)
	conn, err := net.Dial("tcp", address)
	tlc.conn = conn
	tlc.reader = bufio.NewReaderSize(tlc.conn, PacketMaxLen)
	return tlc, err
}

func (tlc *TcpConn) Write(b []byte) (int, error) {
	if !tlc.ok() {
		return 0, syscall.EINVAL
	}
	bin, err := msgEncode(&msg{data: b})
	if err != nil {
		return 0, err
	}
	return tlc.conn.Write(bin)
}

func (tlc *TcpConn) Close() error {
	if !tlc.ok() {
		return syscall.EINVAL
	}
	return tlc.conn.Close()
}
