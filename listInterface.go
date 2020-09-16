package go_ethernet_ip

import (
	"bytes"
	"github.com/loki-os/go-ethernet-ip/typedef"
)

type ListInterfaceItem struct {
	ItemTypeCode typedef.Uint
	ItemLength   typedef.Uint
	ItemData     []byte
}

type ListInterface struct {
	ItemCount typedef.Uint
	Items     []ListInterfaceItem
}

func (l *ListInterface) Decode(data []byte) {
	dataReader := bytes.NewReader(data)
	ReadByte(dataReader, &l.ItemCount)

	for i := typedef.Uint(0); i < l.ItemCount; i++ {
		item := &ListInterfaceItem{}
		ReadByte(dataReader, &item.ItemTypeCode)
		ReadByte(dataReader, &item.ItemLength)
		item.ItemData = make([]byte, item.ItemLength)
		ReadByte(dataReader, &item.ItemData)
		l.Items = append(l.Items, *item)
	}
}

func NewListInterface(context typedef.Ulint) *EncapsulationPacket {
	encapsulationPacket := &EncapsulationPacket{}
	encapsulationPacket.Command = EIPCommandListInterfaces
	encapsulationPacket.SenderContext = context
	return encapsulationPacket
}

func (e *EIPTCP) ListInterface(cb func(interface{}, error)) {
	ctx := CtxGenerator()
	e.router[ctx] = cb

	encapsulationPacket := NewListInterface(ctx)
	b, _ := encapsulationPacket.Encode()

	if e.tcpConn != nil {
		e.sender <- b
	}
}

func (e *EIPUDP) ListInterface() {
	encapsulationPacket := NewListInterface(0)
	b, _ := encapsulationPacket.Encode()

	if e.udpConn != nil {
		e.send(b)
	}
}

func (e *EIPTCP) ListInterfaceDecode(encapsulationPacket *EncapsulationPacket) *ListInterface {
	if len(encapsulationPacket.CommandSpecificData) == 0 {
		return nil
	}

	list := &ListInterface{}
	list.Decode(encapsulationPacket.CommandSpecificData)

	return list
}

func (e *EIPUDP) ListInterfaceDecode(encapsulationPacket *EncapsulationPacket) *ListInterface {
	if len(encapsulationPacket.CommandSpecificData) == 0 {
		return nil
	}

	list := &ListInterface{}
	list.Decode(encapsulationPacket.CommandSpecificData)

	return list
}
