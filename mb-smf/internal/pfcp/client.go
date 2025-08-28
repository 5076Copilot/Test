package pfcp

import (
	"encoding/binary"
	"net"
	"time"
)

type Client struct {
	Remote string
}

func NewClient(remote string) *Client { return &Client{Remote: remote} }

// Minimal PFCP header per 29.244 (very simplified)
type Header struct {
	Version   uint8
	Sflag     bool
	MessageType uint8
	Length    uint16
	Seq       uint32
}

func (h Header) MarshalBinary() []byte {
	b := make([]byte, 8)
	b[0] = (h.Version & 0x07) << 5
	if h.Sflag { b[0] |= 0x10 }
	b[1] = h.MessageType
	binary.BigEndian.PutUint16(b[2:4], h.Length)
	// simplified sequence in 3 bytes + spare
	b[4] = byte((h.Seq>>16)&0xFF)
	b[5] = byte((h.Seq>>8)&0xFF)
	b[6] = byte(h.Seq & 0xFF)
	b[7] = 0
	return b
}

func (c *Client) send(msgType uint8, body []byte) error {
	conn, err := net.Dial("udp", c.Remote)
	if err != nil { return err }
	defer conn.Close()
	h := Header{Version:1, Sflag:true, MessageType: msgType, Length: uint16(len(body)), Seq: uint32(time.Now().UnixNano())}
	pkt := append(h.MarshalBinary(), body...)
	_, err = conn.Write(pkt)
	return err
}

func (c *Client) SetupSession(body []byte) error { return c.send(50, body) }
func (c *Client) ModifySession(body []byte) error { return c.send(52, body) }
func (c *Client) DeleteSession(body []byte) error { return c.send(54, body) }
