package go_ethernet_ip

import (
	"bytes"
	"github.com/loki-os/go-ethernet-ip/typedef"
)

type sendRRDataSpecificData struct {
	InterfaceHandle typedef.Udint
	TimeOut         typedef.Uint
	Packet          *commonPacketFormat
}

func (r *sendRRDataSpecificData) Encode() []byte {
	buffer := new(bytes.Buffer)
	WriteByte(buffer, r.InterfaceHandle)
	WriteByte(buffer, r.TimeOut)
	WriteByte(buffer, r.Packet.Encode())

	return buffer.Bytes()
}

func (r *sendRRDataSpecificData) Decode(data []byte) {
	dataReader := bytes.NewReader(data)
	ReadByte(dataReader, &r.InterfaceHandle)
	ReadByte(dataReader, &r.TimeOut)
	r.Packet = &commonPacketFormat{}
	r.Packet.Decode(dataReader)
}

func NewSendRRData(session typedef.Udint, context typedef.Ulint, cpf *commonPacketFormat) *EncapsulationPacket {
	encapsulationPacket := &EncapsulationPacket{}
	encapsulationPacket.Command = EIPCommandSendRRData
	encapsulationPacket.SessionHandle = session
	encapsulationPacket.SenderContext = context

	sd := &sendRRDataSpecificData{
		InterfaceHandle: 0,
		TimeOut:         0,
		Packet:          cpf,
	}
	encapsulationPacket.CommandSpecificData = sd.Encode()
	encapsulationPacket.Length = typedef.Uint(len(encapsulationPacket.CommandSpecificData))

	return encapsulationPacket
}

func (e *EIPTCP) SendRRData(cpf *commonPacketFormat, cb func(interface{}, error)) {
	ctx := CtxGenerator()
	e.router[ctx] = cb

	encapsulationPacket := NewSendRRData(e.session, ctx, cpf)
	b, _ := encapsulationPacket.Encode()
	e.sender <- b
}

func (e *EIPTCP) SendRRDataDecode(encapsulationPacket *EncapsulationPacket) *sendRRDataSpecificData {
	if len(encapsulationPacket.CommandSpecificData) == 0 {
		return nil
	}

	rrdata := &sendRRDataSpecificData{}
	rrdata.Decode(encapsulationPacket.CommandSpecificData)

	return rrdata
}
