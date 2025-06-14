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

// Game represents the main game state
type Game struct {
	world *World
	slowMode bool  // Slower simulation
	tickCounter int // Count frames to slow down updates
}

// Update proceeds the game state.
func (g *Game) Update() error {
	if g.world == nil {
		g.world = NewWorld()
		rand.Seed(time.Now().UnixNano())
		
		// Add some test entities to see rendering
		g.world.addTestEntities()
		
		log.Println("World initialized with test entities")
	}
	
	// Slow down simulation - update world only every 10 frames (6 FPS instead of 60)
	g.tickCounter++
	if g.tickCounter >= 10 {
		g.tickCounter = 0
		g.world.Update()
		g.world.Tick++
	}
	
	return nil
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
		
		// Show rabbit reproduction status
		if len(g.world.Rabbits) > 0 {
			readyToReproduce := 0
			totalEnergy := 0
			for _, r := range g.world.Rabbits {
				totalEnergy += r.Animal.Energy
				if r.Animal.Energy >= reproduceEnergyThreshold && r.Animal.ReproduceCD == 0 {
					readyToReproduce++
				}
			}
			avgEnergy := totalEnergy / len(g.world.Rabbits)
			debugText += fmt.Sprintf("Avg Energy: %d\n", avgEnergy)
			debugText += fmt.Sprintf("Ready to breed: %d", readyToReproduce)
		}
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