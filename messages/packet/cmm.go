package packet

import (
	"github.com/loki-os/go-ethernet-ip/bufferx"
	"github.com/loki-os/go-ethernet-ip/types"
)

func NewCMM(connectionID types.UDInt, sequenceNumber types.UInt, mr *MessageRouterRequest) *CommonPacketFormat {
	io := bufferx.New(nil)
	io.WL(connectionID)

	io1 := bufferx.New(nil)
	io1.WL(sequenceNumber)
	io1.WL(mr.Encode())

	cpf := NewCommonPacketFormat([]CommonPacketFormatItem{
		{
			TypeID: ItemIDConnectionBased,
			Data:   io.Bytes(),
		},
		{
			TypeID: ItemIDConnectedTransportPacket,
			Data:   io1.Bytes(),
		},
	})

	return cpf
}
