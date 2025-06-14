// Package main implements an ecosystem simulation with grass, rabbits, and foxes.
// The simulation demonstrates predator-prey dynamics and natural population cycles.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"strings"
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
	
	// Drawing modes
	drawMode        string // "none", "rabbit", "fox"
	mousePressed    bool   // Track if mouse is being held down
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
		
		// Initialize drawing mode
		g.drawMode = "none"
		
		log.Println("World initialized with test entities")
		log.Println("Use keys: 1=Draw Rabbits, 2=Draw Foxes, 0=Normal mode")
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
	
	// Number keys to change drawing mode
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.drawMode = "rabbit"
		log.Println("Draw mode: RABBIT (click to place rabbits)")
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.drawMode = "fox"
		log.Println("Draw mode: FOX (click to place foxes)")
	}
	if inpututil.IsKeyJustPressed(ebiten.Key0) {
		g.drawMode = "none"
		log.Println("Draw mode: NONE (normal simulation)")
	}
	
	// S key to save/export data manually
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.saveSimulationData()
	}
	
	// Mouse drawing
	g.handleMouseInput()
	
	// Mouse clicks for buttons
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.handleButtonClick(x, y)
	}
}

// handleMouseInput handles mouse drawing of animals
func (g *Game) handleMouseInput() {
	if g.world == nil || g.drawMode == "none" {
		return
	}
	
	// Check if mouse is pressed
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !g.mousePressed {
			g.mousePressed = true
			g.handleMouseDraw()
		}
	} else {
		g.mousePressed = false
	}
}

// handleMouseDraw places animals at mouse position
func (g *Game) handleMouseDraw() {
	x, y := ebiten.CursorPosition()
	
	// Only draw in the simulation area (not on UI or graph)
	if y >= gameAreaHeight || x < 0 || x >= screenWidth {
		return
	}
	
	// Convert screen coordinates to grid coordinates
	gridX := x / cellSize
	gridY := y / cellSize
	
	// Check bounds
	if gridX < 0 || gridX >= gridWidth || gridY < 0 || gridY >= gridHeight {
		return
	}
	
	pos := Position{gridX, gridY}
	
	// Check if position is already occupied by an animal
	if g.world.Grid[gridX][gridY] == RabbitType || g.world.Grid[gridX][gridY] == FoxType {
		return // Don't overwrite existing animals
	}
	
	switch g.drawMode {
	case "rabbit":
		// Check population limit
		if len(g.world.Rabbits) >= maxRabbits {
			return
		}
		
		// Create new rabbit
		rabbit := &Rabbit{
			Animal: Animal{
				Position:    pos,
				Energy:      80, // Start with good energy
				ReproduceCD: 0,
				Age:         0,
			},
			NewBorn: 60, // Show as yellow briefly
		}
		
		g.world.Rabbits = append(g.world.Rabbits, rabbit)
		g.world.Grid[gridX][gridY] = RabbitType
		log.Printf("Placed rabbit at (%d,%d)", gridX, gridY)
		
	case "fox":
		// Check population limit
		if len(g.world.Foxes) >= maxFoxes {
			return
		}
		
		// Create new fox
		fox := &Fox{
			Animal: Animal{
				Position:    pos,
				Energy:      80, // Start with good energy
				ReproduceCD: 0,
				Age:         0,
			},
		}
		
		g.world.Foxes = append(g.world.Foxes, fox)
		g.world.Grid[gridX][gridY] = FoxType
		log.Printf("Placed fox at (%d,%d)", gridX, gridY)
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
		g.drawMode = "none"
		log.Println("Simulation reset")
	}
	
	// Save button (520, 50, 160, 30) - wider button below pause/play
	if x >= 520 && x <= 680 && y >= 50 && y <= 80 {
		g.saveSimulationData()
	}
}

// saveSimulationData saves both CSV data and JPG screenshot
func (g *Game) saveSimulationData() {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	
	// Save CSV data
	if len(g.populationHistory) > 0 {
		exportPopulationData(g.populationHistory)
	}
	
	// Save current screenshot
	g.saveScreenshot(timestamp)
	
	// Save complete history as image sequence
	g.saveHistorySequence(timestamp)
	
	log.Printf("Saved simulation data with timestamp: %s", timestamp)
}

