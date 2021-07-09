package go_ethernet_ip

import (
	"errors"
	"fmt"
	"github.com/loki-os/go-ethernet-ip/bufferx"
	"github.com/loki-os/go-ethernet-ip/messages/packet"
	"github.com/loki-os/go-ethernet-ip/types"
	"net"
	"sync"
)

type EIPTCP struct {
	config  *Config
	tcpAddr *net.TCPAddr
	tcpConn *net.TCPConn
	session types.UDInt

	established bool
	connID      types.UDInt
	seqNum      types.UInt

	requestLock *sync.Mutex
}

func (t *EIPTCP) reset() {
	t.established = false
}

func (t *EIPTCP) Connect() error {
	t.reset()

	tcpConnection, err := net.DialTCP("tcp", nil, t.tcpAddr)
	if err != nil {
		return err
	}

	err = tcpConnection.SetKeepAlive(true)
	if err != nil {
		return err
	}

	t.tcpConn = tcpConnection

	if err := t.RegisterSession(); err != nil {
		return err
	}

	return nil
}

func (t *EIPTCP) write(data []byte) error {
	_, err := t.tcpConn.Write(data)
	return err
}

func (t *EIPTCP) read() (*packet.Packet, error) {
	buf := make([]byte, 1024*64)
	length, err := t.tcpConn.Read(buf)
	if err != nil {
		return nil, err
	}
	return t.parse(buf[0:length])
}

func (t *EIPTCP) parse(buf []byte) (*packet.Packet, error) {
	if len(buf) < 24 {
		return nil, errors.New("invalid packet, length < 24")
	}
	_packet := new(packet.Packet)
	buffer := bufferx.New(buf)
	buffer.RL(&_packet.Header)
	if buffer.Error() != nil {
		return nil, buffer.Error()
	}
	if _packet.Options != 0 {
		return nil, errors.New("wrong packet with non-zero option")
	}
	if int(_packet.Length) != buffer.Len() {
		return nil, errors.New("wrong packet length")
	}
	_packet.SpecificData = make([]byte, _packet.Length)
	buffer.RL(_packet.SpecificData)
	if buffer.Error() != nil {
		return nil, buffer.Error()
	}
	return _packet, nil
}

func NewTCP(address string, config *Config) (*EIPTCP, error) {
	if config == nil {
		config = DefaultConfig()
	}

	tcpAddress, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", address, config.TCPPort))
	if err != nil {
		return nil, err
	}

	return &EIPTCP{
		requestLock: new(sync.Mutex),
		config:      config,
		tcpAddr:     tcpAddress,
	}, nil
}
