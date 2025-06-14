package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// drawWorld renders the world grid
func (g *Game) drawWorld(screen *ebiten.Image) {
	// First draw grass (bottom layer)
	for pos, grass := range g.world.Grass {
		g.drawGrass(screen, pos, grass.Amount)
	}
	
	// Then draw animals on top
	for _, rabbit := range g.world.Rabbits {
		g.drawRabbit(screen, rabbit.Animal.Position)
	}
	
	for _, fox := range g.world.Foxes {
		g.drawFox(screen, fox.Animal.Position)
	}
}

// drawGrass draws grass at given position with intensity based on amount
func (g *Game) drawGrass(screen *ebiten.Image, pos Position, amount int) {
	x := pos.X * cellSize
	y := pos.Y * cellSize
	
	// Grass color intensity based on amount (0-100)
	intensity := uint8(50 + (amount * 205 / 100)) // 50-255 range
	grassColor := color.RGBA{0, intensity, 0, 255}
	
	// Draw grass cell
	g.fillRect(screen, x, y, cellSize, cellSize, grassColor)
}

// drawRabbit draws a rabbit (smaller white dot so grass is visible underneath)
func (g *Game) drawRabbit(screen *ebiten.Image, pos Position) {
	x := pos.X * cellSize
	y := pos.Y * cellSize
	
	// Find the rabbit at this position to check if it's newborn
	var rabbit *Rabbit
	for _, r := range g.world.Rabbits {
		if r.Animal.Position.X == pos.X && r.Animal.Position.Y == pos.Y {
			rabbit = r
			break
		}
	}
	
	var rabbitColor color.RGBA
	if rabbit != nil && rabbit.NewBorn > 0 {
		// Newborn rabbits are yellow for visibility
		rabbitColor = color.RGBA{255, 255, 0, 255} // Yellow
	} else {
		// Normal rabbits are white
		rabbitColor = color.RGBA{255, 255, 255, 255} // White
	}
	
	// Smaller rabbit so we can see grass underneath
	g.fillRect(screen, x+3, y+3, cellSize-6, cellSize-6, rabbitColor)
}

// drawFox draws a fox (red dot)
func (g *Game) drawFox(screen *ebiten.Image, pos Position) {
	x := pos.X * cellSize
	y := pos.Y * cellSize
	
	foxColor := color.RGBA{255, 0, 0, 255} // Red
	g.fillRect(screen, x+1, y+1, cellSize-2, cellSize-2, foxColor)
}

// fillRect fills a rectangle with given color
func (g *Game) fillRect(screen *ebiten.Image, x, y, width, height int, c color.Color) {
	// Create a small image and fill it
	rect := ebiten.NewImage(width, height)
	rect.Fill(c)
	
	// Draw options
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	
	screen.DrawImage(rect, op)
}