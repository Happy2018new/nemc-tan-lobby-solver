package packet

import "github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"

const (
	TanCreateRoomSuccess              int8 = 0
	TanCreateRoomNeedVipToSetRoomName int8 = 7
)

// TanCreateRoomResponse ..
type TanCreateRoomResponse struct {
	ErrorCode int8
	RoomID    uint32
}

func (*TanCreateRoomResponse) ID() uint16 {
	return IDTanCreateRoomResponse
}

func (*TanCreateRoomResponse) BoundType() uint8 {
	return BoundTypeClient
}

func (t *TanCreateRoomResponse) Marshal(io encoding.IO) {
	io.Int8(&t.ErrorCode)
	if t.ErrorCode != TanCreateRoomSuccess {
		return
	}
	io.Uint32(&t.RoomID)
}
