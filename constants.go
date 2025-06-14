package main

const (
	screenWidth  = 800
	screenHeight = 600
	gridWidth    = 80
	gridHeight   = 60
	cellSize     = 10

	// UI layout
	gameAreaHeight = 400 // Upper area for simulation
	graphHeight    = 150 // Lower area for graph
	graphWidth     = 750 // Graph width
	graphOffsetX   = 25  // Graph left margin
	graphOffsetY   = 420 // Graph top position

	// Graph settings
	maxHistoryPoints = 150 // How many data points to keep

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
	foxMoveChance         = 0.6 // Chance fox moves each tick
	foxEnergyLoss         = 1   // Energy lost per tick every 60 ticks
	rabbitEnergyGain      = 50  // Energy gained from eating rabbit
	foxReproduceThreshold = 70  // Min energy for fox reproduction

	// Fox vision parameters
	foxVisionRange  = 3    // How many cells fox can see in each direction
	foxSmartHunting = true // Whether foxes use smart hunting AI

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

// PopulationData represents population counts at a specific time
type PopulationData struct {
	Tick    int
	Rabbits int
	Foxes   int
	Grass   int
}