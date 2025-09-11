package packet

import "github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"

type TanEnterRoomResponse struct {
	ErrorCode    int8
	PlayerIDList []uint32
	ItemIDs      []uint64
}

func (*TanEnterRoomResponse) ID() uint16 {
	return IDTanEnterRoomResponse
}

func (*TanEnterRoomResponse) BoundType() uint8 {
	return BoundTypeClient
}

func (t *TanEnterRoomResponse) Marshal(io encoding.IO) {
	io.Int8(&t.ErrorCode)
	encoding.FuncSliceUint8Length(io, &t.PlayerIDList, io.Uint32)
	encoding.FuncSliceUint8Length(io, &t.ItemIDs, io.Uint64)
}
