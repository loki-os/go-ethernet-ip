package packet

import (
	"github.com/loki-os/go-ethernet-ip/bufferx"
	"github.com/loki-os/go-ethernet-ip/path"
	"github.com/loki-os/go-ethernet-ip/types"
	"github.com/loki-os/go-ethernet-ip/utils"
)

type UnConnectedSend struct {
	TimeTick           types.USInt
	TimeOutTicks       types.USInt
	MessageRequestSize types.UInt
	MessageRequest     *MessageRouterRequest
	Pad                types.USInt
	RouterPathSize     types.USInt
	Reserved           types.USInt
	RouterPath         []byte
}

func (u *UnConnectedSend) Encode() []byte {
	mr := u.MessageRequest.Encode()

	io := bufferx.New(nil)

	io.WL(u.TimeTick)
	io.WL(u.TimeOutTicks)
	io.WL(types.UInt(len(mr)))
	io.WL(mr)

	if len(mr)%2 == 1 {
		io.WL(uint8(0))
	}

	io.WL(utils.Len(u.RouterPath))
	io.WL(uint8(0))
	io.WL(u.RouterPath)

	return io.Bytes()
}

func UnConnected(slot uint8, timeTick types.USInt, timeOutTicks types.USInt, mr *MessageRouterRequest) *MessageRouterRequest {
	ucs := UnConnectedSend{
		TimeTick:       timeTick,
		TimeOutTicks:   timeOutTicks,
		MessageRequest: mr,
		RouterPath:     path.PortBuild([]byte{slot}, 1, true),
	}

	return NewMessageRouter(0x52, Paths(
		path.LogicalBuild(path.LogicalTypeClassID, 06, true),
		path.LogicalBuild(path.LogicalTypeInstanceID, 01, true),
	), ucs.Encode())
}

func NewUCMM(mr *MessageRouterRequest) *CommonPacketFormat {
	cpf := NewCommonPacketFormat([]CommonPacketFormatItem{
		{
			TypeID: ItemIDUCMM,
			Data:   nil,
		},
		{
			TypeID: ItemIDUnconnectedMessage,
			Data:   mr.Encode(),
		},
	})

	return cpf
}
