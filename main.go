// Package main implements an ecosystem simulation with grass, rabbits, and foxes.
// The simulation demonstrates predator-prey dynamics and natural population cycles.
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
	world       *World
	slowMode    bool // Slower simulation for better observation
	tickCounter int  // Count frames to slow down updates to 6 FPS
}

// Update proceeds the game state.
// Called 60 times per second, but world updates only every 10 frames (6 FPS).
func (g *Game) Update() error {
	if g.world == nil {
		g.world = NewWorld()
		rand.Seed(time.Now().UnixNano())
		
		// Add initial entities to start the simulation
		g.world.addTestEntities()
		
		log.Println("World initialized with test entities")
	}
	
	// Slow down simulation - update world only every 10 frames (6 FPS instead of 60)
	// This makes the simulation easier to observe and follow
	g.tickCounter++
	if g.tickCounter >= 10 {
		g.tickCounter = 0
		g.world.Update()
		g.world.Tick++
	}
	
	return nil
}

// Draw renders the game screen.
// Called every frame (60 FPS) for smooth visual updates.
func (g *Game) Draw(screen *ebiten.Image) {
	// Fill background with black (empty ground)
	screen.Fill(color.RGBA{0, 0, 0, 255})
	
	// Render the world grid with all entities
	if g.world != nil {
		g.drawWorld(screen)
	}
	
	// Display simulation statistics
	debugText := "Ecosystem Simulation\n"
	if g.world != nil {
		debugText += fmt.Sprintf("Tick: %d\n", g.world.Tick)
		debugText += fmt.Sprintf("Grass patches: %d\n", len(g.world.Grass))
		
		// Rabbit population with limit warning
		rabbitCount := len(g.world.Rabbits)
		debugText += fmt.Sprintf("Rabbits: %d", rabbitCount)
		if rabbitCount >= maxRabbits {
			debugText += " (MAX!)"
		}
		debugText += "\n"
		
		// Fox population with limit warning
		foxCount := len(g.world.Foxes)
		debugText += fmt.Sprintf("Foxes: %d", foxCount)
		if foxCount >= maxFoxes {
			debugText += " (MAX!)"
		}
		debugText += "\n"
		
		// Detailed rabbit statistics
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
			debugText += fmt.Sprintf("Rabbit Avg Energy: %d\n", avgEnergy)
			debugText += fmt.Sprintf("Rabbits ready to breed: %d\n", readyToReproduce)
		}
		
		// Fox energy statistics
		if len(g.world.Foxes) > 0 {
			totalFoxEnergy := 0
			for _, f := range g.world.Foxes {
				totalFoxEnergy += f.Animal.Energy
			}
			avgFoxEnergy := totalFoxEnergy / len(g.world.Foxes)
			debugText += fmt.Sprintf("Fox Avg Energy: %d", avgFoxEnergy)
		}
	}
	
	ebitenutil.DebugPrint(screen, debugText)
}

// Layout defines the logical screen size.
// Returns the game's internal resolution regardless of window size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// main initializes and runs the ecosystem simulation.
func main() {
	log.Println("Starting Ecosystem Simulation...")
	
	// Configure the game window
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ecosystem Simulation - Grass, Rabbits, and Foxes")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	
	// Create and run the game
	game := &Game{}
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}