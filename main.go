package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// Game represents the main game state
type Game struct {
	// Będziemy dodawać pola w kolejnych commitach
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Logika gry będzie tutaj
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Wypełnij tło ciemnym zielonym kolorem
	screen.Fill(color.RGBA{34, 139, 34, 255}) // Forest green
	
	// Wyświetl podstawowe info
	ebitenutil.DebugPrint(screen, "Ecosystem Simulation\nPress ESC to quit")
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