package basic_encoding

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

// errVarIntOverflow is an error set if one of the Varint methods encounters a varint that does not terminate
// after 5 or 10 bytes, depending on the data type read into.
var errVarIntOverflow = errors.New("varint overflows integer")

// BasicReader implements an IO that
// allows you read basic data to the
// underlying buffer.
//
// Basic data are as follows:
// number, float and double.
type BasicReader struct {
	r interface {
		io.Reader
		io.ByteReader
	}
}

// NewReader creates a new BasicReader using
// io.ByteReader passed as underlying source
// to read bytes from.
func NewBasicReader(r interface {
	io.Reader
	io.ByteReader
}) *BasicReader {
	return &BasicReader{r: r}
}

// Reader return the underlying reader
// which is r.r.
func (r *BasicReader) Reader() interface {
	io.Reader
	io.ByteReader
} {
	return r.r
}

// Uint8 reads a uint8 from the underlying buffer.
func (r *BasicReader) Uint8(x *uint8) {
	var err error
	*x, err = r.r.ReadByte()
	if err != nil {
		panic(fmt.Sprintf("(r *BasicReader) Uint8: %v", err))
	}
}

// Int8 reads an int8 from the underlying buffer.
func (r *BasicReader) Int8(x *int8) {
	var b uint8
	r.Uint8(&b)
	*x = int8(b)
}

// Uint16 reads a big endian uint16 from the underlying buffer.
func (r *BasicReader) Uint16(x *uint16) {
	b := make([]byte, 2)
	if _, err := r.r.Read(b); err != nil {
		panic(fmt.Sprintf("(r *BasicReader) Uint16: %v", err))
	}
	*x = binary.BigEndian.Uint16(b)
}

// Int16 reads a big endian int16 from the underlying buffer.
func (r *BasicReader) Int16(x *int16) {
	b := make([]byte, 2)
	if _, err := r.r.Read(b); err != nil {
		panic(fmt.Sprintf("(r *BasicReader) Int16: %v", err))
	}
	*x = int16(binary.BigEndian.Uint16(b))
}

// Uint32 reads a big endian uint32 from the underlying buffer.
func (r *BasicReader) Uint32(x *uint32) {
	b := make([]byte, 4)
	if _, err := r.r.Read(b); err != nil {
		panic(fmt.Sprintf("(r *BasicReader) Uint32: %v", err))
	}
	*x = binary.BigEndian.Uint32(b)
}

// Int32 reads a big endian int32 from the underlying buffer.
func (r *BasicReader) Int32(x *int32) {
	b := make([]byte, 4)
	if _, err := r.r.Read(b); err != nil {
		panic(fmt.Sprintf("(r *BasicReader) Int32: %v", err))
	}
	*x = int32(binary.BigEndian.Uint32(b))
}

// Uint64 reads a big endian uint64 from the underlying buffer.
func (r *BasicReader) Uint64(x *uint64) {
	b := make([]byte, 8)
	if _, err := r.r.Read(b); err != nil {
		panic(fmt.Sprintf("(r *BasicReader) Uint64: %v", err))
	}
	*x = binary.BigEndian.Uint64(b)
}

// Int64 reads a big endian int64 from the underlying buffer.
func (r *BasicReader) Int64(x *int64) {
	b := make([]byte, 8)
	if _, err := r.r.Read(b); err != nil {
		panic(fmt.Sprintf("(r *BasicReader) Int64: %v", err))
	}
	*x = int64(binary.BigEndian.Uint64(b))
}

// Float32 reads a big endian float32 from the underlying buffer.
func (r *BasicReader) Float32(x *float32) {
	b := make([]byte, 4)
	if _, err := r.r.Read(b); err != nil {
		panic(fmt.Sprintf("(r *BasicReader) Float32: %v", err))
	}
	*x = math.Float32frombits(binary.BigEndian.Uint32(b))
}

// Float64 reads a big endian float64 from the underlying buffer.
func (r *BasicReader) Float64(x *float64) {
	b := make([]byte, 8)
	if _, err := r.r.Read(b); err != nil {
		panic(fmt.Sprintf("(r *BasicReader) Float64: %v", err))
	}
	*x = math.Float64frombits(binary.BigEndian.Uint64(b))
}

// Varint64 reads up to 10 bytes from the underlying buffer into an int64.
func (r *BasicReader) Varint64(x *int64) {
	var ux uint64
	for i := 0; i < 70; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil {
			panic(fmt.Sprintf("(r *BasicReader) Varint64: %v", err))
		}

		ux |= uint64(b&0x7f) << i
		if b&0x80 == 0 {
			*x = int64(ux >> 1)
			if ux&1 != 0 {
				*x = ^*x
			}
			return
		}
	}
	panic(fmt.Sprintf("(r *BasicReader) Varint64: %v", errVarIntOverflow))
}

// Varuint64 reads up to 10 bytes from the underlying buffer into a uint64.
func (r *BasicReader) Varuint64(x *uint64) {
	var v uint64
	for i := 0; i < 70; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil {
			panic(fmt.Sprintf("(r *BasicReader) Varuint64: %v", err))
		}

		v |= uint64(b&0x7f) << i
		if b&0x80 == 0 {
			*x = v
			return
		}
	}
	panic(fmt.Sprintf("(r *BasicReader) Varuint64: %v", errVarIntOverflow))
}

// Varint32 reads up to 5 bytes from the underlying buffer into an int32.
func (r *BasicReader) Varint32(x *int32) {
	var ux uint32
	for i := 0; i < 35; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil {
			panic(fmt.Sprintf("(r *BasicReader) Varint32: %v", err))
		}

		ux |= uint32(b&0x7f) << i
		if b&0x80 == 0 {
			*x = int32(ux >> 1)
			if ux&1 != 0 {
				*x = ^*x
			}
			return
		}
	}
	panic(fmt.Sprintf("(r *BasicReader) Varint32: %v", errVarIntOverflow))
}

// Varuint32 reads up to 5 bytes from the underlying buffer into a uint32.
func (r *BasicReader) Varuint32(x *uint32) {
	var v uint32
	for i := 0; i < 35; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil {
			panic(fmt.Sprintf("(r *BasicReader) Varint32: %v", err))
		}

		v |= uint32(b&0x7f) << i
		if b&0x80 == 0 {
			*x = v
			return
		}
	}
	panic(fmt.Sprintf("(r *BasicReader) Varuint32: %v", errVarIntOverflow))
}

// Varint16 reads up to 3 bytes from the underlying buffer into an int16.
func (r *BasicReader) Varint16(x *int16) {
	var ux uint16
	for i := 0; i < 21; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil {
			panic(fmt.Sprintf("(r *BasicReader) Varint16: %v", err))
		}

		ux |= uint16(b&0x7f) << i
		if b&0x80 == 0 {
			*x = int16(ux >> 1)
			if ux&1 != 0 {
				*x = ^*x
			}
			return
		}
	}
	panic(fmt.Sprintf("(r *BasicReader) Varint16: %v", errVarIntOverflow))
}

// Varuint16 reads up to 3 bytes from the underlying buffer into a uint16.
func (r *BasicReader) Varuint16(x *uint16) {
	var v uint16
	for i := 0; i < 21; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil {
			panic(fmt.Sprintf("(r *BasicReader) Varint32: %v", err))
		}

		v |= uint16(b&0x7f) << i
		if b&0x80 == 0 {
			*x = v
			return
		}
	}
	panic(fmt.Sprintf("(r *BasicReader) Varuint32: %v", errVarIntOverflow))
}
