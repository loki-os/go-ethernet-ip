package go_ethernet_ip

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

type EIPUDP struct {
	config  *Config
	udpAddr map[*net.UDPAddr]bool
	udpConn *net.UDPConn

	Devices         map[string]*Device
	InterfaceHandle func(string, *ListInterface)
	ServicesHandle  func(string, *ListServices)
}

func (e *EIPUDP) Connect() error {
	if e.udpAddr == nil {
		return errors.New("tcp EIP Object can't call udp function")
	}

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", DefaultConfig().UDPPort))
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

func (e *EIPUDP) Close() {
	e.udpConn.Close()
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
		_l := e.ListIdentityDecode(encapsulationPacket)
		if _l != nil {
			item := _l.Items[0]
			device := &Device{
				IP:           addr.IP,
				VendorID:     item.VendorID,
				DeviceType:   item.DeviceType,
				ProductCode:  item.ProductCode,
				Major:        item.Major,
				Minor:        item.Minor,
				Status:       item.Status,
				SerialNumber: item.SerialNumber,
				ProductName:  string(item.ProductName),
			}

			e.Devices[addr.IP.String()] = device
		}
	case EIPCommandListInterfaces:
		_l := e.ListInterfaceDecode(encapsulationPacket)
		if e.InterfaceHandle != nil {
			e.InterfaceHandle(addr.IP.String(), _l)
		}
	case EIPCommandListServices:
		_l := e.ListServicesDecode(encapsulationPacket)
		if e.ServicesHandle != nil {
			e.ServicesHandle(addr.IP.String(), _l)
		}
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

func (e *EIPUDP) send(message []byte) {
	for udpAddr, _ := range e.udpAddr {
		e.udpConn.WriteTo(message, udpAddr)
	}
}

func NewUDPWithAddress(address string, config *Config) (*EIPUDP, error) {
	if config == nil {
		config = DefaultConfig()
	}

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, config.UDPPort))
	if err != nil {
		return nil, err
	}

	eip := &EIPUDP{}
	eip.config = config
	eip.udpAddr = make(map[*net.UDPAddr]bool)
	eip.udpAddr[udpAddr] = true
	eip.Devices = make(map[string]*Device)
	return eip, nil
}

func NewUDPWithAutoScan(config *Config) (*EIPUDP, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	if config == nil {
		config = DefaultConfig()
	}

	eip := &EIPUDP{}
	eip.config = config
	eip.Devices = make(map[string]*Device)
	eip.udpAddr = make(map[*net.UDPAddr]bool)

	for _, _interface := range netInterfaces {
		if (_interface.Flags & net.FlagUp) != 0 {
			addresses, _ := _interface.Addrs()
			for _, address := range addresses {
				if IPNet, ok := address.(*net.IPNet); ok {
					if IPNet.IP.To4() != nil {
						ip := IPNet.IP.To4()
						ip[3] = 255
						udpAddr, err2 := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, eip.config.UDPPort))
						if err2 == nil {
							eip.udpAddr[udpAddr] = true
						}
					}
				}
			}
		}
	}

	return eip, nil
}
