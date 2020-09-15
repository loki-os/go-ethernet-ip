package go_ethernet_ip

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

type EIP struct {
	config               *config
	tcpAddr              *net.TCPAddr
	tcpConn              *net.TCPConn
	udpAddr              *net.UDPAddr
	udpConn              *net.UDPConn
	sender               chan []byte
	ioCancel             context.CancelFunc
	buffer               []byte
	ListIdentityCallBack func()
}

func (e *EIP) TcpConnect() error {
	var err error
	e.tcpConn, err = net.DialTCP("tcp", nil, e.tcpAddr)
	if err != nil {
		return err
	}

	err = e.tcpConn.SetKeepAlive(true)
	err = e.tcpConn.SetKeepAlivePeriod(time.Second * 10)
	if err != nil {
		return err
	}

	e.tcpConnected()
	return nil
}

func (e *EIP) UdpConnect() error {
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
	go e.udpRead()
	return nil
}

func (e *EIP) udpRead() {
	for {
		data := make([]byte, 1024*64)
		read, remoteAddr, err := e.udpConn.ReadFromUDP(data)
		if err != nil {
			continue
		}
		fmt.Println(read, remoteAddr)
		_, encapsulationPackets, err3 := e.slice(data[0:read])
		if err3 != nil {
			continue
		}

		e.encapsulationParser(encapsulationPackets[0])
	}
}

func (e *EIP) udpSend(message []byte) error {
	if e.udpAddr == nil {
		return errors.New("tcp EIP Object can't call udp function")
	}

	conn, err1 := net.DialUDP("udp", nil, e.udpAddr)
	if err1 != nil {
		return err1
	}

	_, err2 := conn.Write(message)
	if err2 != nil {
		return err2
	}

	return nil
}

func (e *EIP) tcpConnected() {
	if e.config.Connected != nil {
		e.config.Connected()
	}
	ctx, cancel := context.WithCancel(context.Background())
	e.ioCancel = cancel
	go e.tcpWrite(ctx)
	go e.tcpRead(ctx)
}

func (e *EIP) tcpDisconnect(err error) {
	if e.config.Disconnected != nil {
		e.config.Disconnected(err)
	}

	e.ioCancel()
	e.tcpConn.Close()
	e.tcpConn = nil

	if e.config.TCPReconnectionInterval != 0 {
		time.Sleep(e.config.TCPReconnectionInterval)
		if e.config.Reconnecting != nil {
			e.config.Reconnecting()
		}

		err := e.TcpConnect()
		if err != nil {
			panic(err)
		}
	}
}

func (e *EIP) tcpRead(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			go e.tcpDisconnect(err.(error))
		}
	}()

	buf := make([]byte, 1024*64)
	var err error
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var length int
			length, err = e.tcpConn.Read(buf)
			if err != nil {
				panic(err)
			}

			e.buffer = append(e.buffer, buf[0:length]...)
			read, encapsulationPackets, err := e.slice(e.buffer)
			if err != nil {
				panic(err)
			}

			e.buffer = e.buffer[read:]

			for _, encapsulationPacket := range encapsulationPackets {
				e.encapsulationParser(encapsulationPacket)
			}
		}
	}
}

func (e *EIP) encapsulationParser(encapsulationPacket *EncapsulationPacket) {
	switch encapsulationPacket.Command {
	case EIPCommandListIdentity:
		e.ListIdentityDecode(encapsulationPacket)
	case EIPCommandRegisterSession:
		e.RegisterSessionDecode(encapsulationPacket)
	default:
		panic("encapsulation with wrong command")
	}
}

func (e *EIP) tcpWrite(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-e.sender:
			_, _ = e.tcpConn.Write(data)
		}
	}
}

func (e *EIP) slice(data []byte) (uint64, []*EncapsulationPacket, error) {
	if len(data) < 24 {
		return 0, nil, nil
	}

	var result []*EncapsulationPacket

	dataReader := bytes.NewReader(data)
	count := dataReader.Len()

	for dataReader.Len() > 23 {
		_encapsulationPacket := &EncapsulationPacket{}
		ReadByte(dataReader, &_encapsulationPacket.EncapsulationHeader)

		if _encapsulationPacket.Options != 0 {
			return 0, nil, errors.New("wrong package with non-zero option")
		}

		if int(_encapsulationPacket.Length) > dataReader.Len() {
			count += 24
			break
		} else {
			if _encapsulationPacket.Length > 0 {
				_encapsulationPacket.CommandSpecificData = make([]byte, _encapsulationPacket.Length)
				_, e := dataReader.Read(_encapsulationPacket.CommandSpecificData)
				if e != nil {
					panic(e)
				}
			}

			result = append(result, _encapsulationPacket)
		}
	}

	count = count - dataReader.Len()
	return uint64(count), result, nil
}

func NewTCP(addr string, config *config) (*EIP, error) {
	eip := &EIP{}

	if config == nil {
		eip.config = defaultConfig
	} else {
		eip.config = config
	}

	var err error
	eip.tcpAddr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, eip.config.TCPPort))
	if err != nil {
		return nil, err
	}

	eip.sender = make(chan []byte)

	return eip, nil
}

func NewUDP(addr string, config *config) (*EIP, error) {
	eip := &EIP{}

	if config == nil {
		eip.config = defaultConfig
	} else {
		eip.config = config
	}

	var err error
	eip.udpAddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", eip.config.BroadcastAddress, eip.config.UDPPort))
	if err != nil {
		return nil, err
	}

	eip.sender = make(chan []byte)

	return eip, nil
}
