package auth

import "fmt"

// AccessWrapper ..
type AccessWrapper struct {
	RoomID       string
	RoomPasscode string
	Token        string
	Client       *Client
}

// NewAccessWrapper ..
func NewAccessWrapper(client *Client, roomID string, roomPasscode string, token string) *AccessWrapper {
	return &AccessWrapper{
		RoomID:       roomID,
		RoomPasscode: roomPasscode,
		Token:        token,
		Client:       client,
	}
}

// GetRoomID ..
func (aw *AccessWrapper) GetRoomID() string {
	return aw.RoomID
}

// GetRoomPasscode ..
func (aw *AccessWrapper) GetRoomPasscode() string {
	return aw.RoomPasscode
}

// GetAccess ..
func (aw *AccessWrapper) GetAccess() (TanLobbyLoginResponse, error) {
	tanLobbyLoginResp, err := aw.Client.Auth(aw.RoomID, aw.Token)
	if err != nil {
		return TanLobbyLoginResponse{}, fmt.Errorf("GetAccess: %v", err)
	}
	return tanLobbyLoginResp, nil
}
