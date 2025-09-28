package auth

import "fmt"

// AccessWrapper ..
type AccessWrapper struct {
	Token  string
	Client *Client
}

func NewAccessWrapper(client *Client, token string) *AccessWrapper {
	return &AccessWrapper{
		Client: client,
		Token:  token,
	}
}

// GetAccess ..
func (aw *AccessWrapper) GetAccess(roomID string) (TanLobbyLoginResponse, error) {
	tanLobbyLoginResp, err := aw.Client.Auth(roomID, aw.Token)
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
