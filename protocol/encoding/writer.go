package encoding

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"unsafe"

	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding/basic_encoding"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/nbt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/uuid"
)

// Writer implements writing methods for data types
// from Minecraft packets.
//
// Each Packet implementation has one passed to it
// when writing.
//
// Writer implements methods where values are passed
// using a pointer, so that Reader and Writer have a
// synonymous interface and both implement the IO
// interface.
type Writer struct {
	*basic_encoding.BasicWriter
}

// NewWriter creates a new initialised Writer with an
// underlying io.ByteWriter to write to.
func NewWriter(w interface {
	io.Writer
	io.ByteWriter
}, shieldID int32) *Writer {
	return &Writer{basic_encoding.NewBasicWriter(w)}
}

// Bool writes a bool as either 0 or 1 to the underlying buffer.
func (w *Writer) Bool(x *bool) {
	_ = w.Writer().WriteByte(*(*byte)(unsafe.Pointer(x)))
}

// String writes a string, prefixed with a varint16, to the underlying buffer.
func (w *Writer) String(x *string) {
	l := int16(len(*x))
	w.Varint16(&l)
	_, _ = w.Writer().Write([]byte(*x))
}

// Uint8String writes a string, prefixed with a uint8 prefix to the underlying buffer.
func (w *Writer) Uint8String(x *string) {
	l := uint8(len(*x))
	w.Uint8(&l)
	_, _ = w.Writer().Write([]byte(*x))
}

// StringUTF writes a string, prefixed with uint16 prefix to the underlying buffer.
func (w *Writer) StringUTF(x *string) {
	l := uint16(len(*x))
	w.Uint16(&l)
	_, _ = w.Writer().Write([]byte(*x))
}

// Angle writes a rotational float32 as a single byte to the underlying buffer.
func (w *Writer) Angle(x *float32) {
	_ = w.Writer().WriteByte(byte(*x / (360.0 / 256.0)))
}

// UUID writes a UUID to the underlying buffer.
func (w *Writer) UUID(x *uuid.UUID) {
	b := [16]byte(*x)
	_, _ = w.Writer().Write(b[:])
}

// NBT writes a map as NBT to the underlying buffer using the encoding passed.
func (w *Writer) NBT(x *map[string]any, encoding nbt.Encoding) {
	if err := nbt.NewEncoderWithEncoding(w.Writer(), encoding).Encode(*x); err != nil {
		panic(err)
	}
}

// NBTList writes a slice as NBT to the underlying buffer using the encoding passed.
func (w *Writer) NBTList(x *[]any, encoding nbt.Encoding) {
	if err := nbt.NewEncoderWithEncoding(w.Writer(), encoding).Encode(*x); err != nil {
		panic(err)
	}
}

// NBTString writes a string as NBT to the underlying buffer using the encoding passed.
func (w *Writer) NBTString(x *string, encoding nbt.Encoding) {
	if err := nbt.NewEncoderWithEncoding(w.Writer(), encoding).Encode(*x); err != nil {
		panic(err)
	}
}

// Vec4 writes an mgl32.Vec4 as 4 float32s to the underlying buffer.
func (w *Writer) Vec4(x *mgl32.Vec4) {
	w.Float32(&x[0])
	w.Float32(&x[1])
	w.Float32(&x[2])
	w.Float32(&x[3])
}

// Vec3 writes an mgl32.Vec3 as 3 float32s to the underlying buffer.
func (w *Writer) Vec3(x *mgl32.Vec3) {
	w.Float32(&x[0])
	w.Float32(&x[1])
	w.Float32(&x[2])
}

// RGB writes a color.RGBA x as a uint32 0xRRGGBB the underlying buffer.
func (w *Writer) RGB(x *color.RGBA) {
	val := uint32(x.R)<<16 | uint32(x.G)<<8 | uint32(x.B)
	w.Uint32(&val)
}

// RGBA writes a color.RGBA x as a uint32 0xAARRGGBB to the underlying buffer.
func (w *Writer) RGBA(x *color.RGBA) {
	val := uint32(x.A)<<24 | uint32(x.R)<<16 | uint32(x.G)<<8 | uint32(x.B)
	w.Uint32(&val)
}

// RoomTips writes a room tips x to the underlying buffer.
func (w *Writer) RoomTips(x *RoomTips) {
	var bufBytes []byte
	var length uint16

	buf := bytes.NewBuffer(nil)
	writer := NewWriter(buf, 0)
	writer.StringUTF(&x.LevelID)
	writer.Uint8(&x.GameType)
	writer.StringUTF(&x.ConstantTestString)
	writer.Int16(&x.Vioce)
	writer.Uint8(&x.ProtocolID)

	bufBytes, length = buf.Bytes(), uint16(buf.Len())
	w.Uint16(&length)
	w.Bytes(&bufBytes)
}

// Bytes appends a []byte to the underlying buffer.
func (w *Writer) Bytes(x *[]byte) {
	_, _ = w.Writer().Write(*x)
}

// UnknownEnumOption panics with an unknown enum option error.
func (w *Writer) UnknownEnumOption(value any, enum string) {
	w.panicf("unknown value '%v' for enum type '%v'", value, enum)
}

// InvalidValue panics with an invalid value error.
func (w *Writer) InvalidValue(value any, forField, reason string) {
	w.panicf("invalid value '%v' for %v: %v", value, forField, reason)
}

// panicf panics with the format and values passed.
func (w *Writer) panicf(format string, a ...any) {
	panic(fmt.Errorf(format, a...))
}
