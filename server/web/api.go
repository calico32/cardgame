package web

import (
	"time"

	"github.com/gin-gonic/gin"
)

func InitApi(e *gin.RouterGroup) *gin.RouterGroup {
	e.GET("/rooms", GetRooms)
	e.GET("/room/:room", GetRoom)
	e.POST("/room", CreateRoom)
	e.GET("/ws/:room", ServeWS)

	e.GET("/decks", GetDecks)
	e.GET("/deck/:id", GetDeck)

	e.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"time": time.Now().Format(time.RFC3339),
		})
	})

	return e
}
