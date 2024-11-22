package sessionmanager

import (
	"github.com/esonhugh/proxyinbrowser/cmd/server/define"
	"sync"
)

type SafeWebsocketConnMap struct {
	mapper sync.Map
}

func (s *SafeWebsocketConnMap) Get(key string) *define.WebsocketClient {
	if v, e := s.mapper.Load(key); e {
		if rv, ok := v.(*define.WebsocketClient); ok {
			return rv
		}
	}
	return nil
}

func (s *SafeWebsocketConnMap) Set(key string, conn *define.WebsocketClient) {
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
