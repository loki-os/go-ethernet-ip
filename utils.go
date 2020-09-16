package go_ethernet_ip

import (
	"github.com/loki-os/go-ethernet-ip/typedef"
	"math/rand"
	"time"
)

func CtxGenerator() typedef.Ulint {
	rand.Seed(time.Now().UnixNano())
	return typedef.Ulint(rand.Uint64())
}
