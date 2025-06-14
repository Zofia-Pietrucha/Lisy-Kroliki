package main

import (
	"fmt"
	"log"
	"math/rand"
)

type Animal struct {
	Position
	Energy       int
	ReproduceCD  int // Cooldown after reproduction
	Age          int
}

type Rabbit struct {
	Animal
	NewBorn int // Ticks since birth (for visual indication)
}

type Fox struct {
	Animal
}

func (w *World) updateRabbits() {
	// Handle all rabbit updates without reproduction
	for i := len(w.Rabbits) - 1; i >= 0; i-- {
		rabbit := w.Rabbits[i]
		
		rabbit.Animal.Age++
		if rabbit.Animal.ReproduceCD > 0 {
			rabbit.Animal.ReproduceCD--
		}
		
		if rabbit.NewBorn > 0 {
			rabbit.NewBorn--
		}
		
		// Lose energy only every 60 ticks (roughly once per second)
		if w.Tick%60 == 0 {
			rabbit.Animal.Energy -= rabbitEnergyLoss
		}
		
		w.rabbitEatGrass(rabbit)
		
		if rand.Float64() < rabbitMoveChance {
			w.moveRabbit(rabbit)
			w.rabbitEatGrass(rabbit)
		}
		
		if rabbit.Animal.Energy <= 0 {
			w.removeRabbit(i)
		}
	}
	
	// Handle reproduction in separate pass to avoid multiple births per tick
	w.handleRabbitReproduction()
}

func (w *World) rabbitEatGrass(rabbit *Rabbit) {
	pos := rabbit.Animal.Position
	grass, exists := w.Grass[pos]
	
	if exists && grass.Amount >= minGrassToEat {
		rabbit.Animal.Energy += grassEnergyGain
		
		if rabbit.Animal.Energy > 100 {
			rabbit.Animal.Energy = 100
		}
		
		delete(w.Grass, pos)
	}
}

func (w *World) moveRabbit(rabbit *Rabbit) {
	w.Grid[rabbit.Animal.Position.X][rabbit.Animal.Position.Y] = Empty
	
	moves := w.getAdjacentPositions(rabbit.Animal.Position)
	
	validMoves := make([]Position, 0)
	for _, pos := range moves {
		cellType := w.Grid[pos.X][pos.Y]
		if cellType == Empty || cellType == GrassType {
			validMoves = append(validMoves, pos)
		}
	}
	
	if len(validMoves) > 0 {
		newPos := validMoves[rand.Intn(len(validMoves))]
		rabbit.Animal.Position = newPos
	}
	
	w.Grid[rabbit.Animal.Position.X][rabbit.Animal.Position.Y] = RabbitType
}

func (w *World) handleRabbitReproduction() {
	processedPairs := make(map[string]bool)
	
	for _, rabbit := range w.Rabbits {
		if rabbit.Animal.Energy < reproduceEnergyThreshold || rabbit.Animal.ReproduceCD > 0 {
			continue
		}
		
		if rand.Float64() >= reproduceChance {
			continue
		}
		
		partner := w.findAdjacentReproductivePartner(rabbit, processedPairs)
		if partner != nil {
			pairKey1 := fmt.Sprintf("%d,%d-%d,%d", rabbit.Animal.Position.X, rabbit.Animal.Position.Y, partner.Animal.Position.X, partner.Animal.Position.Y)
			pairKey2 := fmt.Sprintf("%d,%d-%d,%d", partner.Animal.Position.X, partner.Animal.Position.Y, rabbit.Animal.Position.X, rabbit.Animal.Position.Y)
			processedPairs[pairKey1] = true
			processedPairs[pairKey2] = true
			
			w.createBabyRabbit(rabbit, partner)
		}
	}
}

func (w *World) findAdjacentReproductivePartner(rabbit *Rabbit, processedPairs map[string]bool) *Rabbit {
	adjacentPositions := w.getAdjacentPositions(rabbit.Animal.Position)
	
	for _, pos := range adjacentPositions {
		if w.Grid[pos.X][pos.Y] == RabbitType {
			partner := w.findRabbitAtPosition(pos)
			if partner != nil && 
			   partner.Animal.Energy >= reproduceEnergyThreshold && 
			   partner.Animal.ReproduceCD == 0 {
				
				pairKey := fmt.Sprintf("%d,%d-%d,%d", rabbit.Animal.Position.X, rabbit.Animal.Position.Y, partner.Animal.Position.X, partner.Animal.Position.Y)
				if !processedPairs[pairKey] {
					return partner
				}
			}
		}
	}
	return nil
}

