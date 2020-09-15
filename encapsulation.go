package go_ethernet_ip

import (
	"bytes"
	"errors"
	"github.com/loki-os/go-ethernet-ip/typedef"
)

type EncapsulationHeader struct {
	Command       typedef.Uint
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
		return nil, errors.New("CommandSpecificData over length")
	}

	buffer := new(bytes.Buffer)
	e.Length = typedef.Uint(len(e.CommandSpecificData))
	WriteByte(buffer, e.EncapsulationHeader)
	WriteByte(buffer, e.CommandSpecificData)
	return buffer.Bytes(), nil
}
