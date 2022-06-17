package card

import "fmt"

var idCounter = -1

// NextId returns an incrementing id with the given prefix, like "c10", "w234", etc.
func NextId(prefix string) string {
	idCounter++
	return fmt.Sprintf("%s%d", prefix, idCounter)
}
