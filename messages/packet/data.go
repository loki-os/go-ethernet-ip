package packet

import (
	"github.com/loki-os/go-ethernet-ip/bufferx"
	"github.com/loki-os/go-ethernet-ip/types"
)

type SpecificData struct {
	InterfaceHandle types.UDInt
	TimeOut         types.UInt
	Packet          *CommonPacketFormat
}

func (r *SpecificData) Encode() []byte {
	io := bufferx.New(nil)
	io.WL(r.InterfaceHandle)
	io.WL(r.TimeOut)
	io.WL(r.Packet.Encode())
	return io.Bytes()
}

func (r *SpecificData) Decode(data []byte) {
	io := bufferx.New(data)
	io.RL(&r.InterfaceHandle)
	io.RL(&r.TimeOut)
	r.Packet = new(CommonPacketFormat)
	r.Packet.Decode(io)
}
