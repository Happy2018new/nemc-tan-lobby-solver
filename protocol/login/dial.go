package login

import (
	"context"
	_ "embed"
	"fmt"
	"math"
	"math/rand/v2"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Happy2018new/nemc-tan-lobby-solver/bunker/auth"
	"github.com/Happy2018new/nemc-tan-lobby-solver/core/nethernet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/core/raknet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/login/signaling"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/packet"
)

const (
	EnableDebug                     = true
	DefaultRaknetServerCollectTimes = time.Second * 5
	DefaultRaknetServerRepeatTimes  = 6
)

// Authenticator ..
type Authenticator interface {
	GetRoomID() string
	GetRoomPasscode() string
	GetAccess() (auth.TanLobbyLoginResponse, error)
	TransferServerList() (auth.TanLobbyTransferServersResponse, error)
}

// Dialer ..
type Dialer struct {
	Authenticator
	clientNetherID uint64
}

// Dial ..
func Dial(authenticator Authenticator) (net.Conn, error) {
	dialer := Dialer{Authenticator: authenticator}
	conn, err := dialer.Dial()
	if err != nil {
		return nil, fmt.Errorf("Dial: %v", err)
	}
	return conn, nil
}

// raknetDialer ..
func (d *Dialer) raknetDialer(
	ctx context.Context,
	address string,
	tanLobbyLoginResp auth.TanLobbyLoginResponse,
) (
	conn net.Conn,
	enc *packet.Encoder,
	dec *packet.Decoder,
	success bool,
) {
	// Create conn
	conn, err := raknet.DialContext(ctx, address)
	if err != nil {
		return nil, nil, nil, false
	}

	// Set encoder and decoder
	enc = packet.NewEncoder(conn)
	dec, err = packet.NewDecoder(conn)
	if err != nil {
		_ = conn.Close()
		return nil, nil, nil, false
	}

	// Send login request
	err = d.writeRaknetPacket(enc, &packet.TanLoginRequest{
		PlayerID:   tanLobbyLoginResp.UserUniqueID,
		Rand:       tanLobbyLoginResp.RaknetRand,
		AESRand:    tanLobbyLoginResp.RaknetAESRand,
		PlayerName: tanLobbyLoginResp.UserPlayerName,
	})
	if err != nil {
		_ = conn.Close()
		return nil, nil, nil, false
	}

	// Handle login response
	pk, err := d.readRaknetPacket(dec)
	if err != nil {
		_ = conn.Close()
		return nil, nil, nil, false
	}
	tanLoginResp, ok := pk.(*packet.TanLoginResponse)
	if !ok {
		_ = conn.Close()
		return nil, nil, nil, false
	}
	if tanLoginResp.ErrorCode != packet.TanLoginSuccess {
		_ = conn.Close()
		return nil, nil, nil, false
	}

	// Enable encryption
	enc.EnableEncryption(tanLobbyLoginResp.EncryptKeyBytes, tanLobbyLoginResp.DecryptKeyBytes)
	dec.EnableEncryption(tanLobbyLoginResp.EncryptKeyBytes, tanLobbyLoginResp.DecryptKeyBytes)

	// Return
	return conn, enc, dec, true
}

