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
	
	g.world.Tick++
	return nil
}

// addTestEntities adds some test entities for visualization
func (g *Game) addTestEntities() {
	// Add some grass patches
	for i := 0; i < 20; i++ {
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		pos := Position{x, y}
		
		g.world.Grass[pos] = &Grass{
			Position: pos,
			Amount:   rand.Intn(101), // 0-100
		}
		g.world.Grid[x][y] = GrassType
	}
	
	// Add some rabbits
	for i := 0; i < 5; i++ {
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		
		rabbit := &Rabbit{
			Position: Position{x, y},
			Energy:   50,
			ReproduceCD: 0,
			Age:      0,
		}
		
		g.world.Rabbits = append(g.world.Rabbits, rabbit)
		g.world.Grid[x][y] = RabbitType
	}
	
	// Add some foxes
	for i := 0; i < 2; i++ {
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		
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
		debugText += fmt.Sprintf("Foxes: %d", len(g.world.Foxes))
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