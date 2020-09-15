package go_ethernet_ip

import (
	"bytes"
	"errors"
	"github.com/loki-os/go-ethernet-ip/typedef"
)

type EIPCommand typedef.Uint

const (
	EIPCommandNOP               EIPCommand = 0x00
	EIPCommandListServices      EIPCommand = 0x04
	EIPCommandListIdentity      EIPCommand = 0x63
	EIPCommandListInterfaces    EIPCommand = 0x64
	EIPCommandRegisterSession   EIPCommand = 0x65
	EIPCommandUnRegisterSession EIPCommand = 0x66
	EIPCommandSendRRData        EIPCommand = 0x6F
	EIPCommandSendUnitData      EIPCommand = 0x70
	//EIPCommandIndicateStatus EIPCommand = 0x72
	//EIPCommandCancel EIPCommand = 0x73
)

var EIPCommandMap map[EIPCommand]bool

func init() {
	EIPCommandMap = make(map[EIPCommand]bool)
	EIPCommandMap[EIPCommandNOP] = true
	EIPCommandMap[EIPCommandListServices] = true
	EIPCommandMap[EIPCommandListIdentity] = true
	EIPCommandMap[EIPCommandListInterfaces] = true
	EIPCommandMap[EIPCommandRegisterSession] = true
	EIPCommandMap[EIPCommandUnRegisterSession] = true
	EIPCommandMap[EIPCommandSendRRData] = true
	EIPCommandMap[EIPCommandSendUnitData] = true
}

func checkCommand(cmd EIPCommand) bool {
	_, ok := EIPCommandMap[cmd]
	return ok
}

type EncapsulationHeader struct {
	Command       EIPCommand
	Length        typedef.Uint
	SessionHandle typedef.Udint
	Status        typedef.Udint
	SenderContext typedef.Ulint
	Options       typedef.Udint
}

type EncapsulationPacket struct {
	EncapsulationHeader
	CommandSpecificData []byte
}

func (e *EncapsulationPacket) Encode() ([]byte, error) {
	//check length !> 65511
	if e.Length > 65511 {
		return nil, errors.New("commandSpecificData over length")
	}

	if !checkCommand(e.Command) {
		return nil, errors.New("command not supported")
	}

	buffer := new(bytes.Buffer)
	e.Length = typedef.Uint(len(e.CommandSpecificData))
	WriteByte(buffer, e.EncapsulationHeader)
	WriteByte(buffer, e.CommandSpecificData)
	return buffer.Bytes(), nil
}