// enterTanLobbyRoom ..
func (d *Dialer) enterTanLobbyRoom(
	tanLobbyLoginResp auth.TanLobbyLoginResponse,
	tanLobbyTransferServersResp auth.TanLobbyTransferServersResponse,
) (pk packet.Packet, raknetAddress string, wrongPasscode bool, success bool) {
	// Prepare
	var raknetServerMu sync.Mutex
	var possibleRaknetServers []string

	// Parse basic info and generate client nether ID
	roomID, err := strconv.ParseUint(d.Authenticator.GetRoomID(), 10, 32)
	if err != nil {
		return nil, "", false, false
	}
	d.clientNetherID = rand.Uint64N(math.MaxUint64)

	// Get available raknet server
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRaknetServerCollectTimes)
	defer cancel()
	for _, value := range tanLobbyTransferServersResp.RaknetServers {
		go func() {
			conn, _, _, success := d.raknetDialer(ctx, value, tanLobbyLoginResp)
			if success {
				raknetServerMu.Lock()
				defer raknetServerMu.Unlock()
				possibleRaknetServers = append(possibleRaknetServers, value)
				_ = conn.Close()
			}
		}()
	}

	// Check possible raknet servers
	<-ctx.Done()
	if EnableDebug {
		fmt.Printf("Find possible available raknet servers: %v\n", possibleRaknetServers)
	}
	if len(possibleRaknetServers) == 0 {
		return nil, "", false, false
	}

	// Find final raknet server
	for _, address := range possibleRaknetServers {
		// Create conn
		conn, enc, dec, success := d.raknetDialer(context.Background(), address, tanLobbyLoginResp)
		if !success {
			continue
		}

		// Enter room
		err := d.writeRaknetPacket(enc, &packet.TanEnterRoomRequest{
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
			_ = conn.Close()
			continue
		}

		// Handle enter room response
		pk, err := d.readRaknetPacket(dec)
		if err != nil {
			_ = conn.Close()
			continue
		}
		tanEnterRoomResp, ok := pk.(*packet.TanEnterRoomResponse)
		if !ok {
			_ = conn.Close()
			continue
		}
		if tanEnterRoomResp.ErrorCode != packet.TanEnterRoomSuccess {
			if tanEnterRoomResp.ErrorCode == packet.TanEnterRoomWrongPasscode {
				return nil, "", true, false
			}
			_ = conn.Close()
			continue
		}

		// Handle ready or kick out packet
		pk, err = d.readRaknetPacket(dec)
		if err != nil {
			_ = conn.Close()
			continue
		}
		switch pk.(type) {
		case *packet.TanNotifyServerReady, *packet.TanKickOutResponse:
			_ = conn.Close()
			return pk, address, false, true
		default:
			_ = conn.Close()
		}
	}

	// Return unsuccessful
	return nil, "", false, false
}

// Dial ..
func (d *Dialer) Dial() (conn net.Conn, err error) {
	var raknetAddress string
	var websocketAddress string
	var remoteNetherNetID uint64

	// First we query room info
	tanLobbyLoginResp, err := d.Authenticator.GetAccess()
	if err != nil {
		return nil, nil
	}
	if !tanLobbyLoginResp.Success {
		return nil, fmt.Errorf("Dial: %v", tanLobbyLoginResp.ErrorInfo)
	}

	// Then get transfer server list
	tanLobbyTransferServersResp, err := d.Authenticator.TransferServerList()
	if err != nil {
		return nil, nil
	}
	if !tanLobbyTransferServersResp.Success {
		return nil, fmt.Errorf("Dial: %v", tanLobbyTransferServersResp.ErrorInfo)
	}

	// Enter tan lobby room
	for range DefaultRaknetServerRepeatTimes {
		pk, addr, wrongPasscode, success := d.enterTanLobbyRoom(tanLobbyLoginResp, tanLobbyTransferServersResp)
		if wrongPasscode {
			return nil, fmt.Errorf("Dial: Provided room passcode is incorrect")
		}
		if !success {
			continue
		}

		switch p := pk.(type) {
		case *packet.TanNotifyServerReady:
			remoteNetherNetID, _ = strconv.ParseUint(p.NetherNetID, 10, 64)
		case *packet.TanKickOutResponse:
			return nil, fmt.Errorf("Dial: The host owner kick you from the room")
		}

		raknetAddress = addr
		break
	}
	if len(raknetAddress) == 0 {
		return nil, fmt.Errorf("Dial: No available raknet server was found")
	}

	// Find websocket server address
	websocketServerIP := strings.Split(raknetAddress, ":")[0]
	for _, value := range tanLobbyTransferServersResp.WebsocketServers {
		if strings.Contains(value, websocketServerIP) {
			websocketAddress = value
			break
		}
	}
	if len(websocketAddress) == 0 {
		return nil, fmt.Errorf("Dial: No available websocket server was found")
	}

	// Connect to websocket server
	wsCTX, wsCTXCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer wsCTXCancel()
	wsConn, err := signaling.Dialer{
		NetworkID: d.clientNetherID,
	}.DialContext(
		wsCTX,
		websocketAddress,
		d.clientNetherID,
		tanLobbyLoginResp.UserUniqueID,
		tanLobbyLoginResp.SignalingSeed,
		tanLobbyLoginResp.SignalingTicket,
	)
	if err != nil {
		return nil, fmt.Errorf("Dial: %v", err)
	}
	defer wsConn.Close()

	// Connect to remote room
	mcCTX, mcCTXCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer mcCTXCancel()
	conn, err = nethernet.Dialer{}.DialContext(
		mcCTX,
		remoteNetherNetID,
		wsConn,
	)
	if err != nil {
		return nil, fmt.Errorf("Dial: %v", err)
	}

	// Return
	return conn, nil
}
