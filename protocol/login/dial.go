package login

import (
	"context"
	_ "embed"
	"fmt"
	"math/rand/v2"
	"net"
	"strconv"
	"time"

	"github.com/Happy2018new/nemc-tan-lobby-solver/bunker/auth"
	"github.com/Happy2018new/nemc-tan-lobby-solver/core/nethernet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/core/raknet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/login/signaling"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/packet"
)

// Authenticator ..
type Authenticator interface {
	GetRoomID() string
	GetRoomPasscode() string
	GetAccess() (auth.TanLobbyLoginResponse, error)
}

// Dialer ..
type Dialer struct {
	Authenticator
	clientNetherID uint64
}

// Dial ..
func Dial(authenticator Authenticator) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	dialer := Dialer{Authenticator: authenticator}
	conn, err := dialer.DialContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("Dial: %v", err)
	}

	return conn, nil
}

// DialContext ..
func DialContext(ctx context.Context, authenticator Authenticator) (net.Conn, error) {
	dialer := Dialer{Authenticator: authenticator}
	conn, err := dialer.DialContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("DialContext: %v", err)
	}
	return conn, nil
}

// enterTanLobbyRoom ..
func (d *Dialer) enterTanLobbyRoom(ctx context.Context, tanLobbyLoginResp auth.TanLobbyLoginResponse) (
	remoteNetherNetID uint64,
	err error,
) {
	// Generate client nether ID and parse basic info
	d.clientNetherID = rand.Uint64()
	roomID, err := strconv.ParseUint(d.Authenticator.GetRoomID(), 10, 32)
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
	err = d.writeRaknetPacket(enc, &packet.TanLoginRequest{
		PlayerID:   tanLobbyLoginResp.UserUniqueID,
		Rand:       tanLobbyLoginResp.RaknetRand,
		AESRand:    tanLobbyLoginResp.RaknetAESRand,
		PlayerName: tanLobbyLoginResp.UserPlayerName,
	})
	if err != nil {
		return 0, fmt.Errorf("enterTanLobbyRoom: %v", err)
	}

	// Handle login response
	pk, err := d.readRaknetPacket(dec)
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
	err = d.writeRaknetPacket(enc, &packet.TanEnterRoomRequest{
		OwnerID:               tanLobbyLoginResp.RoomOwnerID,
		RoomID:                uint32(roomID),
		EnterPassword:         d.Authenticator.GetRoomPasscode(),
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
	pk, err = d.readRaknetPacket(dec)
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
		case packet.TanEnterRoomWrongPasscode:
			return 0, fmt.Errorf("enterTanLobbyRoom: Provided room passcode is incorrect")
		case packet.TanEnterRoomFullOfPeople:
			return 0, fmt.Errorf("enterTanLobbyRoom: Target room (%d) is full of people", roomID)
		default:
			return 0, fmt.Errorf("enterTanLobbyRoom: Enter room failed due to unknown reason (code = %d)", tanEnterRoomResp.ErrorCode)
		}
	}

	// Read incoming packet
	pkChannel := make(chan packet.Packet, 1)
	go func() {
		pk, err = d.readRaknetPacket(dec)
		if err != nil {
			return
		}
		pkChannel <- pk
	}()

	// Handle incoming packet
	select {
	case <-ctx.Done():
		return 0, fmt.Errorf("enterTanLobbyRoom: %v", ctx.Err())
	case pk := <-pkChannel:
		switch p := pk.(type) {
		case *packet.TanNotifyServerReady:
			remoteNetherNetID, err = strconv.ParseUint(p.NetherNetID, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("enterTanLobbyRoom: %v", err)
			}
			return remoteNetherNetID, nil
		case *packet.TanKickOutResponse:
			return 0, fmt.Errorf("enterTanLobbyRoom: The host owner kick you from the room")
		default:
			return 0, fmt.Errorf("enterTanLobbyRoom: Unknown packet received; pk = %#v", pk)
		}
	}
}

// DialContext ..
func (d *Dialer) DialContext(ctx context.Context) (conn net.Conn, err error) {
	// First we query room info
	tanLobbyLoginResp, err := d.Authenticator.GetAccess()
	if err != nil {
		return nil, nil
	}
	if !tanLobbyLoginResp.Success {
		return nil, fmt.Errorf("DialContext: %v", tanLobbyLoginResp.ErrorInfo)
	}

	// Then Enter tan lobby room
	remoteNetherNetID, err := d.enterTanLobbyRoom(ctx, tanLobbyLoginResp)
	if err != nil {
		return nil, fmt.Errorf("DialContext: %v", err)
	}

	// Connect to websocket signaling server
	wsConn, err := signaling.Dialer{
		NetworkID: d.clientNetherID,
	}.DialContext(
		ctx,
		tanLobbyLoginResp.SignalingServerAddress,
		d.clientNetherID,
		tanLobbyLoginResp.UserUniqueID,
		tanLobbyLoginResp.SignalingSeed,
		tanLobbyLoginResp.SignalingTicket,
	)
	if err != nil {
		return nil, fmt.Errorf("DialContext: %v", err)
	}
	defer wsConn.Close()

	// At last we can connect to remote room
	conn, err = nethernet.Dialer{}.DialContext(
		ctx,
		remoteNetherNetID,
		wsConn,
	)
	if err != nil {
		return nil, fmt.Errorf("DialContext: %v", err)
	}

	// Return
	return conn, nil
}
