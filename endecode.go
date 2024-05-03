package gunet

import (
	"encoding/binary"
)

const (
	PacketMaxLen = 343 * 1024 // 最大应用层包长度
)

type msg struct {
	data []byte
}

// 打包
func msgEncode(msg *msg) (bin []byte) {
	if msg.data == nil {
		msg.data = make([]byte, 0)
	}
	packetLen := len(msg.data) + 4
	bin = make([]byte, packetLen)
	// 数据长度
	binary.BigEndian.PutUint32(bin[0:4], uint32(packetLen))
	// proto数据
	copy(bin[4:], msg.data)
	return bin
}

// 解包
func msgDecode(data []byte) *msg {
	newMsg := new(msg)
	// 长度太短
	if len(data) < 6 {
		return nil
	}
	// proto数据
	datas := data[4:]
	newMsg.data = datas
	return newMsg
}
