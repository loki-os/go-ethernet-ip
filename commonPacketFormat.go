package go_ethernet_ip

import (
	"bytes"
	"github.com/loki-os/go-ethernet-ip/typedef"
)

type ItemID typedef.Uint

const (
	ItemIDUCMM                     ItemID = 0x0000
	ItemIDListIdentityResponse     ItemID = 0x000C
	ItemIDConnectionBased          ItemID = 0x00A1
	ItemIDConnectedTransportPacket ItemID = 0x00B1
	ItemIDUnconnectedMessage       ItemID = 0x00B2
	ItemIDListServicesResponse     ItemID = 0x0100
	ItemIDSockaddrInfoO2T          ItemID = 0x8000
	ItemIDSockaddrInfoT2O          ItemID = 0x8001
	ItemIDSequencedAddressItem     ItemID = 0x8002
)

type CommonPacketFormatItem struct {
	TypeID ItemID
	Length typedef.Uint
	Data   []byte
}

func (i *CommonPacketFormatItem) Encode() []byte {
	if i.Length == 0 {
		i.Length = typedef.Uint(len(i.Data))
	}

	buffer := new(bytes.Buffer)
	WriteByte(buffer, i.TypeID)
	WriteByte(buffer, i.Length)
	WriteByte(buffer, i.Data)

	return buffer.Bytes()
}

func (i *CommonPacketFormatItem) New(id ItemID, data []byte) {
	i.TypeID = id
	i.Data = data
	i.Length = typedef.Uint(len(data))
}

func (i *CommonPacketFormatItem) Decode(dataReader *bytes.Reader) {
	ReadByte(dataReader, &i.TypeID)
	ReadByte(dataReader, &i.Length)
	i.Data = make([]byte, i.Length)
	ReadByte(dataReader, &i.Data)
}

type CommonPacketFormat struct {
	ItemCount typedef.Uint
	Items     []CommonPacketFormatItem
}

func (c *CommonPacketFormat) Encode() []byte {
	if c.ItemCount == 0 {
		c.ItemCount = typedef.Uint(len(c.Items))
	}

	buffer := new(bytes.Buffer)

	WriteByte(buffer, c.ItemCount)

	for _, item := range c.Items {
		WriteByte(buffer, item.Encode())
	}

	return buffer.Bytes()
}

func (c *CommonPacketFormat) New(items []CommonPacketFormatItem) {
	c.ItemCount = typedef.Uint(len(c.Items))
	c.Items = items
}

func (c *CommonPacketFormat) Decode(dataReader *bytes.Reader) {
	ReadByte(dataReader, &c.ItemCount)

	for i := typedef.Uint(0); i < c.ItemCount; i++ {
		item := &CommonPacketFormatItem{}
		item.Decode(dataReader)
		c.Items = append(c.Items, *item)
	}
}
