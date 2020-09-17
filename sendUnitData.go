package go_ethernet_ip

import (
	"bytes"
	"github.com/loki-os/go-ethernet-ip/typedef"
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

func NewSendUnitData(session typedef.Udint, context typedef.Ulint, cpf *commonPacketFormat) *EncapsulationPacket {
	encapsulationPacket := &EncapsulationPacket{}
	encapsulationPacket.Command = EIPCommandSendUnitData
	encapsulationPacket.SessionHandle = session
	encapsulationPacket.SenderContext = context

	sd := &sendUnitDataSpecificData{
		InterfaceHandle: 0,
		TimeOut:         0,
		Packet:          cpf,
	}
	encapsulationPacket.CommandSpecificData = sd.Encode()
	encapsulationPacket.Length = typedef.Uint(len(encapsulationPacket.CommandSpecificData))

	return encapsulationPacket
}

func (e *EIPTCP) SendUnitData(cpf *commonPacketFormat, cb func(interface{}, error)) {
	ctx := CtxGenerator()
	e.router[ctx] = cb

	encapsulationPacket := NewSendUnitData(e.session, ctx, cpf)
	b, _ := encapsulationPacket.Encode()
	e.sender <- b
}

func (e *EIPTCP) SendUnitDataDecode(encapsulationPacket *EncapsulationPacket) *sendUnitDataSpecificData {
	if len(encapsulationPacket.CommandSpecificData) == 0 {
		return nil
	}

	unitdata := &sendUnitDataSpecificData{}
	unitdata.Decode(encapsulationPacket.CommandSpecificData)

	return unitdata
}
