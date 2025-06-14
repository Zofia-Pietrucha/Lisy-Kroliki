package main

const (
	screenWidth  = 800
	screenHeight = 600
	gridWidth    = 80
	gridHeight   = 60
	cellSize     = 10

	gameAreaHeight = 400
	graphHeight    = 150
	graphWidth     = 750
	graphOffsetX   = 25
	graphOffsetY   = 420

	maxHistoryPoints = 150

	maxGrassAmount   = 100
	grassGrowthRate  = 2
	grassSpawnChance = 0.01

	rabbitMoveChance = 0.7
	rabbitEnergyLoss = 1
	grassEnergyGain  = 40
	minGrassToEat    = 5

	reproduceEnergyThreshold = 60
	reproductionCooldown     = 180
	reproduceChance          = 0.3

	foxMoveChance         = 0.6
	foxEnergyLoss         = 1
	rabbitEnergyGain      = 50
	foxReproduceThreshold = 70

	foxVisionRange  = 3
	foxSmartHunting = true

	// Population limits prevent overpopulation
	maxRabbits = 50
	maxFoxes   = 15
)

type EntityType int

const (
	Empty EntityType = iota
	GrassType
	RabbitType
	FoxType
)

type Position struct {
	X, Y int
}

type PopulationData struct {
	Tick    int
	Rabbits int
	Foxes   int
	Grass   int
}