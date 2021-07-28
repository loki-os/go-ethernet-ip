package go_ethernet_ip

import (
	"bytes"
	"fmt"
	"github.com/loki-os/go-ethernet-ip/bufferx"
	"github.com/loki-os/go-ethernet-ip/messages/packet"
	"github.com/loki-os/go-ethernet-ip/path"
	"github.com/loki-os/go-ethernet-ip/types"
	"sync"
)

const (
	NULL   types.UInt = 0x00
	BOOL   types.UInt = 0xc1
	SINT   types.UInt = 0xc2
	INT    types.UInt = 0xc3
	DINT   types.UInt = 0xc4
	LINT   types.UInt = 0xc5
	USINT  types.UInt = 0xc6
	UINT   types.UInt = 0xc7
	UDINT  types.UInt = 0xc8
	ULINT  types.UInt = 0xc9
	REAL   types.UInt = 0xca
	LREAL  types.UInt = 0xcb
	STRING types.UInt = 0xfce
)

var TypeMap = map[types.UInt]string{
	NULL:   "NULL",
	BOOL:   "BOOL",
	SINT:   "SINT",
	INT:    "INT",
	DINT:   "DINT",
	LINT:   "LINT",
	USINT:  "USINT",
	UINT:   "UINT",
	UDINT:  "UDINT",
	ULINT:  "ULINT",
	REAL:   "REAL",
	LREAL:  "LREAL",
	STRING: "STRING",
}

type Tag struct {
	Lock *sync.Mutex
	TCP  *EIPTCP

	instanceID types.UDInt
	nameLen    types.UInt
	name       []byte
	Type       types.UInt
	dim1Len    types.UDInt
	dim2Len    types.UDInt
	dim3Len    types.UDInt
	changed    bool

	value    []byte
	wValue   []byte
	Onchange func()
}

func (t *Tag) Read() error {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	res, err := t.TCP.Send(t.readRequest())
	if err != nil {
		return err
	}

	mrres := new(packet.MessageRouterResponse)
	mrres.Decode(res.Packet.Items[1].Data)

	t.readParser(mrres)
	return nil
}

func (t *Tag) readRequest() *packet.MessageRouterRequest {
	io := bufferx.New(nil)
	io.WL(t.count())
	mr := packet.NewMessageRouter(packet.ServiceReadTag, packet.Paths(
		path.LogicalBuild(path.LogicalTypeClassID, 0x6B, true),
		path.LogicalBuild(path.LogicalTypeInstanceID, t.instanceID, true),
	), io.Bytes())
	return mr
}

func (t *Tag) readParser(mr *packet.MessageRouterResponse) {
	io := bufferx.New(mr.ResponseData)

	_t := uint16(0)
	io.RL(&_t)

	if _t == 0x2a0 {
		io.RL(&_t)
	}

	payload := make([]byte, io.Len())
	io.RL(payload)

	if bytes.Compare(t.value, payload) != 0 {
		t.value = payload
		if t.Onchange != nil {
			go t.Onchange()
		}
	}
}

func (t *Tag) Write() error {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	if t.wValue != nil {
		copy(t.wValue, t.value)
	}
	_, err := t.TCP.Send(multiple(t.writeRequest()))
	if err == nil {
		if t.wValue != nil {
			copy(t.value, t.wValue)
			t.wValue = nil
		}
	}
	return err
}

func (t *Tag) writeRequest() []*packet.MessageRouterRequest {
	var result []*packet.MessageRouterRequest
	if 0x8000&t.Type == 0 {
		io := bufferx.New(nil)
		io.WL(t.Type)
		io.WL(t.count())
		io.WL(t.wValue)

		mr := packet.NewMessageRouter(packet.ServiceWriteTag, packet.Paths(
			path.LogicalBuild(path.LogicalTypeClassID, 0x6B, true),
			path.LogicalBuild(path.LogicalTypeInstanceID, t.instanceID, true),
		), io.Bytes())
		result = append(result, mr)
	} else {
		// only string
		io := bufferx.New(nil)
		io.WL(DINT)
		io.WL(types.UInt(1))
		io.WL(types.UDInt(len(t.wValue)))
		mr1 := packet.NewMessageRouter(packet.ServiceWriteTag, packet.Paths(
			path.LogicalBuild(path.LogicalTypeClassID, 0x6B, true),
			path.LogicalBuild(path.LogicalTypeInstanceID, t.instanceID, true),
			path.DataBuild(path.DataTypeANSI, []byte("LEN"), true),
		), io.Bytes())
		result = append(result, mr1)

		io1 := bufferx.New(nil)
		io1.WL(SINT)
		io1.WL(types.UInt(len(t.wValue)))
		io1.WL(t.wValue)
		mr2 := packet.NewMessageRouter(packet.ServiceWriteTag, packet.Paths(
			path.LogicalBuild(path.LogicalTypeClassID, 0x6B, true),
			path.LogicalBuild(path.LogicalTypeInstanceID, t.instanceID, true),
			path.DataBuild(path.DataTypeANSI, []byte("DATA"), true),
		), io1.Bytes())
		result = append(result, mr2)
	}

	return result
}

