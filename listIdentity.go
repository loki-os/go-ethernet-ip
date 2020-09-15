package go_ethernet_ip

import (
	"bytes"
	"encoding/json"
	"github.com/loki-os/go-ethernet-ip/typedef"
	"log"
)

const CommandListIdentity = typedef.Uint(0x63)

type ListIdentityItem struct {
	ItemTypeCode                 typedef.Uint
	ItemLength                   typedef.Uint
	EncapsulationProtocolVersion typedef.Uint
	SinFamily                    typedef.Int
	SinPort                      typedef.Uint
	SinAddr                      typedef.Udint
	SinZero                      typedef.Ulint
	VendorID                     typedef.Uint
	DeviceType                   typedef.Uint
	ProductCode                  typedef.Uint
	Major                        typedef.Usint
	Minor                        typedef.Usint
	Status                       typedef.Word
	SerialNumber                 typedef.Udint
	NameLength                   typedef.Usint
	ProductName                  []byte
	State                        typedef.Usint
}

type ListIdentity struct {
	ItemCount typedef.Uint
	Items     []ListIdentityItem
}

func (l *ListIdentity) Decode(data []byte) {
	dataReader := bytes.NewReader(data)
	ReadByte(dataReader, &l.ItemCount)

	for i := typedef.Uint(0); i < l.ItemCount; i++ {
		item := &ListIdentityItem{}
		ReadByte(dataReader, &item.ItemTypeCode)
		ReadByte(dataReader, &item.ItemLength)
		ReadByte(dataReader, &item.EncapsulationProtocolVersion)
		ReadByteBigEndian(dataReader, &item.SinFamily)
		ReadByteBigEndian(dataReader, &item.SinPort)
		ReadByteBigEndian(dataReader, &item.SinAddr)
		ReadByteBigEndian(dataReader, &item.SinZero)
		ReadByte(dataReader, &item.VendorID)
		ReadByte(dataReader, &item.DeviceType)
		ReadByte(dataReader, &item.ProductCode)
		ReadByte(dataReader, &item.Major)
		ReadByte(dataReader, &item.Minor)
		ReadByte(dataReader, &item.Status)
		ReadByte(dataReader, &item.SerialNumber)
		ReadByte(dataReader, &item.NameLength)
		item.ProductName = make([]byte, item.NameLength)
		ReadByte(dataReader, &item.ProductName)
		ReadByte(dataReader, &item.State)

		l.Items = append(l.Items, *item)
	}
}

func NewListIdentity() *EncapsulationPacket {
	encapsulationPacket := &EncapsulationPacket{}
	encapsulationPacket.Command = CommandListIdentity
	return encapsulationPacket
}

func (e *EIP) ListIdentity() {
	encapsulationPacket := NewListIdentity()
	b, _ := encapsulationPacket.Encode()
	if e.tcpConn != nil {
		e.sender <- b
	}
	if e.udpConn != nil {
		_ = e.udpSend(b)
	}
}

func (e *EIP) ListIdentityDecode(encapsulationPacket *EncapsulationPacket) {
	if len(encapsulationPacket.CommandSpecificData) == 0 {
		return
	}

	list := &ListIdentity{}
	list.Decode(encapsulationPacket.CommandSpecificData)

	b, _ := json.MarshalIndent(list, "", "\t")
	log.Println(string(b))
}
