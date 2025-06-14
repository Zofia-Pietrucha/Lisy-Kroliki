package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	gridWidth    = 80
	gridHeight   = 60
	cellSize     = 10
	
	// Grass parameters
	maxGrassAmount    = 100
	grassGrowthRate   = 2    // How much grass grows per tick
	grassSpawnChance  = 0.01 // Chance for new grass to appear on empty cells
	
	// Rabbit parameters
	rabbitMoveChance = 0.3   // Chance rabbit moves each tick
	rabbitEnergyLoss = 1     // Energy lost per tick
)

// Entity types
type EntityType int

const (
	Empty EntityType = iota
	GrassType
	RabbitType
	FoxType
)

// Position represents coordinates on the grid
type Position struct {
	X, Y int
}

// Grass represents grass at a position
type Grass struct {
	Position
	Amount int // 0-100, where 100 is fully grown
}

// Rabbit represents a rabbit entity
type Rabbit struct {
	Position
	Energy       int
	ReproduceCD  int // Cooldown after reproduction
	Age          int
}

// Fox represents a fox entity
type Fox struct {
	Position
	Energy       int
	ReproduceCD  int // Cooldown after reproduction
	Age          int
}

// World represents the game world
type World struct {
	Grid    [][]EntityType // What type of entity is at each position
	Grass   map[Position]*Grass
	Rabbits []*Rabbit
	Foxes   []*Fox
	Tick    int
}

// NewWorld creates a new world
func NewWorld() *World {
	w := &World{
		Grid:    make([][]EntityType, gridWidth),
		Grass:   make(map[Position]*Grass),
		Rabbits: make([]*Rabbit, 0),
		Foxes:   make([]*Fox, 0),
		Tick:    0,
	}
	
	// Initialize grid
	for x := 0; x < gridWidth; x++ {
		w.Grid[x] = make([]EntityType, gridHeight)
	}
	
	return w
}

// Update updates the world state
func (w *World) Update() {
	w.updateGrass()
	w.updateRabbits()
}

// updateGrass handles grass growth and spawning
func (w *World) updateGrass() {
	// Grow existing grass
	for _, grass := range w.Grass {
		if grass.Amount < maxGrassAmount {
			grass.Amount += grassGrowthRate
			if grass.Amount > maxGrassAmount {
				grass.Amount = maxGrassAmount
			}
		}
	}
	
	// Try to spawn new grass on empty cells
	for attempts := 0; attempts < 10; attempts++ {
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		
		// Check if cell is empty
		if w.Grid[x][y] == Empty {
			// Random chance to spawn grass
			if rand.Float64() < grassSpawnChance {
				pos := Position{x, y}
				w.Grass[pos] = &Grass{
					Position: pos,
					Amount:   grassGrowthRate, // Start with small amount
				}
				w.Grid[x][y] = GrassType
			}
		}
	}
}

// updateRabbits handles rabbit movement and aging
func (w *World) updateRabbits() {
	for i := len(w.Rabbits) - 1; i >= 0; i-- {
		rabbit := w.Rabbits[i]
		
		// Age and lose energy
		rabbit.Age++
		rabbit.Energy -= rabbitEnergyLoss
		
		// Move rabbit randomly
		if rand.Float64() < rabbitMoveChance {
			w.moveRabbit(rabbit)
		}
		
		// Check if rabbit dies
		if rabbit.Energy <= 0 {
			w.removeRabbit(i)
		}
	}
}

// moveRabbit moves a rabbit to a random adjacent position
func (w *World) moveRabbit(rabbit *Rabbit) {
	// Clear current position
	w.Grid[rabbit.Position.X][rabbit.Position.Y] = Empty
	
	// Get possible moves (adjacent cells)
	moves := w.getAdjacentPositions(rabbit.Position)
	
	// Filter for empty positions or grass positions
	validMoves := make([]Position, 0)
	for _, pos := range moves {
		cellType := w.Grid[pos.X][pos.Y]
		if cellType == Empty || cellType == GrassType {
			validMoves = append(validMoves, pos)
		}
	}
	
	// Move to random valid position, or stay if no valid moves
	if len(validMoves) > 0 {
		newPos := validMoves[rand.Intn(len(validMoves))]
		rabbit.Position = newPos
	}
	
	// Update grid with new position
	w.Grid[rabbit.Position.X][rabbit.Position.Y] = RabbitType
}

