package go_ethernet_ip

import (
	"github.com/loki-os/go-ethernet-ip/types"
	"math/rand"
	"time"
)

func contextGenerator() types.ULInt {
	time.Sleep(time.Nanosecond)
	rand.Seed(time.Now().UnixNano())
	return types.ULInt(rand.Int63())
}
