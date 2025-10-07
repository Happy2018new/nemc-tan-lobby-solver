package packet

import "github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"

const (
	TagSurvival uint8 = iota
	TagCreative
	TagAdvance
	TagSimulate
	TagAction
	TagBuilding
	TagCompetitive
	TagCasual
)

// TanSetTagListRequest ..
type TanSetTagListRequest struct {
	TagList []uint8
}

func (*TanSetTagListRequest) ID() uint16 {
	return IDTanSetTagListRequest
}

func (*TanSetTagListRequest) BoundType() uint8 {
	return BoundTypeServer
}

func (t *TanSetTagListRequest) Marshal(io encoding.IO) {
	encoding.FuncSliceUint8Length(io, &t.TagList, io.Uint8)
}