// saveHistorySequence creates images for each point in population history
func (g *Game) saveHistorySequence(timestamp string) {
	if len(g.populationHistory) < 2 {
		log.Println("Not enough history data for sequence")
		return
	}
	
	log.Printf("Creating history sequence with %d frames...", len(g.populationHistory))
	
	// Create directory for sequence
	sequenceDir := fmt.Sprintf("ecosystem_sequence_%s", timestamp)
	err := os.Mkdir(sequenceDir, 0755)
	if err != nil {
		log.Printf("Error creating sequence directory: %v", err)
		return
	}
	
	// Create temporary world state for rendering each frame
	originalHistory := make([]PopulationData, len(g.populationHistory))
	copy(originalHistory, g.populationHistory)
	
	// Generate each frame
	for i := 0; i < len(originalHistory); i++ {
		// Create history up to this point
		currentHistory := originalHistory[:i+1]
		
		// Create frame
		filename := fmt.Sprintf("%s/frame_%03d.jpg", sequenceDir, i)
		g.saveHistoryFrame(filename, currentHistory, originalHistory[i])
		
		// Progress indicator
		if i%10 == 0 || i == len(originalHistory)-1 {
			log.Printf("Generated frame %d/%d", i+1, len(originalHistory))
		}
	}
	
	// Create summary file with instructions
	g.createSequenceSummary(sequenceDir, timestamp, len(originalHistory))
	
	log.Printf("History sequence saved to: %s", sequenceDir)
}

// saveHistoryFrame saves a single frame of the history sequence
func (g *Game) saveHistoryFrame(filename string, historyUpToPoint []PopulationData, currentData PopulationData) {
	// Create image
	bounds := image.Rect(0, 0, screenWidth, screenHeight)
	img := image.NewRGBA(bounds)
	
	// Create temporary screen
	screen := ebiten.NewImage(screenWidth, screenHeight)
	
	// Render frame
	g.renderHistoryFrame(screen, historyUpToPoint, currentData)
	
	// Copy pixels
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			r, g, b, a := screen.At(x, y).RGBA()
			img.Set(x, y, color.RGBA{
				uint8(r >> 8),
				uint8(g >> 8),
				uint8(b >> 8),
				uint8(a >> 8),
			})
		}
	}
	
	// Save file
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error creating frame file: %v", err)
		return
	}
	defer file.Close()
	
	err = jpeg.Encode(file, img, &jpeg.Options{Quality: 80})
	if err != nil {
		log.Printf("Error saving frame: %v", err)
	}
}

// renderHistoryFrame renders a single frame showing history up to a specific point
func (g *Game) renderHistoryFrame(screen *ebiten.Image, historyUpToPoint []PopulationData, currentData PopulationData) {
	// Fill background
	screen.Fill(color.RGBA{0, 0, 0, 255})
	
	// Draw graph area
	g.fillRect(screen, graphOffsetX, graphOffsetY, graphWidth, graphHeight, color.RGBA{20, 20, 20, 255})
	
	// Draw border
	g.fillRect(screen, graphOffsetX, graphOffsetY, graphWidth, 2, color.RGBA{100, 100, 100, 255})
	g.fillRect(screen, graphOffsetX, graphOffsetY+graphHeight-2, graphWidth, 2, color.RGBA{100, 100, 100, 255})
	g.fillRect(screen, graphOffsetX, graphOffsetY, 2, graphHeight, color.RGBA{100, 100, 100, 255})
	g.fillRect(screen, graphOffsetX+graphWidth-2, graphOffsetY, 2, graphHeight, color.RGBA{100, 100, 100, 255})
	
	if len(historyUpToPoint) < 1 {
		return
	}
	
	// Calculate scale
	maxValue := 20
	for _, data := range historyUpToPoint {
		if data.Rabbits > maxValue { maxValue = data.Rabbits }
		if data.Foxes > maxValue { maxValue = data.Foxes }
		if data.Grass/10 > maxValue { maxValue = data.Grass/10 }
	}
	
	// Draw population points up to current point
	g.drawHistoryPoints(screen, historyUpToPoint, "rabbits", maxValue, color.RGBA{255, 255, 255, 255})
	g.drawHistoryPoints(screen, historyUpToPoint, "foxes", maxValue, color.RGBA{255, 0, 0, 255})
	g.drawHistoryPoints(screen, historyUpToPoint, "grass", maxValue, color.RGBA{0, 255, 0, 255})
	
	// Add title and current stats
	title := fmt.Sprintf("Ecosystem Evolution - Tick: %d (Frame %d)", currentData.Tick, len(historyUpToPoint))
	ebitenutil.DebugPrint(screen, title)
	
	// Current population stats
	stats := fmt.Sprintf("Current: Rabbits=%d  Foxes=%d  Grass=%d", 
		currentData.Rabbits, currentData.Foxes, currentData.Grass)
	ebitenutil.DebugPrintAt(screen, stats, 10, 30)
	
	// Legend
	ebitenutil.DebugPrintAt(screen, "■=Rabbits ♦=Foxes +=Grass", 30, screenHeight-20)
	
	// Progress bar
	progressWidth := 200
	progressX := screenWidth - progressWidth - 20
	progressY := 10
	progress := float64(len(historyUpToPoint)) / float64(maxHistoryPoints)
	
	// Progress bar background
	g.fillRect(screen, progressX, progressY, progressWidth, 10, color.RGBA{50, 50, 50, 255})
	// Progress bar fill
	g.fillRect(screen, progressX, progressY, int(float64(progressWidth)*progress), 10, color.RGBA{0, 150, 0, 255})
}

