package signaling

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"container/list"

	"github.com/Happy2018new/nemc-tan-lobby-solver/core/nethernet"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const EnableDebug = true

// Conn ..
type Conn struct {
	refreshMutex  *sync.Mutex
	refreshWaiter *sync.WaitGroup
	isRefreshing  bool

	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelCauseFunc

	dialer      Dialer
	credentials nethernet.Credentials
	signals     chan *nethernet.Signal

	doOnce *sync.Once
}

// NewConn ..
func NewConn(ctx context.Context, conn *websocket.Conn, dialer Dialer) (result *Conn, err error) {
	c := &Conn{
		refreshMutex:  new(sync.Mutex),
		refreshWaiter: new(sync.WaitGroup),
		isRefreshing:  false,
		conn:          conn,
		dialer:        dialer,
		credentials:   nethernet.Credentials{},
		signals:       make(chan *nethernet.Signal),
		doOnce:        new(sync.Once),
	}

	err = c.handleReady(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewConn: %v", err)
	}

	c.ctx, c.cancel = context.WithCancelCause(context.Background())
	go c.read()
	go c.ping()
	go c.autoRefresh()

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

		err := wsjson.Read(context.Background(), c.conn, &message)
		if err != nil {
			if c.refreshMutex.TryLock() {
				c.cancel(err)
				c.refreshMutex.Unlock()
			} else if c.isRefreshing {
				c.refreshWaiter.Done()
			} else {
				c.cancel(err)
			}
			return
		}
		if EnableDebug {
			fmt.Printf("read: Read message %#v\n", message)
		}

		switch message.From {
		case "signalingServer":
			c.refreshMutex.Lock()
			_ = json.Unmarshal([]byte(message.Data), &c.credentials)
			c.refreshMutex.Unlock()
		default:
			signal := new(nethernet.Signal)
			if err = signal.UnmarshalText([]byte(message.Data)); err != nil {
				break
			}
			if signal.NetworkID, err = strconv.ParseUint(message.From, 10, 64); err != nil {
				break
			}
			c.signals <- signal
		}
	}
}

// write ..
func (c *Conn) write(m Message) error {
	c.refreshMutex.Lock()
	defer c.refreshMutex.Unlock()
	return wsjson.Write(context.Background(), c.conn, m)
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

// autoRefresh ..
func (c *Conn) autoRefresh() {
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if c.RefreshConn() != nil {
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
	queue := list.New()
	queueMu := new(sync.Mutex)

	go func() {
		for {
			select {
			case signal := <-c.signals:
				queueMu.Lock()
				_ = queue.PushBack(signal)
				queueMu.Unlock()
			case <-c.ctx.Done():
				n.NotifyError(nethernet.ErrSignalingStopped)
				return
			}
		}
	}()

	go func() {
		for {
			if queue.Len() > 0 {
				queueMu.Lock()
				signals := make([]*nethernet.Signal, 0, queue.Len())
				for queue.Len() > 0 {
					signals = append(signals, queue.Remove(queue.Front()).(*nethernet.Signal))
				}
				queueMu.Unlock()
				for _, signal := range signals {
					n.NotifySignal(signal)
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
	c.refreshMutex.Lock()
	defer c.refreshMutex.Unlock()
	return &c.credentials, nil
}

// NetworkID returns the local network ID of Signaling. It is used by Listener to obtain its local
// network ID.
func (c *Conn) NetworkID() uint64 {
	return c.dialer.NetworkID
}

// PongData ..
func (c *Conn) PongData(d []byte) {}

// Close ..
func (c *Conn) Close() {
	c.refreshMutex.Lock()
	defer c.refreshMutex.Unlock()
	c.doOnce.Do(func() {
		_ = c.conn.Close(websocket.StatusNormalClosure, "")
		c.cancel(fmt.Errorf("Close: Normal close"))
	})
}

// RefreshConn ..
func (c *Conn) RefreshConn() (err error) {
	c.refreshMutex.Lock()
	defer func() {
		c.isRefreshing = false
		c.refreshMutex.Unlock()
	}()

	select {
	case <-c.ctx.Done():
		return fmt.Errorf("RefreshConn: Use of closed network connection")
	default:
	}

	c.refreshWaiter.Add(1)
	c.isRefreshing = true
	_ = c.conn.Close(websocket.StatusNormalClosure, "")
	c.refreshWaiter.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	c.conn, _, err = websocket.Dial(ctx, c.dialer.serverAddress, c.dialer.Options)
	if err != nil {
		c.cancel(fmt.Errorf("RefreshConn: %v", err))
		return fmt.Errorf("RefreshConn: %v", err)
	}

	c.credentials = nethernet.Credentials{}
	if err = c.handleReady(ctx); err != nil {
		c.cancel(fmt.Errorf("RefreshConn: %v", err))
		return fmt.Errorf("RefreshConn: %v", err)
	}

	go c.read()
	return nil
}
