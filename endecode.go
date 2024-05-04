package gunet

import (
	"encoding/binary"
	"errors"
)

const (
	PacketMaxLen = 345 * 1024 // 最大应用层包长度
)

type msg struct {
	data []byte
}

// 打包
func msgEncode(msg *msg) (bin []byte, err error) {
	if msg.data == nil {
		msg.data = make([]byte, 0)
	}
	packetLen := len(msg.data) + 4
	if packetLen > PacketMaxLen {
		return nil, errors.New("packet too big")
	}
	bin = make([]byte, packetLen)
	// 数据长度
	binary.BigEndian.PutUint32(bin[0:4], uint32(packetLen))
	// 数据
	copy(bin[4:], msg.data)
	return bin, nil
}

// 解包
func msgDecode(data []byte) *msg {
	newMsg := new(msg)
	// 长度太短
	if len(data) < 6 {
		return nil
	}
	// 数据
	datas := data[4:]
	newMsg.data = datas
	return newMsg
}
