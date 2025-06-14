package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// drawWorld renders the world grid
func (g *Game) drawWorld(screen *ebiten.Image) {
	// Simplified rendering - no clipping, just avoid drawing below line 400
	
	// First draw grass (bottom layer)
	for pos, grass := range g.world.Grass {
		if pos.Y * cellSize < gameAreaHeight { // Simple check
			g.drawGrass(screen, pos, grass.Amount)
		}
	}
	
	// Then draw animals on top
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

// drawGrassInArea draws grass in the simulation area only
func (g *Game) drawGrassInArea(screen *ebiten.Image, pos Position, amount int) {
	x := pos.X * cellSize
	y := pos.Y * cellSize
	
	// Only draw if within simulation area
	if y >= gameAreaHeight {
		return
	}
	
	// Grass color intensity based on amount (0-100)
	intensity := uint8(50 + (amount * 205 / 100)) // 50-255 range
	grassColor := color.RGBA{0, intensity, 0, 255}
	
	// Draw grass cell
	g.fillRectInArea(screen, x, y, cellSize, cellSize, grassColor)
}

// drawRabbitInArea draws a rabbit in the simulation area only
func (g *Game) drawRabbitInArea(screen *ebiten.Image, pos Position) {
	x := pos.X * cellSize
	y := pos.Y * cellSize
	
	// Only draw if within simulation area
	if y >= gameAreaHeight {
		return
	}
	
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
	g.fillRectInArea(screen, x+3, y+3, cellSize-6, cellSize-6, rabbitColor)
}

// drawFoxInArea draws a fox in the simulation area only
func (g *Game) drawFoxInArea(screen *ebiten.Image, pos Position) {
	x := pos.X * cellSize
	y := pos.Y * cellSize
	
	// Only draw if within simulation area
	if y >= gameAreaHeight {
		return
	}
	
	foxColor := color.RGBA{255, 0, 0, 255} // Red
	g.fillRectInArea(screen, x+1, y+1, cellSize-2, cellSize-2, foxColor)
}

// fillRectInArea fills a rectangle within a specific area
func (g *Game) fillRectInArea(screen *ebiten.Image, x, y, width, height int, c color.Color) {
	// Create a small image and fill it
	rect := ebiten.NewImage(width, height)
	rect.Fill(c)
	
	// Draw options
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	
	screen.DrawImage(rect, op)
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

// drawLine draws a line between two points (optimized)
func (g *Game) drawLine(screen *ebiten.Image, x1, y1, x2, y2 int, lineColor color.RGBA) {
	// Much simpler line drawing - just draw start and end points with a few points in between
	steps := 5 // Much fewer steps
	if steps <= 0 {
		steps = 1
	}
	
	for i := 0; i <= steps; i++ {
		if steps > 0 {
			x := x1 + ((x2-x1)*i)/steps
			y := y1 + ((y2-y1)*i)/steps
			g.fillRect(screen, x, y, 1, 1, lineColor) // Smaller points
		}
	}
}

// drawPopulationGraph renders the population history graph (simplified)
func (g *Game) drawPopulationGraph(screen *ebiten.Image) {
	if len(g.populationHistory) < 1 {
		return
	}
	
	// Draw simple graph background
	g.fillRect(screen, graphOffsetX, graphOffsetY, graphWidth, graphHeight, color.RGBA{20, 20, 20, 255})
	
	// Draw border
	g.fillRect(screen, graphOffsetX, graphOffsetY, graphWidth, 2, color.RGBA{100, 100, 100, 255}) // Top
	g.fillRect(screen, graphOffsetX, graphOffsetY+graphHeight-2, graphWidth, 2, color.RGBA{100, 100, 100, 255}) // Bottom
	g.fillRect(screen, graphOffsetX, graphOffsetY, 2, graphHeight, color.RGBA{100, 100, 100, 255}) // Left
	g.fillRect(screen, graphOffsetX+graphWidth-2, graphOffsetY, 2, graphHeight, color.RGBA{100, 100, 100, 255}) // Right
	
	// Simplified data - use every point but limit complexity
	dataToUse := g.populationHistory
	if len(dataToUse) > 50 { // Limit to 50 points max
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
	
	// Find max values for scaling (simplified)
	maxValue := 20 // Minimum scale
	for _, data := range dataToUse {
		if data.Rabbits > maxValue {
			maxValue = data.Rabbits
		}
		if data.Foxes > maxValue {
			maxValue = data.Foxes
		}
		if data.Grass/10 > maxValue { // Scale grass way down
			maxValue = data.Grass / 10
		}
	}
	
	// Draw population points
	g.drawPopulationPoints(screen, dataToUse, "rabbits", maxValue, color.RGBA{255, 255, 255, 255})
	g.drawPopulationPoints(screen, dataToUse, "foxes", maxValue, color.RGBA{255, 0, 0, 255})
	g.drawPopulationPoints(screen, dataToUse, "grass", maxValue, color.RGBA{0, 255, 0, 255})
}

// drawPopulationPoints draws points instead of lines (much faster)
func (g *Game) drawPopulationPoints(screen *ebiten.Image, history []PopulationData, populationType string, maxValue int, pointColor color.RGBA) {
	if len(history) < 1 || maxValue <= 0 {
		return
	}
	
	// Draw points for population data with vertical offset to avoid overlap
	var yOffset int
	switch populationType {
	case "rabbits":
		yOffset = 0 // Rabbits at exact position
	case "foxes":
		yOffset = -3 // Foxes slightly above
	case "grass":
		yOffset = 3 // Grass slightly below
	}
	
	for i, data := range history {
		var value int
		switch populationType {
		case "rabbits":
			value = data.Rabbits
		case "foxes":
			value = data.Foxes
		case "grass":
			value = data.Grass / 10 // Scale grass down heavily
		default:
			continue
		}
		
		// Skip if value is 0 (makes it clearer when populations die out)
		if value == 0 {
			continue
		}
		
		// Convert to screen coordinates
		x := graphOffsetX + 5 // Start a bit inside the border
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
		
		// Draw different shapes for different populations
		switch populationType {
		case "rabbits":
			// Square for rabbits
			g.fillRect(screen, x-1, y-1, 3, 3, pointColor)
		case "foxes":
			// Diamond for foxes (4 small squares in diamond pattern)
			g.fillRect(screen, x, y-1, 1, 1, pointColor)
			g.fillRect(screen, x-1, y, 1, 1, pointColor)
			g.fillRect(screen, x+1, y, 1, 1, pointColor)
			g.fillRect(screen, x, y+1, 1, 1, pointColor)
		case "grass":
			// Circle-ish for grass (cross pattern)
			g.fillRect(screen, x-1, y, 3, 1, pointColor)
			g.fillRect(screen, x, y-1, 1, 3, pointColor)
		}
	}
}

// drawLegend draws the graph legend (simplified)
func (g *Game) drawLegend(screen *ebiten.Image) {
	// Simple text labels instead of squares
	// This is much faster than drawing rectangles
}