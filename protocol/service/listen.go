package service

import (
	"context"
	_ "embed"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/Happy2018new/nemc-tan-lobby-solver/bunker"
	"github.com/Happy2018new/nemc-tan-lobby-solver/core/nethernet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/core/raknet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/packet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/service/signaling"
)

// ListenConfig ..
type ListenConfig struct {
	bunker.Authenticator
	RoomConfig
	serverNetherID    uint64
	raknetConnection  net.Conn
	netherNetListener *nethernet.Listener
}

// Listen ..
func Listen(roomConfig RoomConfig, authenticator bunker.Authenticator) (
	listenConfig *ListenConfig,
	listener *nethernet.Listener,
	roomID uint32,
	err error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	listenConfig, listener, roomID, err = ListenContext(ctx, roomConfig, authenticator)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("Listen: %v", err)
	}

	return
}

// ListenContext ..
func ListenContext(ctx context.Context, roomConfig RoomConfig, authenticator bunker.Authenticator) (
	listenConfig *ListenConfig,
	listener *nethernet.Listener,
	roomID uint32,
	err error,
) {
	listenConfig = &ListenConfig{
		Authenticator: authenticator,
		RoomConfig:    roomConfig,
	}
	listener, roomID, err = listenConfig.ListenContext(ctx)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("ListenContext: %v", err)
	}
	return
}

// createTanLobbyRoom ..
func (l *ListenConfig) createTanLobbyRoom(
	ctx context.Context,
	tanLobbyCreateResp bunker.TanLobbyCreateResponse,
) (
	conn net.Conn,
	enc *packet.Encoder,
	dec *packet.Decoder,
	roomID uint32,
	err error,
) {
	// Create conn
	conn, err = raknet.DialContext(ctx, tanLobbyCreateResp.RaknetServerAddress)
	if err != nil {
		return nil, nil, nil, 0, fmt.Errorf("createTanLobbyRoom: %v", err)
	}

	// Set encoder and decoder
	enc = packet.NewEncoder(conn)
	dec, err = packet.NewDecoder(conn)
	if err != nil {
		_ = conn.Close()
		return nil, nil, nil, 0, fmt.Errorf("createTanLobbyRoom: %v", err)
	}

	// Send login request
	err = writePacket(enc, &packet.TanLoginRequest{
		PlayerID:   tanLobbyCreateResp.UserUniqueID,
		Rand:       tanLobbyCreateResp.RaknetRand,
		AESRand:    tanLobbyCreateResp.RaknetAESRand,
		PlayerName: tanLobbyCreateResp.UserPlayerName,
	})
	if err != nil {
		_ = conn.Close()
		return nil, nil, nil, 0, fmt.Errorf("createTanLobbyRoom: %v", err)
	}

	// Handle login response
	pk, err := readPacketWithContext(ctx, conn, dec)
	if err != nil {
		_ = conn.Close()
		return nil, nil, nil, 0, fmt.Errorf("createTanLobbyRoom: %v", err)
	}
	tanLoginResp, ok := pk.(*packet.TanLoginResponse)
	if !ok {
		_ = conn.Close()
		return nil, nil, nil, 0, fmt.Errorf("createTanLobbyRoom: Expect the incoming packet is *packet.TanLoginResponse, but got %#v", pk)
	}
	if tanLoginResp.ErrorCode != packet.TanLoginSuccess {
		_ = conn.Close()
		return nil, nil, nil, 0, fmt.Errorf("createTanLobbyRoom: Login failed (code = %d)", tanLoginResp.ErrorCode)
	}

	// Enable encryption
	enc.EnableEncryption(tanLobbyCreateResp.EncryptKeyBytes, tanLobbyCreateResp.DecryptKeyBytes)
	dec.EnableEncryption(tanLobbyCreateResp.EncryptKeyBytes, tanLobbyCreateResp.DecryptKeyBytes)

	// Create room
	err = writePacket(enc, &packet.TanCreateRoomRequest{
		Capacity: l.RoomConfig.MaxPlayerCount,
		Privacy:  l.RoomConfig.RoomPrivacy,
		Name:     "",
		Tips: encoding.RoomTips{
			LevelID:            "World",
			GameType:           0,
			ConstantTestString: "test",
			Vioce:              0,
			ProtocolID:         38,
		},
		ItemIDs:      l.RoomConfig.UsedModItemIDs,
		MinLevel:     0,
		PvP:          l.RoomConfig.AllowPvP,
		TeamID:       0,
		PlayerAuth:   l.RoomConfig.PlayerPermission,
		Password:     l.RoomConfig.RoomPasscode,
		Slogan:       l.RoomConfig.RoomName,
		MapID:        0,
		EnableWebRTC: true,
		OwnerPing:    3,
		PerfLv:       1,
	})
	if err != nil {
		_ = conn.Close()
		return nil, nil, nil, 0, fmt.Errorf("createTanLobbyRoom: %v", err)
	}

	// Read create room response
	pk, err = readPacketWithContext(ctx, conn, dec)
	if err != nil {
		_ = conn.Close()
		return nil, nil, nil, 0, fmt.Errorf("createTanLobbyRoom: %v", err)
	}
	tanCreateRoomResp, ok := pk.(*packet.TanCreateRoomResponse)
	if !ok {
		_ = conn.Close()
		return nil, nil, nil, 0, fmt.Errorf("createTanLobbyRoom: Expect the incoming packet is *packet.TanEnterRoomResponse, but got %#v", pk)
	}

	// Handle create room response
	if tanCreateRoomResp.ErrorCode == packet.TanCreateRoomNeedVipToSetRoomName {
		_ = conn.Close()
		return nil, nil, nil, 0, fmt.Errorf("createTanLobbyRoom: Can only use built in room name (Need VIP to set custom name)")
	}
	if tanCreateRoomResp.ErrorCode != packet.TanCreateRoomSuccess {
		_ = conn.Close()
		return nil, nil, nil, 0, fmt.Errorf("createTanLobbyRoom: Failed to create tan lobby room (code = %d)", tanCreateRoomResp.ErrorCode)
	}

	// Return
	return conn, enc, dec, tanCreateRoomResp.RoomID, nil
}

