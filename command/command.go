package command

import "github.com/loki-os/go-ethernet-ip/types"

type Command types.UInt

const (
	NOP               Command = 0x00
	ListServices      Command = 0x04
	ListIdentity      Command = 0x63
	ListInterfaces    Command = 0x64
	RegisterSession   Command = 0x65
	UnRegisterSession Command = 0x66
	SendRRData        Command = 0x6F
	SendUnitData      Command = 0x70
	IndicateStatus    Command = 0x72
	Cancel            Command = 0x73
)

var commandMap = map[Command]string{
	NOP:               "NOP",
	ListServices:      "ListServices",
	ListIdentity:      "ListIdentity",
	ListInterfaces:    "ListInterfaces",
	RegisterSession:   "RegisterSession",
	UnRegisterSession: "UnRegisterSession",
	SendRRData:        "SendRRData",
	SendUnitData:      "SendUnitData",
	IndicateStatus:    "IndicateStatus",
	Cancel:            "Cancel",
}

func CheckValid(command Command) bool {
	_, ok := commandMap[command]
	return ok
}
