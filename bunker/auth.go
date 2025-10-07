package bunker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// PhoenixSkinInfo ..
type PhoenixSkinInfo struct {
	ItemID          string `json:"entity_id"`
	SkinDownloadURL string `json:"res_url"`
	SkinIsSlim      bool   `json:"is_slim"`
}

// TanLobbyLoginRequest ..
type TanLobbyLoginRequest struct {
	FBToken            string `json:"login_token"`
	ProvidedPEAuthData string `json:"provided_pe_auth_data"`
	ProvidedSaAuthData string `json:"provided_sa_auth_data"`
	RoomID             string `json:"room_id"`
}

// TanLobbyLoginResponse ..
type TanLobbyLoginResponse struct {
	Success   bool   `json:"success"`
	ErrorInfo string `json:"error_info"`

	UserUniqueID   uint32          `json:"user_unique_id"`
	UserPlayerName string          `json:"user_player_name"`
	BotLevel       int             `json:"growth_level"`
	BotSkin        PhoenixSkinInfo `json:"skin_info"`
	BotComponent   map[string]*int `json:"outfit_info,omitempty"`

	RoomOwnerID        uint32   `json:"room_owner_id"`
	RoomModDisplayName []string `json:"room_mod_display_name"`
	RoomModDownloadURL []string `json:"room_mod_download_url"`
	RoomModEncryptKey  [][]byte `json:"room_mod_encrypt_key"`

	RaknetServerAddress string `json:"raknet_server_address"`
	RaknetRand          []byte `json:"raknet_rand"`
	RaknetAESRand       []byte `json:"raknet_aes_rand"`
	EncryptKeyBytes     []byte `json:"encrypt_key_bytes"`
	DecryptKeyBytes     []byte `json:"decrypt_key_bytes"`

	SignalingServerAddress string `json:"signaling_server_address"`
	SignalingSeed          []byte `json:"signaling_seed"`
	SignalingTicket        []byte `json:"signaling_ticket"`
}

func (client *Client) Auth(roomID string) (TanLobbyLoginResponse, error) {
	// Pack request
	request := TanLobbyLoginRequest{
		FBToken:            client.FBToken,
		ProvidedPEAuthData: client.ProvidedPEAuthData,
		ProvidedSaAuthData: client.ProvidedSaAuthData,
		RoomID:             roomID,
	}
	requestJsonBytes, _ := json.Marshal(request)

	// Post request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/phoenix/tan_lobby_login", client.AuthServer),
		"application/json",
		bytes.NewBuffer(requestJsonBytes),
	)
	if err != nil {
		return TanLobbyLoginResponse{}, fmt.Errorf("Auth: %v", err)
	}

	// Parse response
	tanLobbyLoginResp, err := parseHttpResponse[TanLobbyLoginResponse](resp)
	if err != nil {
		return TanLobbyLoginResponse{}, fmt.Errorf("Auth: %v", err)
	}

	// Return
	return tanLobbyLoginResp, nil
}

// TanLobbyCreateRequest ..
type TanLobbyCreateRequest struct {
	FBToken            string `json:"login_token"`
	ProvidedPEAuthData string `json:"provided_pe_auth_data"`
	ProvidedSaAuthData string `json:"provided_sa_auth_data"`
}

// TanLobbyCreateResponse ..
type TanLobbyCreateResponse struct {
	Success   bool   `json:"success"`
	ErrorInfo string `json:"error_info"`

	UserUniqueID   uint32 `json:"user_unique_id"`
	UserPlayerName string `json:"user_player_name"`

	RaknetServerAddress string `json:"raknet_server_address"`
	RaknetRand          []byte `json:"raknet_rand"`
	RaknetAESRand       []byte `json:"raknet_aes_rand"`
	EncryptKeyBytes     []byte `json:"encrypt_key_bytes"`
	DecryptKeyBytes     []byte `json:"decrypt_key_bytes"`

	SignalingServerAddress string `json:"signaling_server_address"`
	SignalingSeed          []byte `json:"signaling_seed"`
	SignalingTicket        []byte `json:"signaling_ticket"`
}

func (client *Client) TanLobbyCreate() (TanLobbyCreateResponse, error) {
	// Pack request
	request := TanLobbyCreateRequest{
		FBToken:            client.FBToken,
		ProvidedPEAuthData: client.ProvidedPEAuthData,
		ProvidedSaAuthData: client.ProvidedSaAuthData,
	}
	requestJsonBytes, _ := json.Marshal(request)

	// Post request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/phoenix/tan_lobby_create", client.AuthServer),
		"application/json",
		bytes.NewBuffer(requestJsonBytes),
	)
	if err != nil {
		return TanLobbyCreateResponse{}, fmt.Errorf("TanLobbyCreate: %v", err)
	}

	// Parse response
	tanLobbyCreateResp, err := parseHttpResponse[TanLobbyCreateResponse](resp)
	if err != nil {
		return TanLobbyCreateResponse{}, fmt.Errorf("TanLobbyCreate: %v", err)
	}

	// Return
	return tanLobbyCreateResp, nil
}

// TanLobbyDebugRequest ..
type TanLobbyDebugRequest struct {
	FBToken       string `json:"login_token"`
	LoginResponse string `json:"login_response"`
	RaknetRand    []byte `json:"raknet_rand"`
}

// TanLobbyDebugResponse ..
type TanLobbyDebugResponse struct {
	Success         bool   `json:"success"`
	ErrorInfo       string `json:"error_info"`
	EncryptKeyBytes []byte `json:"encrypt_key_bytes"`
	DecryptKeyBytes []byte `json:"decrypt_key_bytes"`
}

func (client *Client) TanLobbyDebug(loginResponse string, raknetRand []byte) (TanLobbyDebugResponse, error) {
	// Pack request
	request := TanLobbyDebugRequest{
		FBToken:       client.FBToken,
		LoginResponse: loginResponse,
		RaknetRand:    raknetRand,
	}
	requestJsonBytes, _ := json.Marshal(request)

	// Post request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/phoenix/tan_lobby_debug", client.AuthServer),
		"application/json",
		bytes.NewBuffer(requestJsonBytes),
	)
	if err != nil {
		return TanLobbyDebugResponse{}, fmt.Errorf("TanLobbyDebug: %v", err)
	}

	// Parse response
	tanLobbyDebugResp, err := parseHttpResponse[TanLobbyDebugResponse](resp)
	if err != nil {
		return TanLobbyDebugResponse{}, fmt.Errorf("TanLobbyDebug: %v", err)
	}

	// Return
	return tanLobbyDebugResp, nil
}
