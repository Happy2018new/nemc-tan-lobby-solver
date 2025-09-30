package service

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
