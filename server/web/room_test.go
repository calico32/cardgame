package web

import (
	"cardgame/game"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func makePublicRoom(t *testing.T, api *gin.Engine) *game.Room {
	t.Helper()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/room", nil)
	api.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "should be able to create public room")

	type response struct {
		Room *game.Room `json:"room"`
	}
	var r response
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
	assert.NotNil(t, r.Room, "should receive room")

	return r.Room
}

func makePrivateRoom(t *testing.T, api *gin.Engine, password string) *game.Room {
	t.Helper()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/room", nil)
	req.Header.Add("X-Password", password)
	api.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "should be able to create private room")

	type response struct {
		Room *game.Room `json:"room"`
	}
	var r response
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
	assert.NotNil(t, r.Room, "should receive room")

	return r.Room
}

func TestPublicRoom(t *testing.T) {
	api := initTestApi(t)
	rm := makePublicRoom(t, api)

	assert.Contains(t, game.Rooms, rm.Id, "should contain public room")

	type response struct {
		Room  *game.Room `json:"room"`
		Error string     `json:"error"`
	}
	var r response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/room/"+rm.Id, nil)
	api.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "should be able to get public room")
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
	assert.Equal(t, rm.Id, r.Room.Id, "should be able to get correct room")

	delete(game.Rooms, rm.Id)
}

func TestRoomList(t *testing.T) {
	api := initTestApi(t)
	pubRoom := makePublicRoom(t, api)
	privRoom := makePrivateRoom(t, api, "correct horse battery staple")

	type response struct {
		Rooms []*game.Room `json:"rooms"`
	}
	var r response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/rooms", nil)
	api.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "should be able to get public room list")
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
	assert.Contains(t, r.Rooms, pubRoom, "should contain public room")
	assert.NotContains(t, r.Rooms, privRoom, "should not contain private room")

	delete(game.Rooms, pubRoom.Id)
	delete(game.Rooms, privRoom.Id)
}

func TestPrivateRoom(t *testing.T) {
	api := initTestApi(t)
	password := "correct horse battery staple"
	rm := makePrivateRoom(t, api, password)

	assert.Contains(t, game.Rooms, rm.Id, "should contain private room")

	type response struct {
		Error string     `json:"error"`
		Room  *game.Room `json:"room"`
	}
	var r response

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/room/"+rm.Id, nil)
	api.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code, "should not be able to get private room without password")
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/room/"+rm.Id, nil)
	req.Header.Add("X-Password", "wrong password")
	api.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code, "should not be able to get private room with wrong password")
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/room/"+rm.Id, nil)
	req.Header.Add("X-Password", password)
	api.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "should be able to get private room with correct password")
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))

	delete(game.Rooms, rm.Id)
}
