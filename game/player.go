package game

import (
	"math/rand"
	"time"
)

type Player struct {
	Health     int
	Attack     int
	Position   Vec2
	UUID       int
	Username   string
	NextMove   Vec2
	NextAttack int
}

func NewPlayer(username string) *Player {
	return &Player{
		Health:     100,
		Attack:     1,
		Position:   generateRandomVec2(),
		UUID:       generateUUID(),
		Username:   username,
		NextMove:   Vec2{X: 0, Y: 0},
		NextAttack: 0,
	}
}

func generateUUID() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Int()
}
