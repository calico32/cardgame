package web

import (
	"cardgame/build"
	"time"

	"github.com/gin-gonic/gin"
)

func InitApi(e *gin.RouterGroup) *gin.RouterGroup {
	e.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"commit":    build.Commit(),
			"version":   build.Version(),
			"branch":    build.Branch(),
			"buildTime": build.Time().Format(time.RFC3339),
		})
	})

	e.GET("/rooms", GetRooms)
	e.GET("/room/:room", GetRoom)
	e.POST("/room", CreateRoom)
	e.GET("/ws/:room", ServeWS)

	e.GET("/decks", GetDecks)
	e.GET("/deck/:id", GetDeck)

	e.GET("/me", GetUser)
	e.POST("/me", CreateUser)
	e.PUT("/me", UpdateUser)
	e.DELETE("/me", DeleteUser)

	e.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"time": time.Now().Format(time.RFC3339),
		})
	})

	return e
}
