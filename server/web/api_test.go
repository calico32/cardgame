package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func initTestApi(t *testing.T) *gin.Engine {
	t.Helper()
	e := gin.New()
	InitApi(e.Group("/api"))
	return e
}

func TestApiTime(t *testing.T) {
	api := initTestApi(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api", nil)
	api.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	type response struct {
		Time string `json:"time"`
	}
	var r response
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
	assert.NotEmpty(t, r.Time)
	// time matches RFC3339 format (from https://regex101.com/r/qH0sU7/1)
	assert.Regexp(t, `^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)$`, r.Time)
}
