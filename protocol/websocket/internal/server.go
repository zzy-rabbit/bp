package internal

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/zzy-rabbit/bp/protocol/websocket/api"
	"net/http"
	"sync"
	"time"
)

type server struct {
	mutex    sync.Mutex
	conns    map[string]api.IConn
	callback api.OnConnCallbackFunc
	upgrade  *websocket.Upgrader
	mux      *http.ServeMux
	httpSvr  *http.Server
	service  *service
}

func (s *service) NewServer(ctx context.Context, addr string, callback api.OnConnCallbackFunc) api.IServer {
	svr := &server{
		conns:    make(map[string]api.IConn),
		callback: callback,
		upgrade: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		mux:     http.NewServeMux(),
		service: s,
	}
	svr.mux.HandleFunc("/", svr.handler)
	svr.httpSvr = &http.Server{Addr: addr, Handler: svr.mux}
	go func() {
		for {
			err := svr.httpSvr.ListenAndServe()
			if err != nil {
				time.Sleep(time.Second * 3)
				continue
			}
		}
	}()
	return s
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		return
	}
	req := api.Request{
		Headers: r.Header,
	}
	ctx := context.Background()
	c := s.service.NewConnection(ctx, conn)
	s.conns[c.RemoteAddr(ctx).String()] = c
	s.callback(ctx, c, req)
}