// getAdjacentPositions returns valid adjacent positions (8-directional)
func (w *World) getAdjacentPositions(pos Position) []Position {
	adjacent := make([]Position, 0, 8)
	
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue // Skip current position
			}
			
			newX := pos.X + dx
			newY := pos.Y + dy
			
			// Check bounds
			if newX >= 0 && newX < gridWidth && newY >= 0 && newY < gridHeight {
				adjacent = append(adjacent, Position{newX, newY})
			}
		}
	}
	
	return adjacent
}

// removeRabbit removes a rabbit from the world
func (w *World) removeRabbit(index int) {
	rabbit := w.Rabbits[index]
	
	// Clear grid position
	w.Grid[rabbit.Position.X][rabbit.Position.Y] = Empty
	
	// Remove from slice
	w.Rabbits = append(w.Rabbits[:index], w.Rabbits[index+1:]...)
}

// Game represents the main game state
type Game struct {
	world *World
}

// Update proceeds the game state.
func (g *Game) Update() error {
	if g.world == nil {
		g.world = NewWorld()
		rand.Seed(time.Now().UnixNano())
		
		// Add some test entities to see rendering
		g.addTestEntities()
		
		log.Println("World initialized with test entities")
	}
	
	// Update world every tick
	g.world.Update()
	g.world.Tick++
	return nil
}

// addTestEntities adds some test entities for visualization
func (g *Game) addTestEntities() {
	// Add some grass patches (mniej niż wcześniej)
	for i := 0; i < 10; i++ {
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		pos := Position{x, y}
		
		g.world.Grass[pos] = &Grass{
			Position: pos,
			Amount:   rand.Intn(51) + 25, // 25-75 (średnio rozwiniętą trawę)
		}
		g.world.Grid[x][y] = GrassType
	}
	
	// Add some rabbits
	for i := 0; i < 3; i++ {
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		
		// Make sure position is not occupied
		if g.world.Grid[x][y] == Empty {
			rabbit := &Rabbit{
				Position: Position{x, y},
				Energy:   50,
				ReproduceCD: 0,
				Age:      0,
			}
			
			g.world.Rabbits = append(g.world.Rabbits, rabbit)
			g.world.Grid[x][y] = RabbitType
		}
	}
	
	// Add some foxes
	for i := 0; i < 2; i++ {
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		
		// Make sure position is not occupied
		if g.world.Grid[x][y] == Empty {
			fox := &Fox{
				Position: Position{x, y},
				Energy:   50,
				ReproduceCD: 0,
				Age:      0,
			}
			
			g.world.Foxes = append(g.world.Foxes, fox)
			g.world.Grid[x][y] = FoxType
		}
	}
}

// Draw draws the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Wypełnij tło czarnym kolorem (ziemia bez trawy)
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Black
	
	// Draw world grid
	if g.world != nil {
		g.drawWorld(screen)
	}
	
	// Debug info
	debugText := "Ecosystem Simulation\n"
	if g.world != nil {
		debugText += fmt.Sprintf("Tick: %d\n", g.world.Tick)
		debugText += fmt.Sprintf("Grass patches: %d\n", len(g.world.Grass))
		debugText += fmt.Sprintf("Rabbits: %d\n", len(g.world.Rabbits))
		debugText += fmt.Sprintf("Foxes: %d\n", len(g.world.Foxes))
		
		// Show first rabbit info if any
		if len(g.world.Rabbits) > 0 {
			r := g.world.Rabbits[0]
			debugText += fmt.Sprintf("Rabbit0: E=%d A=%d", r.Energy, r.Age)
		}
	}
	
	ebitenutil.DebugPrint(screen, debugText)
}

// drawWorld renders the world grid
func (g *Game) drawWorld(screen *ebiten.Image) {
	// Draw grass
	for pos, grass := range g.world.Grass {
		g.drawGrass(screen, pos, grass.Amount)
	}
	
	// Draw rabbits
	for _, rabbit := range g.world.Rabbits {
		g.drawRabbit(screen, rabbit.Position)
	}
	
	// Draw foxes
	for _, fox := range g.world.Foxes {
		g.drawFox(screen, fox.Position)
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

// drawRabbit draws a rabbit (white dot)
func (g *Game) drawRabbit(screen *ebiten.Image, pos Position) {
	x := pos.X * cellSize
	y := pos.Y * cellSize
	
	rabbitColor := color.RGBA{255, 255, 255, 255} // White
	g.fillRect(screen, x+2, y+2, cellSize-4, cellSize-4, rabbitColor)
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

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	log.Println("Starting Ecosystem Simulation...")
	
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ecosystem Simulation")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	
	game := &Game{}
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}