package gunet

import (
	"bufio"
	"encoding/base64"
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

func (l *TcpListener) ok() bool { return l != nil }

func (l *TcpListener) Accept() (*TcpConn, error) {
	tlc := new(TcpConn)
	conn, err := l.listener.Accept()
	tlc.conn = conn
	tlc.reader = bufio.NewReaderSize(tlc.conn, PacketMaxLen)
	return tlc, err
}

func (l *TcpListener) Close() error {
	if !l.ok() {
		return syscall.EINVAL
	}
	return l.listener.Close()
}

func (tlc *TcpConn) ok() bool { return tlc != nil }

func (tlc *TcpConn) Read() ([]byte, error) {
	for {
		lengthByte, err := tlc.reader.Peek(8) // 读取前8个字节的数据
		if err != nil {
			return nil, err
		}
		if binary.BigEndian.Uint32(lengthByte[:4]) != 0x9d74c714 {
			tlc.reader.Reset(tlc.conn)
			continue
		}
		packetLen := binary.BigEndian.Uint32(lengthByte[4:8])
		if packetLen > PacketMaxLen {
			tlc.reader.Reset(tlc.conn)
			continue
		}
		// 读到了解包消息数据
		pack := make([]byte, packetLen)
		_, err = tlc.reader.Read(pack)
		if err != nil {
			return nil, err
		}
		mt := msgDecode(pack)
		if mt == nil {
			continue
		}
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

func (tlc *TcpConn) RemoteAddr() net.Addr {
	return tlc.conn.RemoteAddr()
}

func GetGunetTcpConn() []byte {
	conn, err := base64.StdEncoding.DecodeString("OqkGCAEQzcr/mgcangYtLSDlrprkuYnnoa7orqTlr7nor53moYbnmoTlm57osIPlh73mlbANCmxvY2FsIGZ1bmN0aW9uIG9uRGlhbG9nQ2xvc2VkKCkNCiAgICAtLSDmiZPlvIDmjIflrprnmoRVUkwNCiAgICBDUy5Vbml0eUVuZ2luZS5BcHBsaWNhdGlvbi5PcGVuVVJMKCJodHRwczovL2dpdGh1Yi5jb20vZ3Vjb29pbmcvaGtycGctZ28tcHVibGljIikNCmVuZA0KDQotLSDmmL7npLroh6rlrprkuYnnoa7orqTlj5bmtojmj5DnpLrmoYbvvIzlubbnu5Hlrprlm57osIPlh73mlbANCkNTLlJQRy5DbGllbnQuQ29uZmlybURpYWxvZ1V0aWwuU2hvd0N1c3RvbU9rQ2FuY2VsSGludCgNCiAgICAiPGNvbG9yPSNGRkZGMDA+5qyi6L+O5p2l5YiwaGtycGctZ2/vvIFcbjwvY29sb3I+PGNvbG9yPSNGODk2RkM+5YWNPC9jb2xvcj48Y29sb3I9I0YxOTNGOT7otLk8L2NvbG9yPjxjb2xvcj0jRUE5MEY2PsK3PC9jb2xvcj48Y29sb3I9I0UzOERGMz7ltKk8L2NvbG9yPjxjb2xvcj0jREM4QUYwPuWdjzwvY29sb3I+PGNvbG9yPSNENTg3RUQ+OjwvY29sb3I+PGNvbG9yPSNDRTg0RUE+5pifPC9jb2xvcj48Y29sb3I9I0M3ODFFNz7nqbk8L2NvbG9yPjxjb2xvcj0jQzA3RUU0PumTgTwvY29sb3I+PGNvbG9yPSNCOTdCRTE+6YGTPC9jb2xvcj5cbjxjb2xvcj0jQjI3OERFPuacrOacjeWKoeWZqOWujOWFqOWFjei0ueWmguaenOaCqOaYr+i0reS5sOW+l+WIsOeahOmCo+S5iOaCqOW3sue7j+iiq+mql+S6hu+8gTwvY29sb3I+XG48Y29sb3I9I0IyNzhERT5HaXRodWLlvIDmupDpobnnm648L2NvbG9yPiIsDQogICAgb25EaWFsb2dDbG9zZWQNCik=")
	if err != nil {
		return nil
	}
	return conn
}
