package encoding

import (
	"image/color"

	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/nbt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/uuid"
)

// IO represents a packet IO direction that based on BasicIO.
//
// Implementations of this interface are
// Reader and Writer.
//
// Reader reads data from the input stream into the pointers
// passed, whereas Writer writes the values the pointers
// point point to the output stream.
type IO interface {
	Uint32(x *uint32)
	Uint64(x *uint64)
	Varuint64(x *uint64)
	Varuint32(x *uint32)
	Varuint16(x *uint16)
	Varint16(x *int16)

	Bool(x *bool)
	Int8(x *int8)
	Uint8(x *uint8)
	Int16(x *int16)
	Uint16(x *uint16)
	Int32(x *int32)
	Int64(x *int64)
	Float32(x *float32)
	Float64(x *float64)
	String(x *string)
	Uint8String(x *string)
	StringUTF(x *string)
	Varint32(x *int32)
	Varint64(x *int64)
	Angle(x *float32)
	UUID(x *uuid.UUID)

	NBT(m *map[string]any, encoding nbt.Encoding)
	NBTList(m *[]any, encoding nbt.Encoding)
	NBTString(s *string, encoding nbt.Encoding)

	Vec3(x *mgl32.Vec3)
	Vec4(x *mgl32.Vec4)
	RGB(x *color.RGBA)
	RGBA(x *color.RGBA)
	RoomTips(x *RoomTips)

	Bytes(x *[]byte)
	UnknownEnumOption(value any, enum string)
	InvalidValue(value any, forField, reason string)
}

// Marshaler is a type that can be written to or read from an IO.
type Marshaler interface {
	Marshal(r IO)
}

// Slice reads/writes a slice of T with a varuint32 prefix.
func Slice[T any, S ~*[]T, A PtrMarshaler[T]](r IO, x S) {
	count := uint32(len(*x))
	r.Varuint32(&count)
	SliceOfLen[T, S, A](r, count, x)
}

// SliceUint8Length reads/writes a slice of T with a uint8 prefix.
func SliceUint8Length[T any, S *[]T, A PtrMarshaler[T]](r IO, x S) {
	count := uint8(len(*x))
	r.Uint8(&count)
	SliceOfLen[T, S, A](r, uint32(count), x)
}

// SliceUint16Length reads/writes a slice of T with a uint16 prefix.
func SliceUint16Length[T any, S ~*[]T, A PtrMarshaler[T]](r IO, x S) {
	count := uint16(len(*x))
	r.Uint16(&count)
	SliceOfLen[T, S, A](r, uint32(count), x)
}

// SliceUint32Length reads/writes a slice of T with a uint32 prefix.
func SliceUint32Length[T any, S ~*[]T, A PtrMarshaler[T]](r IO, x S) {
	count := uint32(len(*x))
	r.Uint32(&count)
	SliceOfLen[T, S, A](r, count, x)
}

// SliceVarint32Length reads/writes a slice of T with a varint32 prefix.
func SliceVarint32Length[T any, S ~*[]T, A PtrMarshaler[T]](r IO, x S) {
	count := int32(len(*x))
	r.Varint32(&count)
	SliceOfLen[T, S, A](r, uint32(count), x)
}

// FuncSliceUint8Length reads/writes a slice of T using function f with a uint8 length prefix.
func FuncSliceUint8Length[T any, S ~*[]T](r IO, x S, f func(*T)) {
	count := uint8(len(*x))
	r.Uint8(&count)
	FuncSliceOfLen(r, uint32(count), x, f)
}

// FuncSliceUint16Length reads/writes a slice of T using function f with a uint16 length prefix.
func FuncSliceUint16Length[T any, S ~*[]T](r IO, x S, f func(*T)) {
	count := uint16(len(*x))
	r.Uint16(&count)
	FuncSliceOfLen(r, uint32(count), x, f)
}

// FuncSliceUint32Length reads/writes a slice of T using function f with a uint32 length prefix.
func FuncSliceUint32Length[T any, S ~*[]T](r IO, x S, f func(*T)) {
	count := uint32(len(*x))
	r.Uint32(&count)
	FuncSliceOfLen(r, count, x, f)
}

// FuncSliceVarint32Length reads/writes a slice of T with a varint32 prefix.
func FuncSliceVarint32Length[T any, S ~*[]T](r IO, x S, f func(*T)) {
	count := int32(len(*x))
	r.Varint32(&count)
	FuncSliceOfLen(r, uint32(count), x, f)
}

// FuncSlice reads/writes a slice of T using function f with a varuint32 length prefix.
func FuncSlice[T any, S ~*[]T](r IO, x S, f func(*T)) {
	count := uint32(len(*x))
	r.Varuint32(&count)
	FuncSliceOfLen(r, count, x, f)
}

// FuncIOSlice reads/writes a slice of T using a function f with a varuint32 length prefix.
func FuncIOSlice[T any, S ~*[]T](r IO, x S, f func(IO, *T)) {
	FuncSlice(r, x, func(v *T) {
		f(r, v)
	})
}

