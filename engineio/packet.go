package engineioclient

// 来源： https://github.com/shekhei/go-engine.io-client/blob/master/client_test.go
type PacketType string

type Packet struct {
	Type PacketType
	Data []byte
}

// copy from "github.com/googollee/go-engine.io/parser"
const (
	OPEN    PacketType = "open"
	CLOSE   PacketType = "close"
	PING    PacketType = "ping"
	PONG    PacketType = "pong"
	MESSAGE PacketType = "message"
	UPGRADE PacketType = "upgrade"
	NOOP    PacketType = "noop"
)

/*const (
	Open    PacketType = '0'
	Close   PacketType = '1'
	Ping    PacketType = '2'
	Pong    PacketType = '3'
	Message PacketType = '4'
	Upgrade PacketType = '5'
	Noop    PacketType = '6'
) */

func NewEmptyPacket(pType PacketType) *Packet {
	return &Packet{pType, make([]byte, 0)}
}

func NewPacket(pType PacketType, content []byte) *Packet {
	return &Packet{pType, content}
}

// added by sunny
func NewPacket2(content []byte) *Packet {
	//return &Packet{parser.MESSAGE, content}
	return &Packet{MESSAGE, content}
}

func NewClosePacket() *Packet {
	//return &Packet{Close, nil}
	return &Packet{CLOSE, nil}
}

func BytesToPacket(from []byte) (to *Packet) {
	return &Packet{Type: PacketType(from[0]), Data: from[1:]}
}

func PacketToBytes(from *Packet) (to []byte) {
	/*to = make([]byte, len(from.Data)+1)
	copy(to[1:], from.Data)
	to[0] = byte(from.Type) */
	//log.Println(string(to))
	return
}
