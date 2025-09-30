package signaling

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"container/list"

	"github.com/Happy2018new/nemc-tan-lobby-solver/core/nethernet"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const (
	DefaultRefreshRetryTimes = 5
	EnableDebug              = false
	PrintRefreshInfo         = true
)

// Conn ..
type Conn struct {
	globalMutex     *sync.Mutex
	runningReadFunc *atomic.Int32

	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelCauseFunc

	dialer      Dialer
	credentials nethernet.Credentials
	messages    chan Message

	doOnce *sync.Once
}

// NewConn ..
func NewConn(ctx context.Context, conn *websocket.Conn, dialer Dialer) (result *Conn, err error) {
	c := &Conn{
		globalMutex:     new(sync.Mutex),
		runningReadFunc: new(atomic.Int32),
		conn:            conn,
		dialer:          dialer,
		credentials:     nethernet.Credentials{},
		messages:        make(chan Message),
		doOnce:          new(sync.Once),
	}

	err = c.handleReady(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewConn: %v", err)
	}

	c.runningReadFunc.Add(1)
	c.ctx, c.cancel = context.WithCancelCause(context.Background())

	go c.read()
	go c.ping()
	go c.checkConn()
	go c.autoRefresh(dialer.RefreshDuration)

	return c, nil
}

// handleReady ..
func (c *Conn) handleReady(ctx context.Context) error {
	var message Message

	err := wsjson.Read(ctx, c.conn, &message)
	if err != nil {
		_ = c.conn.Close(websocket.StatusNormalClosure, "")
		return fmt.Errorf("handleReady: %v", err)
	}
	if EnableDebug {
		fmt.Printf("handleReady: Read message %#v\n", message)
	}

	if message.From != "signalingServer" {
		_ = c.conn.Close(websocket.StatusNormalClosure, "")
		return fmt.Errorf("handleReady: First message is not credentials info; message = %#v", message)
	}

	err = json.Unmarshal([]byte(message.Data), &c.credentials)
	if err != nil {
		_ = c.conn.Close(websocket.StatusNormalClosure, "")
		return fmt.Errorf("handleReady: %v", err)
	}

	return nil
}

// read ..
func (c *Conn) read() {
	for {
		var message Message

		err := wsjson.Read(c.ctx, c.conn, &message)
		if err != nil {
			c.runningReadFunc.Add(-1)
			return
		}
		if EnableDebug {
			fmt.Printf("read: Read message %#v\n", message)
		}

		c.messages <- message
	}
}

// write ..
func (c *Conn) write(m Message) error {
	c.globalMutex.Lock()
	defer c.globalMutex.Unlock()

	err := wsjson.Write(c.ctx, c.conn, m)
	if err != nil {
		return fmt.Errorf("write: %v", err)
	}

	return nil
}

// ping ..
func (c *Conn) ping() {
	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = c.write(Message{Type: MessageTypeClientRequestPing})
		case <-c.ctx.Done():
			return
		}
	}
}

// checkConn ..
func (c *Conn) checkConn() {
	for {
		c.globalMutex.Lock()
		if c.runningReadFunc.Load() == 0 {
			c.globalMutex.Unlock()
			c.Close()
			return
		}
		c.globalMutex.Unlock()
		time.Sleep(time.Second / 20)
	}
}