func (t *Tag) SetInt32(i int32) {
	t.changed = true
	io := bufferx.New(nil)
	io.WL(i)
	t.wValue = io.Bytes()
}

func (t *Tag) SetString(i string) {
	t.changed = true
	io := bufferx.New(nil)
	io.WL([]byte(i))
	t.wValue = io.Bytes()
}

func (t *Tag) dims() types.USInt {
	return types.USInt((0x6000 & t.Type) >> 13)
}

func (t *Tag) TypeString() string {
	var _type string
	if 0x8000&t.Type == 0 {
		_type = "atomic"
	} else {
		_type = "struct"
	}

	return fmt.Sprintf("%#04x(%6s) | %s | %d dims", uint16(t.Type), TypeMap[0xFFF&t.Type], _type, (0x6000&t.Type)>>13)
}

func (t *Tag) Name() string {
	return string(t.name)
}

func (t *Tag) count() types.UInt {
	a := types.UInt(1)
	if t.dim1Len > 0 {
		a = types.UInt(t.dim1Len)
	}
	b := types.UInt(1)
	if t.dim2Len > 0 {
		b = types.UInt(t.dim2Len)
	}
	c := types.UInt(1)
	if t.dim3Len > 0 {
		c = types.UInt(t.dim3Len)
	}
	return a * b * c
}

func (t *Tag) Int32() int32 {
	io := bufferx.New(t.value)
	var val int32
	io.RL(&val)
	return val
}

func (t *Tag) String() string {
	io := bufferx.New(t.value)
	_len := types.UDInt(0)
	io.RL(&_len)
	val := make([]byte, _len)
	io.RL(&val)
	return string(val)
}

func (t *Tag) XInt32() int32 {
	var _value []byte
	if len(t.wValue) > 0 {
		_value = t.wValue
	} else {
		_value = t.value
	}
	io := bufferx.New(_value)
	var val int32
	io.RL(&val)
	return val
}

func (t *Tag) XString() string {
	var _value []byte
	if len(t.wValue) > 0 {
		_value = t.wValue
	} else {
		_value = t.value
	}
	io := bufferx.New(_value)
	_len := types.UDInt(0)
	io.RL(&_len)
	val := make([]byte, _len)
	io.RL(&val)
	return string(val)
}

func multiple(mrs []*packet.MessageRouterRequest) *packet.MessageRouterRequest {
	if len(mrs) == 1 {
		return mrs[0]
	} else {
		io := bufferx.New(nil)
		io.WL(types.UInt(len(mrs)))
		offset := 2 * (len(mrs) + 1) // offset0 = 上一个(2) + 所有offset的长度的长度综合 2xN
		io.WL(types.UInt(offset))
		for i := range mrs {
			if i != len(mrs)-1 {
				offset += len(mrs[i].Encode())
				io.WL(types.UInt(offset))
			}
		}
		for i := range mrs {
			io.WL(mrs[i].Encode())
		}
		return packet.NewMessageRouter(packet.ServiceMultipleServicePacket, packet.Paths(
			path.LogicalBuild(path.LogicalTypeClassID, 0x02, true),
			path.LogicalBuild(path.LogicalTypeInstanceID, 0x01, true),
		), io.Bytes())
	}
}

func (t *EIPTCP) AllTags() (map[string]*Tag, error) {
	result := make(map[string]*Tag)
	return t.allTags(result, 0)
}

