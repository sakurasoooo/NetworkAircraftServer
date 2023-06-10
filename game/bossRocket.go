package game

// bossRocket.go
type BossRocket struct {
	Health   int
	Attack   int
	Position Vec2
	UUID     int
	Parent   int
	NextMove Vec2
	Target   int
}

func NewBossRocket(position Vec2, parent int, target int) *BossRocket {
	return &BossRocket{
		Health:   100,
		Attack:   20,
		Position: position,
		UUID:     generateUUID(),
		Parent:   parent,
		NextMove: Vec2{X: 0, Y: 0},
		Target:   target,
	}
}
