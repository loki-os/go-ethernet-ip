package packet

import (
	"errors"
	"github.com/loki-os/go-ethernet-ip/bufferx"
	"github.com/loki-os/go-ethernet-ip/command"
	"github.com/loki-os/go-ethernet-ip/types"
)

type Header struct {
	Command       command.Command
	Length        types.UInt
	SessionHandle types.UDInt
	Status        types.UDInt
	SenderContext types.ULInt
	Options       types.UDInt
}

type Packet struct {
	Header
	SpecificData []byte
}

func (p *Packet) Encode() ([]byte, error) {
	if p.Length > 65511 {
		return nil, errors.New("specific data over length 65511")
	}

	if !command.CheckValid(p.Command) {
		return nil, errors.New("command not supported")
	}

	buffer := bufferx.New(nil)
	buffer.WL(p.Header)
	buffer.WL(p.SpecificData)
	if buffer.Error() != nil {
		return nil, buffer.Error()
	}

	return buffer.Bytes(), nil
}
