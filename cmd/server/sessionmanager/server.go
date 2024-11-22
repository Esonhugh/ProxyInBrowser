package sessionmanager

import (
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"

	define "github.com/esonhugh/proxyinbrowser/cmd/server/define"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var WebsocketConnMap = SafeWebsocketConnMap{mapper: sync.Map{}}

//go:embed bundle.js
var fileContent string

func RunServer(buffer io.Writer) {
	gin.DefaultWriter = buffer
	gin.DisableConsoleColor()

	router := gin.Default()

	_ = router.SetTrustedProxies(nil) // disable
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024 * 4,
		WriteBufferSize: 1024 * 4,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	router.GET("/cert", func(c *gin.Context) {
		c.FileAttachment("cert/cert.pem", "cert.pem")
	})

	router.GET("/:id/init", func(c *gin.Context) {
		c.Header("Content-Type", "application/javascript")
		content := strings.ReplaceAll(
			strings.ReplaceAll(
				fileContent, "localhost:9999", c.Request.Host),
			"__ENDPOINT__", c.Param("id"))
		c.String(http.StatusOK, content)
	})

	router.Any("/:id/infos", func(c *gin.Context) {
		req, err := httputil.DumpRequest(c.Request, true)
		if err != nil {
			log.Error("Try dump request error", err)
		}
		log.Info(string(req))
	})

	// Handle WebSocket connections
	router.GET("/:id", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			// panic(err)
			log.Errorf("%s, error while Upgrading websocket connection", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		connId := uuid.New().String()
		client := define.NewWebSocketClient(conn)
		WebsocketConnMap.Set(connId, client)
		defer WebsocketConnMap.Delete(connId)

		log.Infof("receive new connection! alloc new session id: %v", connId)
		l := log.WithField("session", connId)

		// init!
		_, p, err := conn.ReadMessage()
		if err != nil {
			l.Errorf("Init read message failed")
		}
		var msg map[string]string
		if err = json.Unmarshal(p, &msg); err != nil {
			l.Errorf("Init message unmarshal failed. reason: %s, data: %v", err.Error(), string(p))
		} else {
			for k, v := range msg {
				l.Infof("%s: %s", k, v)
			}
		}

		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				// panic(err)
				l.Errorf("%s, error while reading message", err.Error())
				c.AbortWithError(http.StatusInternalServerError, err)
				break
			}
			p = define.Decode(p)

			var test define.RelayCommandResp
			if json.Unmarshal(p, &test) == nil {
				client.RelayChan <- test
				continue
			}
			time.Sleep(1000 * time.Millisecond)
		}

	})

	router.Run(":9999")
}
