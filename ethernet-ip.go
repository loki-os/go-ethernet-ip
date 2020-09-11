package go_ethernet_ip

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/loki-os/go-ethernet-ip/typedef"
	"net"
	"time"
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
	//check exist
	if e == nil {
		return nil, errors.New("EIP not initialized")
	}

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

type EIP struct {
	config   *config
	tcpAddr  *net.TCPAddr
	tcpConn  *net.TCPConn
	sender   chan []byte
	ioCancel context.CancelFunc
	buffer   []byte
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

func (e *EIP) tcpConnected() {
	if e.config.Connected != nil {
		e.config.Connected()
	}
	ctx, cancel := context.WithCancel(context.Background())
	e.ioCancel = cancel
	go e.write(ctx)
	go e.read(ctx)
}

func (e *EIP) disconnect(err error) {
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

func (e *EIP) read(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			go e.disconnect(err.(error))
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
	case 0x63:

	default:
		panic("encapsulation with wrong command")
	}
}

func (e *EIP) write(ctx context.Context) {
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

				result = append(result, _encapsulationPacket)
			}
		}
	}

	count = count - dataReader.Len()
	return uint64(count), result, nil
}

func New(addr string, config *config) (*EIP, error) {
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

	return eip, nil
}
