package go_ethernet_ip

import (
	"errors"
	"github.com/loki-os/go-ethernet-ip/typedef"
	"time"
)

func NewSendUnitData(session typedef.Udint, context typedef.Ulint, cpf *CommonPacketFormat, timeout typedef.Uint) *EncapsulationPacket {
	encapsulationPacket := &EncapsulationPacket{}
	encapsulationPacket.Command = EIPCommandSendUnitData
	encapsulationPacket.SessionHandle = session
	encapsulationPacket.SenderContext = context

	sd := &SendDataSpecificData{
		InterfaceHandle: 0,
		TimeOut:         timeout,
		Packet:          cpf,
	}
	encapsulationPacket.CommandSpecificData = sd.Encode()
	encapsulationPacket.Length = typedef.Uint(len(encapsulationPacket.CommandSpecificData))

	return encapsulationPacket
}

func (e *EIPTCP) SendUnitData(cpf *CommonPacketFormat, timeout typedef.Uint) (*SendDataSpecificData, error) {
	ctx := CtxGenerator()
	e.receiver[ctx] = make(chan *EncapsulationPacket)

	encapsulationPacket := NewSendUnitData(e.session, ctx, cpf, timeout)
	b, _ := encapsulationPacket.Encode()
	e.sender <- b

	for {
		select {
		case <-time.After(e.config.TCPTimeout):
			return nil, errors.New("tcp timeout")
		case received := <-e.receiver[ctx]:
			return e.SendUnitDataDecode(received), nil
		}
	}
}

func (e *EIPTCP) SendUnitDataDecode(encapsulationPacket *EncapsulationPacket) *SendDataSpecificData {
	if len(encapsulationPacket.CommandSpecificData) == 0 {
		return nil
	}

	unitdata := &SendDataSpecificData{}
	unitdata.Decode(encapsulationPacket.CommandSpecificData)

	return unitdata
}
