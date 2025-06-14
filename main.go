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
	rabbitMoveChance = 0.7   // Chance rabbit moves each tick (zwiększone z 0.5)
	rabbitEnergyLoss = 1     // Energy lost per tick every 60 ticks (co sekundę)
	grassEnergyGain  = 40    // Energy gained from eating grass (dużo więcej)
	minGrassToEat    = 5     // Minimum grass amount to be edible (bardzo niskie)
	
	// Reproduction parameters
	reproduceEnergyThreshold = 60  // Min energy to reproduce (zmniejszone z 70)
	reproductionCooldown     = 180 // Ticks before can reproduce again (3 seconds zamiast 5)
	reproduceChance          = 0.3 // Chance to try reproduction each tick (zwiększone z 0.1)
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
	NewBorn      int // Ticks since birth (for visual indication)
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
	// First, handle all rabbit updates without reproduction
	for i := len(w.Rabbits) - 1; i >= 0; i-- {
		rabbit := w.Rabbits[i]
		
		// Age and reduce reproduction cooldown
		rabbit.Age++
		if rabbit.ReproduceCD > 0 {
			rabbit.ReproduceCD--
		}
		
		// Reduce newborn indicator
		if rabbit.NewBorn > 0 {
			rabbit.NewBorn--
		}
		
		// Lose energy only every 60 ticks (roughly once per second)
		if w.Tick%60 == 0 {
			rabbit.Energy -= rabbitEnergyLoss
		}
		
		// Try to eat grass at current position
		w.rabbitEatGrass(rabbit)
		
		// Move rabbit randomly
		if rand.Float64() < rabbitMoveChance {
			w.moveRabbit(rabbit)
			// Try to eat grass at new position too
			w.rabbitEatGrass(rabbit)
		}
		
		// Check if rabbit dies
		if rabbit.Energy <= 0 {
			w.removeRabbit(i)
		}
	}
	
	// Then handle reproduction in separate pass to avoid multiple births per tick
	w.handleRabbitReproduction()
}

// handleRabbitReproduction handles all rabbit reproduction in one pass
func (w *World) handleRabbitReproduction() {
	processedPairs := make(map[string]bool) // Track processed rabbit pairs
	
	for _, rabbit := range w.Rabbits {
		// Skip if this rabbit can't reproduce
		if rabbit.Energy < reproduceEnergyThreshold || rabbit.ReproduceCD > 0 {
			continue
		}
		
		// Skip reproduction attempt with some probability
		if rand.Float64() >= reproduceChance {
			continue
		}
		
		// Look for adjacent partner
		partner := w.findAdjacentReproductivePartner(rabbit, processedPairs)
		if partner != nil {
			// Mark this pair as processed (both directions)
			pairKey1 := fmt.Sprintf("%d,%d-%d,%d", rabbit.Position.X, rabbit.Position.Y, partner.Position.X, partner.Position.Y)
			pairKey2 := fmt.Sprintf("%d,%d-%d,%d", partner.Position.X, partner.Position.Y, rabbit.Position.X, rabbit.Position.Y)
			processedPairs[pairKey1] = true
			processedPairs[pairKey2] = true
			
			// Create baby
			w.createBabyRabbit(rabbit, partner)
		}
	}
}

// findAdjacentReproductivePartner finds a suitable partner for reproduction
func (w *World) findAdjacentReproductivePartner(rabbit *Rabbit, processedPairs map[string]bool) *Rabbit {
	adjacentPositions := w.getAdjacentPositions(rabbit.Position)
	
	for _, pos := range adjacentPositions {
		if w.Grid[pos.X][pos.Y] == RabbitType {
			partner := w.findRabbitAtPosition(pos)
			if partner != nil && 
			   partner.Energy >= reproduceEnergyThreshold && 
			   partner.ReproduceCD == 0 {
				
				// Check if this pair was already processed this tick
				pairKey := fmt.Sprintf("%d,%d-%d,%d", rabbit.Position.X, rabbit.Position.Y, partner.Position.X, partner.Position.Y)
				if !processedPairs[pairKey] {
					return partner
				}
			}
		}
	}
	return nil
}

// findRabbitAtPosition finds rabbit at given position
func (w *World) findRabbitAtPosition(pos Position) *Rabbit {
	for _, rabbit := range w.Rabbits {
		if rabbit.Position.X == pos.X && rabbit.Position.Y == pos.Y {
			return rabbit
		}
	}
	return nil
}

