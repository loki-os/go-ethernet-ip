package sendUnitData

import "github.com/loki-os/go-ethernet-ip/messages/packet"

func Decode(_packet *packet.Packet) (*packet.SpecificData, error) {
	result := new(packet.SpecificData)
	result.Decode(_packet.SpecificData)
	return result, nil
}
