package main

import "math/rand"

type Grass struct {
	Position
	Amount int // 0-100, where 100 is fully grown
}

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
		
		if w.Grid[x][y] == Empty {
			if rand.Float64() < grassSpawnChance {
				pos := Position{x, y}
				w.Grass[pos] = &Grass{
					Position: pos,
					Amount:   grassGrowthRate,
				}
				w.Grid[x][y] = GrassType
			}
		}
	}
}