package internal

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/zzy-rabbit/bp/protocol/websocket/api"
	uniformApi "github.com/zzy-rabbit/bp/tool/uniform/api"
	"github.com/zzy-rabbit/xtools/xerror"
	"net"
	"sync"
	"time"
)

type connection struct {
	mutex    sync.RWMutex
	conn     *websocket.Conn
	IUniform uniformApi.IPlugin `xplugin:"bp.tool.uniform"`
	async    *async
	callback api.OnReceiveCallbackFunc
	sendChan chan []byte
}

func (s *service) NewConnection(ctx context.Context, conn *websocket.Conn) api.IConn {
	c := &connection{
		conn:  conn,
		async: NewSync(ctx),
		callback: func(ctx context.Context, frame uniformApi.Frame) {
		},
		IUniform: s.IUniform,
		sendChan: make(chan []byte, 1024),
	}
	c.startReceiveMonitor(ctx)
	c.startSendMonitor(ctx)
	return c
}

func (c *connection) RemoteAddr(ctx context.Context) net.Addr {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.conn.RemoteAddr()
}

func (c *connection) Send(ctx context.Context, frame uniformApi.Frame, timeout time.Duration) (uniformApi.Frame, xerror.IError) {
	err := c.Post(ctx, frame)
	if err != nil {
		return uniformApi.Frame{}, err
	}
	wait := c.async.wait(ctx, frame, timeout)
	err = <-wait.errors
	return wait.response, err
}

func (c *connection) Post(ctx context.Context, frame uniformApi.Frame) xerror.IError {
	content, err := c.IUniform.Marshal(ctx, &frame)
	if err != nil {
		return err
	}
	c.sendChan <- content
	return nil
}

func (c *connection) SetCallback(ctx context.Context, callback api.OnReceiveCallbackFunc) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.callback = callback
	return
}

func (c *connection) startSendMonitor(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case content := <-c.sendChan:
				err := c.conn.WriteMessage(websocket.TextMessage, content)
				if err != nil {
					continue
				}
			}
		}
	}()
}

func (c *connection) startReceiveMonitor(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				typ, message, err := c.conn.ReadMessage()
				if err != nil {
					continue
				}
				switch typ {
				case websocket.TextMessage:
					frame, err := c.IUniform.Unmarshal(ctx, message)
					if err != nil {
						continue
					}
					if exist := c.async.receive(ctx, frame); exist {
						continue
					}
					c.mutex.RLock()
					go c.callback(ctx, frame)
					c.mutex.RUnlock()
					continue
				case websocket.BinaryMessage:
					continue
				case websocket.CloseMessage:
					_ = c.conn.Close()
					return
				case websocket.PingMessage:
					_ = c.conn.WriteControl(websocket.PongMessage, message, time.Now().Add(time.Second))
					continue
				case websocket.PongMessage:
					continue
				default:
					continue
				}
			}
		}
	}()
}
