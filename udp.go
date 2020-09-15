package go_ethernet_ip

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
)

type EIPUDP struct {
	config  *config
	udpAddr *net.UDPAddr
	udpConn *net.UDPConn
}

func (e *EIPUDP) Connect() error {
	if e.udpAddr == nil {
		return errors.New("tcp EIP Object can't call udp function")
	}

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", defaultConfig.UDPPort))
	if err != nil {
		return err
	}

	udpConn, err2 := net.ListenUDP("udp", udpAddr)
	if err2 != nil {
		return err2
	}

	e.udpConn = udpConn
	go e.read()
	return nil
}

func (e *EIPUDP) read() {
	for {
		data := make([]byte, 1024*64)
		read, remoteAddr, err := e.udpConn.ReadFromUDP(data)
		if err != nil {
			continue
		}
		encapsulationPacket, err2 := e.decode(data[0:read])
		if err2 != nil {
			continue
		}

		e.encapsulationParser(encapsulationPacket, remoteAddr)
	}
}

func (e *EIPUDP) encapsulationParser(encapsulationPacket *EncapsulationPacket, addr *net.UDPAddr) {
	switch encapsulationPacket.Command {
	case EIPCommandListIdentity:
		tcp := (*net.TCPAddr)(addr)
		log.Printf("%+v\n", tcp)
		log.Printf("%+v\n", e.ListIdentityDecode(encapsulationPacket))
	default:
	}
}

func (e *EIPUDP) decode(data []byte) (*EncapsulationPacket, error) {
	if len(data) < 24 {
		return nil, errors.New("wrong package length")
	}

	dataReader := bytes.NewReader(data)

	_encapsulationPacket := &EncapsulationPacket{}
	ReadByte(dataReader, &_encapsulationPacket.EncapsulationHeader)

	if _encapsulationPacket.Options != 0 {
		return nil, errors.New("wrong package with non-zero option")
	}

	if int(_encapsulationPacket.Length) != dataReader.Len() {
		return nil, errors.New("wrong package length")
	} else {
		if _encapsulationPacket.Length > 0 {
			_encapsulationPacket.CommandSpecificData = make([]byte, dataReader.Len())
			ReadByte(dataReader, &_encapsulationPacket.CommandSpecificData)
		}
		return _encapsulationPacket, nil
	}
}

func (e *EIPUDP) send(message []byte) error {
	_, err2 := e.udpConn.WriteTo(message, e.udpAddr)
	if err2 != nil {
		return err2
	}

	return nil
}

func NewUDPWithBroadcastAddress(addr string, config *config) (*EIPUDP, error) {
	eip := &EIPUDP{}

	if config == nil {
		eip.config = defaultConfig
	} else {
		eip.config = config
	}

	var err error
	eip.udpAddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, eip.config.UDPPort))
	if err != nil {
		return nil, err
	}

	return eip, nil
}