// createBabyRabbit creates new rabbit from two parents
func (w *World) createBabyRabbit(parent1, parent2 *Rabbit) {
	// Find empty adjacent position for baby
	adjacentPositions := w.getAdjacentPositions(parent1.Position)
	
	for _, pos := range adjacentPositions {
		if w.Grid[pos.X][pos.Y] == Empty {
			// Create baby rabbit
			baby := &Rabbit{
				Position:    pos,
				Energy:      60, // Born with decent energy
				ReproduceCD: reproductionCooldown, // Can't reproduce immediately
				Age:         0,
				NewBorn:     180, // Visual indicator for 30 seconds (180 ticks at 6 FPS)
			}
			
			w.Rabbits = append(w.Rabbits, baby)
			w.Grid[pos.X][pos.Y] = RabbitType
			
			// Parents lose energy and get cooldown
			parent1.Energy -= 20
			parent2.Energy -= 20
			parent1.ReproduceCD = reproductionCooldown
			parent2.ReproduceCD = reproductionCooldown
			
			// Log birth for visibility
			log.Printf("New rabbit born at (%d,%d)! Total rabbits: %d", pos.X, pos.Y, len(w.Rabbits))
			
			return
		}
	}
}

// rabbitEatGrass makes rabbit eat grass at its current position
func (w *World) rabbitEatGrass(rabbit *Rabbit) {
	pos := rabbit.Position
	grass, exists := w.Grass[pos]
	
	if exists && grass.Amount >= minGrassToEat {
		// Rabbit eats the grass
		rabbit.Energy += grassEnergyGain
		
		// Cap energy at reasonable level
		if rabbit.Energy > 100 {
			rabbit.Energy = 100
		}
		
		// Remove grass completely (rabbit ate it all)
		delete(w.Grass, pos)
		// Don't update grid here - let moveRabbit handle it properly
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
	slowMode bool  // Slower simulation
	tickCounter int // Count frames to slow down updates
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
	
	// Slow down simulation - update world only every 10 frames (6 FPS instead of 60)
	g.tickCounter++
	if g.tickCounter >= 10 {
		g.tickCounter = 0
		g.world.Update()
		g.world.Tick++
	}
	
	return nil
}

// addTestEntities adds some test entities for visualization
func (g *Game) addTestEntities() {
	// Add some grass patches (więcej trawy)
	for i := 0; i < 30; i++ {
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		pos := Position{x, y}
		
		g.world.Grass[pos] = &Grass{
			Position: pos,
			Amount:   rand.Intn(51) + 50, // 50-100 (dojrzała trawa)
		}
		g.world.Grid[x][y] = GrassType
	}
	
	// Add some rabbits (w grupach żeby się częściej spotykały)
	for group := 0; group < 3; group++ {
		// Random center for group
		centerX := rand.Intn(gridWidth-10) + 5
		centerY := rand.Intn(gridHeight-10) + 5
		
		// Add 3-4 rabbits around this center
		for i := 0; i < 3+rand.Intn(2); i++ {
			x := centerX + rand.Intn(6) - 3 // Within 3 cells of center
			y := centerY + rand.Intn(6) - 3
			
			// Make sure within bounds and not occupied
			if x >= 0 && x < gridWidth && y >= 0 && y < gridHeight && g.world.Grid[x][y] == Empty {
				rabbit := &Rabbit{
					Position: Position{x, y},
					Energy:   80, // Start with more energy
					ReproduceCD: 0,
					Age:      0,
					NewBorn:  0, // Not newborn
				}
				
				g.world.Rabbits = append(g.world.Rabbits, rabbit)
				g.world.Grid[x][y] = RabbitType
			}
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
		
		// Show rabbit reproduction status
		if len(g.world.Rabbits) > 0 {
			readyToReproduce := 0
			totalEnergy := 0
			for _, r := range g.world.Rabbits {
				totalEnergy += r.Energy
				if r.Energy >= reproduceEnergyThreshold && r.ReproduceCD == 0 {
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

// drawWorld renders the world grid
func (g *Game) drawWorld(screen *ebiten.Image) {
	// First draw grass (bottom layer)
	for pos, grass := range g.world.Grass {
		g.drawGrass(screen, pos, grass.Amount)
	}
	
	// Then draw animals on top
	for _, rabbit := range g.world.Rabbits {
		g.drawRabbit(screen, rabbit.Position)
	}
	
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

// drawRabbit draws a rabbit (smaller white dot so grass is visible underneath)
func (g *Game) drawRabbit(screen *ebiten.Image, pos Position) {
	x := pos.X * cellSize
	y := pos.Y * cellSize
	
	// Find the rabbit at this position to check if it's newborn
	var rabbit *Rabbit
	for _, r := range g.world.Rabbits {
		if r.Position.X == pos.X && r.Position.Y == pos.Y {
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