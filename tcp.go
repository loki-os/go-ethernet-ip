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

type EIPTCP struct {
	config   *Config
	tcpAddr  *net.TCPAddr
	tcpConn  *net.TCPConn
	sender   chan []byte
	receiver map[typedef.Ulint]chan *EncapsulationPacket

	ioCancel context.CancelFunc
	buffer   []byte
	session  typedef.Udint

	Connected    func()
	Disconnected func(error)
	Reconnecting func()
}

func (e *EIPTCP) Connect() error {
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

	return e.connected()
}

func (e *EIPTCP) Close() {
	e.ioCancel()
	e.tcpConn.Close()
	e.tcpConn = nil
}

func (e *EIPTCP) connected() error {
	ctx, cancel := context.WithCancel(context.Background())
	e.ioCancel = cancel
	go e.write(ctx)
	go e.read(ctx)

	if e.config.AutoSession {
		return e.RegisterSession()
	}

	return nil
}

func (e *EIPTCP) write(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-e.sender:
			_, _ = e.tcpConn.Write(data)
		}
	}
}

func (e *EIPTCP) read(ctx context.Context) {
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
				e.disconnect(err)
			}

			e.buffer = append(e.buffer, buf[0:length]...)
			read, encapsulationPackets, err := e.slice(e.buffer)
			if err != nil {
				panic(err)
			}

			e.buffer = e.buffer[read:]

			for _, encapsulationPacket := range encapsulationPackets {
				channel, ok := e.receiver[encapsulationPacket.SenderContext]
				if ok {
					channel <- encapsulationPacket
					close(channel)
					delete(e.receiver, encapsulationPacket.SenderContext)
				}
			}
		}
	}
}

func (e *EIPTCP) slice(data []byte) (uint64, []*EncapsulationPacket, error) {
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

func (e *EIPTCP) disconnect(err error) {
	if e.Disconnected != nil {
		e.Disconnected(err)
	}

	e.ioCancel()
	e.tcpConn.Close()
	e.tcpConn = nil

	if e.config.TCPReconnectionInterval != 0 {
		time.Sleep(e.config.TCPReconnectionInterval)
		if e.Reconnecting != nil {
			e.Reconnecting()
		}

		err := e.Connect()
		if err != nil {
			panic(err)
		}
	}
}

func NewTcpWithAddress(addr string, config *Config) (*EIPTCP, error) {
	if config == nil {
		config = DefaultConfig()
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, config.TCPPort))
	if err != nil {
		return nil, err
	}

	eip := &EIPTCP{}
	eip.config = config
	eip.tcpAddr = tcpAddr
	eip.sender = make(chan []byte)
	eip.receiver = make(map[typedef.Ulint]chan *EncapsulationPacket)

	return eip, nil
}
