package web

import (
	"cardgame/game"

	"github.com/gin-gonic/gin"
)

func GetRooms(c *gin.Context) {
	rooms := []*game.Room{}
	for _, r := range game.Rooms {
		if r.IsPrivate() {
			continue
		}
		rooms = append(rooms, r)
	}

	c.JSON(200, gin.H{
		"rooms": rooms,
		"count": gin.H{
			"public":  len(rooms),
			"private": len(game.Rooms) - len(rooms),
			"total":   len(game.Rooms),
		},
	})
}

func GetRoom(c *gin.Context) {
	id := c.Param("room")
	password := c.Request.Header.Get("X-Password")
	r, ok := game.Rooms[id]
	if !ok || (r.IsPrivate() && password == "") {
		c.AbortWithStatusJSON(404, gin.H{"error": "room not found"})
		return
	}
	if r.IsPrivate() && !r.TryPassword(password) {
		c.AbortWithStatusJSON(403, gin.H{"error": "invalid password"})
		return
	}

	c.JSON(200, gin.H{"room": r})
}

func CreateRoom(c *gin.Context) {
	password := c.Request.Header.Get("X-Password")
	r := game.NewRoom(password)

	c.JSON(200, gin.H{"room": r})
}
