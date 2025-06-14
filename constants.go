package main

const (
	screenWidth  = 800
	screenHeight = 600
	gridWidth    = 80
	gridHeight   = 60
	cellSize     = 10

	// Grass parameters
	maxGrassAmount   = 100
	grassGrowthRate  = 2    // How much grass grows per tick
	grassSpawnChance = 0.01 // Chance for new grass to appear on empty cells

	// Rabbit parameters
	rabbitMoveChance = 0.7 // Chance rabbit moves each tick
	rabbitEnergyLoss = 1   // Energy lost per tick every 60 ticks
	grassEnergyGain  = 40  // Energy gained from eating grass
	minGrassToEat    = 5   // Minimum grass amount to be edible

	// Reproduction parameters
	reproduceEnergyThreshold = 60  // Min energy to reproduce
	reproductionCooldown     = 180 // Ticks before can reproduce again
	reproduceChance          = 0.3 // Chance to try reproduction each tick

	// Fox parameters
	foxMoveChance         = 0.6 // Chance fox moves each tick (zwiększone z 0.4)
	foxEnergyLoss         = 1   // Energy lost per tick every 60 ticks
	rabbitEnergyGain      = 50  // Energy gained from eating rabbit
	foxReproduceThreshold = 70  // Min energy for fox reproduction (zmniejszone z 80)

	// Population limits (prevent overpopulation)
	maxRabbits = 50 // Maximum rabbit population
	maxFoxes   = 15 // Maximum fox population
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