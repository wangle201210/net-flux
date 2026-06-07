package stream

import (
	"encoding/binary"
	"errors"
)

// ReadStream 定义了一个读取数据流的接口，提供了多种读取数据的方法
type ReadStream interface {
	// Size 返回数据流的总大小
	Size() int
	// Length 返回数据流的长度
	Length() int
	// Left 返回数据流剩余未读取的数据长度
	Left() int
	// Reset 重置数据流，使用提供的字节切片作为新的数据源
	Reset([]byte)
	// Data 返回数据流中的全部数据
	Data() []byte
	// ReadByte 从数据流中读取一个字节
	ReadByte() (b byte, err error)
	// ReadUint16 从数据流中读取一个uint16类型的值
	ReadUint16() (b uint16, err error)
	// ReadUint32 从数据流中读取一个uint32类型的值
	ReadUint32() (b uint32, err error)
	// ReadUint64 从数据流中读取一个uint64类型的值
	ReadUint64() (b uint64, err error)
	// ReadBuff 从数据流中读取指定大小的字节切片
	ReadBuff(size int) (b []byte, err error)
	// ReadString 从数据流中读取一个字符串
	ReadString() (b string, err error)
	// CopyBuff 将数据流中的数据复制到提供的字节切片中
	CopyBuff(b []byte) error
}

// WriteStream 定义了一个写入流的接口，提供了一系列写入数据的方法
type WriteStream interface {
	// Size 返回写入流的总大小（容量）
	Size() int
	// Length 返回写入流中已使用的长度
	Length() int
	// Left 返回写入流剩余可用的空间大小
	Left() int
	// Reset 重置写入流，使用提供的字节切片作为新的缓冲区
	Reset([]byte)
	// Data 返回写入流中当前的数据
	Data() []byte
	// WriteByte 将一个字节写入流中
	WriteByte(b byte) error
	// WriteUint16 将一个16位无符号整数写入流中
	WriteUint16(b uint16) error
	// WriteUint32 将一个32位无符号整数写入流中
	WriteUint32(b uint32) error
	// WriteUint64 将一个64位无符号整数写入流中
	WriteUint64(b uint64) error
	// WriteBuff 将一个字节切片写入流中
	WriteBuff(b []byte) error
	// WriteString 将一个字符串写入流中
	WriteString(b string) error
}

type BigEndianStreamImpl struct {
	pos  int
	buff []byte
}

type LittleEndianStreamImpl struct {
	pos  int
	buff []byte
}

// NewBigEndianStream 创建一个新的大端字节流实现
// 该函数接收一个字节数组作为缓冲区，并返回一个指向BigEndianStreamImpl结构体的指针
// 参数:
//
//	buff: 用于初始化字节流的字节数组
//
// 返回值:
//
//	*BigEndianStreamImpl: 指向新创建的大端字节流实现结构体的指针
func NewBigEndianStream(buff []byte) *BigEndianStreamImpl {
	return &BigEndianStreamImpl{
		buff: buff,
	}
}

// NewLittleEndianStream 创建一个新的小端字节流实现
// 参数:
//
//	buff: 字节数组，用于初始化字节流
//
// 返回值:
//
//	*LittleEndianStreamImpl: 返回一个小端字节流实现的指针
func NewLittleEndianStream(buff []byte) *LittleEndianStreamImpl {
	return &LittleEndianStreamImpl{ // 创建并返回一个新的 LittleEndianStreamImpl 实例
		buff: buff, // 使用传入的字节数组初始化 buff 字段
	}
}

func (impl *BigEndianStreamImpl) Size() int { return len(impl.buff) }

func (impl *BigEndianStreamImpl) Length() int { return impl.pos }

func (impl *BigEndianStreamImpl) Data() []byte { return impl.buff }

func (impl *BigEndianStreamImpl) Left() int { return len(impl.buff) - impl.pos }

func (impl *BigEndianStreamImpl) Reset(buff []byte) { impl.pos = 0; impl.buff = buff }