func (t *EIPTCP) allTags(tagMap map[string]*Tag, instanceID types.UDInt) (map[string]*Tag, error) {
	paths := packet.Paths(
		path.LogicalBuild(path.LogicalTypeClassID, 0x6B, true),
		path.LogicalBuild(path.LogicalTypeInstanceID, instanceID, true),
	)

	io := bufferx.New(nil)
	io.WL(types.UInt(3))
	io.WL(types.UInt(1))
	io.WL(types.UInt(2))
	io.WL(types.UInt(8))

	mr := packet.NewMessageRouter(packet.ServiceGetInstanceAttributeList, paths, io.Bytes())

	res, err := t.Send(mr)
	if err != nil {
		return nil, err
	}

	mrres := new(packet.MessageRouterResponse)
	mrres.Decode(res.Packet.Items[1].Data)

	io1 := bufferx.New(mrres.ResponseData)
	for io1.Len() > 0 {
		tag := new(Tag)
		tag.TCP = t
		tag.Lock = new(sync.Mutex)

		io1.RL(&tag.instanceID)
		io1.RL(&tag.nameLen)
		tag.name = make([]byte, tag.nameLen)
		io1.RL(tag.name)
		io1.RL(&tag.Type)
		io1.RL(&tag.dim1Len)
		io1.RL(&tag.dim2Len)
		io1.RL(&tag.dim3Len)

		tagMap[tag.Name()] = tag
		instanceID = tag.instanceID
	}

	if mrres.GeneralStatus == 0x06 {
		return t.allTags(tagMap, instanceID+1)
	}

	return tagMap, nil
}

type TagGroup struct {
	tags map[types.UDInt]*Tag
	Tcp  *EIPTCP
}

func NewTagGroup() *TagGroup {
	return &TagGroup{tags: make(map[types.UDInt]*Tag)}
}

func (tg *TagGroup) Add(tag *Tag) {
	if tg.Tcp == nil {
		tg.Tcp = tag.TCP
	} else {
		if tg.Tcp != tag.TCP {
			return
		}
	}
	tg.tags[tag.instanceID] = tag
}

func (tg *TagGroup) Remove(tag *Tag) {
	delete(tg.tags, tag.instanceID)
}

func (tg *TagGroup) Read() error {
	if len(tg.tags) == 0 {
		return nil
	}

	if len(tg.tags) == 1 {
		for _, v := range tg.tags {
			return v.Read()
		}
	}

	var list []types.UDInt
	var mrs []*packet.MessageRouterRequest

	for i := range tg.tags {
		tg.tags[i].Lock.Lock()
		list = append(list, tg.tags[i].instanceID)
		mrs = append(mrs, tg.tags[i].readRequest())
	}

	defer func() {
		for i := range tg.tags {
			tg.tags[i].Lock.Unlock()
		}
	}()

	_sb := multiple(mrs)
	res, err := tg.Tcp.Send(_sb)
	if err != nil {
		return err
	}

	rmr := &packet.MessageRouterResponse{}
	rmr.Decode(res.Packet.Items[1].Data)

	io1 := bufferx.New(rmr.ResponseData)
	count := types.UInt(0)
	io1.RL(&count)

	if int(count) != len(list) {
		return nil
	}

	var offsets []types.UInt
	for i := types.UInt(0); i < count; i++ {
		one := types.UInt(0)
		io1.RL(&one)
		offsets = append(offsets, one)
	}
	for i2 := range list {
		mr := &packet.MessageRouterResponse{}
		if (i2 + 1) != len(offsets) {
			mr.Decode(rmr.ResponseData[offsets[i2]:offsets[i2+1]])
		} else {
			mr.Decode(rmr.ResponseData[offsets[i2]:])
		}
		tg.tags[list[i2]].readParser(mr)
	}

	return nil
}

func (tg *TagGroup) Write() error {
	var list []types.UDInt
	var mrs []*packet.MessageRouterRequest

	for i := range tg.tags {
		tg.tags[i].Lock.Lock()
		if tg.tags[i].changed {
			list = append(list, tg.tags[i].instanceID)
			mrs = append(mrs, tg.tags[i].writeRequest()...)
			tg.tags[i].changed = false
		}
	}

	defer func() {
		for i := range tg.tags {
			tg.tags[i].Lock.Unlock()
		}
	}()

	if len(list) == 0 {
		return nil
	}

	_, err := tg.Tcp.Send(multiple(mrs))
	if err != nil {
		return err
	}
	for i := range tg.tags {
		if tg.tags[i].changed {
			if tg.tags[i].wValue != nil {
				copy(tg.tags[i].value, tg.tags[i].wValue)
				tg.tags[i].wValue = nil
			}
		}
	}

	return nil
}
