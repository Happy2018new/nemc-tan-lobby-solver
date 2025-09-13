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
	DefaultRaknetServerRepeatTimes  = 30
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
	tanLobbyLoginResp *auth.TanLobbyLoginResponse
	clientNetherID    uint64
}

// DialContext dials a Minecraft connection to the address passed over the network passed. The network is
// typically "raknet". A Conn is returned which may be used to receive packets from and send packets to.
// If a connection is not established before the context passed is cancelled, DialContext returns an error.
// DialContext uses a zero value of Dialer to initiate the connection.
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
	defer func() {
		if conn != nil && !success {
			_ = conn.Close()
		}
	}()

	// Set encoder and decoder
	enc = packet.NewEncoder(conn)
	dec, err = packet.NewDecoder(conn)
	if err != nil {
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
		return nil, nil, nil, false
	}

	// Handle login response
	pk, err := d.readRaknetPacket(dec)
	if err != nil {
		return nil, nil, nil, false
	}
	tanLoginResp, ok := pk.(*packet.TanLoginResponse)
	if !ok {
		return nil, nil, nil, false
	}
	if tanLoginResp.ErrorCode != packet.TanLoginSuccess {
		return nil, nil, nil, false
	}

	// Enable encryption
	enc.EnableEncryption(tanLobbyLoginResp.EncryptKeyBytes, tanLobbyLoginResp.DecryptKeyBytes)
	dec.EnableEncryption(tanLobbyLoginResp.EncryptKeyBytes, tanLobbyLoginResp.DecryptKeyBytes)

	// Return
	success = true
	return conn, enc, dec, success
}

// enterTanLobbyRoom ..
func (d *Dialer) enterTanLobbyRoom(
	tanLobbyLoginResp auth.TanLobbyLoginResponse,
	tanLobbyTransferServersResp auth.TanLobbyTransferServersResponse,
) (pk packet.Packet, raknetAddress string, success bool) {
	// Prepare
	var raknetServerMu sync.Mutex
	var possibleRaknetServers []string

	// Parse basic info and generate client nether ID
	roomID, err := strconv.ParseUint(d.Authenticator.GetRoomID(), 10, 32)
	if err != nil {
		return nil, "", false
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
		return nil, "", false
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
			return pk, address, true
		default:
			_ = conn.Close()
		}
	}

	// Return unsuccessful
	return nil, "", false
}

// DialContext dials a Minecraft connection to the address passed over the network passed. The network is
// typically "raknet". A Conn is returned which may be used to receive packets from and send packets to.
// If a connection is not established before the context passed is cancelled, DialContext returns an error.
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
	d.tanLobbyLoginResp = &tanLobbyLoginResp

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
		pk, addr, success := d.enterTanLobbyRoom(tanLobbyLoginResp, tanLobbyTransferServersResp)
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
	wsConn, err := signaling.Dialer{}.DialContext(
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
	conn, err = nethernet.Dialer{ConnectionID: d.clientNetherID}.DialContext(
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
