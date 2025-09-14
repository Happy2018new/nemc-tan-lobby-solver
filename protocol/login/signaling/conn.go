package signaling

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Happy2018new/nemc-tan-lobby-solver/core/nethernet"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const EnableDebug = false

type Conn struct {
	conn *websocket.Conn
	ctx  context.Context
	d    Dialer

	credentials atomic.Pointer[nethernet.Credentials]
	ready       chan struct{}

	once sync.Once

	signals chan *nethernet.Signal
}

func (c *Conn) Signal(signal *nethernet.Signal) error {
	return c.write(Message{
		Type: MessageTypeClientSendSignal,
		To:   json.Number(strconv.FormatUint(signal.NetworkID, 10)),
		Data: signal.String(),
	})
}

func (c *Conn) ReadSignal() (*nethernet.Signal, error) {
	select {
	case s := <-c.signals:
		return s, nil
	case <-c.ctx.Done():
		return nil, context.Cause(c.ctx)
	}
}

func (c *Conn) Credentials(ctx context.Context) (*nethernet.Credentials, error) {
	select {
	case <-c.ctx.Done():
		return nil, context.Cause(c.ctx)
	case <-c.ready:
		return c.credentials.Load(), nil
	}
}

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

func (c *Conn) read(cancel context.CancelCauseFunc) {
	for {
		var message Message

		err := wsjson.Read(context.Background(), c.conn, &message)
		if err != nil {
			cancel(err)
			return
		}
		if EnableDebug {
			fmt.Printf("read: Read message %#v\n", message)
		}

		switch message.From {
		case "signalingServer":
			var credentials nethernet.Credentials
			if err := json.Unmarshal([]byte(message.Data), &credentials); err != nil {
				continue
			}
			c.credentials.Store(&credentials)
			close(c.ready)
		default:
			s := &nethernet.Signal{}
			if err := s.UnmarshalText([]byte(message.Data)); err != nil {
				continue
			}
			var err error
			s.NetworkID, err = strconv.ParseUint(message.From, 10, 64)
			if err != nil {
				continue
			}
			c.signals <- s
		}
	}
}

func (c *Conn) write(m Message) error {
	return wsjson.Write(context.Background(), c.conn, m)
}

func (c *Conn) Close() (err error) {
	c.once.Do(func() {
		err = c.conn.Close(websocket.StatusNormalClosure, "")
	})
	return err
}

func (c *Conn) NetworkID() uint64 {
	return c.d.NetworkID
}

func (c *Conn) Notify(n nethernet.Notifier) (stop func()) {
	go func() {
		for {
			c, err := c.ReadSignal()
			if err != nil {
				n.NotifyError(nethernet.ErrSignalingStopped)
				return
			}
			n.NotifySignal(c)
		}
	}()
	return func() {
		c.Close()
	}
}

func (c *Conn) PongData(d []byte) {}
