package main

import "math/rand"

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
	w.updateFoxes()
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

// addTestEntities adds some test entities for visualization
func (w *World) addTestEntities() {
	// Add some grass patches (więcej trawy)
	for i := 0; i < 30; i++ {
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		pos := Position{x, y}
		
		w.Grass[pos] = &Grass{
			Position: pos,
			Amount:   rand.Intn(51) + 50, // 50-100 (dojrzała trawa)
		}
		w.Grid[x][y] = GrassType
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
			if x >= 0 && x < gridWidth && y >= 0 && y < gridHeight && w.Grid[x][y] == Empty {
				rabbit := &Rabbit{
					Animal: Animal{
						Position:    Position{x, y},
						Energy:      80, // Start with more energy
						ReproduceCD: 0,
						Age:         0,
					},
					NewBorn: 0, // Not newborn
				}
				
				w.Rabbits = append(w.Rabbits, rabbit)
				w.Grid[x][y] = RabbitType
			}
		}
	}
	
	// Add some foxes (w grupach jak króliki)
	for group := 0; group < 2; group++ {
		// Random center for fox group
		centerX := rand.Intn(gridWidth-6) + 3
		centerY := rand.Intn(gridHeight-6) + 3
		
		// Add 2-3 foxes around this center  
		for i := 0; i < 2+rand.Intn(2); i++ {
			x := centerX + rand.Intn(4) - 2 // Within 2 cells of center
			y := centerY + rand.Intn(4) - 2
			
			// Make sure within bounds and not occupied
			if x >= 0 && x < gridWidth && y >= 0 && y < gridHeight && w.Grid[x][y] == Empty {
				fox := &Fox{
					Animal: Animal{
						Position:    Position{x, y},
						Energy:      80, // Start with high energy (ready to reproduce)
						ReproduceCD: 0,
						Age:         0,
					},
				}
				
				w.Foxes = append(w.Foxes, fox)
				w.Grid[x][y] = FoxType
			}
		}
	}
}