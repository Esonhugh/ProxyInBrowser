package sessionmanager

import (
	"sync"

	"github.com/gorilla/websocket"
)

type SafeWebsocketConnMap struct {
	mapper sync.Map
}

func (s *SafeWebsocketConnMap) Get(key string) *websocket.Conn {
	if v, e := s.mapper.Load(key); e {
		if rv, ok := v.(*websocket.Conn); ok {
			return rv
		}
	}
	return nil
}

func (s *SafeWebsocketConnMap) Set(key string, conn *websocket.Conn) {
	s.mapper.Store(key, conn)
}

func (s *SafeWebsocketConnMap) Delete(key string) {
	s.mapper.Delete(key)
}

func (s *SafeWebsocketConnMap) List() []string {
	var ret []string
	s.mapper.Range(func(key, value interface{}) bool {
		if rv, ok := key.(string); ok {
			ret = append(ret, rv)
			return true
		} else {
			return false
		}
	})
	return ret
}
