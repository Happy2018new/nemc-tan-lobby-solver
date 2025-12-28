package service

import (
	"context"
	_ "embed"
	"fmt"
	"math/rand/v2"
	"net"
	"strconv"
	"time"

	"github.com/Happy2018new/nemc-tan-lobby-solver/bunker"
	"github.com/Happy2018new/nemc-tan-lobby-solver/core/nethernet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/core/raknet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/packet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/service/signaling"
)

// Dialer ..
type Dialer struct {
	bunker.Authenticator
	RoomID         string
	RoomPasscode   string
	clientNetherID uint64
}

// Dial ..
func Dial(roomID string, roomPasscode string, authenticator bunker.Authenticator) (net.Conn, bunker.TanLobbyLoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	conn, tanLobbyLoginResp, err := DialContext(ctx, roomID, roomPasscode, authenticator)
	if err != nil {
		return nil, bunker.TanLobbyLoginResponse{}, fmt.Errorf("Dial: %v", err)
	}

	return conn, tanLobbyLoginResp, nil
}

// DialContext ..
func DialContext(
	ctx context.Context,
	roomID string,
	roomPasscode string,
	authenticator bunker.Authenticator,
) (net.Conn, bunker.TanLobbyLoginResponse, error) {
	dialer := Dialer{
		Authenticator: authenticator,
		RoomID:        roomID,
		RoomPasscode:  roomPasscode,
	}
	conn, tanLobbyLoginResp, err := dialer.DialContext(ctx)
	if err != nil {
		return nil, bunker.TanLobbyLoginResponse{}, fmt.Errorf("DialContext: %v", err)
	}
	return conn, tanLobbyLoginResp, nil
}

// enterTanLobbyRoom ..
func (d *Dialer) enterTanLobbyRoom(ctx context.Context, tanLobbyLoginResp bunker.TanLobbyLoginResponse) (
	remoteNetherNetID uint64,
	err error,
) {
	// Generate client nether ID and parse basic info
	d.clientNetherID = rand.Uint64()
	roomID, err := strconv.ParseUint(d.RoomID, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("enterTanLobbyRoom: %v", err)
	}

	// Create conn
	conn, err := raknet.DialContext(ctx, tanLobbyLoginResp.RaknetServerAddress)
	if err != nil {
		return 0, fmt.Errorf("enterTanLobbyRoom: %v", err)
	}
	defer conn.Close()

	// Set encoder and decoder
	enc := packet.NewEncoder(conn)
	dec, err := packet.NewDecoder(conn)
	if err != nil {
		return 0, fmt.Errorf("enterTanLobbyRoom: %v", err)
	}

	// Send login request
	err = writePacket(enc, &packet.TanLoginRequest{
		PlayerID:   tanLobbyLoginResp.UserUniqueID,
		Rand:       tanLobbyLoginResp.RaknetRand,
		AESRand:    tanLobbyLoginResp.RaknetAESRand,
		PlayerName: tanLobbyLoginResp.UserPlayerName,
	})
	if err != nil {
		return 0, fmt.Errorf("enterTanLobbyRoom: %v", err)
	}

	// Handle login response
	pk, err := readPacketWithContext(ctx, conn, dec)
	if err != nil {
		return 0, fmt.Errorf("enterTanLobbyRoom: %v", err)
	}
	tanLoginResp, ok := pk.(*packet.TanLoginResponse)
	if !ok {
		return 0, fmt.Errorf("enterTanLobbyRoom: Expect the incoming packet is *packet.TanLoginResponse, but got %#v", pk)
	}
	if tanLoginResp.ErrorCode != packet.TanLoginSuccess {
		return 0, fmt.Errorf("enterTanLobbyRoom: Login failed (code = %d)", tanLoginResp.ErrorCode)
	}

	// Enable encryption
	enc.EnableEncryption(tanLobbyLoginResp.EncryptKeyBytes, tanLobbyLoginResp.DecryptKeyBytes)
	dec.EnableEncryption(tanLobbyLoginResp.EncryptKeyBytes, tanLobbyLoginResp.DecryptKeyBytes)

	// Enter room
	err = writePacket(enc, &packet.TanEnterRoomRequest{
		OwnerID:               tanLobbyLoginResp.RoomOwnerID,
		RoomID:                uint32(roomID),
		EnterPassword:         d.RoomPasscode,
		EnterTeamID:           0,
		EnterToken:            0,
		FollowTeamID:          0,
		NetherNetID:           fmt.Sprintf("%d", d.clientNetherID),
		SupportWebRTCCompress: true,
	})
	if err != nil {
		return 0, fmt.Errorf("enterTanLobbyRoom: %v", err)
	}

	// Handle enter room response
	pk, err = readPacketWithContext(ctx, conn, dec)
	if err != nil {
		return 0, fmt.Errorf("enterTanLobbyRoom: %v", err)
	}
	tanEnterRoomResp, ok := pk.(*packet.TanEnterRoomResponse)
	if !ok {
		return 0, fmt.Errorf("enterTanLobbyRoom: Expect the incoming packet is *packet.TanEnterRoomResponse, but got %#v", pk)
	}
	if tanEnterRoomResp.ErrorCode != packet.TanEnterRoomSuccess {
		switch tanEnterRoomResp.ErrorCode {
		case packet.TanEnterRoomNotFound:
			return 0, fmt.Errorf("enterTanLobbyRoom: Target room (%d) is closed", roomID)
		case packet.TanEnterRoomAccessDenied:
			return 0, fmt.Errorf("enterTanLobbyRoom: Access denied (Possibly is that you are not the friend of the owner)")
		case packet.TanEnterRoomFullOfPeople:
			return 0, fmt.Errorf("enterTanLobbyRoom: Target room (%d) is full of people", roomID)
		case packet.TanEnterRoomWrongPasscode:
			return 0, fmt.Errorf("enterTanLobbyRoom: Provided room passcode is incorrect")
		default:
			return 0, fmt.Errorf("enterTanLobbyRoom: Enter room failed due to unknown reason (code = %d)", tanEnterRoomResp.ErrorCode)
		}
	}

	// Read and handle incoming packet
	for {
		pk, err = readPacketWithContext(ctx, conn, dec)
		if err != nil {
			return 0, fmt.Errorf("enterTanLobbyRoom: %v", err)
		}
		switch p := pk.(type) {
		case *packet.TanNotifyServerReady:
			remoteNetherNetID, err = strconv.ParseUint(p.NetherNetID, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("enterTanLobbyRoom: %v", err)
			}
			return remoteNetherNetID, nil
		case *packet.TanKickOutResponse:
			return 0, fmt.Errorf("enterTanLobbyRoom: The host owner kick you from the room")
		case *packet.TanNewGuestResponse, *packet.TanSetTagListResponse:
		default:
			return 0, fmt.Errorf("enterTanLobbyRoom: Unknown packet received; pk = %#v", pk)
		}
	}
}

