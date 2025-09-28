package service

import "github.com/Happy2018new/nemc-tan-lobby-solver/bunker/auth"

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

// Authenticator ..
type Authenticator interface {
	GetAccess(roomID string) (auth.TanLobbyLoginResponse, error)
	GetCreate() (auth.TanLobbyCreateResponse, error)
}

// RoomConfig ..
type RoomConfig struct {
	RoomName         string
	RoomPasscode     string
	RoomPrivacy      uint8
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
		MaxPlayerCount:   maxPlayerCount,
		UsedModItemIDs:   nil,
		PlayerPermission: playerPermission,
		AllowPvP:         true,
	}
}
