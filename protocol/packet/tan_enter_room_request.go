package packet

import "github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"

// TanEnterRoomRequest ..
type TanEnterRoomRequest struct {
	OwnerID               uint32
	RoomID                uint32
	EnterPassword         string
	EnterTeamID           uint64
	EnterToken            uint32
	FollowTeamID          uint64
	NetherNetID           string
	SupportWebRTCCompress bool
}

func (*TanEnterRoomRequest) ID() uint16 {
	return IDTanEnterRoomRequest
}

func (*TanEnterRoomRequest) BoundType() uint8 {
	return BoundTypeServer
}

func (t *TanEnterRoomRequest) Marshal(io encoding.IO) {
	io.Uint32(&t.OwnerID)
	io.Uint32(&t.RoomID)
	io.StringUTF(&t.EnterPassword)
	io.Uint64(&t.EnterTeamID)
	io.Uint32(&t.EnterToken)
	io.Uint64(&t.FollowTeamID)
	io.StringUTF(&t.NetherNetID)
	io.Bool(&t.SupportWebRTCCompress)
}
