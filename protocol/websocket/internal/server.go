package internal

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/zzy-rabbit/bp/protocol/websocket/api"
	logApi "github.com/zzy-rabbit/bp/tool/log/api"
	"net/http"
	"sync"
	"time"
)

type server struct {
	ILogger logApi.IPlugin `xplugin:"bp.tool.log"`
	mutex   sync.Mutex
	conns   map[string]api.IConn
	upgrade *websocket.Upgrader
	mux     *http.ServeMux
	httpSvr *http.Server
	service *service
}

func (s *service) NewServer(ctx context.Context, addr string) api.IServer {
	svr := &server{
		ILogger: s.ILogger,
		conns:   make(map[string]api.IConn),
		upgrade: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		mux:     http.NewServeMux(),
		service: s,
	}
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
	return svr
}

func (s *server) Handler(ctx context.Context, url string, callback api.OnConnCallbackFunc) {
	s.mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
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
		callback(ctx, c, req)
	})
}

//func (s *server) handler(w http.ResponseWriter, r *http.Request) {
//	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
//	if err != nil {
//		return
//	}
//	req := api.Request{
//		Headers: r.Header,
//	}
//	ctx := context.Background()
//	c := s.service.NewConnection(ctx, conn)
//	s.conns[c.RemoteAddr(ctx).String()] = c
//	s.callback(ctx, c, req)
//}