// ListenContext ..
func (l *ListenConfig) ListenContext(ctx context.Context) (listener *nethernet.Listener, roomID uint32, err error) {
	// Prepare
	var enc *packet.Encoder
	var dec *packet.Decoder
	l.serverNetherID = rand.Uint64()

	// First we query room info
	tanLobbyCreateResp, err := l.Authenticator.GetCreate()
	if err != nil {
		return nil, 0, fmt.Errorf("ListenContext: %v", err)
	}
	if !tanLobbyCreateResp.Success {
		return nil, 0, fmt.Errorf("ListenContext: %v", tanLobbyCreateResp.ErrorInfo)
	}

	// Create tan lobby room
	l.raknetConnection, enc, dec, roomID, err = l.createTanLobbyRoom(ctx, tanLobbyCreateResp)
	if err != nil {
		return nil, 0, fmt.Errorf("ListenContext: %v", err)
	}
	go func() {
		for {
			pk, err := readPacket(dec)
			if err != nil {
				l.CloseRoom()
				return
			}
			if _, ok := pk.(*packet.TanNewGuestResponse); !ok {
				continue
			}
			writePacket(enc, &packet.TanNotifyServerReady{
				ServerAddress:         "127.0.0.1|19132",
				ServerRaknetGuid:      "",
				RTCRoomID:             fmt.Sprintf("%d", roomID),
				NetherNetID:           fmt.Sprintf("%d", l.serverNetherID),
				WebRTCCompressEnabled: true,
			})
		}
	}()

	// Connect to websocket signaling server
	wsConnection, err := signaling.Dialer{
		Authenticator:     l.Authenticator,
		RefreshTime:       l.RoomRefreshTime,
		G79UserUID:        tanLobbyCreateResp.UserUniqueID,
		ServerBaseAddress: tanLobbyCreateResp.SignalingServerAddress,
		ClientNetherNetID: l.serverNetherID,
	}.DialContext(
		ctx,
		tanLobbyCreateResp.SignalingSeed,
		tanLobbyCreateResp.SignalingTicket,
	)
	if err != nil {
		_ = l.raknetConnection.Close()
		return nil, 0, fmt.Errorf("ListenContext: %v", err)
	}

	// Init listen config
	listenConfig := nethernet.ListenConfig{
		ConnContext: func(parent context.Context, conn *nethernet.Conn) context.Context {
			ctx, cancel := context.WithTimeout(parent, time.Second*30)
			go func() {
				<-ctx.Done()
				cancel()
			}()
			return ctx
		},
	}

	// Create listener
	l.netherNetListener, err = listenConfig.Listen(wsConnection)
	if err != nil {
		_ = l.raknetConnection.Close()
		wsConnection.Close()
		return nil, 0, fmt.Errorf("ListenContext: %v", err)
	}

	// Return
	return l.netherNetListener, roomID, nil
}

// CloseRoom ..
func (l *ListenConfig) CloseRoom() {
	if l.raknetConnection != nil {
		_ = l.raknetConnection.Close()
	}
	if l.netherNetListener != nil {
		_ = l.netherNetListener.Close()
	}
}
