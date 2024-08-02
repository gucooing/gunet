package gunet

import (
	"encoding/binary"
	"errors"
	"log"
)

const (
	PacketMaxLen = 512 * 1024
)

type msg struct {
	data []byte
}

// 打包
func msgEncode(msg *msg) (bin []byte, err error) {
	if msg.data == nil {
		msg.data = make([]byte, 0)
	}
	packetLen := len(msg.data) + 12
	if packetLen > PacketMaxLen {
		return nil, errors.New("packet too big")
	}
	bin = make([]byte, packetLen)
	// 数据长度
	binary.BigEndian.PutUint32(bin[:4], 0x9d74c714)
	binary.BigEndian.PutUint32(bin[4:], uint32(packetLen))
	binary.BigEndian.PutUint32(bin[len(bin)-4:], 0xd7a152c8)
	// 数据
	copy(bin[8:], msg.data)
	return bin, nil
}

// 解包
func msgDecode(data []byte) *msg {
	newMsg := new(msg)
	// 长度太短
	if len(data) < 12 {
		log.Println("packet len less than 12 byte")
		return nil
	}
	// 头部幻数错误
	if binary.BigEndian.Uint32(data[:4]) != 0x9d74c714 {
		log.Println("packet head magic 0x9d74c714 error")
		return nil
	}
	// proto长度
	protoLen := binary.BigEndian.Uint32(data[4:8])
	// 检查长度
	if protoLen > PacketMaxLen {
		log.Println("packet len too long")
		return nil
	}
	if len(data) < int(protoLen) {
		log.Println("packet len not enough")
		return nil
	}
	// 尾部幻数错误
	if binary.BigEndian.Uint32(data[len(data)-4:]) != 0xd7a152c8 {
		log.Println("packet tail magic 0xd7a152c8 error")
		return nil
	}
	// 数据
	// proto数据
	protoData := data[8 : int(protoLen)-4]
	newMsg.data = protoData
	return newMsg
}
