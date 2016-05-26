package engineioclient2

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/googollee/go-engine.io/message"
	"github.com/googollee/go-engine.io/parser"
	"github.com/googollee/go-engine.io/transport"
	eiowebsocket "github.com/googollee/go-engine.io/websocket"
	"github.com/gorilla/websocket"
)

type Event struct {
	Type string
	Data []byte
}

var closeErrorType = reflect.TypeOf((*websocket.CloseError)(nil))

type Client struct {
	Conn   transport.Client
	sender chan *Packet
	//receiver chan *Packet
	Event chan *Event
	//done     chan bool
	opened   bool
	eventMap map[string]interface{}
	lastPing time.Time
}

func NewClient(conn transport.Client) *Client {
	return &Client{
		Conn:   conn,
		opened: false,
		sender: make(chan *Packet, 10),
		//receiver: make(chan *Packet, 5),
		Event: make(chan *Event, 5),
		//done:     make(chan bool),
		eventMap: make(map[string]interface{}),
		lastPing: time.Now(),
	}
}

func Dial(urlStr string, extras map[string]string) (client *Client, err error) {
	if nil == extras {
		extras = make(map[string]string)
	}
	if nil != extras {
		extras["transport"] = "websocket"
	}
	querystring := ""
	for k, v := range extras {
		querystring += "&" + k + "=" + v
	}
	urlStr += "/socket.io/?EIO=3" + querystring
	//log.Println(urlStr)
	client = NewClient(nil)
	go func() {
		// listens for open event, and then make it open
		for {
			wc := client.Conn
			if wc != nil {
				decoder, err := wc.NextReader()
				if err != nil {
					wc.Close()
					client.Conn = nil
					client.onPacket(NewClosePacket())
					continue
				}
				res, err := ioutil.ReadAll(decoder)
				decoder.Close()
				if nil != err {
					t := reflect.TypeOf(err)
					log.Printf("Reading next Message Error, %s", err)
					if closeErrorType == t {
						wc.Close()
						client.Conn = nil
						// handle reconnections
						client.onPacket(NewClosePacket())
					} else {
						log.Printf("Reading Message %s", err.Error())
						wc.Close()
						client.Conn = nil
						client.onPacket(NewClosePacket())
					}

				} else {
					//log.Printf("Reading next Message Success, %s:%s", decoder.Type(), string(res))
					//log.Printf("%s", string(res))
					//packet := BytesToPacket(res)
					//log.Printf("Decoding packet %q", packet.Type)
					packet := &Packet{Type: decoder.Type(), Data: res}
					client.onPacket(packet)
				}
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
	go func() {
		// reconnection ticker
		ticker := time.NewTicker(1 * time.Second)
		ping := time.NewTicker(20 * time.Second)
		for {
			select {
			case _ = <-ticker.C:
				// do reconnection here if client is not connected
				//log.Printf("reconnecting ticker")
				if nil == client.Conn {
					//log.Printf("reconnecting")
					u, _ := url.Parse(urlStr)
					u.Scheme = "ws"
					req, err := http.NewRequest("GET", u.String(), nil)
					conn, err := eiowebsocket.NewClient(req)
					if nil != err {
						//log.Printf("reconnecting fucked up %s\n", err.Error())
						continue
					} else {
						//log.Printf("reconnecting succeeded")
						client.Conn = conn
					}
				}
			case _ = <-ping.C:
				client.sender <- &Packet{Type: parser.PONG, Data: []byte{}}
			case p := <-client.sender:
				wc := client.Conn
				if wc == nil {
					continue
				}
				//log.Printf("write begin %v:%s", p.Type, string(p.Data))
				w, err := wc.NextWriter(message.MessageText, p.Type)
				if nil != err {
					log.Printf("WriteMessage %s\n", err.Error())
				} else {
					_, err := w.Write(p.Data)
					w.Close()
					if nil != err {
						log.Printf("WriteMessage %s\n", err.Error())
					}
				}
			}
		}
	}()
	return
}

func (c *Client) onPacket(packet *Packet) {
	//log.Printf("Got a packet!!! yeah!!")
	event := &Event{Data: packet.Data}
	switch packet.Type {
	case parser.OPEN:
		c.opened = true
		event.Type = "Open"
		//log.Printf("%s:%s", packet.Type, string(packet.Data))
	case parser.CLOSE:
		c.opened = false
		event.Type = "Close"
	case parser.MESSAGE:
		//log.Printf("Got a message event!!! yeah!!")
		//log.Printf("Got a Message %s", string(event.Data))
		event.Type = "Message"
	case parser.PING:
		event.Type = "Ping"
		now := time.Now()
		//diff := now.Sub(c.lastPing)
		//log.Printf("Got a ping event, replying now delta:%d", diff)
		c.lastPing = now
		c.sender <- &Packet{Type: parser.PONG, Data: packet.Data}
	case parser.PONG:
		event.Type = "Pong"
	default:
		event = nil
	}
	if nil != event {
		//c.Emit(event.Type, event.Data)
		c.Event <- event
	}
}

func (c *Client) On(event string, callback interface{}) {
	c.eventMap[event] = callback
}

/*func (c *Client) Emit(event string, args ...interface{}) {
	if f, ok := c.eventMap[event]; ok {
		fv := reflect.ValueOf(f)
		if fv.Kind() != reflect.Func {
			return
		}
		ft := fv.Type()
		if ft.NumIn() == 0 {
			fv.Call([]reflect.Value{})
			return
		}
		parm := make([]reflect.Type, ft.NumIn())
		for i, n := 0, ft.NumIn(); i < n; i++ {
			parm[i] = ft.In(i)
		}

		a := make([]reflect.Value, len(args))
		for i, arg := range args {
			v := reflect.ValueOf(arg)
			if parm[i].Kind() != reflect.Ptr {
				if v.IsValid() {
					a[i] = v
				} else {
					a[i] = reflect.Zero(parm[i])
				}
			}
		}
		fv.Call(a)
	}
} */

func (c *Client) Emit(event string, v interface{}) error {
	obj := []interface{}{event, v}
	result, err := json.Marshal(obj)
	if nil != err {
		return err
	}
	b := make([]byte, len(result)+1)
	copy(b[1:], result)
	b[0] = byte(_EVENT) + '0'

	c.SendPacket(NewPacket(b))
	return nil
}

func (c *Client) SendMessage(obj interface{}) error {
	result, err := json.Marshal(obj)
	if nil != err {
		return err
	}
	log.Printf(string(result))
	c.sender <- NewPacket(result)
	return nil
}

func (c *Client) SendPacket(packet *Packet) {
	c.sender <- packet
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