func (w *World) findRabbitAtPosition(pos Position) *Rabbit {
	for _, rabbit := range w.Rabbits {
		if rabbit.Animal.Position.X == pos.X && rabbit.Animal.Position.Y == pos.Y {
			return rabbit
		}
	}
	return nil
}

func (w *World) createBabyRabbit(parent1, parent2 *Rabbit) {
	if len(w.Rabbits) >= maxRabbits {
		return
	}
	
	adjacentPositions := w.getAdjacentPositions(parent1.Animal.Position)
	
	for _, pos := range adjacentPositions {
		if w.Grid[pos.X][pos.Y] == Empty {
			baby := &Rabbit{
				Animal: Animal{
					Position:    pos,
					Energy:      60,
					ReproduceCD: reproductionCooldown,
					Age:         0,
				},
				NewBorn: 180, // Visual indicator for 30 seconds (180 ticks at 6 FPS)
			}
			
			w.Rabbits = append(w.Rabbits, baby)
			w.Grid[pos.X][pos.Y] = RabbitType
			
			parent1.Animal.Energy -= 20
			parent2.Animal.Energy -= 20
			parent1.Animal.ReproduceCD = reproductionCooldown
			parent2.Animal.ReproduceCD = reproductionCooldown
			
			log.Printf("New rabbit born at (%d,%d)! Total rabbits: %d", pos.X, pos.Y, len(w.Rabbits))
			
			return
		}
	}
}

func (w *World) removeRabbit(index int) {
	rabbit := w.Rabbits[index]
	
	w.Grid[rabbit.Animal.Position.X][rabbit.Animal.Position.Y] = Empty
	
	w.Rabbits = append(w.Rabbits[:index], w.Rabbits[index+1:]...)
}

func (w *World) updateFoxes() {
	for i := len(w.Foxes) - 1; i >= 0; i-- {
		fox := w.Foxes[i]
		
		fox.Animal.Age++
		if fox.Animal.ReproduceCD > 0 {
			fox.Animal.ReproduceCD--
		}
		
		// Lose energy only every 60 ticks
		if w.Tick%60 == 0 {
			fox.Animal.Energy -= foxEnergyLoss
		}
		
		w.foxHuntRabbit(fox)
		
		if rand.Float64() < foxMoveChance {
			if w.smartHunting {
				w.moveFoxSmart(fox)
			} else {
				w.moveFoxHunting(fox)
			}
			w.foxHuntRabbit(fox)
		}
		
		if fox.Animal.Energy >= foxReproduceThreshold && fox.Animal.ReproduceCD == 0 {
			if rand.Float64() < reproduceChance*1.5 {
				w.tryFoxReproduction(fox)
			}
		}
		
		if fox.Animal.Energy <= 0 {
			w.removeFox(i)
		}
	}
}

// moveFoxSmart uses enhanced vision to hunt rabbits more effectively
func (w *World) moveFoxSmart(fox *Fox) {
	w.Grid[fox.Animal.Position.X][fox.Animal.Position.Y] = Empty
	
	targetRabbit := w.findNearestRabbit(fox.Animal.Position)
	
	var newPos Position
	if targetRabbit != nil {
		newPos = w.moveTowardsTarget(fox.Animal.Position, *targetRabbit)
		log.Printf("Fox at (%d,%d) spotted rabbit at (%d,%d), moving towards it", 
			fox.Animal.Position.X, fox.Animal.Position.Y, targetRabbit.X, targetRabbit.Y)
	} else {
		moves := w.getAdjacentPositions(fox.Animal.Position)
		validMoves := make([]Position, 0)
		for _, pos := range moves {
			cellType := w.Grid[pos.X][pos.Y]
			if cellType == Empty || cellType == GrassType || cellType == RabbitType {
				validMoves = append(validMoves, pos)
			}
		}
		
		if len(validMoves) > 0 {
			newPos = validMoves[rand.Intn(len(validMoves))]
		} else {
			newPos = fox.Animal.Position
		}
	}
	
	fox.Animal.Position = newPos
	w.Grid[fox.Animal.Position.X][fox.Animal.Position.Y] = FoxType
}

