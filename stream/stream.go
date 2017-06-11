package stream

import (
	"encoding/binary"
	"fmt"
)

var ErrBuffOverflow = fmt.Errorf("DSProtocol: buff is too small to io")

type ReadStream interface {
	ReadByte() (b byte, err error)
	ReadUint16() (b uint16, err error)
	ReadUint32() (b uint32, err error)
	ReadUint64() (b uint64, err error)
	ReadBuff(size int) (b []byte, err error)
	CopyBuff(b []byte) error
}

type WriteStream interface {
	WriteByte(b byte) error
	WriteUint16(b uint16) error
	WriteUint32(b uint32) error
	WriteUint64(b uint64) error
	WriteBuff(b []byte) error
}
type CommonStream interface {
	Pos() int
	Size() int
	Left() int
	Reset([]byte)
	Data() []byte
	DataSelect(from, to int) []byte
}

type Stream interface {
	ReadStream
	WriteStream
	CommonStream
}

//BigEndianStream
type BEStream struct {
	pos  int
	buff []byte
}

//LittleEndianStream
type LEStream struct {
	pos  int
	buff []byte
}

func NewBEStream(buff []byte) Stream {
	return &BEStream{
		buff: buff,
	}
}

func NewLEStream(buff []byte) Stream {
	return &LEStream{
		buff: buff,
	}
}

func (impl *BEStream) Pos() int { return impl.pos }

func (impl *BEStream) Size() int { return len(impl.buff) }

func (impl *BEStream) Data() []byte { return impl.buff }

func (impl *BEStream) DataSelect(from, to int) []byte {
	if from < 0 || to < 0 || to > len(impl.buff) || from > len(impl.buff) {
		panic("out of index")
	}
	return impl.buff[from:to]
}

func (impl *BEStream) Left() int { return len(impl.buff) - impl.pos }

func (impl *BEStream) Reset(buff []byte) { impl.pos = 0; impl.buff = buff }

func (impl *BEStream) ReadByte() (b byte, err error) {
	if impl.Left() < 1 {
		return 0, ErrBuffOverflow
	}
	b = impl.buff[impl.pos]
	impl.pos += 1
	return b, nil
}

func (impl *BEStream) ReadUint16() (b uint16, err error) {
	if impl.Left() < 2 {
		return 0, ErrBuffOverflow
	}
	b = binary.BigEndian.Uint16(impl.buff[impl.pos:])
	impl.pos += 2
	return b, nil
}

func (impl *BEStream) ReadUint32() (b uint32, err error) {
	if impl.Left() < 4 {
		return 0, ErrBuffOverflow
	}
	b = binary.BigEndian.Uint32(impl.buff[impl.pos:])
	impl.pos += 4
	return b, nil
}

func (impl *BEStream) ReadUint64() (b uint64, err error) {
	if impl.Left() < 8 {
		return 0, ErrBuffOverflow
	}
	b = binary.BigEndian.Uint64(impl.buff[impl.pos:])
	impl.pos += 8
	return b, nil
}

func (impl *BEStream) ReadBuff(size int) (buff []byte, err error) {
	if impl.Left() < size {
		return nil, ErrBuffOverflow
	}
	buff = make([]byte, size, size)
	copy(buff, impl.buff[impl.pos:impl.pos+size])
	impl.pos += size
	return buff, nil
}

func (impl *BEStream) CopyBuff(b []byte) error {
	if impl.Left() < len(b) {
		return ErrBuffOverflow
	}
	copy(b, impl.buff[impl.pos:impl.pos+len(b)])
	return nil
}

func (impl *BEStream) WriteByte(b byte) error {
	if impl.Left() < 1 {
		return ErrBuffOverflow
	}
	impl.buff[impl.pos] = b
	impl.pos += 1
	return nil
}

func (impl *BEStream) WriteUint16(b uint16) error {
	if impl.Left() < 2 {
		return ErrBuffOverflow
	}
	binary.BigEndian.PutUint16(impl.buff[impl.pos:], b)
	impl.pos += 2
	return nil
}

