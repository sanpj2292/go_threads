package main

import (
	"image/color"
	"log"
	"sync"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth, screenHeight = 640, 320
	boidCount                 = 500
	viewRadius                = 13    // No of pixels that each boid can see
	adjRate                   = 0.015 // How jerky or smooth our simulation is
)

var (
	green   = color.RGBA{10, 255, 50, 255}
	boids   [boidCount]*Boid
	boidMap [screenWidth + 1][screenHeight + 1]int
	lock    = sync.RWMutex{}
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, boid := range boids {
		xpos := boid.Position.x
		ypos := boid.Position.y

		screen.Set(int(xpos+1), int(ypos), green)
		screen.Set(int(xpos-1), int(ypos), green)
		screen.Set(int(xpos), int(ypos+1), green)
		screen.Set(int(xpos), int(ypos-1), green)
	}
}

func (g *Game) Layout(_, _ int) (w, h int) {
	return screenWidth, screenHeight
}

func main() {
	// pre-filling
	for i, row := range boidMap {
		for j := range row {
			boidMap[i][j] = -1
		}
	}
	for i := 0; i < boidCount; i++ {
		createBoid(i)
	}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Boid simulation")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
