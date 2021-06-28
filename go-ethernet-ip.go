package go_ethernet_ip

func Init() {
	tcp, err := NewTCP("192.168.1.1", nil)
	if err != nil {
		panic(err)
	}

	tcp.established = true
	tcp.AllTags()
}