// findNearestRabbit looks for rabbits within fox vision range
func (w *World) findNearestRabbit(foxPos Position) *Position {
	var nearestRabbit *Position
	minDistance := foxVisionRange + 1
	
	for dx := -foxVisionRange; dx <= foxVisionRange; dx++ {
		for dy := -foxVisionRange; dy <= foxVisionRange; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			
			x := foxPos.X + dx
			y := foxPos.Y + dy
			
			if x < 0 || x >= gridWidth || y < 0 || y >= gridHeight {
				continue
			}
			
			if w.Grid[x][y] == RabbitType {
				distance := abs(dx) + abs(dy)
				if distance < minDistance {
					minDistance = distance
					pos := Position{x, y}
					nearestRabbit = &pos
				}
			}
		}
	}
	
	return nearestRabbit
}

// moveTowardsTarget calculates the best move towards a target position
func (w *World) moveTowardsTarget(current, target Position) Position {
	moves := w.getAdjacentPositions(current)
	
	var bestMove Position = current
	bestDistance := abs(current.X - target.X) + abs(current.Y - target.Y)
	
	for _, move := range moves {
		cellType := w.Grid[move.X][move.Y]
		if cellType == Empty || cellType == GrassType || cellType == RabbitType {
			distance := abs(move.X - target.X) + abs(move.Y - target.Y)
			if distance < bestDistance {
				bestDistance = distance
				bestMove = move
			}
		}
	}
	
	return bestMove
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (w *World) moveFoxHunting(fox *Fox) {
	w.Grid[fox.Animal.Position.X][fox.Animal.Position.Y] = Empty
	
	moves := w.getAdjacentPositions(fox.Animal.Position)
	
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
		newPos := validMoves[rand.Intn(len(validMoves))]
		fox.Animal.Position = newPos
	}
	
	w.Grid[fox.Animal.Position.X][fox.Animal.Position.Y] = FoxType
}

func (w *World) foxHuntRabbit(fox *Fox) {
	pos := fox.Animal.Position
	
	for i, rabbit := range w.Rabbits {
		if rabbit.Animal.Position.X == pos.X && rabbit.Animal.Position.Y == pos.Y {
			fox.Animal.Energy += rabbitEnergyGain
			
			if fox.Animal.Energy > 150 {
				fox.Animal.Energy = 150
			}
			
			w.removeRabbit(i)
			
			log.Printf("Fox hunted rabbit at (%d,%d)! Rabbits left: %d", pos.X, pos.Y, len(w.Rabbits))
			return
		}
	}
}

func (w *World) moveFox(fox *Fox) {
	w.Grid[fox.Animal.Position.X][fox.Animal.Position.Y] = Empty
	
	moves := w.getAdjacentPositions(fox.Animal.Position)
	
	validMoves := make([]Position, 0)
	for _, pos := range moves {
		cellType := w.Grid[pos.X][pos.Y]
		if cellType == Empty || cellType == GrassType || cellType == RabbitType {
			validMoves = append(validMoves, pos)
		}
	}
	
	if len(validMoves) > 0 {
		newPos := validMoves[rand.Intn(len(validMoves))]
		fox.Animal.Position = newPos
	}
	
	w.Grid[fox.Animal.Position.X][fox.Animal.Position.Y] = FoxType
}

func (w *World) tryFoxReproduction(fox *Fox) {
	if len(w.Foxes) >= maxFoxes {
		return
	}
	
	adjacentPositions := w.getAdjacentPositions(fox.Animal.Position)
	
	for _, pos := range adjacentPositions {
		if w.Grid[pos.X][pos.Y] == FoxType {
			partner := w.findFoxAtPosition(pos)
			if partner != nil && 
			   partner.Animal.Energy >= foxReproduceThreshold && 
			   partner.Animal.ReproduceCD == 0 {
				
				for _, babyPos := range w.getAdjacentPositions(fox.Animal.Position) {
					if w.Grid[babyPos.X][babyPos.Y] == Empty {
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

func (w *World) findFoxAtPosition(pos Position) *Fox {
	for _, fox := range w.Foxes {
		if fox.Animal.Position.X == pos.X && fox.Animal.Position.Y == pos.Y {
			return fox
		}
	}
	return nil
}

func (w *World) removeFox(index int) {
	fox := w.Foxes[index]
	
	log.Printf("Fox died at (%d,%d) with energy %d! Foxes left: %d", fox.Animal.Position.X, fox.Animal.Position.Y, fox.Animal.Energy, len(w.Foxes)-1)
	
	w.Grid[fox.Animal.Position.X][fox.Animal.Position.Y] = Empty
	
	w.Foxes = append(w.Foxes[:index], w.Foxes[index+1:]...)
}