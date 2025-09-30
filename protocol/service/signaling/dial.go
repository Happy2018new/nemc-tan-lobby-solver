package signaling

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/Happy2018new/nemc-tan-lobby-solver/bunker"
	"github.com/coder/websocket"
)

const (
	RefreshTimeDisbale = time.Duration(0)
	RefreshTimeDefault = time.Minute * 30
)

// Dialer ..
type Dialer struct {
	bunker.Authenticator
	Options           *websocket.DialOptions
	RefreshTime       time.Duration
	G79UserUID        uint32
	ServerBaseAddress string
	ClientNetherNetID uint64
}

// DialContext ..
func (d Dialer) DialContext(ctx context.Context, signalingSeed []byte, signalingTicket []byte) (*Conn, error) {
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
	if d.ClientNetherNetID == 0 {
		d.ClientNetherNetID = rand.Uint64()
	}

	finalAddress := fmt.Sprintf(
		"ws://%s/%d/%d/%s/%s",
		d.ServerBaseAddress,
		d.ClientNetherNetID,
		d.G79UserUID,
		base64.URLEncoding.EncodeToString(signalingSeed),
		base64.URLEncoding.EncodeToString(signalingTicket),
	)
	c, _, err := websocket.Dial(ctx, finalAddress, d.Options)
	if err != nil {
		return nil, err
	}

	conn, err := NewConn(ctx, c, d)
	if err != nil {
		return nil, fmt.Errorf("DialContext: %v", err)
	}
	return conn, nil
}
