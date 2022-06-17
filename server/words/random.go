package words

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Words(list []string, n int) []string {
	var out []string
	for len(out) < n {
		idx := rand.Intn(len(list))
		out = append(out, list[idx])
	}
	return out
}
