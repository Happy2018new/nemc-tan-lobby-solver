package bunker

import "fmt"

// Authenticator ..
type Authenticator interface {
	GetAccess(roomID string) (TanLobbyLoginResponse, error)
	GetCreate() (TanLobbyCreateResponse, error)
	GetRefresh() (TanLobbyRefreshResponse, error)
}

// AccessWrapper ..
type AccessWrapper struct {
	client *Client
}

func NewAccessWrapper(authServer string, token string) *AccessWrapper {
	return &AccessWrapper{
		client: NewClient(authServer, token),
	}
}

// GetAccess ..
func (aw *AccessWrapper) GetAccess(roomID string) (TanLobbyLoginResponse, error) {
	tanLobbyLoginResp, err := aw.client.Auth(roomID)
	if err != nil {
		return TanLobbyLoginResponse{}, fmt.Errorf("GetAccess: %v", err)
	}
	return tanLobbyLoginResp, nil
}

// GetCreate ..
func (aw *AccessWrapper) GetCreate() (TanLobbyCreateResponse, error) {
	tanLobbyCreateResp, err := aw.client.TanLobbyCreate()
	if err != nil {
		return TanLobbyCreateResponse{}, fmt.Errorf("GetCreate: %v", err)
	}
	return tanLobbyCreateResp, nil
}

// GetRefresh ..
func (aw *AccessWrapper) GetRefresh() (TanLobbyRefreshResponse, error) {
	tanLobbyRefreshResp, err := aw.client.TanLobbyRefresh()
	if err != nil {
		return TanLobbyRefreshResponse{}, fmt.Errorf("GetRefresh: %v", err)
	}
	return tanLobbyRefreshResp, nil
}
