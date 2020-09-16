package go_ethernet_ip

import (
	"bytes"
	"github.com/loki-os/go-ethernet-ip/typedef"
)

type ListServicesItem struct {
	ItemTypeCode typedef.Uint
	ItemLength   typedef.Uint
	Version      typedef.Uint
	Flags        typedef.Uint
	Name         []byte
}

type ListServices struct {
	ItemCount typedef.Uint
	Items     []ListServicesItem
}

func (l *ListServices) Decode(data []byte) {
	dataReader := bytes.NewReader(data)
	ReadByte(dataReader, &l.ItemCount)

	for i := typedef.Uint(0); i < l.ItemCount; i++ {
		item := &ListServicesItem{}
		ReadByte(dataReader, &item.ItemTypeCode)
		ReadByte(dataReader, &item.ItemLength)
		ReadByte(dataReader, &item.Version)
		ReadByte(dataReader, &item.Flags)
		item.Name = make([]byte, 16)
		ReadByte(dataReader, &item.Name)
		l.Items = append(l.Items, *item)
	}
}

func NewListServices(context typedef.Ulint) *EncapsulationPacket {
	encapsulationPacket := &EncapsulationPacket{}
	encapsulationPacket.Command = EIPCommandListServices
	encapsulationPacket.SenderContext = context
	return encapsulationPacket
}

func (e *EIPTCP) ListServices(cb func(interface{}, error)) {
	ctx := CtxGenerator()
	e.router[ctx] = cb

	encapsulationPacket := NewListServices(ctx)
	b, _ := encapsulationPacket.Encode()

	if e.tcpConn != nil {
		e.sender <- b
	}
}

func (e *EIPUDP) ListServices() {
	encapsulationPacket := NewListServices(0)
	b, _ := encapsulationPacket.Encode()

	if e.udpConn != nil {
		e.send(b)
	}
}

func (e *EIPTCP) ListServicesDecode(encapsulationPacket *EncapsulationPacket) *ListServices {
	if len(encapsulationPacket.CommandSpecificData) == 0 {
		return nil
	}

	list := &ListServices{}
	list.Decode(encapsulationPacket.CommandSpecificData)

	return list
}

func (e *EIPUDP) ListServicesDecode(encapsulationPacket *EncapsulationPacket) *ListServices {
	if len(encapsulationPacket.CommandSpecificData) == 0 {
		return nil
	}

	list := &ListServices{}
	list.Decode(encapsulationPacket.CommandSpecificData)

	return list
}
