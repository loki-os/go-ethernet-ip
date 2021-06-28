package sendRRData

import (
	"github.com/loki-os/go-ethernet-ip/command"
	"github.com/loki-os/go-ethernet-ip/messages/packet"
	"github.com/loki-os/go-ethernet-ip/types"
)

func New(session types.UDInt, context types.ULInt, cpf *packet.CommonPacketFormat, timeout types.UInt) (*packet.Packet, error) {
	specificData := &packet.SpecificData{
		InterfaceHandle: 0,
		TimeOut:         timeout,
		Packet:          cpf,
	}
	specificDataBytes := specificData.Encode()
	return &packet.Packet{
		Header: packet.Header{
			Command:       command.SendRRData,
			Length:        types.UInt(len(specificDataBytes)),
			SessionHandle: session,
			Status:        0,
			SenderContext: context,
			Options:       0,
		},
		SpecificData: specificDataBytes,
	}, nil
}
