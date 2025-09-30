package signaling

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/coder/websocket"
)

// Dialer ..
type Dialer struct {
	Options       *websocket.DialOptions
	NetworkID     uint64
	serverAddress string
}

// DialContext ..
func (d Dialer) DialContext(
	ctx context.Context,
	serverBaseAddress string,
	clientNetherID uint64,
	g79UserUID uint32,
	signalingSeed []byte,
	signalingTicket []byte,
) (*Conn, error) {
	if d.Options == nil {
		d.Options = &websocket.DialOptions{}
	}
	if d.Options.HTTPClient == nil {
		d.Options.HTTPClient = &http.Client{}
	}
	if d.Options.HTTPHeader == nil {
		d.Options.HTTPHeader = make(http.Header)
		d.Options.HTTPHeader.Set("Authorization", "NeteaseSignalingAuthToken")
	}
	if d.NetworkID == 0 {
		d.NetworkID = rand.Uint64()
	}

	d.serverAddress = fmt.Sprintf(
		"ws://%s/%d/%d/%s/%s",
		serverBaseAddress,
		clientNetherID,
		g79UserUID,
		base64.URLEncoding.EncodeToString(signalingSeed),
		base64.URLEncoding.EncodeToString(signalingTicket),
	)
	c, _, err := websocket.Dial(ctx, d.serverAddress, d.Options)
	if err != nil {
		return nil, err
	}

	conn, err := NewConn(ctx, c, d)
	if err != nil {
		return nil, fmt.Errorf("DialContext: %v", err)
	}
	return conn, nil
}
