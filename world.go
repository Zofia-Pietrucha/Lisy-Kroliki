package main

import "math/rand"

type World struct {
	Grid    [][]EntityType
	Grass   map[Position]*Grass
	Rabbits []*Rabbit
	Foxes   []*Fox
	Tick    int
	smartHunting bool
}

func NewWorld() *World {
	w := &World{
		Grid:    make([][]EntityType, gridWidth),
		Grass:   make(map[Position]*Grass),
		Rabbits: make([]*Rabbit, 0),
		Foxes:   make([]*Fox, 0),
		Tick:    0,
		smartHunting: foxSmartHunting,
	}
	
	for x := 0; x < gridWidth; x++ {
		w.Grid[x] = make([]EntityType, gridHeight)
	}
	
	return w
}

func (w *World) Update() {
	w.updateGrass()
	w.updateRabbits()
	w.updateFoxes()
}

func (w *World) getAdjacentPositions(pos Position) []Position {
	adjacent := make([]Position, 0, 8)
	
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			
			newX := pos.X + dx
			newY := pos.Y + dy
			
			if newX >= 0 && newX < gridWidth && newY >= 0 && newY < gridHeight {
				adjacent = append(adjacent, Position{newX, newY})
			}
		}
	}
	
	return adjacent
}

func (w *World) addTestEntities() {
	for i := 0; i < 30; i++ {
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		pos := Position{x, y}
		
		w.Grass[pos] = &Grass{
			Position: pos,
			Amount:   rand.Intn(51) + 50,
		}
		w.Grid[x][y] = GrassType
	}
	
	for group := 0; group < 3; group++ {
		centerX := rand.Intn(gridWidth-10) + 5
		centerY := rand.Intn(gridHeight-10) + 5
		
		for i := 0; i < 3+rand.Intn(2); i++ {
			x := centerX + rand.Intn(6) - 3
			y := centerY + rand.Intn(6) - 3
			
			if x >= 0 && x < gridWidth && y >= 0 && y < gridHeight && w.Grid[x][y] == Empty {
				rabbit := &Rabbit{
					Animal: Animal{
						Position:    Position{x, y},
						Energy:      80,
						ReproduceCD: 0,
						Age:         0,
					},
					NewBorn: 0,
				}
				
				w.Rabbits = append(w.Rabbits, rabbit)
				w.Grid[x][y] = RabbitType
			}
		}
	}
	
	for group := 0; group < 2; group++ {
		centerX := rand.Intn(gridWidth-6) + 3
		centerY := rand.Intn(gridHeight-6) + 3
		
		for i := 0; i < 2+rand.Intn(2); i++ {
			x := centerX + rand.Intn(4) - 2
			y := centerY + rand.Intn(4) - 2
			
			if x >= 0 && x < gridWidth && y >= 0 && y < gridHeight && w.Grid[x][y] == Empty {
				fox := &Fox{
					Animal: Animal{
						Position:    Position{x, y},
						Energy:      80,
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