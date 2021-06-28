package registerSession

import (
	"github.com/loki-os/go-ethernet-ip/bufferx"
	"github.com/loki-os/go-ethernet-ip/command"
	"github.com/loki-os/go-ethernet-ip/messages/packet"
	"github.com/loki-os/go-ethernet-ip/types"
)

type specificData struct {
	ProtocolVersion types.UInt
	OptionsFlags    types.UInt
}

func (s *specificData) Encode() ([]byte, error) {
	buffer := bufferx.New(nil)
	buffer.WL(s.ProtocolVersion)
	buffer.WL(s.OptionsFlags)
	if buffer.Error() != nil {
		return nil, buffer.Error()
	}
	return buffer.Bytes(), nil
}

func New(context types.ULInt) (*packet.Packet, error) {
	specificData := specificData{
		ProtocolVersion: 1,
		OptionsFlags:    0,
	}
	specificDataBytes, err := specificData.Encode()
	if err != nil {
		return nil, err
	}

	return &packet.Packet{
		Header: packet.Header{
			Command:       command.RegisterSession,
			Length:        4,
			SessionHandle: 0,
			Status:        0,
			SenderContext: context,
			Options:       0,
		},
		SpecificData: specificDataBytes,
	}, nil
}
