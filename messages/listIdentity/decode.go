package listIdentity

import (
	"github.com/loki-os/go-ethernet-ip/bufferx"
	"github.com/loki-os/go-ethernet-ip/messages/packet"
	"github.com/loki-os/go-ethernet-ip/types"
)

type ListIdentityItem struct {
	ItemTypeCode                 types.UInt
	ItemLength                   types.UInt
	EncapsulationProtocolVersion types.UInt
	SinFamily                    types.Int
	SinPort                      types.UInt
	SinAddr                      types.UDInt
	SinZero                      types.ULInt
	VendorID                     types.UInt
	DeviceType                   types.UInt
	ProductCode                  types.UInt
	Major                        types.USInt
	Minor                        types.USInt
	Status                       types.Word
	SerialNumber                 types.UDInt
	NameLength                   types.USInt
	ProductName                  []byte
	State                        types.USInt
}

type ListIdentity struct {
	ItemCount types.UInt
	Items     []ListIdentityItem
}

func Decode(packet *packet.Packet) (*ListIdentity, error) {
	result := new(ListIdentity)
	io := bufferx.New(packet.SpecificData)
	io.RL(&result.ItemCount)

	for i := types.UInt(0); i < result.ItemCount; i++ {
		item := &ListIdentityItem{}
		io.RL(&item.ItemTypeCode)
		io.RL(&item.ItemLength)
		io.RL(&item.EncapsulationProtocolVersion)
		io.RB(&item.SinFamily)
		io.RB(&item.SinPort)
		io.RB(&item.SinAddr)
		io.RB(&item.SinZero)
		io.RL(&item.VendorID)
		io.RL(&item.DeviceType)
		io.RL(&item.ProductCode)
		io.RL(&item.Major)
		io.RL(&item.Minor)
		io.RL(&item.Status)
		io.RL(&item.SerialNumber)
		io.RL(&item.NameLength)
		item.ProductName = make([]byte, item.NameLength)
		io.RL(&item.ProductName)
		io.RL(&item.State)
		result.Items = append(result.Items, *item)
	}

	return result, nil
}
