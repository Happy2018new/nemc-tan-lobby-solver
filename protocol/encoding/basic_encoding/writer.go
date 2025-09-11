package basic_encoding

import (
	"encoding/binary"
	"io"
	"math"
)

// BasicWriter implements an IO that
// allows you write basic data to the
// underlying buffer.
//
// Basic data are as follows:
// number, float and double.
type BasicWriter struct {
	w interface {
		io.Writer
		io.ByteWriter
	}
}

// NewReader creates a new BasicReader using
// io.ByteWriter passed as underlying source
// to read bytes from.
func NewBasicWriter(w interface {
	io.Writer
	io.ByteWriter
}) *BasicWriter {
	return &BasicWriter{w: w}
}

// Writer return the underlying writer
// which is w.w.
func (w *BasicWriter) Writer() interface {
	io.Writer
	io.ByteWriter
} {
	return w.w
}

// Uint8 writes a uint8 to the underlying buffer.
func (w *BasicWriter) Uint8(x *uint8) {
	_ = w.w.WriteByte(*x)
}

// Int8 writes an int8 to the underlying buffer.
func (w *BasicWriter) Int8(x *int8) {
	_ = w.w.WriteByte(byte(*x) & 0xff)
}

// Uint16 writes a big endian uint16 to the underlying buffer.
func (w *BasicWriter) Uint16(x *uint16) {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, *x)
	_, _ = w.w.Write(data)
}

// Int16 writes a big endian int16 to the underlying buffer.
func (w *BasicWriter) Int16(x *int16) {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, uint16(*x))
	_, _ = w.w.Write(data)
}

// Uint32 writes a big endian uint32 to the underlying buffer.
func (w *BasicWriter) Uint32(x *uint32) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, *x)
	_, _ = w.w.Write(data)
}

// Int32 writes a big endian int32 to the underlying buffer.
func (w *BasicWriter) Int32(x *int32) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(*x))
	_, _ = w.w.Write(data)
}

// Uint64 writes a big endian uint64 to the underlying buffer.
func (w *BasicWriter) Uint64(x *uint64) {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, *x)
	_, _ = w.w.Write(data)
}

// Int64 writes a big endian int64 to the underlying buffer.
func (w *BasicWriter) Int64(x *int64) {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(*x))
	_, _ = w.w.Write(data)
}

// Float32 writes a big endian float32 to the underlying buffer.
func (w *BasicWriter) Float32(x *float32) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, math.Float32bits(*x))
	_, _ = w.w.Write(data)
}

// Float64 writes a big endian float64 to the underlying buffer.
func (w *BasicWriter) Float64(x *float64) {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, math.Float64bits(*x))
	_, _ = w.w.Write(data)
}

// Varint64 writes an int64 as 1-10 bytes to the underlying buffer.
func (w *BasicWriter) Varint64(x *int64) {
	u := *x
	ux := uint64(u) << 1
	if u < 0 {
		ux = ^ux
	}
	for ux >= 0x80 {
		_ = w.w.WriteByte(byte(ux) | 0x80)
		ux >>= 7
	}
	_ = w.w.WriteByte(byte(ux))
}

// Varuint64 writes a uint64 as 1-10 bytes to the underlying buffer.
func (w *BasicWriter) Varuint64(x *uint64) {
	u := *x
	for u >= 0x80 {
		_ = w.w.WriteByte(byte(u) | 0x80)
		u >>= 7
	}
	_ = w.w.WriteByte(byte(u))
}

// Varint32 writes an int32 as 1-5 bytes to the underlying buffer.
func (w *BasicWriter) Varint32(x *int32) {
	u := *x
	ux := uint32(u) << 1
	if u < 0 {
		ux = ^ux
	}
	for ux >= 0x80 {
		_ = w.w.WriteByte(byte(ux) | 0x80)
		ux >>= 7
	}
	_ = w.w.WriteByte(byte(ux))
}

// Varuint32 writes a uint32 as 1-5 bytes to the underlying buffer.
func (w *BasicWriter) Varuint32(x *uint32) {
	u := *x
	for u >= 0x80 {
		_ = w.w.WriteByte(byte(u) | 0x80)
		u >>= 7
	}
	_ = w.w.WriteByte(byte(u))
}

// Varint16 writes an int16 as 1-3 bytes to the underlying buffer.
func (w *BasicWriter) Varint16(x *int16) {
	u := *x
	ux := uint16(u) << 1
	if u < 0 {
		ux = ^ux
	}
	for ux >= 0x80 {
		_ = w.w.WriteByte(byte(ux) | 0x80)
		ux >>= 7
	}
	_ = w.w.WriteByte(byte(ux))
}

// Varuint16 writes a uint16 as 1-3 bytes to the underlying buffer.
func (w *BasicWriter) Varuint16(x *uint16) {
	u := *x
	for u >= 0x80 {
		_ = w.w.WriteByte(byte(u) | 0x80)
		u >>= 7
	}
	_ = w.w.WriteByte(byte(u))
}
