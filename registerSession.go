package go_ethernet_ip

import (
	"bytes"
	"github.com/loki-os/go-ethernet-ip/typedef"
	"log"
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

func (e *EIPTCP) RegisterSession() {
	encapsulationPacket := NewRegisterSession(0)
	b, _ := encapsulationPacket.Encode()
	e.sender <- b
}

func (e *EIPTCP) RegisterSessionDecode(encapsulationPacket *EncapsulationPacket) {
	log.Printf("%+v\n", encapsulationPacket)
}
