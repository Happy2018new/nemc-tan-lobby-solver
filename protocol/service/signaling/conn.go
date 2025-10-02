package signaling

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Happy2018new/nemc-tan-lobby-solver/core/nethernet"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const (
	PrintDebugInfo   = false
	PrintRefreshInfo = true
)

// Conn ..
type Conn struct {
	mu     *sync.Mutex
	dialer Dialer

	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelCauseFunc

	credentials chan nethernet.Credentials
	signals     chan *nethernet.Signal
	doOnce      *sync.Once
}

// NewConn ..
func NewConn(ctx context.Context, conn *websocket.Conn, dialer Dialer) (result *Conn, err error) {
	c := &Conn{
		mu:          new(sync.Mutex),
		dialer:      dialer,
		conn:        conn,
		credentials: make(chan nethernet.Credentials, 1),
		signals:     make(chan *nethernet.Signal),
		doOnce:      new(sync.Once),
	}
	c.ctx, c.cancel = context.WithCancelCause(context.Background())

	go c.read()
	go c.ping()

	select {
	case <-ctx.Done():
		c.Close(fmt.Errorf("NewConn: %v", ctx.Err()))
		return nil, fmt.Errorf("NewConn: %v", ctx.Err())
	case <-c.ctx.Done():
		return nil, fmt.Errorf("NewConn: %v", c.ctx.Err())
	case credential := <-c.credentials:
		c.credentials <- credential
		go c.autoRefresh(dialer.RefreshTime)
		return c, nil
	}
}

// read ..
func (c *Conn) read() {
	for {
		var message Message

		if err := wsjson.Read(c.ctx, c.conn, &message); err != nil {
			c.Close(fmt.Errorf("read: %v", err))
			return
		}
		if PrintDebugInfo {
			fmt.Printf("read: Read message %#v\n", message)
		}

		switch message.From {
		case "signalingServer":
			var credential nethernet.Credentials
			if err := json.Unmarshal([]byte(message.Data), &credential); err != nil {
				c.Close(fmt.Errorf("read: %v", err))
				return
			}
			select {
			case c.credentials <- credential:
			default:
				c.Close(fmt.Errorf("read: Should never happened"))
				return
			}
		default:
			var signal nethernet.Signal
			var err error
			if err = signal.UnmarshalText([]byte(message.Data)); err != nil {
				c.Close(fmt.Errorf("read: %v", err))
				return
			}
			if signal.NetworkID, err = strconv.ParseUint(message.From, 10, 64); err != nil {
				c.Close(fmt.Errorf("read: %v", err))
				return
			}
			c.signals <- &signal
		}
	}
}

// write ..
func (c *Conn) write(m Message) error {
	if err := wsjson.Write(c.ctx, c.conn, m); err != nil {
		c.Close(fmt.Errorf("write: %v", err))
		return fmt.Errorf("write: %v", err)
	}
	return nil
}

// ping ..
func (c *Conn) ping() {
	ticker := time.NewTicker(time.Second * 50)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if c.write(Message{Type: MessageTypeClientRequestPing}) != nil {
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// refreshCredentials ..
func (c *Conn) refreshCredentials() (err error) {
	if PrintDebugInfo || PrintRefreshInfo {
		fmt.Printf(
			"[%s] refreshCredentials: Start refresh\n",
			time.Now().Format("2006-01-02 15:04:05"),
		)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.credentials:
	default:
		err = fmt.Errorf("refreshCredentials: Should never happened")
		c.Close(err)
		return
	}

	c.write(Message{Type: MessageTypeClientRequestCredential})
	ctx, cancel := context.WithTimeout(c.ctx, time.Second*30)
	defer cancel()

	select {
	case credential := <-c.credentials:
		c.credentials <- credential
		return nil
	case <-ctx.Done():
		err = fmt.Errorf("refreshCredentials: Refresh timeout")
		c.Close(err)
		return
	case <-c.ctx.Done():
		return fmt.Errorf("refreshCredentials: %v", c.ctx.Err())
	}
}

// autoRefresh ..
func (c *Conn) autoRefresh(refreshTime time.Duration) {
	if refreshTime == RefreshTimeDisbale {
		return
	}

	ticker := time.NewTicker(refreshTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if c.refreshCredentials() != nil {
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// Signal sends a Signal to a remote network referenced by [Signal.NetworkID].
func (c *Conn) Signal(signal *nethernet.Signal) error {
	err := c.write(Message{
		Type: MessageTypeClientSendSignal,
		To:   json.Number(strconv.FormatUint(signal.NetworkID, 10)),
		Data: signal.String(),
	})
	if err != nil {
		return fmt.Errorf("Signal: %v", err)
	}
	return nil
}

// Notify registers a Notifier to receive notifications for signals and errors. It returns
// a function to stop receiving notifications on Notifier. Once the stopping function is called,
// ErrSignalingStopped will be notified to the Notifier, and the underlying negotiator should
// handle the error by closing or returning.
func (c *Conn) Notify(n nethernet.Notifier) (stop func()) {
	go func() {
		for {
			select {
			case signal := <-c.signals:
				n.NotifySignal(signal)
			case <-c.ctx.Done():
				n.NotifyError(nethernet.ErrSignalingStopped)
				return
			}
		}
	}()
	return func() {
		c.Close(fmt.Errorf("Notify: Use of closed network connection"))
	}
}

// Credentials blocks until Credentials are received by Signaling, and returns them. If Signaling
// does not support returning Credentials, it will return nil. Credentials are typically received
// from a WebSocket connection. The [context.Context] may be used to cancel the blocking.
func (c *Conn) Credentials(ctx context.Context) (*nethernet.Credentials, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case credential := <-c.credentials:
		c.credentials <- credential
		return &credential, nil
	case <-c.ctx.Done():
		return nil, fmt.Errorf("Credentials: %v", c.ctx.Err())
	}
}

// NetworkID returns the local network ID of Signaling. It is used by Listener to obtain its local
// network ID.
func (c *Conn) NetworkID() uint64 {
	return c.dialer.ClientNetherNetID
}

// PongData ..
func (c *Conn) PongData(d []byte) {}

// Close ..
func (c *Conn) Close(err error) {
	c.doOnce.Do(func() {
		_ = c.conn.Close(websocket.StatusNormalClosure, "")
		c.cancel(fmt.Errorf("Close: %v", err))
	})
}
