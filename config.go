package go_ethernet_ip

import "time"

type config struct {
	TCPPort                 uint16
	UDPPort                 uint16
	TCPReconnectionInterval time.Duration
}

var defaultConfig *config

func init() {
	defaultConfig = &config{}
	defaultConfig.TCPPort = 0xAF12
	defaultConfig.UDPPort = 0xAF12
	defaultConfig.TCPReconnectionInterval = 0
}
