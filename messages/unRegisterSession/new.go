package unRegisterSession

import (
	"github.com/loki-os/go-ethernet-ip/command"
	"github.com/loki-os/go-ethernet-ip/messages/packet"
	"github.com/loki-os/go-ethernet-ip/types"
)

func New(session types.UDInt, context types.ULInt) (*packet.Packet, error) {
	return &packet.Packet{
		Header: packet.Header{
			Command:       command.UnRegisterSession,
			Length:        0,
			SessionHandle: session,
			Status:        0,
			SenderContext: context,
			Options:       0,
		},
		SpecificData: nil,
	}, nil
}
