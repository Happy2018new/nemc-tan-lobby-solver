package encoding

import (
	"fmt"
	"image/color"
	"io"
	"unsafe"

	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding/basic_encoding"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/nbt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/uuid"
)

// Reader implements reading operations for
// reading types from Minecraft packets.
//
// Each Packet implementation
// has one passed to it.
//
// Reader's uses should always be encapsulated
// with a deferred recovery.
// Reader panics on invalid data.
type Reader struct {
	*basic_encoding.BasicReader
}

// NewReader creates a new Reader using the
// io.ByteReader passed as underlying source
// to read bytes from.
func NewReader(r interface {
	io.Reader
	io.ByteReader
}) *Reader {
	return &Reader{basic_encoding.NewBasicReader(r)}
}

// Bool reads a bool from the underlying buffer.
func (r *Reader) Bool(x *bool) {
	u, err := r.Reader().ReadByte()
	if err != nil {
		r.panic(err)
	}
	*x = *(*bool)(unsafe.Pointer(&u))
}

// String reads a string from the underlying buffer.
func (r *Reader) String(x *string) {
	var length int16
	r.Varint16(&length)
	l := int(length)
	data := make([]byte, l)
	if _, err := r.Reader().Read(data); err != nil {
		r.panic(err)
	}
	*x = *(*string)(unsafe.Pointer(&data))
}

// StringUint8 reads a string that length is uint8 prefix from the underlying buffer.
func (r *Reader) Uint8String(x *string) {
	var length uint8
	r.Uint8(&length)
	l := int(length)
	data := make([]byte, l)
	if _, err := r.Reader().Read(data); err != nil {
		r.panic(err)
	}
	*x = *(*string)(unsafe.Pointer(&data))
}

// StringUTF reads a string that length is uint16 prefix from the underlying buffer.
func (r *Reader) StringUTF(x *string) {
	var length uint16
	r.Uint16(&length)
	l := int(length)
	data := make([]byte, l)
	if _, err := r.Reader().Read(data); err != nil {
		r.panic(err)
	}
	*x = *(*string)(unsafe.Pointer(&data))
}

// Angle reads a rotational float32 from a single byte.
func (r *Reader) Angle(x *float32) {
	var v uint8
	r.Uint8(&v)
	*x = float32(v) * (360.0 / 256.0)
}

// UUID reads a uuid.UUID from the underlying buffer.
func (r *Reader) UUID(x *uuid.UUID) {
	b := make([]byte, 16)
	if _, err := r.Reader().Read(b); err != nil {
		r.panic(err)
	}
	*x = uuid.UUID(b)
}

// NBT reads a compound tag into a map from the underlying buffer.
func (r *Reader) NBT(m *map[string]any, encoding nbt.Encoding) {
	dec := nbt.NewDecoderWithEncoding(r.Reader(), encoding)
	dec.AllowZero = true

	*m = make(map[string]any)
	if err := dec.Decode(m); err != nil {
		r.panic(err)
	}
}

// NBTList reads a list of NBT tags from the underlying buffer.
func (r *Reader) NBTList(m *[]any, encoding nbt.Encoding) {
	if err := nbt.NewDecoderWithEncoding(r.Reader(), encoding).Decode(m); err != nil {
		r.panic(err)
	}
}

// NBTString reads a string tag into a string from the underlying buffer.
func (r *Reader) NBTString(s *string, encoding nbt.Encoding) {
	dec := nbt.NewDecoderWithEncoding(r.Reader(), encoding)
	dec.AllowZero = true

	*s = ""
	if err := dec.Decode(s); err != nil {
		r.panic(err)
	}
}

// Vec3 reads three float32s into an mgl32.Vec3 from the underlying buffer.
func (r *Reader) Vec3(x *mgl32.Vec3) {
	r.Float32(&x[0])
	r.Float32(&x[1])
	r.Float32(&x[2])
}

// Vec4 reads four float32s into an mgl32.Vec4 from the underlying buffer.
func (r *Reader) Vec4(x *mgl32.Vec4) {
	r.Float32(&x[0])
	r.Float32(&x[1])
	r.Float32(&x[2])
	r.Float32(&x[3])
}

// RGB reads a color.RGBA x from a 0xRRGGBB uint32.
func (r *Reader) RGB(x *color.RGBA) {
	var v uint32
	r.Uint32(&v)
	*x = color.RGBA{
		R: byte((v >> 16) & 0xff),
		G: byte((v >> 8) & 0xff),
		B: byte(v & 0xff),
		A: 255,
	}
}

// RGBA reads a color.RGBA x from a 0xAARRGGBB uint32.
func (r *Reader) RGBA(x *color.RGBA) {
	var v uint32
	r.Uint32(&v)
	*x = color.RGBA{
		A: byte((v >> 24) & 0xff),
		R: byte((v >> 16) & 0xff),
		G: byte((v >> 8) & 0xff),
		B: byte(v & 0xff),
	}
}

// Bytes reads the leftover bytes into a byte slice.
func (r *Reader) Bytes(p *[]byte) {
	var err error
	*p, err = io.ReadAll(r.Reader())
	if err != nil {
		r.panic(err)
	}
}

// UnknownEnumOption panics with an unknown enum option error.
func (r *Reader) UnknownEnumOption(value any, enum string) {
	r.panicf("unknown value '%v' for enum type '%v'", value, enum)
}

// InvalidValue panics with an error indicating that the value passed is not valid for a specific field.
func (r *Reader) InvalidValue(value any, forField, reason string) {
	r.panicf("invalid value '%v' for %v: %v", value, forField, reason)
}

// panicf panics with the format and values passed and assigns the error created to the Reader.
func (r *Reader) panicf(format string, a ...any) {
	panic(fmt.Errorf(format, a...))
}

// panic panics with the error passed, similarly to panicf.
func (r *Reader) panic(err error) {
	panic(err)
}
