package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// TanLobbyLoginRequest ..
type TanLobbyLoginRequest struct {
	FBToken string `json:"login_token"`
	RoomID  string `json:"room_id"`
}

// TanLobbyLoginResponse ..
type TanLobbyLoginResponse struct {
	Success   bool   `json:"success"`
	ErrorInfo string `json:"error_info"`

	RoomOwnerID    uint32 `json:"room_owner_id"`
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

func (client *Client) Auth(roomID string, fbtoken string) (TanLobbyLoginResponse, error) {
	// Pack request
	request := TanLobbyLoginRequest{
		FBToken: fbtoken,
		RoomID:  roomID,
	}
	requestJsonBytes, _ := json.Marshal(request)

	// Post request
	resp, err := client.client.Post(
		fmt.Sprintf("%s/api/phoenix/tan_lobby_login", client.AuthServer),
		"application/json",
		bytes.NewBuffer(requestJsonBytes),
	)
	if err != nil {
		return TanLobbyLoginResponse{}, fmt.Errorf("Auth: %v", err)
	}

	// Parse response
	tanLobbyLoginResp, err := assertAndParse[TanLobbyLoginResponse](resp)
	if err != nil {
		return TanLobbyLoginResponse{}, fmt.Errorf("Auth: %v", err)
	}

	// Return
	return tanLobbyLoginResp, nil
}
