package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	Position Vector2d
	Velocity Vector2d
	id       int
}

func (b *Boid) CalculateAccelerate() Vector2d {
	upper, lower := b.Position.AddV(viewRadius), b.Position.AddV(-viewRadius)
	avgPosition, avgVel, separation := Vector2d{}, Vector2d{}, Vector2d{}
	count := 0.0

	lock.RLock()
	for i := math.Max(lower.x, 0); i < math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j < math.Min(upper.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				if dist := boids[otherBoidId].Position.Distance(b.Position); dist < viewRadius {
					count++
					avgVel = avgVel.Add(boids[otherBoidId].Velocity)
					avgPosition = avgPosition.Add(boids[otherBoidId].Position)
					separation = separation.Add(b.Position.Subtract(boids[otherBoidId].Position).DivisionV(dist))
				}
			}
		}
	}
	lock.RUnlock()

	accel := Vector2d{x: b.BorderBounce(b.Position.x, screenWidth), y: b.BorderBounce(b.Position.y, screenHeight)}
	if count > 0 {
		avgVel = avgVel.DivisionV(count)
		avgPosition = avgPosition.DivisionV(count)
		accelAlignment := avgVel.Subtract(b.Velocity).MultiplyV(adjRate)
		accelCohesion := avgPosition.Subtract(b.Position).MultiplyV(adjRate)
		accelSep := separation.MultiplyV(adjRate)
		accel = accel.Add(accelAlignment).Add(accelCohesion).Add(accelSep)
	}

	return accel
}

func (b *Boid) BorderBounce(pos, maxBorderPos float64) float64 {
	if pos < viewRadius {
		// When position is close to (0,0)
		return 1 / pos
	} else if pos > (maxBorderPos - viewRadius) {
		// When position is close to (screenWidth, screenHeight)
		return 1 / (pos - maxBorderPos)
	}
	return 0
}

func (b *Boid) MoveOne() {
	acceleration := b.CalculateAccelerate()
	lock.Lock()

	b.Velocity = b.Velocity.Add(acceleration).Limit(-1.0, 1.0)

	boidMap[int(b.Position.x)][int(b.Position.y)] = -1
	b.Position = b.Position.Add(b.Velocity)
	boidMap[int(b.Position.x)][int(b.Position.y)] = b.id

	// next := b.Position.Add(b.Velocity)
	// if next.x >= screenWidth || next.x < 0 {
	// 	b.Velocity = Vector2d{x: -b.Velocity.x, y: b.Velocity.y}
	// }
	// if next.y >= screenHeight || next.y < 0 {
	// 	b.Velocity = Vector2d{x: b.Velocity.x, y: -b.Velocity.y}
	// }

	lock.Unlock()
}

func (b *Boid) Start() {
	for {
		b.MoveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func createBoid(i int) {
	b := Boid{
		Position: Vector2d{rand.Float64() * screenWidth, rand.Float64() * screenHeight},
		Velocity: Vector2d{(rand.Float64() * 2) - 1.0, (rand.Float64() * 2) - 1.0},
		id:       i,
	}
	boids[i] = &b
	boidMap[int(b.Position.x)][int(b.Position.y)] = b.id
	go b.Start()
}
