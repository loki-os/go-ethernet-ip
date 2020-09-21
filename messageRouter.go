package go_ethernet_ip

import (
	"bytes"
	"github.com/loki-os/go-ethernet-ip/typedef"
)

type MessageRouterRequest struct {
	Service         typedef.Usint
	RequestPathSize typedef.Usint
	RequestPath     []byte
	RequestData     []byte
}

func (m *MessageRouterRequest) Encode() []byte {
	if m.RequestPathSize == 0 {
		m.RequestPathSize = typedef.Usint(len(m.RequestPath) / 2)
	}

	buffer := new(bytes.Buffer)
	WriteByte(buffer, m.Service)
	WriteByte(buffer, m.RequestPathSize)
	WriteByte(buffer, m.RequestPath)
	WriteByte(buffer, m.RequestData)

	return buffer.Bytes()
}

func (m *MessageRouterRequest) New(service typedef.Usint, path []byte, data []byte) {
	m.Service = service
	m.RequestPathSize = typedef.Usint(len(path) / 2)
	m.RequestPath = path
	m.RequestData = data
}

type MessageRouterResponse struct {
	ReplyService           typedef.Usint
	Reserved               typedef.Usint
	GeneralStatus          typedef.Usint
	SizeOfAdditionalStatus typedef.Usint
	AdditionalStatus       []byte
	ResponseData           []byte
}

func (m *MessageRouterResponse) Decode(data []byte) {
	dataReader := bytes.NewReader(data)
	ReadByte(dataReader, &m.ReplyService)
	ReadByte(dataReader, &m.Reserved)
	ReadByte(dataReader, &m.GeneralStatus)
	ReadByte(dataReader, &m.SizeOfAdditionalStatus)
	m.AdditionalStatus = make([]byte, m.SizeOfAdditionalStatus*2)
	ReadByte(dataReader, &m.AdditionalStatus)
	m.ResponseData = make([]byte, dataReader.Len())
	ReadByte(dataReader, &m.ResponseData)
}
