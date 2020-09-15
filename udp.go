package go_ethernet_ip

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

type EIPUDP struct {
	config  *config
	udpAddr []*net.UDPAddr
	udpConn *net.UDPConn
	Devices map[string]*Device
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
	for _, _addr := range e.udpAddr {
		e.udpConn.WriteTo(message, _addr)
	}
}

func NewUDPWithAddress(addr string, config *config) (*EIPUDP, error) {
	eip := &EIPUDP{}

	if config == nil {
		eip.config = defaultConfig
	} else {
		eip.config = config
	}
	eip.Devices = make(map[string]*Device)

	var err error
	var _udpAddr *net.UDPAddr
	_udpAddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, eip.config.UDPPort))
	if err != nil {
		return nil, err
	}

	eip.udpAddr = []*net.UDPAddr{_udpAddr}

	return eip, nil
}

func NewUDPWithAutoScan(config *config) (*EIPUDP, error) {
	eip := &EIPUDP{}

	if config == nil {
		eip.config = defaultConfig
	} else {
		eip.config = config
	}
	eip.Devices = make(map[string]*Device)

	eip.udpAddr = []*net.UDPAddr{}

	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok {
					if ipnet.IP.To4() != nil {
						ip := ipnet.IP.To4()
						ip[3] = 255
						_udpAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, eip.config.UDPPort))
						eip.udpAddr = append(eip.udpAddr, _udpAddr)
					}
				}
			}
		}
	}

	return eip, nil
}