func (impl *BEStream) WriteUint32(b uint32) error {
	if impl.Left() < 4 {
		return ErrBuffOverflow
	}
	binary.BigEndian.PutUint32(impl.buff[impl.pos:], b)
	impl.pos += 4
	return nil
}

func (impl *BEStream) WriteUint64(b uint64) error {
	if impl.Left() < 8 {
		return ErrBuffOverflow
	}
	binary.BigEndian.PutUint64(impl.buff[impl.pos:], b)
	impl.pos += 8
	return nil
}

func (impl *BEStream) WriteBuff(buff []byte) error {
	if impl.Left() < len(buff) {
		return ErrBuffOverflow
	}
	copy(impl.buff[impl.pos:], buff)
	impl.pos += len(buff)
	return nil
}
func (impl *LEStream) Pos() int { return impl.pos }

func (impl *LEStream) Size() int { return len(impl.buff) }

func (impl *LEStream) Data() []byte { return impl.buff }

func (impl *LEStream) DataSelect(from, to int) []byte {
	if from < 0 || to < 0 || to > len(impl.buff) || from > len(impl.buff) {
		panic("out of index")
	}
	return impl.buff[from:to]
}

func (impl *LEStream) Left() int { return len(impl.buff) - impl.pos }

func (impl *LEStream) Reset(buff []byte) { impl.pos = 0; impl.buff = buff }

func (impl *LEStream) ReadByte() (b byte, err error) {
	if impl.Left() < 1 {
		return 0, ErrBuffOverflow
	}
	b = impl.buff[impl.pos]
	impl.pos += 1
	return b, nil
}

func (impl *LEStream) ReadUint16() (b uint16, err error) {
	if impl.Left() < 2 {
		return 0, ErrBuffOverflow
	}
	b = binary.LittleEndian.Uint16(impl.buff[impl.pos:])
	impl.pos += 2
	return b, nil
}

func (impl *LEStream) ReadUint32() (b uint32, err error) {
	if impl.Left() < 4 {
		return 0, ErrBuffOverflow
	}
	b = binary.LittleEndian.Uint32(impl.buff[impl.pos:])
	impl.pos += 4
	return b, nil
}

func (impl *LEStream) ReadUint64() (b uint64, err error) {
	if impl.Left() < 8 {
		return 0, ErrBuffOverflow
	}
	b = binary.LittleEndian.Uint64(impl.buff[impl.pos:])
	impl.pos += 8
	return b, nil
}

func (impl *LEStream) ReadBuff(size int) (buff []byte, err error) {
	if impl.Left() < size {
		return nil, ErrBuffOverflow
	}
	buff = make([]byte, size, size)
	copy(buff, impl.buff[impl.pos:impl.pos+size])
	impl.pos += size
	return buff, nil
}

func (impl *LEStream) CopyBuff(b []byte) error {
	if impl.Left() < len(b) {
		return ErrBuffOverflow
	}
	copy(b, impl.buff[impl.pos:impl.pos+len(b)])
	return nil
}

func (impl *LEStream) WriteByte(b byte) error {
	if impl.Left() < 1 {
		return ErrBuffOverflow
	}
	impl.buff[impl.pos] = b
	impl.pos += 1
	return nil
}

func (impl *LEStream) WriteUint16(b uint16) error {
	if impl.Left() < 2 {
		return ErrBuffOverflow
	}
	binary.LittleEndian.PutUint16(impl.buff[impl.pos:], b)
	impl.pos += 2
	return nil
}

func (impl *LEStream) WriteUint32(b uint32) error {
	if impl.Left() < 4 {
		return ErrBuffOverflow
	}
	binary.LittleEndian.PutUint32(impl.buff[impl.pos:], b)
	impl.pos += 4
	return nil
}

func (impl *LEStream) WriteUint64(b uint64) error {
	if impl.Left() < 8 {
		return ErrBuffOverflow
	}
	binary.LittleEndian.PutUint64(impl.buff[impl.pos:], b)
	impl.pos += 8
	return nil
}

func (impl *LEStream) WriteBuff(buff []byte) error {
	if impl.Left() < len(buff) {
		return ErrBuffOverflow
	}
	copy(impl.buff[impl.pos:], buff)
	impl.pos += len(buff)
	return nil
}
