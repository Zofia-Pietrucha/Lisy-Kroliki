package main

import (
	"fmt"
	"log"
	"math/rand"
)

// Animal represents base animal properties
type Animal struct {
	Position
	Energy       int
	ReproduceCD  int // Cooldown after reproduction
	Age          int
}

// Rabbit represents a rabbit entity
type Rabbit struct {
	Animal
	NewBorn int // Ticks since birth (for visual indication)
}

// Fox represents a fox entity
type Fox struct {
	Animal
	// Fox-specific fields can be added here later
}

// updateRabbits handles rabbit movement and aging
func (w *World) updateRabbits() {
	// First, handle all rabbit updates without reproduction
	for i := len(w.Rabbits) - 1; i >= 0; i-- {
		rabbit := w.Rabbits[i]
		
		// Age and reduce reproduction cooldown
		rabbit.Animal.Age++
		if rabbit.Animal.ReproduceCD > 0 {
			rabbit.Animal.ReproduceCD--
		}
		
		// Reduce newborn indicator
		if rabbit.NewBorn > 0 {
			rabbit.NewBorn--
		}
		
		// Lose energy only every 60 ticks (roughly once per second)
		if w.Tick%60 == 0 {
			rabbit.Animal.Energy -= rabbitEnergyLoss
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
		if rabbit.Animal.Energy <= 0 {
			w.removeRabbit(i)
		}
	}
	
	// Then handle reproduction in separate pass to avoid multiple births per tick
	w.handleRabbitReproduction()
}

// rabbitEatGrass makes rabbit eat grass at its current position
func (w *World) rabbitEatGrass(rabbit *Rabbit) {
	pos := rabbit.Animal.Position
	grass, exists := w.Grass[pos]
	
	if exists && grass.Amount >= minGrassToEat {
		// Rabbit eats the grass
		rabbit.Animal.Energy += grassEnergyGain
		
		// Cap energy at reasonable level
		if rabbit.Animal.Energy > 100 {
			rabbit.Animal.Energy = 100
		}
		
		// Remove grass completely (rabbit ate it all)
		delete(w.Grass, pos)
		// Don't update grid here - let moveRabbit handle it properly
	}
}

// moveRabbit moves a rabbit to a random adjacent position
func (w *World) moveRabbit(rabbit *Rabbit) {
	// Clear current position
	w.Grid[rabbit.Animal.Position.X][rabbit.Animal.Position.Y] = Empty
	
	// Get possible moves (adjacent cells)
	moves := w.getAdjacentPositions(rabbit.Animal.Position)
	
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
		rabbit.Animal.Position = newPos
	}
	
	// Update grid with new position
	w.Grid[rabbit.Animal.Position.X][rabbit.Animal.Position.Y] = RabbitType
}

// handleRabbitReproduction handles all rabbit reproduction in one pass
func (w *World) handleRabbitReproduction() {
	processedPairs := make(map[string]bool) // Track processed rabbit pairs
	
	for _, rabbit := range w.Rabbits {
		// Skip if this rabbit can't reproduce
		if rabbit.Animal.Energy < reproduceEnergyThreshold || rabbit.Animal.ReproduceCD > 0 {
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
			pairKey1 := fmt.Sprintf("%d,%d-%d,%d", rabbit.Animal.Position.X, rabbit.Animal.Position.Y, partner.Animal.Position.X, partner.Animal.Position.Y)
			pairKey2 := fmt.Sprintf("%d,%d-%d,%d", partner.Animal.Position.X, partner.Animal.Position.Y, rabbit.Animal.Position.X, rabbit.Animal.Position.Y)
			processedPairs[pairKey1] = true
			processedPairs[pairKey2] = true
			
			// Create baby
			w.createBabyRabbit(rabbit, partner)
		}
	}
}

