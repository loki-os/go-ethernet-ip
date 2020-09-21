package go_ethernet_ip

import (
	"bytes"
	"errors"
	"github.com/loki-os/go-ethernet-ip/typedef"
	"time"
)

type sendRRDataSpecificData struct {
	InterfaceHandle typedef.Udint
	TimeOut         typedef.Uint
	Packet          *CommonPacketFormat
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
	r.Packet = &CommonPacketFormat{}
	r.Packet.Decode(dataReader)
}

func NewSendRRData(session typedef.Udint, context typedef.Ulint, cpf *CommonPacketFormat, timeout typedef.Uint) *EncapsulationPacket {
	encapsulationPacket := &EncapsulationPacket{}
	encapsulationPacket.Command = EIPCommandSendRRData
	encapsulationPacket.SessionHandle = session
	encapsulationPacket.SenderContext = context

	sd := &sendRRDataSpecificData{
		InterfaceHandle: 0,
		TimeOut:         timeout,
		Packet:          cpf,
	}
	encapsulationPacket.CommandSpecificData = sd.Encode()
	encapsulationPacket.Length = typedef.Uint(len(encapsulationPacket.CommandSpecificData))

	return encapsulationPacket
}

func (e *EIPTCP) SendRRData(cpf *CommonPacketFormat, timeout typedef.Uint) (*sendRRDataSpecificData, error) {
	ctx := CtxGenerator()
	e.receiver[ctx] = make(chan *EncapsulationPacket)

	encapsulationPacket := NewSendRRData(e.session, ctx, cpf, timeout)
	b, _ := encapsulationPacket.Encode()
	e.sender <- b

	for {
		select {
		case <-time.After(e.config.TCPTimeout):
			return nil, errors.New("tcp timeout")
		case received := <-e.receiver[ctx]:
			return e.SendRRDataDecode(received), nil
		}
	}
}

func (e *EIPTCP) SendRRDataDecode(encapsulationPacket *EncapsulationPacket) *sendRRDataSpecificData {
	if len(encapsulationPacket.CommandSpecificData) == 0 {
		return nil
	}

	rrdata := &sendRRDataSpecificData{}
	rrdata.Decode(encapsulationPacket.CommandSpecificData)

	return rrdata
}
