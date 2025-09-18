package packet

import "github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"

// TanNewGuestResponse ..
type TanNewGuestResponse struct {
	PlayerID              uint32
	NetherNetID           string
	SupportWebRTCCompress bool
}

func (*TanNewGuestResponse) ID() uint16 {
	return IDTanNewGuestResponse
}

func (*TanNewGuestResponse) BoundType() uint8 {
	return BoundTypeClient
}

func (t *TanNewGuestResponse) Marshal(io encoding.IO) {
	io.Uint32(&t.PlayerID)
	io.StringUTF(&t.NetherNetID)
	io.Bool(&t.SupportWebRTCCompress)
}
