# Isaac – gra typu roguelike (Go + Pixel)

Projekt przedstawia prostą grę typu **roguelike / twin-stick shooter**, inspirowaną *The Binding of Isaac*. Gra została napisana w języku **Go** z wykorzystaniem biblioteki **Pixel** do obsługi grafiki, okna oraz wejścia z klawiatury.

---

## Opis gry

Gracz steruje postacią poruszającą się po losowo generowanych pokojach.  
Celem jest:
- eliminowanie przeciwników,
- zbieranie itemów,
- pokonanie bossów,
- przechodzenie na kolejne piętra.

Gra posiada minimapę, różne typy pomieszczeń, system przeciwników oraz bossów.

---

## Typy pomieszczeń

- **StartRoom** – pokój startowy
- **NormalRoom** – standardowy pokój z przeciwnikami
- **ItemRoom** – pokój z losowym itemem
- **BossRoom** – pokój z bossem i przejściem na kolejne piętro

Mapa jest generowana losowo przy każdym uruchomieniu gry.

---

## Przeciwnicy i bossowie

### Przeciwnicy
- **Fly** – szybki przeciwnik strzelający pociskami
- **Slime (Follower)** – wolniejszy przeciwnik podążający za graczem

### Bossowie
- **Gurdy** – boss pierwszego piętra  
- **Hush** – boss drugiego piętra  

Każdy boss posiada:
- własne punkty życia,
- pasek HP,
- unikalne wzorce ataków.

---

## Itemy

Itemy losowane są na piętro i mogą zwiększać statystyki gracza:

- **Tears Up** – szybsze strzelanie
- **Damage Up** – większe obrażenia
- **Range Up** – większy zasięg łez
- **Health Up** – dodatkowe serce
- **Piercing** (item po bossie) – łzy przebijają przeciwników

Itemy można podnosić klawiszem **E**.

---

## Sterowanie

| Klawisz | Akcja |
|------|------|
| W A S D | Ruch postaci |
| ↑ ↓ ← → | Strzelanie |
| E | Podnoszenie itemów |
| Zamknięcie okna | Wyjście z gry |

---

## Mechaniki gry

- zamykane drzwi w pokojach z przeciwnikami,
- minimapa pokazująca układ pomieszczeń,
- system kolizji (ściany, przeciwnicy, bossowie),
- przejście na kolejne piętro po pokonaniu bossa,
- koniec gry po 2 piętrach.

---

## Technologie

- **Go**
- **Pixel** (`github.com/faiface/pixel`)
- **pixelgl**
- **imdraw**
- **goroutines + mutexy** (logika przeciwników, bossów, pocisków)

---

## Uruchomienie projektu

1. Zainstaluj Go (>= 1.20)
2. Zainstaluj Pixel:
   ```bash
   go get github.com/faiface/pixel
   go get github.com/faiface/pixel/pixelgl
