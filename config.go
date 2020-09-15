package go_ethernet_ip

import "time"

type config struct {
	TCPPort                 uint16
	UDPPort                 uint16
	BroadcastAddress        string
	TCPReconnectionInterval time.Duration

	//runtime callback
	Connected    func()
	Disconnected func(error)
	Reconnecting func()
}

var defaultConfig *config

func init() {
	defaultConfig = &config{}
	defaultConfig.TCPPort = 0xAF12
	defaultConfig.UDPPort = 0xAF12
	defaultConfig.BroadcastAddress = "255.255.255.255"
	defaultConfig.TCPReconnectionInterval = 0
}
