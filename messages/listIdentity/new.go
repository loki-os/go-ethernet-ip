package listIdentity

import (
	"github.com/loki-os/go-ethernet-ip/command"
	"github.com/loki-os/go-ethernet-ip/messages/packet"
	"github.com/loki-os/go-ethernet-ip/types"
)

func New(context types.ULInt) (*packet.Packet, error) {
	return &packet.Packet{
		Header: packet.Header{
			Command:       command.ListIdentity,
			Length:        0,
			SessionHandle: 0,
			Status:        0,
			SenderContext: context,
			Options:       0,
		},
		SpecificData: nil,
	}, nil
}
