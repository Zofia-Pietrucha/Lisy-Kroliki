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
		log.Println("World initialized")
	}
	
	g.world.Tick++
	return nil
}

// Draw draws the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Wypełnij tło ciemnym zielonym kolorem
	screen.Fill(color.RGBA{34, 139, 34, 255}) // Forest green
	
	// Draw grid (będziemy rysować w następnym commit)
	
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