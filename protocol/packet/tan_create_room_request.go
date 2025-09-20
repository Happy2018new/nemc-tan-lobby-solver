package packet

import "github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"

// BuiltInRoomName ..
var BuiltInRoomName = []string{
	"来和我一起玩吧！",
	"一起玩生存",
	"一起做建筑",
	"自由创造",
	"探索冒险",
	"动作游戏爱好者",
	"精美建筑欣赏",
	"交个新朋友",
	"寻找通关搭档",
	"寻找跑酷达人",
	"来和我PK吧",
	"有人来挑战我吗？",
	"来一场速度比拼！",
	"轻松休闲，快乐游戏",
	"一起聊聊天吧",
}

// TanCreateRoomRequest ..
type TanCreateRoomRequest struct {
	Capacity     uint8
	Privacy      uint8
	Name         string
	Tips         encoding.RoomTips
	ItemIDs      []uint64
	MinLevel     uint32
	PvP          bool
	TeamID       uint64
	PlayerAuth   uint32
	Password     string
	Slogan       string
	MapID        uint64
	EnableWebRTC bool
	OwnerPing    uint8
	PerfLv       uint8
}

func (*TanCreateRoomRequest) ID() uint16 {
	return IDTanCreateRoomRequest
}

func (*TanCreateRoomRequest) BoundType() uint8 {
	return BoundTypeServer
}

func (t *TanCreateRoomRequest) Marshal(io encoding.IO) {
	io.Uint8(&t.Capacity)
	io.Uint8(&t.Privacy)
	io.StringUTF(&t.Name)
	io.RoomTips(&t.Tips)
	encoding.FuncSliceUint8Length(io, &t.ItemIDs, io.Uint64)
	io.Uint32(&t.MinLevel)
	io.Bool(&t.PvP)
	io.Uint64(&t.TeamID)
	io.Uint32(&t.PlayerAuth)
	io.StringUTF(&t.Password)
	io.StringUTF(&t.Slogan)
	io.Uint64(&t.MapID)
	io.Bool(&t.EnableWebRTC)
	io.Uint8(&t.OwnerPing)
	io.Uint8(&t.PerfLv)
}
