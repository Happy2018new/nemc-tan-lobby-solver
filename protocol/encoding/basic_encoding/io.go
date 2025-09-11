package basic_encoding

// BasicIO represents a packet IO direction,
// but only used to read/write basic data type,
// which are number, float and double.
//
// Implementations of this interface are
// Reader and Writer.
//
// Reader reads data from the input stream into the pointers passed,
// whereas Writer writes the values the pointers point point to the
// output stream.
type BasicIO interface {
	Uint8(x *uint8)
	Int8(x *int8)
	Uint16(x *uint16)
	Int16(x *int16)
	Uint32(x *uint32)
	Int32(x *int32)
	Uint64(x *uint64)
	Int64(x *int64)

	Varint64(x *int64)
	Varuint64(x *uint64)
	Varint32(x *int32)
	Varuint32(x *uint32)
	Varint16(x *int16)
	Varuint16(x *uint16)

	Float32(x *float32)
	Float64(x *float64)
}
