package path

import (
	"github.com/loki-os/go-ethernet-ip/bufferx"
	"github.com/loki-os/go-ethernet-ip/types"
)

type SegmentType types.USInt

const (
	SegmentTypePort      SegmentType = 0 << 5
	SegmentTypeLogical   SegmentType = 1 << 5
	SegmentTypeNetwork   SegmentType = 2 << 5
	SegmentTypeSymbolic  SegmentType = 3 << 5
	SegmentTypeData      SegmentType = 4 << 5
	SegmentTypeDataType1 SegmentType = 5 << 5
	SegmentTypeDataType2 SegmentType = 6 << 5
)

func Paths(arg ...[]byte) []byte {
	io := bufferx.New(nil)
	for i := 0; i < len(arg); i++ {
		io.WL(arg[i])
	}
	return io.Bytes()
}

type DataTypes types.USInt

const (
	DataTypeSimple DataTypes = 0x0
	DataTypeANSI   DataTypes = 0x11
)

type LogicalType types.USInt

const (
	LogicalTypeClassID     LogicalType = 0 << 2
	LogicalTypeInstanceID  LogicalType = 1 << 2
	LogicalTypeMemberID    LogicalType = 2 << 2
	LogicalTypeConnPoint   LogicalType = 3 << 2
	LogicalTypeAttributeID LogicalType = 4 << 2
	LogicalTypeSpecial     LogicalType = 5 << 2
	LogicalTypeServiceID   LogicalType = 6 << 2
)

func DataBuild(tp DataTypes, data []byte, padded bool) []byte {
	io := bufferx.New(nil)

	firstByte := uint8(SegmentTypeData) | uint8(tp)
	io.WL(firstByte)

	length := uint8(len(data))
	io.WL(length)

	io.WL(data)

	if padded && io.Len()%2 == 1 {
		io.WL(uint8(0))
	}

	return io.Bytes()
}

func LogicalBuild(tp LogicalType, address types.UDInt, padded bool) []byte {
	format := uint8(0)

	if address <= 255 {
		format = 0
	} else if address > 255 && address <= 65535 {
		format = 1
	} else {
		format = 2
	}

	io := bufferx.New(nil)
	firstByte := uint8(SegmentTypeLogical) | uint8(tp) | format
	io.WL(firstByte)

	if address > 255 && address <= 65535 && padded {
		io.WL(uint8(0))
	}

	if address <= 255 {
		io.WL(uint8(address))
	} else if address > 255 && address <= 65535 {
		io.WL(uint16(address))
	} else {
		io.WL(address)
	}

	return io.Bytes()
}

func PortBuild(link []byte, portID uint16, padded bool) []byte {
	extendedLinkTag := len(link) > 1
	extendedPortTag := !(portID < 15)

	io := bufferx.New(nil)

	firstByte := uint8(SegmentTypePort)
	if extendedLinkTag {
		firstByte = firstByte | 0x10
	}

	if !extendedPortTag {
		firstByte = firstByte | uint8(portID)
	} else {
		firstByte = firstByte | 0xf
	}

	io.WL(firstByte)

	if extendedLinkTag {
		io.WL(uint8(len(link)))
	}

	if extendedPortTag {
		io.WL(portID)
	}

	io.WL(link)

	if padded && io.Len()%2 == 1 {
		io.WL(uint8(0))
	}

	return io.Bytes()
}
