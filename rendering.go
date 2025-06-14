package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) drawWorld(screen *ebiten.Image) {
	// Draw grass (bottom layer)
	for pos, grass := range g.world.Grass {
		if pos.Y * cellSize < gameAreaHeight {
			g.drawGrass(screen, pos, grass.Amount)
		}
	}
	
	// Draw animals on top
	for _, rabbit := range g.world.Rabbits {
		if rabbit.Animal.Position.Y * cellSize < gameAreaHeight {
			g.drawRabbit(screen, rabbit.Animal.Position)
		}
	}
	
	for _, fox := range g.world.Foxes {
		if fox.Animal.Position.Y * cellSize < gameAreaHeight {
			g.drawFox(screen, fox.Animal.Position)
		}
	}
}

func (g *Game) drawGrassInArea(screen *ebiten.Image, pos Position, amount int) {
	x := pos.X * cellSize
	y := pos.Y * cellSize
	
	if y >= gameAreaHeight {
		return
	}
	
	// Grass color intensity based on amount (0-100)
	intensity := uint8(50 + (amount * 205 / 100))
	grassColor := color.RGBA{0, intensity, 0, 255}
	
	g.fillRectInArea(screen, x, y, cellSize, cellSize, grassColor)
}

func (g *Game) drawRabbitInArea(screen *ebiten.Image, pos Position) {
	x := pos.X * cellSize
	y := pos.Y * cellSize
	
	if y >= gameAreaHeight {
		return
	}
	
	var rabbit *Rabbit
	for _, r := range g.world.Rabbits {
		if r.Animal.Position.X == pos.X && r.Animal.Position.Y == pos.Y {
			rabbit = r
			break
		}
	}
	
	var rabbitColor color.RGBA
	if rabbit != nil && rabbit.NewBorn > 0 {
		rabbitColor = color.RGBA{255, 255, 0, 255} // Yellow for newborns
	} else {
		rabbitColor = color.RGBA{255, 255, 255, 255} // White
	}
	
	// Smaller rabbit so we can see grass underneath
	g.fillRectInArea(screen, x+3, y+3, cellSize-6, cellSize-6, rabbitColor)
}

func (g *Game) drawFoxInArea(screen *ebiten.Image, pos Position) {
	x := pos.X * cellSize
	y := pos.Y * cellSize
	
	if y >= gameAreaHeight {
		return
	}
	
	foxColor := color.RGBA{255, 0, 0, 255}
	g.fillRectInArea(screen, x+1, y+1, cellSize-2, cellSize-2, foxColor)
}

func (g *Game) fillRectInArea(screen *ebiten.Image, x, y, width, height int, c color.Color) {
	rect := ebiten.NewImage(width, height)
	rect.Fill(c)
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	
	screen.DrawImage(rect, op)
}

func (g *Game) drawGrass(screen *ebiten.Image, pos Position, amount int) {
	x := pos.X * cellSize
	y := pos.Y * cellSize
	
	// Grass color intensity based on amount (0-100)
	intensity := uint8(50 + (amount * 205 / 100))
	grassColor := color.RGBA{0, intensity, 0, 255}
	
	g.fillRect(screen, x, y, cellSize, cellSize, grassColor)
}

func (g *Game) drawRabbit(screen *ebiten.Image, pos Position) {
	x := pos.X * cellSize
	y := pos.Y * cellSize
	
	var rabbit *Rabbit
	for _, r := range g.world.Rabbits {
		if r.Animal.Position.X == pos.X && r.Animal.Position.Y == pos.Y {
			rabbit = r
			break
		}
	}
	
	var rabbitColor color.RGBA
	if rabbit != nil && rabbit.NewBorn > 0 {
		rabbitColor = color.RGBA{255, 255, 0, 255} // Yellow for newborns
	} else {
		rabbitColor = color.RGBA{255, 255, 255, 255} // White
	}
	
	// Smaller rabbit so we can see grass underneath
	g.fillRect(screen, x+3, y+3, cellSize-6, cellSize-6, rabbitColor)
}

func (g *Game) drawFox(screen *ebiten.Image, pos Position) {
	x := pos.X * cellSize
	y := pos.Y * cellSize
	
	foxColor := color.RGBA{255, 0, 0, 255}
	g.fillRect(screen, x+1, y+1, cellSize-2, cellSize-2, foxColor)
}

func (g *Game) fillRect(screen *ebiten.Image, x, y, width, height int, c color.Color) {
	rect := ebiten.NewImage(width, height)
	rect.Fill(c)
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	
	screen.DrawImage(rect, op)
}

