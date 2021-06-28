package utils

import "github.com/loki-os/go-ethernet-ip/types"

func Len(data []byte) types.USInt {
	_len := len(data)
	if _len%2 == 1 {
		_len += 1
	}
	return types.USInt(_len / 2)
}
