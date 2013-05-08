package netwk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/fmstephe/matching_engine/trade"
	"net"
)

type Listener struct {
	conn   *net.UDPConn
	submit chan interface{}
}

func NewListener(port string, submit chan interface{}) (*Listener, error) {
	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return &Listener{conn: conn, submit: submit}, nil
}

func (l *Listener) Listen() {
	for {
		s := make([]byte, trade.SizeofOrderData)
		n, _, err := l.conn.ReadFromUDP(s)
		if err != nil {
			println("Listener - UDP Read: ", err.Error())
			continue
		}
		if n != trade.SizeofOrderData {
			println(fmt.Sprintf("Listener: Error incorrect number of bytes. Expecting %d, found %d submit %v", trade.SizeofOrderData, n, s))
			continue
		}
		od := &trade.OrderData{}
		buf := bytes.NewBuffer(s)
		err = binary.Read(buf, binary.BigEndian, od)
		if err != nil {
			println("Listener - to []byte: ", err.Error())
			continue
		}
		l.submit <- od
	}
}