// drawHistoryPoints draws points for a specific population type in history frame
func (g *Game) drawHistoryPoints(screen *ebiten.Image, history []PopulationData, populationType string, maxValue int, pointColor color.RGBA) {
	if len(history) < 1 || maxValue <= 0 {
		return
	}
	
	var yOffset int
	switch populationType {
	case "rabbits": yOffset = 0
	case "foxes": yOffset = -3  
	case "grass": yOffset = 3
	}
	
	for i, data := range history {
		var value int
		switch populationType {
		case "rabbits": value = data.Rabbits
		case "foxes": value = data.Foxes
		case "grass": value = data.Grass / 10
		default: continue
		}
		
		if value == 0 { continue }
		
		// Calculate position (spread across full width)
		x := graphOffsetX + 5 + ((i * (graphWidth - 10)) / maxHistoryPoints)
		y := graphOffsetY + graphHeight - 5 - ((value * (graphHeight - 10)) / maxValue) + yOffset
		
		// Clamp y
		if y < graphOffsetY + 5 { y = graphOffsetY + 5 }
		if y > graphOffsetY + graphHeight - 5 { y = graphOffsetY + graphHeight - 5 }
		
		// Draw point with different shapes
		switch populationType {
		case "rabbits":
			g.fillRect(screen, x-1, y-1, 3, 3, pointColor)
		case "foxes":
			g.fillRect(screen, x, y-1, 1, 1, pointColor)
			g.fillRect(screen, x-1, y, 1, 1, pointColor)
			g.fillRect(screen, x+1, y, 1, 1, pointColor)
			g.fillRect(screen, x, y+1, 1, 1, pointColor)
		case "grass":
			g.fillRect(screen, x-1, y, 3, 1, pointColor)
			g.fillRect(screen, x, y-1, 1, 3, pointColor)
		}
	}
}

// createSequenceSummary creates a text file with information about the sequence
func (g *Game) createSequenceSummary(sequenceDir, timestamp string, frameCount int) {
	summaryFile := fmt.Sprintf("%s/README.txt", sequenceDir)
	file, err := os.Create(summaryFile)
	if err != nil {
		log.Printf("Error creating summary file: %v", err)
		return
	}
	defer file.Close()
	
	summary := fmt.Sprintf(`Ecosystem Simulation History Sequence
=====================================

Generated: %s
Frames: %d
Duration: %d ticks
Frame rate: One frame per 5 seconds of simulation

Files:
- frame_000.jpg to frame_%03d.jpg: Individual frames showing population evolution
- README.txt: This file

To create a video from these frames, you can use tools like:

FFmpeg (example):
ffmpeg -framerate 10 -i frame_%%03d.jpg -c:v libx264 -pix_fmt yuv420p ecosystem_evolution.mp4

VirtualDub, Adobe Premiere, or other video editing software can also be used.

Each frame shows:
- Population graph growing over time
- Current population numbers
- Progress bar showing simulation progress
- Legend for graph symbols

This sequence captures the complete evolution of your ecosystem simulation!
`, time.Now().Format("2006-01-02 15:04:05"), frameCount, 
   g.populationHistory[len(g.populationHistory)-1].Tick, frameCount-1)
	
	file.WriteString(summary)
	log.Printf("Created sequence summary: %s", summaryFile)
}

