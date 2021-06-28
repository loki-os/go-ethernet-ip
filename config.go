package go_ethernet_ip

import "github.com/loki-os/go-ethernet-ip/types"

type Config struct {
	TCPPort     uint16
	UDPPort     uint16
	Slot        uint8
	TimeTick    types.USInt
	TimeTickOut types.USInt
}

func DefaultConfig() *Config {
	cfg := &Config{}
	cfg.TCPPort = 0xAF12
	cfg.UDPPort = 0xAF12
	cfg.Slot = 0
	cfg.TimeTick = types.USInt(3)
	cfg.TimeTickOut = types.USInt(250)
	return cfg
}
