package service

import (
	"time"

	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/service/signaling"
)

const (
	RoomPrivacyEveryoneCanSee uint8 = iota
	RoomPrivacyOnlyFriendsCanSee
)

const (
	PlayerPermissionVisitor uint32 = iota
	PlayerPermissionMember
	PlayerPermissionOperator
	PlayerPermissionCustom
)

// RoomConfig ..
type RoomConfig struct {
	RoomName         string
	RoomPasscode     string
	RoomPrivacy      uint8
	RoomTagList      []uint8
	RoomRefreshTime  time.Duration
	MaxPlayerCount   uint8
	UsedModItemIDs   []uint64
	PlayerPermission uint32
	AllowPvP         bool
}

// DefaultRoomConfig ..
func DefaultRoomConfig(roomName string, roomPasscode string, maxPlayerCount uint8, playerPermission uint32) RoomConfig {
	return RoomConfig{
		RoomName:         roomName,
		RoomPasscode:     roomPasscode,
		RoomPrivacy:      RoomPrivacyEveryoneCanSee,
		RoomTagList:      nil,
		RoomRefreshTime:  signaling.RefreshTimeDefault,
		MaxPlayerCount:   maxPlayerCount,
		UsedModItemIDs:   nil,
		PlayerPermission: playerPermission,
		AllowPvP:         true,
	}
}
