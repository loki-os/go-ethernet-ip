package bufferx

import (
	"bytes"
	"encoding/binary"
)

type BufferX struct {
	*bytes.Buffer
	err error
}

func (b *BufferX) WL(target interface{}) {
	b.err = binary.Write(b, binary.LittleEndian, target)
}

func (b *BufferX) WB(target interface{}) {
	b.err = binary.Write(b, binary.BigEndian, target)
}

func (b *BufferX) RL(target interface{}) {
	b.err = binary.Read(b.Buffer, binary.LittleEndian, target)
}

func (b *BufferX) RB(target interface{}) {
	b.err = binary.Read(b.Buffer, binary.BigEndian, target)
}

func (b *BufferX) Error() error {
	return b.err
}

func New(data []byte) *BufferX {
	var buffer *bytes.Buffer
	if data == nil {
		buffer = new(bytes.Buffer)
	} else {
		buffer = bytes.NewBuffer(data)
	}
	return &BufferX{
		Buffer: buffer,
	}
}
