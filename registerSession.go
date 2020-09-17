package go_ethernet_ip

import (
	"bytes"
	"errors"
	"github.com/loki-os/go-ethernet-ip/typedef"
	"time"
)

type registerSessionSpecificData struct {
	ProtocolVersion typedef.Uint
	OptionsFlags    typedef.Uint
}

func (r *registerSessionSpecificData) Encode() []byte {
	buffer := new(bytes.Buffer)
	WriteByte(buffer, r.ProtocolVersion)
	WriteByte(buffer, r.OptionsFlags)
	return buffer.Bytes()
}

func NewRegisterSession(context typedef.Ulint) *EncapsulationPacket {
	encapsulationPacket := &EncapsulationPacket{}
	encapsulationPacket.Command = EIPCommandRegisterSession
	encapsulationPacket.Length = 4
	encapsulationPacket.SenderContext = context

	sd := &registerSessionSpecificData{
		ProtocolVersion: 1,
		OptionsFlags:    0,
	}
	encapsulationPacket.CommandSpecificData = sd.Encode()

	return encapsulationPacket
}

func (e *EIPTCP) RegisterSession() error {
	ctx := CtxGenerator()
	e.receiver[ctx] = make(chan *EncapsulationPacket)

	encapsulationPacket := NewRegisterSession(ctx)
	b, _ := encapsulationPacket.Encode()
	e.sender <- b

	for {
		select {
		case <-time.After(e.config.TCPTimeout):
			return errors.New("tcp timeout")
		case received := <-e.receiver[ctx]:
			e.RegisterSessionDecode(received)
			return nil
		}
	}
}

func (e *EIPTCP) RegisterSessionDecode(encapsulationPacket *EncapsulationPacket) {
	e.session = encapsulationPacket.SessionHandle
	if e.Connected != nil {
		e.Connected()
	}
}