// DialContext ..
func (d *Dialer) DialContext(ctx context.Context) (conn net.Conn, authResp bunker.TanLobbyLoginResponse, err error) {
	// First we query room info
	tanLobbyLoginResp, err := d.Authenticator.GetAccess(d.RoomID)
	if err != nil {
		return nil, bunker.TanLobbyLoginResponse{}, fmt.Errorf("DialContext: %v", err)
	}
	if !tanLobbyLoginResp.Success {
		return nil, bunker.TanLobbyLoginResponse{}, fmt.Errorf("DialContext: %v", tanLobbyLoginResp.ErrorInfo)
	}

	// Then Enter tan lobby room
	remoteNetherNetID, err := d.enterTanLobbyRoom(ctx, tanLobbyLoginResp)
	if err != nil {
		return nil, bunker.TanLobbyLoginResponse{}, fmt.Errorf("DialContext: %v", err)
	}

	// Connect to websocket signaling server
	wsConnection, err := signaling.Dialer{
		Authenticator: d.Authenticator,
		RefreshTime:   signaling.RefreshTimeDisbale,
		NetherNetID:   d.clientNetherID,
	}.DialContext(
		ctx,
		tanLobbyLoginResp.SignalingServerAddress,
		tanLobbyLoginResp.UserUniqueID,
		tanLobbyLoginResp.SignalingSeed,
		tanLobbyLoginResp.SignalingTicket,
	)
	if err != nil {
		return nil, bunker.TanLobbyLoginResponse{}, fmt.Errorf("DialContext: %v", err)
	}
	defer wsConnection.Close(fmt.Errorf("DialContext: Close as expected (this error can ignore)"))

	// At last we can connect to remote room
	conn, err = nethernet.Dialer{}.DialContext(
		ctx,
		remoteNetherNetID,
		wsConnection,
	)
	if err != nil {
		return nil, bunker.TanLobbyLoginResponse{}, fmt.Errorf("DialContext: %v", err)
	}

	// Return
	return conn, tanLobbyLoginResp, nil
}
