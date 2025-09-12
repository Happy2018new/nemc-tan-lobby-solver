package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// TanLobbyLoginRequest ..
type TanLobbyLoginRequest struct {
	FBToken      string `json:"login_token"`
	RoomID       string `json:"room_id"`
	RoomPasscode string `json:"room_passcode"`
}

// TanLobbyLoginResponse ..
type TanLobbyLoginResponse struct {
	Success   bool   `json:"success"`
	ErrorInfo string `json:"error_info"`

	RoomOwnerID    uint32 `json:"room_owner_id"`
	UserUniqueID   uint32 `json:"user_unique_id"`
	UserPlayerName string `json:"user_player_name"`

	RaknetRand      []byte `json:"raknet_rand"`
	RaknetAESRand   []byte `json:"raknet_aes_rand"`
	EncryptKeyBytes []byte `json:"encrypt_key_bytes"`
	DecryptKeyBytes []byte `json:"decrypt_key_bytes"`

	SignalingSeed   []byte `json:"signaling_seed"`
	SignalingTicket []byte `json:"signaling_ticket"`
}

func (client *Client) Auth(roomID string, roomPasscode string, fbtoken string) (TanLobbyLoginResponse, error) {
	// Pack request
	request := TanLobbyLoginRequest{
		FBToken:      fbtoken,
		RoomID:       roomID,
		RoomPasscode: roomPasscode,
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

// TanLobbyTransferServersRequest ..
type TanLobbyTransferServersRequest struct {
	FBToken string `json:"login_token"`
}

// TanLobbyTransferServersResponse ..
type TanLobbyTransferServersResponse struct {
	Success          bool     `json:"success"`
	ErrorInfo        string   `json:"error_info"`
	RaknetServers    []string `json:"raknet_servers"`
	WebsocketServers []string `json:"websocket_servers"`
}

func (client *Client) TransferServerList(fbtoken string) (TanLobbyTransferServersResponse, error) {
	// Pack request
	request := TanLobbyTransferServersRequest{
		FBToken: fbtoken,
	}
	requestJsonBytes, _ := json.Marshal(request)

	// Post request
	resp, err := client.client.Post(
		fmt.Sprintf("%s/api/phoenix/tan_lobby_transfer_server", client.AuthServer),
		"application/json",
		bytes.NewBuffer(requestJsonBytes),
	)
	if err != nil {
		return TanLobbyTransferServersResponse{}, fmt.Errorf("TransferServerList: %v", err)
	}

	// Parse response
	tanLobbyTransferServersResp, err := assertAndParse[TanLobbyTransferServersResponse](resp)
	if err != nil {
		return TanLobbyTransferServersResponse{}, fmt.Errorf("TransferServerList: %v", err)
	}

	// Return
	return tanLobbyTransferServersResp, nil
}
