package go_ethernet_ip

import "time"

type Config struct {
	TCPPort                 uint16
	UDPPort                 uint16
	TCPTimeout              time.Duration
	TCPReconnectionInterval time.Duration
	AutoSession             bool
}

func DefaultConfig() *Config {
	cfg := &Config{}
	cfg.TCPPort = 0xAF12
	cfg.UDPPort = 0xAF12
	cfg.TCPTimeout = time.Second
	cfg.TCPReconnectionInterval = 0
	cfg.AutoSession = true
	return cfg
}
