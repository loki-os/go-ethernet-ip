package packet

import "github.com/loki-os/go-ethernet-ip/types"

const (
	ServiceGetInstanceAttributeList types.USInt = 0x55
	ServiceGetAttributes            types.USInt = 0x03
	ServiceGetAttributeAll          types.USInt = 0x01
	ServiceGetAttributeSingle       types.USInt = 0x0e
	ServiceReset                    types.USInt = 0x05
	ServiceStart                    types.USInt = 0x06
	ServiceStop                     types.USInt = 0x07
	ServiceCreate                   types.USInt = 0x08
	ServiceDelete                   types.USInt = 0x09
	ServiceMultipleServicePacket    types.USInt = 0x0a
	ServiceApplyAttributes          types.USInt = 0x0d
	ServiceSetAttributeSingle       types.USInt = 0x10
	ServiceFindNext                 types.USInt = 0x11
	ServiceReadTag                  types.USInt = 0x4c
	ServiceWriteTag                 types.USInt = 0x4d
	ServiceReadTagFragmented        types.USInt = 0x52
	ServiceWriteTagFragmented       types.USInt = 0x53
	ServiceReadModifyWriteTag       types.USInt = 0x4e
	ServiceForwardOpen              types.USInt = 0x54
	ServiceForwardOpenLarge         types.USInt = 0x5b
	ServiceForwardClose             types.USInt = 0x4e
)
