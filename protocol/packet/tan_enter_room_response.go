package packet

import "github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"

const (
	TanEnterRoomSuccess       int8 = 0
	TanEnterRoomNotFound      int8 = 1
	TanEnterRoomFullOfPeople  int8 = 4
	TanEnterRoomWrongPasscode int8 = 14
)

// TanEnterRoomResponse ..
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
	if t.ErrorCode != TanEnterRoomSuccess {
		return
	}
	encoding.FuncSliceUint8Length(io, &t.PlayerIDList, io.Uint32)
	encoding.FuncSliceUint8Length(io, &t.ItemIDs, io.Uint64)
}
