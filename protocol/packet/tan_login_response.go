package packet

import "github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"

type TanLoginResponse struct {
	ErrorCode int8
}

func (*TanLoginResponse) ID() uint16 {
	return IDTanLoginResponse
}

func (*TanLoginResponse) BoundType() uint8 {
	return BoundTypeClient
}

func (t *TanLoginResponse) Marshal(io encoding.IO) {
	io.Int8(&t.ErrorCode)
}
