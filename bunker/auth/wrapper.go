package auth

import "fmt"

// AccessWrapper ..
type AccessWrapper struct {
	RoomID       string
	RoomPasscode string
	Token        string
	Client       *Client
}

// NewClientAccessWrapper ..
func NewClientAccessWrapper(client *Client, roomID string, roomPasscode string, token string) *AccessWrapper {
	return &AccessWrapper{
		RoomID:       roomID,
		RoomPasscode: roomPasscode,
		Token:        token,
		Client:       client,
	}
}

func NewServerAccessWrapper(client *Client, token string) *AccessWrapper {
	return &AccessWrapper{
		Client: client,
		Token:  token,
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

// GetCreate ..
func (aw *AccessWrapper) GetCreate() (TanLobbyCreateResponse, error) {
	tanLobbyCreateResp, err := aw.Client.TanLobbyCreate(aw.Token)
	if err != nil {
		return TanLobbyCreateResponse{}, fmt.Errorf("GetAccess: %v", err)
	}
	return tanLobbyCreateResp, nil
}
