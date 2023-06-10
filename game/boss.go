package game

import (
	"math"
	"math/rand"
)

type Vec2 struct {
	X float32
	Y float32
}

type Boss struct {
	Health      int
	Attack      int
	Position    Vec2
	UUID        int
	AttackTimer int
	NextMove    Vec2
	NextAttack  int
}

func NewBoss() *Boss {

	return &Boss{
		Health:      100,
		Attack:      10,
		Position:    generateRandomVec2(),
		UUID:        generateUUID(),
		AttackTimer: 25,
		NextMove:    Vec2{X: 0, Y: 0},
		NextAttack:  0,
	}
}

// fucntion to generate a random Vec2 limit in -10 <X <10 , -4.5<Y< 4.5
func generateRandomVec2() Vec2 {
	return Vec2{X: rand.Float32()*20 - 10, Y: rand.Float32()*9 - 4.5}
}

// define next direction vector2
var nextDirection = Vec2{X: -1, Y: 0}

// update the boss position, choose a random direction to move in, not change the direction if the boss is at the edge of the screen
// choose a random direction when the boss is at the edge of the screen, range in -10 <X <10 , -4.5<Y< 4.5
func (b *Boss) Update() {
	// print
	// print Direction
	if b.Position.X < 10 && b.Position.X > -10 && b.Position.Y < 4.5 && b.Position.Y > -4.5 {
		b.Position.X += nextDirection.X * 0.1
		b.Position.Y += nextDirection.Y * 0.1
	} else {
		nextDirection = generateRandomVec2()
		// normalize the vector
		nextDirection.Normalize()
		b.Position.X += nextDirection.X * 0.1
		b.Position.Y += nextDirection.Y * 0.1
	}

	// update the attack timer
	if b.AttackTimer > 0 {
		b.AttackTimer -= 1
	} else {
		b.AttackTimer = 25
		b.NextAttack += 1
	}
}

// Normalize normalizes the vector to unit length
func (v *Vec2) Normalize() {
	length := math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
	v.X /= float32(length)
	v.Y /= float32(length)
}

// calculate the distance between the another vector2
func (v *Vec2) Distance(v2 Vec2) float32 {
	return float32(math.Sqrt(float64((v.X-v2.X)*(v.X-v2.X) + (v.Y-v2.Y)*(v.Y-v2.Y))))
}

// calculate the direction between the another vector2, normalize the vector
func (v *Vec2) Direction(v2 Vec2) Vec2 {
	direction := Vec2{X: v2.X - v.X, Y: v2.Y - v.Y}
	direction.Normalize()
	return direction
}
