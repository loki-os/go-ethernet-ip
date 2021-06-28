package packet

import (
	"github.com/loki-os/go-ethernet-ip/bufferx"
	"github.com/loki-os/go-ethernet-ip/types"
)

type ItemID types.UInt

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
	Length types.UInt
	Data   []byte
}

func (i *CommonPacketFormatItem) Encode() []byte {
	if i.Length == 0 {
		i.Length = types.UInt(len(i.Data))
	}
	io := bufferx.New(nil)
	io.WL(i.TypeID)
	io.WL(i.Length)
	io.WL(i.Data)
	return io.Bytes()
}

func (i *CommonPacketFormatItem) Decode(io *bufferx.BufferX) {
	io.RL(&i.TypeID)
	io.RL(&i.Length)
	i.Data = make([]byte, i.Length)
	io.RL(&i.Data)
}

type CommonPacketFormat struct {
	ItemCount types.UInt
	Items     []CommonPacketFormatItem
}

func (c *CommonPacketFormat) Encode() []byte {
	if c.ItemCount == 0 {
		c.ItemCount = types.UInt(len(c.Items))
	}

	io := bufferx.New(nil)
	io.WL(c.ItemCount)

	for _, item := range c.Items {
		io.WL(item.Encode())
	}

	return io.Bytes()
}

func (c *CommonPacketFormat) Decode(io *bufferx.BufferX) {
	io.RL(&c.ItemCount)

	for i := types.UInt(0); i < c.ItemCount; i++ {
		item := &CommonPacketFormatItem{}
		item.Decode(io)
		c.Items = append(c.Items, *item)
	}
}

func NewCommonPacketFormat(items []CommonPacketFormatItem) *CommonPacketFormat {
	return &CommonPacketFormat{
		ItemCount: types.UInt(len(items)),
		Items:     items,
	}
}
