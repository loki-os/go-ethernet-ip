package nop

import (
	"github.com/loki-os/go-ethernet-ip/command"
	"github.com/loki-os/go-ethernet-ip/messages/packet"
	"github.com/loki-os/go-ethernet-ip/types"
)

func New(data []byte) (*packet.Packet, error) {
	return &packet.Packet{
		Header: packet.Header{
			Command:       command.NOP,
			Length:        types.UInt(len(data)),
			SessionHandle: 0,
			Status:        0,
			SenderContext: 0,
			Options:       0,
		},
		SpecificData: data,
	}, nil
}
