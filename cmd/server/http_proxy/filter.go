package http_proxy

import (
	"fmt"
	"net/http"
)

func blockHostFilter(blockedHostnames []string, next http.HandlerFunc) http.HandlerFunc {
	hosts := map[string]struct{}{}
	for _, host := range blockedHostnames {
		hosts[host] = struct{}{}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := hosts[r.URL.Host]; ok {
			msg := fmt.Sprintf("access to %s blocked by proxy", r.URL.Host)
			http.Error(w, msg, http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
