package go_ethernet_ip

import (
	"errors"
	"github.com/loki-os/go-ethernet-ip/typedef"
)

func NewUnRegisterSession(session typedef.Udint, context typedef.Ulint) *EncapsulationPacket {
	encapsulationPacket := &EncapsulationPacket{}
	encapsulationPacket.Command = EIPCommandUnRegisterSession
	encapsulationPacket.SenderContext = context
	encapsulationPacket.SessionHandle = session

	return encapsulationPacket
}

func (e *EIPTCP) UnRegisterSession() {
	ctx := CtxGenerator()
	encapsulationPacket := NewUnRegisterSession(e.session, ctx)
	b, _ := encapsulationPacket.Encode()
	e.sender <- b
}

func (e *EIPTCP) UnRegisterSessionDecode(encapsulationPacket *EncapsulationPacket) {
	if e.session == encapsulationPacket.SessionHandle {
		e.disconnect(errors.New("unRegister session"))
	}
}