// findAdjacentReproductivePartner finds a suitable partner for reproduction
func (w *World) findAdjacentReproductivePartner(rabbit *Rabbit, processedPairs map[string]bool) *Rabbit {
	adjacentPositions := w.getAdjacentPositions(rabbit.Animal.Position)
	
	for _, pos := range adjacentPositions {
		if w.Grid[pos.X][pos.Y] == RabbitType {
			partner := w.findRabbitAtPosition(pos)
			if partner != nil && 
			   partner.Animal.Energy >= reproduceEnergyThreshold && 
			   partner.Animal.ReproduceCD == 0 {
				
				// Check if this pair was already processed this tick
				pairKey := fmt.Sprintf("%d,%d-%d,%d", rabbit.Animal.Position.X, rabbit.Animal.Position.Y, partner.Animal.Position.X, partner.Animal.Position.Y)
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
		if rabbit.Animal.Position.X == pos.X && rabbit.Animal.Position.Y == pos.Y {
			return rabbit
		}
	}
	return nil
}

// createBabyRabbit creates new rabbit from two parents
func (w *World) createBabyRabbit(parent1, parent2 *Rabbit) {
	// Check population limit
	if len(w.Rabbits) >= maxRabbits {
		return // Too many rabbits already
	}
	
	// Find empty adjacent position for baby
	adjacentPositions := w.getAdjacentPositions(parent1.Animal.Position)
	
	for _, pos := range adjacentPositions {
		if w.Grid[pos.X][pos.Y] == Empty {
			// Create baby rabbit
			baby := &Rabbit{
				Animal: Animal{
					Position:    pos,
					Energy:      60, // Born with decent energy
					ReproduceCD: reproductionCooldown, // Can't reproduce immediately
					Age:         0,
				},
				NewBorn: 180, // Visual indicator for 30 seconds (180 ticks at 6 FPS)
			}
			
			w.Rabbits = append(w.Rabbits, baby)
			w.Grid[pos.X][pos.Y] = RabbitType
			
			// Parents lose energy and get cooldown
			parent1.Animal.Energy -= 20
			parent2.Animal.Energy -= 20
			parent1.Animal.ReproduceCD = reproductionCooldown
			parent2.Animal.ReproduceCD = reproductionCooldown
			
			// Log birth for visibility
			log.Printf("New rabbit born at (%d,%d)! Total rabbits: %d", pos.X, pos.Y, len(w.Rabbits))
			
			return
		}
	}
}

// removeRabbit removes a rabbit from the world
func (w *World) removeRabbit(index int) {
	rabbit := w.Rabbits[index]
	
	// Clear grid position
	w.Grid[rabbit.Animal.Position.X][rabbit.Animal.Position.Y] = Empty
	
	// Remove from slice
	w.Rabbits = append(w.Rabbits[:index], w.Rabbits[index+1:]...)
}

// updateFoxes handles fox movement, hunting and aging
func (w *World) updateFoxes() {
	for i := len(w.Foxes) - 1; i >= 0; i-- {
		fox := w.Foxes[i]
		
		// Age and reduce reproduction cooldown
		fox.Animal.Age++
		if fox.Animal.ReproduceCD > 0 {
			fox.Animal.ReproduceCD--
		}
		
		// Lose energy only every 60 ticks
		if w.Tick%60 == 0 {
			fox.Animal.Energy -= foxEnergyLoss
		}
		
		// Try to hunt rabbit at current position first
		w.foxHuntRabbit(fox)
		
		// Move fox (actively hunt - prefer moving toward rabbits)
		if rand.Float64() < foxMoveChance {
			w.moveFoxHunting(fox)
			// Try to hunt at new position too
			w.foxHuntRabbit(fox)
		}
		
		// Try reproduction if well-fed
		if fox.Animal.Energy >= foxReproduceThreshold && fox.Animal.ReproduceCD == 0 {
			if rand.Float64() < reproduceChance*1.5 { // Foxes reproduce easier
				w.tryFoxReproduction(fox)
			}
		}
		
		// Check if fox dies
		if fox.Animal.Energy <= 0 {
			w.removeFox(i)
		}
	}
}

// moveFoxHunting moves fox, preferring positions with rabbits nearby
func (w *World) moveFoxHunting(fox *Fox) {
	// Clear current position
	w.Grid[fox.Animal.Position.X][fox.Animal.Position.Y] = Empty
	
	// Get possible moves (adjacent cells)
	moves := w.getAdjacentPositions(fox.Animal.Position)
	
	// Separate moves into: rabbit positions and other valid positions
	rabbitMoves := make([]Position, 0)
	validMoves := make([]Position, 0)
	
	for _, pos := range moves {
		cellType := w.Grid[pos.X][pos.Y]
		if cellType == RabbitType {
			rabbitMoves = append(rabbitMoves, pos)
		} else if cellType == Empty || cellType == GrassType {
			validMoves = append(validMoves, pos)
		}
	}
	
	// Prefer moving to rabbit positions (hunting!)
	if len(rabbitMoves) > 0 {
		newPos := rabbitMoves[rand.Intn(len(rabbitMoves))]
		fox.Animal.Position = newPos
	} else if len(validMoves) > 0 {
		// No rabbits nearby, move randomly
		newPos := validMoves[rand.Intn(len(validMoves))]
		fox.Animal.Position = newPos
	}
	// If no valid moves, stay in place
	
	// Update grid with new position
	w.Grid[fox.Animal.Position.X][fox.Animal.Position.Y] = FoxType
}

// foxHuntRabbit makes fox hunt rabbit at its current position
func (w *World) foxHuntRabbit(fox *Fox) {
	pos := fox.Animal.Position
	
	// Look for rabbit at this position
	for i, rabbit := range w.Rabbits {
		if rabbit.Animal.Position.X == pos.X && rabbit.Animal.Position.Y == pos.Y {
			// Fox catches rabbit!
			fox.Animal.Energy += rabbitEnergyGain
			
			// Cap energy
			if fox.Animal.Energy > 150 {
				fox.Animal.Energy = 150
			}
			
			// Remove rabbit
			w.removeRabbit(i)
			
			log.Printf("Fox hunted rabbit at (%d,%d)! Rabbits left: %d", pos.X, pos.Y, len(w.Rabbits))
			return
		}
	}
}

// moveFox moves a fox to a random adjacent position
func (w *World) moveFox(fox *Fox) {
	// Clear current position
	w.Grid[fox.Animal.Position.X][fox.Animal.Position.Y] = Empty
	
	// Get possible moves (adjacent cells)
	moves := w.getAdjacentPositions(fox.Animal.Position)
	
	// Filter for empty positions, grass positions, or rabbit positions (can hunt)
	validMoves := make([]Position, 0)
	for _, pos := range moves {
		cellType := w.Grid[pos.X][pos.Y]
		if cellType == Empty || cellType == GrassType || cellType == RabbitType {
			validMoves = append(validMoves, pos)
		}
	}
	
	// Move to random valid position, or stay if no valid moves
	if len(validMoves) > 0 {
		newPos := validMoves[rand.Intn(len(validMoves))]
		fox.Animal.Position = newPos
	}
	
	// Update grid with new position
	w.Grid[fox.Animal.Position.X][fox.Animal.Position.Y] = FoxType
}

// tryFoxReproduction attempts fox reproduction with nearby fox
func (w *World) tryFoxReproduction(fox *Fox) {
	// Check population limit
	if len(w.Foxes) >= maxFoxes {
		return // Too many foxes already
	}
	
	adjacentPositions := w.getAdjacentPositions(fox.Animal.Position)
	
	for _, pos := range adjacentPositions {
		if w.Grid[pos.X][pos.Y] == FoxType {
			partner := w.findFoxAtPosition(pos)
			if partner != nil && 
			   partner.Animal.Energy >= foxReproduceThreshold && 
			   partner.Animal.ReproduceCD == 0 {
				
				// Find empty position for baby fox
				for _, babyPos := range w.getAdjacentPositions(fox.Animal.Position) {
					if w.Grid[babyPos.X][babyPos.Y] == Empty {
						// Create baby fox
						baby := &Fox{
							Animal: Animal{
								Position:    babyPos,
								Energy:      60,
								ReproduceCD: reproductionCooldown,
								Age:         0,
							},
						}
						
						w.Foxes = append(w.Foxes, baby)
						w.Grid[babyPos.X][babyPos.Y] = FoxType
						
						// Parents lose energy and get cooldown
						fox.Animal.Energy -= 30
						partner.Animal.Energy -= 30
						fox.Animal.ReproduceCD = reproductionCooldown
						partner.Animal.ReproduceCD = reproductionCooldown
						
						log.Printf("New fox born at (%d,%d)! Total foxes: %d", babyPos.X, babyPos.Y, len(w.Foxes))
						return
					}
				}
			}
		}
	}
}

// findFoxAtPosition finds fox at given position
func (w *World) findFoxAtPosition(pos Position) *Fox {
	for _, fox := range w.Foxes {
		if fox.Animal.Position.X == pos.X && fox.Animal.Position.Y == pos.Y {
			return fox
		}
	}
	return nil
}

// removeFox removes a fox from the world
func (w *World) removeFox(index int) {
	fox := w.Foxes[index]
	
	// Log fox death for debugging
	log.Printf("Fox died at (%d,%d) with energy %d! Foxes left: %d", fox.Animal.Position.X, fox.Animal.Position.Y, fox.Animal.Energy, len(w.Foxes)-1)
	
	// Clear grid position
	w.Grid[fox.Animal.Position.X][fox.Animal.Position.Y] = Empty
	
	// Remove from slice
	w.Foxes = append(w.Foxes[:index], w.Foxes[index+1:]...)
}