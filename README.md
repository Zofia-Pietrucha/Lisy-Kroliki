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
- [x] System trawy (wzrost i pojawianie się)
- [x] Króliki z podstawowym ruchem
- [x] Króliki jedzą trawę
- [x] Rozmnażanie królików
- [x] Refactoring na Animal + Rabbit/Fox struktury
- [x] Lisy i polowanie
- [x] Balansowanie ekosystemu
- [x] Finalne dopracowanie i dokumentacja

## Projekt zakończony!

Wszystkie wymagania zostały spełnione:
✅ Graficzna symulacja w Go z Ebiten
✅ Współistnienie trawy, królików i lisów
✅ Wzrost i regeneracja trawy
✅ Ruch i żywienie królików
✅ Rozmnażanie zwierząt z cooldownem
✅ Polowanie lisów na króliki
✅ System energii i śmierci z głodu
✅ Dynamika populacji i balans ekosystemu
✅ Działanie pod Linux
✅ Dokumentacja i kod źródłowy
