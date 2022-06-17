package web

import (
	"cardgame/game"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	EnableCompression: true,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: validate origin
		return true
	},
}

func upgrade(c *gin.Context) (*websocket.Conn, error) {
	return upgrader.Upgrade(c.Writer, c.Request, nil)
}

func ServeWS(c *gin.Context) {
	roomId := c.Param("room")
	password := c.Query("password") // js doesn't allow custom headers, so we use query
	r, ok := game.Rooms[roomId]
	if !ok || (r.IsPrivate() && password == "") {
		c.JSON(404, gin.H{"error": "room not found"})
		return
	}
	if r.IsPrivate() && !r.TryPassword(password) {
		c.JSON(403, gin.H{"error": "invalid password"})
		return
	}

	conn, err := upgrade(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	game.NewPlayer(conn, r)
}
