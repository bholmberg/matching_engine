package coordinator

import (
	"fmt"
	"github.com/fmstephe/matching_engine/msg"
	"github.com/fmstephe/matching_engine/msg/msgutil"
	"io"
	"net"
	"os"
	"time"
)

const RESEND_MILLIS = time.Duration(10) * time.Millisecond

type stdResponder struct {
	unacked *msgutil.Set
	writer  io.WriteCloser
	msgHelper
}

func newResponder(writer io.WriteCloser) *stdResponder {
	return &stdResponder{unacked: msgutil.NewSet(), writer: writer}
}

func (r *stdResponder) Run() {
	defer r.shutdown()
	t := time.NewTimer(RESEND_MILLIS)
	for {
		select {
		case resp := <-r.msgs:
			switch {
			case resp.Direction == msg.IN && resp.Route == msg.ACK:
				r.handleInAck(resp)
			case resp.Direction == msg.OUT && (resp.Status != msg.NORMAL || resp.Route == msg.APP || resp.Route == msg.ACK):
				r.writeResponse(resp)
			case resp.Direction == msg.IN && resp.Route == msg.APP && resp.Status == msg.NORMAL:
				continue
			case resp.Route == msg.SHUTDOWN:
				return
			default:
				panic(fmt.Sprintf("Unhandleable response %v", resp))
			}
		case <-t.C:
			r.resend()
			t = time.NewTimer(RESEND_MILLIS)
		}
	}
}

func (r *stdResponder) handleInAck(ca *msg.Message) {
	r.unacked.Remove(ca)
}

func (r *stdResponder) writeResponse(resp *msg.Message) {
	resp.Direction = msg.IN
	r.addToUnacked(resp)
	r.write(resp)
}

func (r *stdResponder) addToUnacked(resp *msg.Message) {
	if resp.Route == msg.APP {
		r.unacked.Add(resp)
	}
}

func (r *stdResponder) resend() {
	r.unacked.Do(func(m *msg.Message) {
		r.write(m)
	})
}

func (r *stdResponder) write(resp *msg.Message) {
	b := make([]byte, msg.SizeofMessage)
	resp.WriteTo(b)
	n, err := r.writer.Write(b)
	if err != nil {
		r.handleError(resp, err, msg.WRITE_ERROR)
	}
	if n != msg.SizeofMessage {
		r.handleError(resp, err, msg.SMALL_WRITE_ERROR)
	}
}

func (r *stdResponder) handleError(resp *msg.Message, err error, s msg.MsgStatus) {
	em := &msg.Message{}
	*em = *resp
	em.Status = s
	println(resp.String(), err.Error())
	if e, ok := err.(net.Error); ok && !e.Temporary() {
		os.Exit(1)
	}
}

func (r *stdResponder) shutdown() {
	r.writer.Close()
}