func (impl *BigEndianStreamImpl) ReadByte() (b byte, err error) {
	if impl.Left() < 1 {
		return 0, errors.New("buff is too small to io")
	}
	b = impl.buff[impl.pos]
	impl.pos += 1
	return b, nil
}

func (impl *BigEndianStreamImpl) ReadUint16() (b uint16, err error) {
	if impl.Left() < 2 {
		return 0, errors.New("buff is too small to io")
	}
	b = binary.BigEndian.Uint16(impl.buff[impl.pos:])
	impl.pos += 2
	return b, nil
}

func (impl *BigEndianStreamImpl) ReadUint32() (b uint32, err error) {
	if impl.Left() < 4 {
		return 0, errors.New("buff is too small to io")
	}
	b = binary.BigEndian.Uint32(impl.buff[impl.pos:])
	impl.pos += 4
	return b, nil
}

func (impl *BigEndianStreamImpl) ReadUint64() (b uint64, err error) {
	if impl.Left() < 8 {
		return 0, errors.New("buff is too small to io")
	}
	b = binary.BigEndian.Uint64(impl.buff[impl.pos:])
	impl.pos += 8
	return b, nil
}

func (impl *BigEndianStreamImpl) ReadBuff(size int) (buff []byte, err error) {
	if impl.Left() < size {
		return nil, errors.New("buff is too small to io")
	}
	buff = make([]byte, size, size)
	copy(buff, impl.buff[impl.pos:impl.pos+size])
	impl.pos += size
	return buff, nil
}

func (impl *BigEndianStreamImpl) Read(p []byte) (int, error) {
	l := len(p)
	buf, err := impl.ReadBuff(l)
	if err == nil {
		copy(p, buf)
		return len(buf), err
	}
	return -1, err
}

