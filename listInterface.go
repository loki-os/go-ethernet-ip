package go_ethernet_ip

import (
	"bytes"
	"errors"
	"github.com/loki-os/go-ethernet-ip/typedef"
	"time"
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

func (e *EIPTCP) ListInterface() (*ListInterface, error) {
	ctx := CtxGenerator()
	e.receiver[ctx] = make(chan *EncapsulationPacket)

	encapsulationPacket := NewListInterface(ctx)
	b, _ := encapsulationPacket.Encode()

	if e.tcpConn != nil {
		e.sender <- b
	}

	for {
		select {
		case <-time.After(e.config.TCPTimeout):
			return nil, errors.New("tcp timeout")
		case received := <-e.receiver[ctx]:
			return e.ListInterfaceDecode(received), nil
		}
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
