package web

import (
	"cardgame/deck"

	"github.com/gin-gonic/gin"
)

func GetDecks(c *gin.Context) {
	d := deck.Decks()
	c.JSON(200, gin.H{"decks": d})
}

func GetDeck(c *gin.Context) {
	id := c.Param("id")
	if d, ok := deck.Decks()[id]; ok {
		c.Header("Cache-Control", "public, max-age=31536000") // 1 year
		c.JSON(200, gin.H{"deck": d})
	} else {
		c.JSON(404, gin.H{"error": "deck not found"})
	}
}
