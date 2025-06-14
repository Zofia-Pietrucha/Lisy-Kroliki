# Ecosystem Simulation

Symulacja ekosystemu z trawą, królikami i lisami napisana w Go z użyciem biblioteki Ebiten.

## Wymagania

- Go 1.21+
- Linux (projekt testowany pod Linux)

## Uruchamianie

```bash
# Pobranie zależności
go mod download

# Uruchomienie
go run main.go

# Budowanie
go build -o ecosystem-sim main.go

# Budowanie zoptymalizowanej wersji dla Linux
go build -ldflags="-s -w" -o ecosystem-sim-optimized main.go
```

## Status

- [x] Setup projektu
- [x] Podstawowe okno gry
- [x] System świata i entity
- [x] Renderowanie
- [ ] Logika symulacji
