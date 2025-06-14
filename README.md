# Ecosystem Simulation

Symulacja ekosystemu z trawą, królikami i lisami napisana w Go z użyciem biblioteki Ebiten.

## Opis projektu

Projekt implementuje symulację ekosystemu zawierającego:

- **Trawę** - rośnie na pustych polach, regeneruje się z czasem
- **Króliki** - poruszają się, jedzą trawę, rozmnażają się
- **Lisy** - polują na króliki, rozmnażają się po udanych polowaniach

## Mechaniki symulacji

### Trawa

- Pojawia się losowo na pustych polach (1% szansy na tick)
- Rośnie do maksymalnej wartości 100 punktów
- Różne odcienie zieleni w zależności od dojrzałości

### Króliki (białe/żółte punkty)

- Poruszają się losowo po planszy (70% szansy na ruch)
- Jedzą trawę gdy mają min. 5 punktów (zyskują 35 energii)
- Tracą 1 energię co sekundę
- Rozmnażają się gdy mają 65+ energii i spotykają partnera
- Nowo narodzone króliki są żółte przez ~30 sekund
- Maksymalna populacja: 50 osobników

### Lisy (czerwone punkty)

- Poruszają się aktywnie w poszukiwaniu królików (60% szansy na ruch)
- Polują na króliki (zyskują 50 energii za każdego)
- Preferują ruchy w kierunku królików w sąsiednich polach
- Rozmnażają się gdy mają 70+ energii
- Maksymalna populacja: 15 osobników

### Dynamika ekosystemu

- Naturalna konkurencja: więcej trawy → więcej królików → więcej lisów → mniej królików
- Cykliczne zmiany populacji przypominające rzeczywiste ekosystemy
- Limity populacji zapobiegają nierealistycznym eksplozjom

## Wymagania

- Go 1.21+
- Linux (projekt testowany pod Linux)
- Biblioteki systemowe dla Ebiten (zazwyczaj dostępne domyślnie)

## Uruchamianie

```bash
# Pobranie zależności
go mod download

# Uruchomienie symulacji
go run .

# Budowanie
go build -o ecosystem-sim .

# Budowanie zoptymalizowanej wersji dla Linux
go build -ldflags="-s -w" -o ecosystem-sim-optimized .
```

## Struktura projektu

```
ecosystem-sim/
├── main.go         # Główna aplikacja i interfejs
├── constants.go    # Parametry symulacji
├── world.go        # Logika świata i inicjalizacja
├── animals.go      # Logika królików i lisów
├── grass.go        # System trawy
├── rendering.go    # Funkcje rysowania
├── go.mod          # Definicja modułu
└── README.md       # Dokumentacja
```

## Parametry symulacji

Wszystkie parametry można łatwo dostosować w pliku `constants.go`:

- Szybkość poruszania się zwierząt
- Tempo wzrostu trawy
- Wymagania energetyczne
- Progi reprodukcji
- Limity populacji

## Kontrolki

- **Spacja** - pauza/wznowienie symulacji
- **1** - tryb rysowania królików (kliknij myszą żeby postawić)
- **2** - tryb rysowania lisów (kliknij myszą żeby postawić)
- **0** - tryb normalny (bez rysowania)
- **S** - zapisz dane populacji do pliku CSV
- **Mysz** - kliknij przyciski Pause/Play/Reset lub rysuj zwierzęta

## Eksport danych

Symulacja automatycznie zapisuje dane populacji do pliku CSV po zamknięciu programu. Plik zawiera:

- Znaczniki czasowe każdego pomiaru
- Liczby królików, lisów i trawy w czasie
- Metadane symulacji (parametry, statystyki)
- Format gotowy do analizy w Excel lub innych narzędziach

Można też zapisać dane ręcznie klawiszem **S** podczas symulacji.

## Obserwacje z symulacji

1. **Cykle populacyjne** - populacje oscylują w naturalnych cyklach
2. **Konkurencja o zasoby** - króliki konkurują o trawę
3. **Relacje drapieżnik-ofiara** - lisy kontrolują populację królików
4. **Równowaga ekologiczna** - system dąży do stabilnego stanu

## Rozwój projektu

Projekt przeszedł przez następujące etapy:

1. Setup projektu i podstawowe okno
2. Struktury danych i renderowanie
3. System trawy
4. Logika królików (ruch, jedzenie, reprodukcja)
5. Refactoring do modularnej struktury
6. Implementacja lisów i polowania
7. Balansowanie i optymalizacja

## Możliwe rozszerzenia

- Większy zasięg widzenia dla lisów
- System tropienia zapachu
- Różne typy trawy o różnej wartości odżywczej
- Sezonowość (zima/lato wpływające na wzrost trawy)
- Choroby i epidemie
- Migracje zwierząt
- Interfejs użytkownika do zmiany parametrów w czasie rzeczywistym
