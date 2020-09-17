package go_ethernet_ip

import (
	"bytes"
	"github.com/loki-os/go-ethernet-ip/typedef"
)

type commonPacketFormatItem struct {
	TypeID typedef.Uint
	Length typedef.Uint
	Data   []byte
}

func (i *commonPacketFormatItem) Encode() []byte {
	buffer := new(bytes.Buffer)
	WriteByte(buffer, i.TypeID)
	WriteByte(buffer, i.Length)
	WriteByte(buffer, i.Data)

	return buffer.Bytes()
}

func (i *commonPacketFormatItem) Decode(dataReader *bytes.Reader) {
	ReadByte(dataReader, &i.TypeID)
	ReadByte(dataReader, &i.Length)
	i.Data = make([]byte, i.Length)
	ReadByte(dataReader, &i.Data)
}

type commonPacketFormat struct {
	ItemCount typedef.Uint
	Items     []commonPacketFormatItem
}

func (c *commonPacketFormat) Encode() []byte {
	buffer := new(bytes.Buffer)
	WriteByte(buffer, c.ItemCount)
	for _, item := range c.Items {
		WriteByte(buffer, item.Encode())
	}

	return buffer.Bytes()
}

func (c *commonPacketFormat) Decode(dataReader *bytes.Reader) {
	ReadByte(dataReader, &c.ItemCount)

	for i := typedef.Uint(0); i < c.ItemCount; i++ {
		item := &commonPacketFormatItem{}
		item.Decode(dataReader)
		c.Items = append(c.Items, *item)
	}
}