func (impl *BigEndianStreamImpl) ReadString() (string, error) {
	len, err := impl.ReadUint16()
	if err != nil {
		return "", err
	}
	buff, err := impl.ReadBuff(int(len))
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

func (impl *BigEndianStreamImpl) CopyBuff(b []byte) error {
	if impl.Left() < len(b) {
		return errors.New("buff is too small to io")
	}
	copy(b, impl.buff[impl.pos:impl.pos+len(b)])
	return nil
}

func (impl *BigEndianStreamImpl) WriteByte(b byte) error {
	if impl.Left() < 1 {
		return errors.New("buff is too small to io")
	}
	impl.buff[impl.pos] = b
	impl.pos += 1
	return nil
}

func (impl *BigEndianStreamImpl) WriteUint16(b uint16) error {
	if impl.Left() < 2 {
		return errors.New("buff is too small to io")
	}
	binary.BigEndian.PutUint16(impl.buff[impl.pos:], b)
	impl.pos += 2
	return nil
}

func (impl *BigEndianStreamImpl) WriteUint32(b uint32) error {
	if impl.Left() < 4 {
		return errors.New("buff is too small to io")
	}
	binary.BigEndian.PutUint32(impl.buff[impl.pos:], b)
	impl.pos += 4
	return nil
}

func (impl *BigEndianStreamImpl) WriteUint64(b uint64) error {
	if impl.Left() < 8 {
		return errors.New("buff is too small to io")
	}
	binary.BigEndian.PutUint64(impl.buff[impl.pos:], b)
	impl.pos += 8
	return nil
}

func (impl *BigEndianStreamImpl) WriteBuff(buff []byte) error {
	left := impl.Left()
	length := len(buff)
	if left < length {
		return errors.New("buff is too small to io")
	}
	if buff == nil || impl.Left() < len(buff) {
		return errors.New("buff is too small to io")
	}
	copy(impl.buff[impl.pos:], buff)
	impl.pos += len(buff)
	return nil
}

func (impl *BigEndianStreamImpl) Write(b []byte) (n int, err error) {
	err = impl.WriteBuff(b)
	return len(b), err
}

func (impl *BigEndianStreamImpl) WriteString(b string) error {
	len := len(b)
	if err := impl.WriteUint16(uint16(len)); err != nil {
		return err
	}
	if err := impl.WriteBuff([]byte(b)); err != nil {
		return err
	}
	return nil
}

func (impl *LittleEndianStreamImpl) Size() int { return len(impl.buff) }

func (impl *LittleEndianStreamImpl) Data() []byte { return impl.buff }

func (impl *LittleEndianStreamImpl) Left() int { return len(impl.buff) - impl.pos }

func (impl *LittleEndianStreamImpl) Reset(buff []byte) { impl.pos = 0; impl.buff = buff }

func (impl *LittleEndianStreamImpl) ReadByte() (b byte, err error) {
	if impl.Left() < 1 {
		return 0, errors.New("buff is too small to io")
	}
	b = impl.buff[impl.pos]
	impl.pos += 1
	return b, nil
}

func (impl *LittleEndianStreamImpl) ReadUint16() (b uint16, err error) {
	if impl.Left() < 2 {
		return 0, errors.New("buff is too small to io")
	}
	b = binary.LittleEndian.Uint16(impl.buff[impl.pos:])
	impl.pos += 2
	return b, nil
}

func (impl *LittleEndianStreamImpl) ReadUint32() (b uint32, err error) {
	if impl.Left() < 4 {
		return 0, errors.New("buff is too small to io")
	}
	b = binary.LittleEndian.Uint32(impl.buff[impl.pos:])
	impl.pos += 4
	return b, nil
}

func (impl *LittleEndianStreamImpl) ReadUint64() (b uint64, err error) {
	if impl.Left() < 8 {
		return 0, errors.New("buff is too small to io")
	}
	b = binary.LittleEndian.Uint64(impl.buff[impl.pos:])
	impl.pos += 8
	return b, nil
}

func (impl *LittleEndianStreamImpl) ReadBuff(size int) (buff []byte, err error) {
	if impl.Left() < size {
		return nil, errors.New("buff is too small to io")
	}
	buff = make([]byte, size, size)
	copy(buff, impl.buff[impl.pos:impl.pos+size])
	impl.pos += size
	return buff, nil
}

func (impl *LittleEndianStreamImpl) ReadString() (string, error) {
	len, err := impl.ReadUint16()
	if err != nil {
		return "", err
	}
	buff, err := impl.ReadBuff(int(len))
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

func (impl *LittleEndianStreamImpl) CopyBuff(b []byte) error {
	if impl.Left() < len(b) {
		return errors.New("buff is too small to io")
	}
	copy(b, impl.buff[impl.pos:impl.pos+len(b)])
	return nil
}

func (impl *LittleEndianStreamImpl) WriteByte(b byte) error {
	if impl.Left() < 1 {
		return errors.New("buff is too small to io")
	}
	impl.buff[impl.pos] = b
	impl.pos += 1
	return nil
}

func (impl *LittleEndianStreamImpl) WriteUint16(b uint16) error {
	if impl.Left() < 2 {
		return errors.New("buff is too small to io")
	}
	binary.LittleEndian.PutUint16(impl.buff[impl.pos:], b)
	impl.pos += 2
	return nil
}

func (impl *LittleEndianStreamImpl) WriteUint32(b uint32) error {
	if impl.Left() < 4 {
		return errors.New("buff is too small to io")
	}
	binary.LittleEndian.PutUint32(impl.buff[impl.pos:], b)
	impl.pos += 4
	return nil
}

func (impl *LittleEndianStreamImpl) WriteUint64(b uint64) error {
	if impl.Left() < 8 {
		return errors.New("buff is too small to io")
	}
	binary.LittleEndian.PutUint64(impl.buff[impl.pos:], b)
	impl.pos += 8
	return nil
}

func (impl *LittleEndianStreamImpl) WriteBuff(buff []byte) error {
	if impl.Left() < len(buff) {
		return errors.New("buff is too small to io")
	}
	copy(impl.buff[impl.pos:], buff)
	impl.pos += len(buff)
	return nil
}

func (impl *LittleEndianStreamImpl) WriteString(b string) error {
	len := len(b)
	if err := impl.WriteUint16(uint16(len)); err != nil {
		return err
	}
	if err := impl.WriteBuff([]byte(b)); err != nil {
		return err
	}
	return nil
}
