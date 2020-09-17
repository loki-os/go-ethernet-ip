package go_ethernet_ip

import (
	"bytes"
	"errors"
	"github.com/loki-os/go-ethernet-ip/typedef"
	"time"
)

type sendUnitDataSpecificData struct {
	InterfaceHandle typedef.Udint
	TimeOut         typedef.Uint
	Packet          *commonPacketFormat
}

func (r *sendUnitDataSpecificData) Encode() []byte {
	buffer := new(bytes.Buffer)
	WriteByte(buffer, r.InterfaceHandle)
	WriteByte(buffer, r.TimeOut)
	WriteByte(buffer, r.Packet.Encode())

	return buffer.Bytes()
}

func (r *sendUnitDataSpecificData) Decode(data []byte) {
	dataReader := bytes.NewReader(data)
	ReadByte(dataReader, &r.InterfaceHandle)
	ReadByte(dataReader, &r.TimeOut)
	r.Packet = &commonPacketFormat{}
	r.Packet.Decode(dataReader)
}

func NewSendUnitData(session typedef.Udint, context typedef.Ulint, cpf *commonPacketFormat, timeout typedef.Uint) *EncapsulationPacket {
	encapsulationPacket := &EncapsulationPacket{}
	encapsulationPacket.Command = EIPCommandSendUnitData
	encapsulationPacket.SessionHandle = session
	encapsulationPacket.SenderContext = context

	sd := &sendUnitDataSpecificData{
		InterfaceHandle: 0,
		TimeOut:         timeout,
		Packet:          cpf,
	}
	encapsulationPacket.CommandSpecificData = sd.Encode()
	encapsulationPacket.Length = typedef.Uint(len(encapsulationPacket.CommandSpecificData))

	return encapsulationPacket
}

func (e *EIPTCP) SendUnitData(cpf *commonPacketFormat, timeout typedef.Uint) (*sendUnitDataSpecificData, error) {
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

func (e *EIPTCP) SendUnitDataDecode(encapsulationPacket *EncapsulationPacket) *sendUnitDataSpecificData {
	if len(encapsulationPacket.CommandSpecificData) == 0 {
		return nil
	}

	unitdata := &sendUnitDataSpecificData{}
	unitdata.Decode(encapsulationPacket.CommandSpecificData)

	return unitdata
}