// saveScreenshot captures and saves the current game screen as JPG
func (g *Game) saveScreenshot(timestamp string) {
	// Create a new image to capture the screen
	bounds := image.Rect(0, 0, screenWidth, screenHeight)
	img := image.NewRGBA(bounds)
	
	// Create a temporary screen image
	screen := ebiten.NewImage(screenWidth, screenHeight)
	
	// Render the game to the temporary screen
	g.renderToImage(screen)
	
	// Read pixels from the screen
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			r, g, b, a := screen.At(x, y).RGBA()
			img.Set(x, y, color.RGBA{
				uint8(r >> 8),
				uint8(g >> 8), 
				uint8(b >> 8),
				uint8(a >> 8),
			})
		}
	}
	
	// Create file
	filename := fmt.Sprintf("ecosystem_screenshot_%s.jpg", timestamp)
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error creating screenshot file: %v", err)
		return
	}
	defer file.Close()
	
	// Save as JPEG
	err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	if err != nil {
		log.Printf("Error saving screenshot: %v", err)
		return
	}
	
	log.Printf("Screenshot saved: %s", filename)
}

// renderToImage renders the game state to an image (without UI elements)
func (g *Game) renderToImage(screen *ebiten.Image) {
	// Fill background with black
	screen.Fill(color.RGBA{0, 0, 0, 255})
	
	// Draw simulation area
	if g.world != nil {
		g.drawWorld(screen)
	}
	
	// Draw population graph
	if g.world != nil {
		g.drawPopulationGraph(screen)
	}
	
	// Add title and basic info (without debug details)
	title := fmt.Sprintf("Ecosystem Simulation - Tick: %d", g.world.Tick)
	if g.paused {
		title += " (PAUSED)"
	}
	ebitenutil.DebugPrint(screen, title)
	
	// Add population info at bottom
	if g.world != nil {
		info := fmt.Sprintf("Rabbits: %d  Foxes: %d  Grass: %d", 
			len(g.world.Rabbits), len(g.world.Foxes), len(g.world.Grass))
		ebitenutil.DebugPrintAt(screen, info, 10, screenHeight-20)
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
		
		// Drawing mode info
		debugText += fmt.Sprintf("Draw Mode: %s\n", strings.ToUpper(g.drawMode))
		debugText += "Controls: SPACE=Pause 1=Rabbit 2=Fox 0=None S=Save+Screenshot"
	}
	
	ebitenutil.DebugPrint(screen, debugText)
	
	// Draw button labels using DebugPrint
	ebitenutil.DebugPrintAt(screen, "PAUSE", 535, 20)
	ebitenutil.DebugPrintAt(screen, "PLAY", 630, 20)
	ebitenutil.DebugPrintAt(screen, "RESET", 715, 20)
	ebitenutil.DebugPrintAt(screen, "SAVE DATA + SCREENSHOT", 535, 60)
	
	// Draw cursor indicator in drawing mode
	if g.drawMode != "none" {
		g.drawCursor(screen)
	}
	
	// Draw simple graph legend
	ebitenutil.DebugPrintAt(screen, "■=Rabbits ♦=Foxes +=Grass", 30, 580)
}

