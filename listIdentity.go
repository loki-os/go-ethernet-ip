package go_ethernet_ip

import "github.com/loki-os/go-ethernet-ip/typedef"

type ListIdentity struct {
	ItemCount typedef.Uint
}

func NewListIdentity() *EncapsulationPacket {
	encapsulationPacket := &EncapsulationPacket{}
	encapsulationPacket.Command = 0x63

	return encapsulationPacket
}

//func DecodeListidentity(data []byte) *ListIdentity {
//
//}
