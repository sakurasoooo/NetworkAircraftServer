package game

type PlayerRocket struct {
	Health   int
	Attack   int
	Position Vec2
	UUID     int
	Parent   int
	NextMove Vec2
}

func NewPlayerRocket(position Vec2, parent int) *PlayerRocket {
	return &PlayerRocket{
		Health:   100,
		Attack:   10,
		Position: position,
		UUID:     generateUUID(),
		Parent:   parent,
		NextMove: Vec2{X: 0, Y: 0},
	}
}
