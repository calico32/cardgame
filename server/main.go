package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"cardgame/deck"
	"cardgame/game"
	"cardgame/web"
)

var mode string

func main() {
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		godotenv.Load()

		room := game.NewRoom("")
		delete(game.Rooms, room.Id)
		room.Id = "r_debug"
		game.Rooms[room.Id] = room
	}

	deck.InitDecks("./data/decks")

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.SetTrustedProxies([]string{})

	c := cors.DefaultConfig()
	c.AllowAllOrigins = true
	c.AllowHeaders = []string{"Origin", "Content-Type", "X-Password"}
	r.Use(cors.New(c))

	web.InitApi(r.Group("/api"))

	r.Use(func(c *gin.Context) {
		c.AbortWithStatusJSON(404, gin.H{"error": "not found"})
	})

	r.Run()
}