// drawCursor shows what will be placed at mouse position
func (g *Game) drawCursor(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	
	// Only show cursor in simulation area
	if y >= gameAreaHeight || x < 0 || x >= screenWidth {
		return
	}
	
	// Convert to grid position
	gridX := x / cellSize
	gridY := y / cellSize
	
	// Check bounds
	if gridX < 0 || gridX >= gridWidth || gridY < 0 || gridY >= gridHeight {
		return
	}
	
	// Draw preview at grid position
	drawX := gridX * cellSize
	drawY := gridY * cellSize
	
	var cursorColor color.RGBA
	switch g.drawMode {
	case "rabbit":
		cursorColor = color.RGBA{255, 255, 255, 128} // Semi-transparent white
	case "fox":
		cursorColor = color.RGBA{255, 0, 0, 128} // Semi-transparent red
	}
	
	// Draw preview rectangle
	g.fillRect(screen, drawX+1, drawY+1, cellSize-2, cellSize-2, cursorColor)
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
	
	// Save button (new, wider button below the others)
	g.fillRect(screen, 520, 50, 160, 30, color.RGBA{200, 200, 100, 255}) // Yellow
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
	
	// Defer function to export data when program exits
	defer func() {
		if game.world != nil && len(game.populationHistory) > 5 { // Only if we have meaningful data
			timestamp := time.Now().Format("2006-01-02_15-04-05")
			
			log.Println("Saving final simulation data...")
			exportPopulationData(game.populationHistory)
			
			log.Println("Creating complete history sequence...")
			game.saveHistorySequence(timestamp)
			
			log.Println("Simulation data export complete!")
		}
	}()
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// exportPopulationData saves population history to CSV file
func exportPopulationData(history []PopulationData) {
	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("ecosystem_data_%s.csv", timestamp)
	
	// Create file
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error creating CSV file: %v", err)
		return
	}
	defer file.Close()
	
	// Write file header with simulation info
	file.WriteString("# Ecosystem Simulation Data Export\n")
	file.WriteString(fmt.Sprintf("# Generated: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	file.WriteString(fmt.Sprintf("# Duration: %d ticks (%d data points)\n", history[len(history)-1].Tick, len(history)))
	file.WriteString(fmt.Sprintf("# Data recorded every 5 seconds (30 ticks)\n"))
	file.WriteString("# \n")
	file.WriteString("# Simulation parameters:\n")
	file.WriteString(fmt.Sprintf("# - Grid size: %dx%d\n", gridWidth, gridHeight))
	file.WriteString(fmt.Sprintf("# - Max rabbits: %d\n", maxRabbits))
	file.WriteString(fmt.Sprintf("# - Max foxes: %d\n", maxFoxes))
	file.WriteString(fmt.Sprintf("# - Grass growth rate: %d per tick\n", grassGrowthRate))
	file.WriteString(fmt.Sprintf("# - Grass spawn chance: %.3f per tick\n", grassSpawnChance))
	file.WriteString("# \n")
	
	// Write CSV header
	_, err = file.WriteString("Tick,Rabbits,Foxes,Grass,Timestamp\n")
	if err != nil {
		log.Printf("Error writing CSV header: %v", err)
		return
	}
	
	// Write data rows
	startTime := time.Now().Add(-time.Duration(len(history)) * 5 * time.Second) // Approximate start time
	for i, data := range history {
		rowTime := startTime.Add(time.Duration(i) * 5 * time.Second) // 5 seconds between data points
		line := fmt.Sprintf("%d,%d,%d,%d,%s\n", 
			data.Tick, 
			data.Rabbits, 
			data.Foxes, 
			data.Grass,
			rowTime.Format("15:04:05"))
		
		_, err = file.WriteString(line)
		if err != nil {
			log.Printf("Error writing CSV data: %v", err)
			return
		}
	}
	
	log.Printf("Population data exported to: %s", filename)
	log.Printf("Exported %d data points covering %d ticks", len(history), history[len(history)-1].Tick)
	
	// Calculate and log some statistics
	if len(history) > 1 {
		maxRabbits := 0
		maxFoxes := 0
		maxGrass := 0
		for _, data := range history {
			if data.Rabbits > maxRabbits { maxRabbits = data.Rabbits }
			if data.Foxes > maxFoxes { maxFoxes = data.Foxes }
			if data.Grass > maxGrass { maxGrass = data.Grass }
		}
		log.Printf("Peak populations: Rabbits=%d, Foxes=%d, Grass=%d", maxRabbits, maxFoxes, maxGrass)
	}
}