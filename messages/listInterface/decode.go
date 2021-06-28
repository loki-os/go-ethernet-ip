package listInterface

import (
	"github.com/loki-os/go-ethernet-ip/bufferx"
	"github.com/loki-os/go-ethernet-ip/messages/packet"
	"github.com/loki-os/go-ethernet-ip/types"
)

type ListInterfaceItem struct {
	ItemTypeCode types.UInt
	ItemLength   types.UInt
	ItemData     []byte
}

type ListInterface struct {
	ItemCount types.UInt
	Items     []ListInterfaceItem
}

func Decode(packet *packet.Packet) (*ListInterface, error) {
	result := new(ListInterface)
	io := bufferx.New(packet.SpecificData)
	io.RL(&result.ItemCount)

	for i := types.UInt(0); i < result.ItemCount; i++ {
		item := &ListInterfaceItem{}
		io.RL(&item.ItemTypeCode)
		io.RL(&item.ItemLength)
		item.ItemData = make([]byte, item.ItemLength)
		io.RL(&item.ItemData)
		result.Items = append(result.Items, *item)
	}

	return result, nil
}
