package listServices

import (
	"github.com/loki-os/go-ethernet-ip/bufferx"
	"github.com/loki-os/go-ethernet-ip/messages/packet"
	"github.com/loki-os/go-ethernet-ip/types"
)

type ListServicesItem struct {
	ItemTypeCode types.UInt
	ItemLength   types.UInt
	Version      types.UInt
	Flags        types.UInt
	Name         []byte
}

type ListServices struct {
	ItemCount types.UInt
	Items     []ListServicesItem
}

func Decode(packet *packet.Packet) (*ListServices, error) {
	result := new(ListServices)
	io := bufferx.New(packet.SpecificData)
	io.RL(&result.ItemCount)

	for i := types.UInt(0); i < result.ItemCount; i++ {
		item := &ListServicesItem{}
		io.RL(&item.ItemTypeCode)
		io.RL(&item.ItemLength)
		io.RL(&item.Version)
		io.RL(&item.Flags)
		item.Name = make([]byte, 16)
		io.RL(&item.Name)
		result.Items = append(result.Items, *item)
	}

	return result, nil
}
