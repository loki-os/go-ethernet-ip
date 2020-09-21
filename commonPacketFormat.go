package go_ethernet_ip

import (
	"bytes"
	"github.com/loki-os/go-ethernet-ip/typedef"
)

type CommonPacketFormatItem struct {
	TypeID typedef.Uint
	Length typedef.Uint
	Data   []byte
}

func (i *CommonPacketFormatItem) Encode() []byte {
	buffer := new(bytes.Buffer)
	WriteByte(buffer, i.TypeID)
	WriteByte(buffer, i.Length)
	WriteByte(buffer, i.Data)

	return buffer.Bytes()
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
