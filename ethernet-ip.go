package go_ethernet_ip

import (
	"github.com/loki-os/go-ethernet-ip/typedef"
	"net"
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
}
