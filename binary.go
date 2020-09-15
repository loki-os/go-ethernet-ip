package go_ethernet_ip

import (
	"encoding/binary"
	"io"
)

func WriteByte(writer io.Writer, target interface{}) {
	binary.Write(writer, binary.LittleEndian, target)
}

func ReadByte(reader io.Reader, target interface{}) {
	binary.Read(reader, binary.LittleEndian, target)
}

func ReadByteBigEndian(reader io.Reader, target interface{}) {
	binary.Read(reader, binary.BigEndian, target)
}
