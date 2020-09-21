package go_ethernet_ip

import (
	"encoding/json"
	"github.com/loki-os/go-ethernet-ip/typedef"
	"log"
	"net"
	"time"
)

type Device struct {
	IP           net.IP
	VendorID     typedef.Uint
	DeviceType   typedef.Uint
	ProductCode  typedef.Uint
	Major        typedef.Usint
	Minor        typedef.Usint
	Status       typedef.Word
	SerialNumber typedef.Udint
	ProductName  string
	tcp          *EIPTCP
}

func (d *Device) String() string {
	b, _ := json.MarshalIndent(*d, "", "\t")
	return string(b)
}

func (d *Device) Connect(config *Config) error {
	if d.tcp != nil {
		return nil
	}
	tcp, e := NewTcpWithAddress(d.IP.String(), config)
	if e != nil {
		return e
	}

	e1 := tcp.Connect()
	if e1 != nil {
		return e1
	}

	d.tcp = tcp
	return nil
}

func (d *Device) ListIdentity() (*ListIdentity, error) {
	e := d.Connect(nil)
	if e != nil {
		return nil, e
	}

	return d.tcp.ListIdentity()
}

func (d *Device) ListServices() (*ListServices, error) {
	e := d.Connect(nil)
	if e != nil {
		return nil, e
	}

	return d.tcp.ListServices()
}

func (d *Device) ListInterface() (*ListInterface, error) {
	e := d.Connect(nil)
	if e != nil {
		return nil, e
	}

	return d.tcp.ListInterface()
}

func (d *Device) SendRRData(cpf *CommonPacketFormat, timeout typedef.Uint) (*SendDataSpecificData, error) {
	e := d.Connect(nil)
	if e != nil {
		return nil, e
	}

	return d.tcp.SendRRData(cpf, timeout)
}

func (d *Device) SendUnitData(cpf *CommonPacketFormat, timeout typedef.Uint) (*SendDataSpecificData, error) {
	e := d.Connect(nil)
	if e != nil {
		return nil, e
	}

	return d.tcp.SendUnitData(cpf, timeout)
}

func (d *Device) Disconnect() {
	d.tcp.UnRegisterSession()
	d.tcp.Close()
}

func NewDevice(ip string) (*Device, error) {
	cfg := DefaultConfig()
	cfg.AutoSession = false
	tcp, e := NewTcpWithAddress(ip, cfg)

	if e != nil {
		return nil, e
	}

	e1 := tcp.Connect()
	if e1 != nil {
		return nil, e1
	}

	listIdentity, e2 := tcp.ListIdentity()
	if e2 != nil {
		return nil, e2
	}

	item := listIdentity.Items[0]
	device := &Device{
		IP:           tcp.tcpAddr.IP,
		VendorID:     item.VendorID,
		DeviceType:   item.DeviceType,
		ProductCode:  item.ProductCode,
		Major:        item.Major,
		Minor:        item.Minor,
		Status:       item.Status,
		SerialNumber: item.SerialNumber,
		ProductName:  string(item.ProductName),
	}

	tcp.Close()
	return device, nil
}

func GetLanDevices(scanTimeout time.Duration) (map[string]*Device, error) {
	udp, e := NewUDPWithAutoScan(nil)
	if e != nil {
		return nil, e
	}

	e1 := udp.Connect()
	if e1 != nil {
		log.Println(e1)
		return nil, e1
	}

	udp.ListIdentity()
	time.Sleep(scanTimeout)

	udp.Close()
	return udp.Devices, nil
}
