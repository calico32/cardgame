package web

import (
	"cardgame/game"

	"github.com/gin-gonic/gin"
)

func GetRooms(c *gin.Context) {
	rooms := []*game.Room{}
	for _, r := range game.HubMain.Rooms {
		if r.IsPrivate() {
			continue
		}
		rooms = append(rooms, r)
	}

	c.JSON(200, gin.H{
		"rooms": rooms,
		"count": gin.H{
			"public":  len(rooms),
			"private": len(game.HubMain.Rooms) - len(rooms),
			"total":   len(game.HubMain.Rooms),
		},
	})
}

func GetRoom(c *gin.Context) {
	id := c.Param("room")
	password := c.Request.Header.Get("X-Password")
	r, ok := game.HubMain.Rooms[id]
	if !ok || (r.IsPrivate() && password == "") {
		c.AbortWithStatusJSON(404, gin.H{"error": "room not found"})
		return
	}
	if r.IsPrivate() && !r.CheckPassword(password) {
		c.AbortWithStatusJSON(403, gin.H{"error": "invalid password"})
		return
	}

	c.JSON(200, gin.H{"room": r})
}

func CreateRoom(c *gin.Context) {
	password := c.Request.Header.Get("X-Password")
	r := game.HubMain.NewRoom(password)

	c.JSON(200, gin.H{"room": r})
}
