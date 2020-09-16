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
var EIPError map[typedef.Udint]string

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

	EIPError = make(map[typedef.Udint]string)
	EIPError[0x0000] = "success"
	EIPError[0x0001] = "the sender issued an invalid or unsupported encapsulation command."
	EIPError[0x0002] = "insufficient memory resources in the receiver to handle the command. This is not an application error. Instead, it only results if the encapsulation layer cannot obtain memory resources that it needs."
	EIPError[0x0064] = "poorly formed or incorrect data in the data portion of the encapsulation message."
	EIPError[0x0003] = "an originator used an invalid session handle when sending an encapsulation message to the target."
	EIPError[0x0065] = "the target received a message of invalid length."
	EIPError[0x0069] = "unsupported encapsulation protocol revision."
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
