package packet

import "github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"

// TanKickOutResponse ..
type TanKickOutResponse struct{}

func (*TanKickOutResponse) ID() uint16 {
	return IDTanKickOutResponse
}

func (*TanKickOutResponse) BoundType() uint8 {
	return BoundTypeClient
}

func (t *TanKickOutResponse) Marshal(io encoding.IO) {}
