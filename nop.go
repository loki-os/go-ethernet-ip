package go_ethernet_ip

import "github.com/loki-os/go-ethernet-ip/typedef"

func NewNOP(data []byte) *EncapsulationPacket {
	encapsulationPacket := &EncapsulationPacket{}
	encapsulationPacket.Command = EIPCommandNOP
	encapsulationPacket.Length = typedef.Uint(len(data))
	encapsulationPacket.CommandSpecificData = data

	return encapsulationPacket
}