// refreshConn ..
func (c *Conn) refreshConn() (err error) {
	var finalAddress string

	for range DefaultRefreshRetryTimes {
		tanLobbyRefreshResp, err := c.dialer.GetRefresh()
		if err != nil {
			continue
		}
		if !tanLobbyRefreshResp.Success {
			c.cancel(fmt.Errorf("refreshConn: %v", tanLobbyRefreshResp.ErrorInfo))
			return fmt.Errorf("refreshConn: %v", tanLobbyRefreshResp.ErrorInfo)
		}
		finalAddress = fmt.Sprintf(
			"ws://%s/%d/%d/%s/%s",
			c.dialer.ServerBaseAddress,
			c.dialer.ClientNetherNetID,
			c.dialer.G79UserUID,
			base64.URLEncoding.EncodeToString(tanLobbyRefreshResp.SignalingSeed),
			base64.URLEncoding.EncodeToString(tanLobbyRefreshResp.SignalingTicket),
		)
		break
	}
	if len(finalAddress) == 0 {
		c.cancel(fmt.Errorf("refreshConn: %v", err))
		return fmt.Errorf("refreshConn: %v", err)
	}

	c.globalMutex.Lock()
	defer c.globalMutex.Unlock()
	if EnableDebug || PrintRefreshInfo {
		fmt.Println(time.Now(), "refreshConn: Start refresh")
	}

	_ = c.conn.Close(websocket.StatusNormalClosure, "")
	for {
		if c.runningReadFunc.Load() == 0 {
			break
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	c.conn, _, err = websocket.Dial(ctx, finalAddress, c.dialer.Options)
	if err != nil {
		c.cancel(fmt.Errorf("refreshConn: %v", err))
		return fmt.Errorf("refreshConn: %v", err)
	}

	c.credentials = nethernet.Credentials{}
	if err = c.handleReady(ctx); err != nil {
		c.cancel(fmt.Errorf("refreshConn: %v", err))
		return fmt.Errorf("refreshConn: %v", err)
	}

	c.runningReadFunc.Add(1)
	go c.read()
	return nil
}

// autoRefresh ..
func (c *Conn) autoRefresh(refreshDuration time.Duration) {
	if refreshDuration == RefreshDurationDisable {
		return
	}

	ticker := time.NewTicker(refreshDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if c.refreshConn() != nil {
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// Signal sends a Signal to a remote network referenced by [Signal.NetworkID].
func (c *Conn) Signal(signal *nethernet.Signal) error {
	return c.write(Message{
		Type: MessageTypeClientSendSignal,
		To:   json.Number(strconv.FormatUint(signal.NetworkID, 10)),
		Data: signal.String(),
	})
}

// Notify registers a Notifier to receive notifications for signals and errors. It returns
// a function to stop receiving notifications on Notifier. Once the stopping function is called,
// ErrSignalingStopped will be notified to the Notifier, and the underlying negotiator should
// handle the error by closing or returning.
func (c *Conn) Notify(n nethernet.Notifier) (stop func()) {
	mu := new(sync.Mutex)
	queue := list.New()

	go func() {
		for {
			select {
			case message := <-c.messages:
				mu.Lock()
				_ = queue.PushBack(message)
				mu.Unlock()
			case <-c.ctx.Done():
				n.NotifyError(nethernet.ErrSignalingStopped)
				return
			}
		}
	}()

	go func() {
		for {
			if queue.Len() > 0 {
				// Read message from queue
				mu.Lock()
				messages := make([]Message, 0, queue.Len())
				for queue.Len() > 0 {
					messages = append(messages, queue.Remove(queue.Front()).(Message))
				}
				mu.Unlock()
				// Handle message
				for _, message := range messages {
					switch message.From {
					case "signalingServer":
						c.globalMutex.Lock()
						_ = json.Unmarshal([]byte(message.Data), &c.credentials)
						c.globalMutex.Unlock()
					default:
						signal := new(nethernet.Signal)
						err := signal.UnmarshalText([]byte(message.Data))
						if err != nil {
							break
						}
						signal.NetworkID, err = strconv.ParseUint(message.From, 10, 64)
						if err != nil {
							break
						}
						n.NotifySignal(signal)
					}
				}
			}
			select {
			case <-c.ctx.Done():
				return
			default:
				time.Sleep(time.Second / 20)
			}
		}
	}()

	return func() {
		c.Close()
	}
}

// Credentials blocks until Credentials are received by Signaling, and returns them. If Signaling
// does not support returning Credentials, it will return nil. Credentials are typically received
// from a WebSocket connection. The [context.Context] may be used to cancel the blocking.
func (c *Conn) Credentials(ctx context.Context) (*nethernet.Credentials, error) {
	c.globalMutex.Lock()
	defer c.globalMutex.Unlock()
	return &c.credentials, nil
}

// NetworkID returns the local network ID of Signaling. It is used by Listener to obtain its local
// network ID.
func (c *Conn) NetworkID() uint64 {
	return c.dialer.ClientNetherNetID
}

// PongData ..
func (c *Conn) PongData(d []byte) {}

// Close ..
func (c *Conn) Close() {
	c.globalMutex.Lock()
	defer c.globalMutex.Unlock()
	c.doOnce.Do(func() {
		_ = c.conn.Close(websocket.StatusNormalClosure, "")
		c.cancel(fmt.Errorf("Close: Use of closed network connection"))
	})
}
