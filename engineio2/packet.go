package engineioclient2

import (
	// "encoding/json"
	//"log"
	"github.com/googollee/go-engine.io/parser"
)

type PacketType byte

type Packet struct {
	Type parser.PacketType
	Data []byte
}

const (
	Open    PacketType = '0'
	Close   PacketType = '1'
	Ping    PacketType = '2'
	Pong    PacketType = '3'
	Message PacketType = '4'
	Upgrade PacketType = '5'
	Noop    PacketType = '6'
)

func NewEmptyPacket(pType parser.PacketType) *Packet {
	return &Packet{pType, make([]byte, 0)}
}

func NewPacket(content []byte) *Packet {
	return &Packet{parser.MESSAGE, content}
}

func NewClosePacket() *Packet {
	return &Packet{parser.CLOSE, nil}
}

func PacketToBytes(from *Packet) (to []byte) {
	to = make([]byte, len(from.Data)+1)
	copy(to[1:], from.Data)
	//to[0] = byte(from.Type)
	//log.Println(string(to))
	return
}