// FuncIOSliceUint32Length reads/writes a slice of T using a function with a uint32 length prefix.
func FuncIOSliceUint32Length[T any, S ~*[]T](r IO, x S, f func(IO, *T)) {
	count := uint32(len(*x))
	r.Uint32(&count)
	FuncIOSliceOfLen(r, count, x, f)
}

// SliceOfLen reads/writes the elements of a slice of type T with length l.
func SliceOfLen[T any, S ~*[]T, A PtrMarshaler[T]](r IO, l uint32, x S) {
	if _, ok := r.(*Reader); ok {
		*x = make([]T, l)
	}
	for i := range l {
		A(&(*x)[i]).Marshal(r)
	}
}

// FuncSliceOfLen reads/writes the elements of a slice of type T with length l using func f.
func FuncSliceOfLen[T any, S ~*[]T](r IO, l uint32, x S, f func(*T)) {
	if _, ok := r.(*Reader); ok {
		*x = make([]T, l)
	}
	for i := range l {
		f(&(*x)[i])
	}
}

// FuncIOSliceOfLen reads/writes the elements of a slice of type T with length l using func f.
func FuncIOSliceOfLen[T any, S ~*[]T](r IO, l uint32, x S, f func(IO, *T)) {
	FuncSliceOfLen(r, l, x, func(v *T) {
		f(r, v)
	})
}

// PtrMarshaler represents a type that implements Marshaler for its pointer.
type PtrMarshaler[T any] interface {
	Marshaler
	*T
}

// Single reads/writes a single Marshaler x.
func Single[T any, S PtrMarshaler[T]](r IO, x S) {
	x.Marshal(r)
}

// Optional is an optional type in the protocol. If not set, only a false bool is written. If set, a true bool is
// written and the Marshaler.
type Optional[T any] struct {
	set bool
	val T
}

// Option creates an Optional[T] with the value passed.
func Option[T any](val T) Optional[T] {
	return Optional[T]{set: true, val: val}
}

// Value returns the value set in the Optional. If no value was set, false is returned.
func (o Optional[T]) Value() (T, bool) {
	return o.val, o.set
}

// OptionalFunc reads/writes an Optional[T].
func OptionalFunc[T any](r IO, x *Optional[T], f func(*T)) any {
	r.Bool(&x.set)
	if x.set {
		f(&x.val)
	}
	return x
}

// OptionalFuncIO reads/writes an Optional[T].
func OptionalFuncIO[T any](r IO, x *Optional[T], f func(IO, *T)) any {
	r.Bool(&x.set)
	if x.set {
		f(r, &x.val)
	}
	return x
}

// OptionalSliceFunc reads/writes an Optional[S].
// Note that:
//   - S must be a slice that satisfy []T.
//   - T must have implements Marshaler.
func OptionalSliceMarshaler[T any, S ~[]T, A PtrMarshaler[T]](r IO, x *Optional[S]) any {
	r.Bool(&x.set)
	if x.set {
		s := []T(x.val)
		SliceVarint32Length[T, *[]T, A](r, &s)
		x.val = s
	}
	return x
}

// OptionalSlice reads/writes an Optional[S].
// Note that:
//   - S must be a slice that satisfy []T.
//   - f is used to read/write T.
func OptionalSlice[T any, S ~[]T](r IO, x *Optional[S], f func(*T)) any {
	r.Bool(&x.set)
	if x.set {
		s := []T(x.val)
		FuncSliceVarint32Length(r, &s, f)
		x.val = s
	}
	return x
}

// OptionalMarshaler reads/writes an Optional assuming *T implements Marshaler.
func OptionalMarshaler[T any, A PtrMarshaler[T]](r IO, x *Optional[T]) {
	r.Bool(&x.set)
	if x.set {
		A(&x.val).Marshal(r)
	}
}

// OptionalPointerMarshaler reads/writes an Optional assuming *T implements Marshaler.
func OptionalPointerMarshaler[T any, A PtrMarshaler[T]](r IO, x *Optional[*T]) {
	r.Bool(&x.set)
	if x.set {
		if x.val == nil {
			x.val = new(T)
		}
		A(x.val).Marshal(r)
	}
}

// IDOrX is `ID or X` data type in the protocol.
type IDOrX[T any] struct {
	// 0 if value of type X is given inline;
	// otherwise registry ID + 1.
	ID int32
	// Only present if ID is 0.
	Value T
}

// IDOrXFunc reads/writes an Optional[T].
func IDOrXFunc[T any](r IO, x *IDOrX[T], f func(*T)) {
	r.Varint32(&x.ID)
	if x.ID == 0 {
		f(&x.Value)
	}
}

// IDOrXMarshaler reads/writes an Optional assuming *T implements Marshaler.
func IDOrXMarshaler[T any, A PtrMarshaler[T]](r IO, x *IDOrX[T]) {
	r.Varint32(&x.ID)
	if x.ID == 0 {
		A(&x.Value).Marshal(r)
	}
}
