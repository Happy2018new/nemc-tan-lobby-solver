package service

import "github.com/Happy2018new/nemc-tan-lobby-solver/bunker/auth"

// Authenticator ..
type Authenticator interface {
	GetRoomID() string
	GetRoomPasscode() string
	GetAccess() (auth.TanLobbyLoginResponse, error)
	GetCreate() (auth.TanLobbyCreateResponse, error)
}
