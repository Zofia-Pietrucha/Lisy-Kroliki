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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game represents the main game state
type Game struct {
	world           *World
	slowMode        bool // Slower simulation for better observation
	tickCounter     int  // Count frames to slow down updates to 6 FPS
	paused          bool // Simulation pause state
	populationHistory []PopulationData // Historical population data for graphing
	recordCounter   int  // Counter for recording data points
}

// Update proceeds the game state.
// Called 60 times per second, but world updates only every 10 frames (6 FPS).
func (g *Game) Update() error {
	if g.world == nil {
		g.world = NewWorld()
		rand.Seed(time.Now().UnixNano())
		
		// Add initial entities to start the simulation
		g.world.addTestEntities()
		
		// Initialize population history
		g.populationHistory = make([]PopulationData, 0, maxHistoryPoints)
		g.recordPopulationData()
		
		log.Println("World initialized with test entities")
	}
	
	// Handle input
	g.handleInput()
	
	// Only update world if not paused
	if !g.paused {
		// Slow down simulation - update world only every 10 frames (6 FPS instead of 60)
		// This makes the simulation easier to observe and follow
		g.tickCounter++
		if g.tickCounter >= 10 {
			g.tickCounter = 0
			g.world.Update()
			g.world.Tick++
			
			// Record population data every 30 ticks (once per 5 seconds) instead of every second
			g.recordCounter++
			if g.recordCounter >= 30 {
				g.recordCounter = 0
				g.recordPopulationData()
			}
		}
	}
	
	return nil
}

// handleInput processes user input for pause/play controls
func (g *Game) handleInput() {
	// Space bar to toggle pause
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.paused = !g.paused
		if g.paused {
			log.Println("Simulation paused")
		} else {
			log.Println("Simulation resumed")
		}
	}
	
	// Mouse clicks for buttons
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.handleButtonClick(x, y)
	}
}

// handleButtonClick handles clicks on pause/play buttons
func (g *Game) handleButtonClick(x, y int) {
	// Pause button (520, 10, 80, 30)
	if x >= 520 && x <= 600 && y >= 10 && y <= 40 {
		g.paused = true
		log.Println("Simulation paused")
	}
	
	// Play button (610, 10, 80, 30)
	if x >= 610 && x <= 690 && y >= 10 && y <= 40 {
		g.paused = false
		log.Println("Simulation resumed")
	}
	
	// Reset button (700, 10, 80, 30)
	if x >= 700 && x <= 780 && y >= 10 && y <= 40 {
		g.world = NewWorld()
		g.world.addTestEntities()
		g.populationHistory = make([]PopulationData, 0, maxHistoryPoints)
		g.recordPopulationData()
		g.paused = false
		log.Println("Simulation reset")
	}
}

// recordPopulationData adds current population counts to history
func (g *Game) recordPopulationData() {
	if g.world == nil {
		return
	}
	
	data := PopulationData{
		Tick:    g.world.Tick,
		Rabbits: len(g.world.Rabbits),
		Foxes:   len(g.world.Foxes),
		Grass:   len(g.world.Grass),
	}
	
	g.populationHistory = append(g.populationHistory, data)
	
	// Keep only last maxHistoryPoints
	if len(g.populationHistory) > maxHistoryPoints {
		g.populationHistory = g.populationHistory[1:]
	}
}

// Draw renders the game screen.
// Called every frame (60 FPS) for smooth visual updates.
func (g *Game) Draw(screen *ebiten.Image) {
	// Fill background with black (empty ground)
	screen.Fill(color.RGBA{0, 0, 0, 255})
	
	// Draw simulation area (top part)
	if g.world != nil {
		g.drawWorld(screen)
	}
	
	// Draw control buttons
	g.drawControlButtons(screen)
	
	// Draw population graph (bottom part) - FIX: Actually call the function!
	if g.world != nil {
		g.drawPopulationGraph(screen)
	}
	
	// Display simulation statistics (top-left corner)
	debugText := "Ecosystem Simulation"
	if g.paused {
		debugText += " (PAUSED)"
	}
	debugText += "\n"
	
	if g.world != nil {
		debugText += fmt.Sprintf("Tick: %d\n", g.world.Tick)
		debugText += fmt.Sprintf("Grass: %d\n", len(g.world.Grass))
		
		// Rabbit population with limit warning
		rabbitCount := len(g.world.Rabbits)
		debugText += fmt.Sprintf("Rabbits: %d", rabbitCount)
		if rabbitCount >= maxRabbits {
			debugText += " (MAX!)"
		}
		debugText += "\n"
		
		// Fox population with limit warning and energy info
		foxCount := len(g.world.Foxes)
		debugText += fmt.Sprintf("Foxes: %d", foxCount)
		if foxCount >= maxFoxes {
			debugText += " (MAX!)"
		}
		if foxCount == 0 {
			debugText += " (EXTINCT!)"
		}
		debugText += "\n"
		
		// Show average energies to understand why animals might be dying
		if rabbitCount > 0 {
			totalRabbitEnergy := 0
			for _, r := range g.world.Rabbits {
				totalRabbitEnergy += r.Animal.Energy
			}
			avgRabbitEnergy := totalRabbitEnergy / rabbitCount
			debugText += fmt.Sprintf("Avg Rabbit Energy: %d\n", avgRabbitEnergy)
		}
		
		if foxCount > 0 {
			totalFoxEnergy := 0
			for _, f := range g.world.Foxes {
				totalFoxEnergy += f.Animal.Energy
			}
			avgFoxEnergy := totalFoxEnergy / foxCount
			debugText += fmt.Sprintf("Avg Fox Energy: %d\n", avgFoxEnergy)
		}
		
		debugText += "Controls: SPACE=Pause"
	}
	
	ebitenutil.DebugPrint(screen, debugText)
	
	// Draw button labels using DebugPrint
	ebitenutil.DebugPrintAt(screen, "PAUSE", 535, 20)
	ebitenutil.DebugPrintAt(screen, "PLAY", 630, 20)
	ebitenutil.DebugPrintAt(screen, "RESET", 715, 20)
	
	// Draw simple graph legend
	ebitenutil.DebugPrintAt(screen, "■=Rabbits ♦=Foxes +=Grass", 30, 580)
}

// drawControlButtons draws pause/play/reset buttons
func (g *Game) drawControlButtons(screen *ebiten.Image) {
	// Pause button
	if g.paused {
		g.fillRect(screen, 520, 10, 80, 30, color.RGBA{100, 100, 100, 255})
	} else {
		g.fillRect(screen, 520, 10, 80, 30, color.RGBA{200, 100, 100, 255})
	}
	
	// Play button
	if !g.paused {
		g.fillRect(screen, 610, 10, 80, 30, color.RGBA{100, 100, 100, 255})
	} else {
		g.fillRect(screen, 610, 10, 80, 30, color.RGBA{100, 200, 100, 255})
	}
	
	// Reset button
	g.fillRect(screen, 700, 10, 80, 30, color.RGBA{100, 100, 200, 255})
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