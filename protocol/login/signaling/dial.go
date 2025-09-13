package signaling

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/Happy2018new/nemc-tan-lobby-solver/core/nethernet"
	"github.com/coder/websocket"
)

// Dialer ..
type Dialer struct {
	Options   *websocket.DialOptions
	NetworkID uint64
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

	finalAddress := fmt.Sprintf(
		"ws://%s/%d/%d/%s/%s",
		serverBaseAddress,
		clientNetherID,
		g79UserUID,
		base64.URLEncoding.EncodeToString(signalingSeed),
		base64.URLEncoding.EncodeToString(signalingTicket),
	)
	c, _, err := websocket.Dial(ctx, finalAddress, d.Options)
	if err != nil {
		return nil, err
	}

	conn := &Conn{
		conn:    c,
		d:       d,
		signals: make(chan *nethernet.Signal),
		ready:   make(chan struct{}),
	}
	var cancel context.CancelCauseFunc
	conn.ctx, cancel = context.WithCancelCause(context.Background())

	go conn.read(cancel)
	go conn.ping()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-conn.ready:
		return conn, nil
	}
}