func (g *Game) drawLine(screen *ebiten.Image, x1, y1, x2, y2 int, lineColor color.RGBA) {
	steps := 5
	if steps <= 0 {
		steps = 1
	}
	
	for i := 0; i <= steps; i++ {
		if steps > 0 {
			x := x1 + ((x2-x1)*i)/steps
			y := y1 + ((y2-y1)*i)/steps
			g.fillRect(screen, x, y, 1, 1, lineColor)
		}
	}
}

// drawPopulationGraph renders the population history graph
func (g *Game) drawPopulationGraph(screen *ebiten.Image) {
	if len(g.populationHistory) < 1 {
		return
	}
	
	g.fillRect(screen, graphOffsetX, graphOffsetY, graphWidth, graphHeight, color.RGBA{20, 20, 20, 255})
	
	// Draw border
	g.fillRect(screen, graphOffsetX, graphOffsetY, graphWidth, 2, color.RGBA{100, 100, 100, 255})
	g.fillRect(screen, graphOffsetX, graphOffsetY+graphHeight-2, graphWidth, 2, color.RGBA{100, 100, 100, 255})
	g.fillRect(screen, graphOffsetX, graphOffsetY, 2, graphHeight, color.RGBA{100, 100, 100, 255})
	g.fillRect(screen, graphOffsetX+graphWidth-2, graphOffsetY, 2, graphHeight, color.RGBA{100, 100, 100, 255})
	
	// Simplify data for performance
	dataToUse := g.populationHistory
	if len(dataToUse) > 50 {
		step := len(dataToUse) / 50
		simplified := make([]PopulationData, 0, 50)
		for i := 0; i < len(dataToUse); i += step {
			simplified = append(simplified, dataToUse[i])
		}
		dataToUse = simplified
	}
	
	if len(dataToUse) < 1 {
		return
	}
	
	// Find max values for scaling
	maxValue := 20
	for _, data := range dataToUse {
		if data.Rabbits > maxValue {
			maxValue = data.Rabbits
		}
		if data.Foxes > maxValue {
			maxValue = data.Foxes
		}
		if data.Grass/10 > maxValue {
			maxValue = data.Grass / 10
		}
	}
	
	g.drawPopulationPoints(screen, dataToUse, "rabbits", maxValue, color.RGBA{255, 255, 255, 255})
	g.drawPopulationPoints(screen, dataToUse, "foxes", maxValue, color.RGBA{255, 0, 0, 255})
	g.drawPopulationPoints(screen, dataToUse, "grass", maxValue, color.RGBA{0, 255, 0, 255})
}

func (g *Game) drawPopulationPoints(screen *ebiten.Image, history []PopulationData, populationType string, maxValue int, pointColor color.RGBA) {
	if len(history) < 1 || maxValue <= 0 {
		return
	}
	
	// Vertical offset to avoid overlap
	var yOffset int
	switch populationType {
	case "rabbits":
		yOffset = 0
	case "foxes":
		yOffset = -3
	case "grass":
		yOffset = 3
	}
	
	for i, data := range history {
		var value int
		switch populationType {
		case "rabbits":
			value = data.Rabbits
		case "foxes":
			value = data.Foxes
		case "grass":
			value = data.Grass / 10
		default:
			continue
		}
		
		if value == 0 {
			continue
		}
		
		x := graphOffsetX + 5
		if len(history) > 1 {
			x = graphOffsetX + 5 + ((i * (graphWidth - 10)) / (len(history) - 1))
		}
		
		y := graphOffsetY + graphHeight - 5 - ((value * (graphHeight - 10)) / maxValue) + yOffset
		if y < graphOffsetY + 5 {
			y = graphOffsetY + 5
		}
		if y > graphOffsetY + graphHeight - 5 {
			y = graphOffsetY + graphHeight - 5
		}
		
		// Different shapes for different populations
		switch populationType {
		case "rabbits":
			g.fillRect(screen, x-1, y-1, 3, 3, pointColor)
		case "foxes":
			// Diamond shape
			g.fillRect(screen, x, y-1, 1, 1, pointColor)
			g.fillRect(screen, x-1, y, 1, 1, pointColor)
			g.fillRect(screen, x+1, y, 1, 1, pointColor)
			g.fillRect(screen, x, y+1, 1, 1, pointColor)
		case "grass":
			// Cross shape
			g.fillRect(screen, x-1, y, 3, 1, pointColor)
			g.fillRect(screen, x, y-1, 1, 3, pointColor)
		}
	}
}
