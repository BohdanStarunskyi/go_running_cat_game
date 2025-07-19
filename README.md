# Go Cat Game

A simple, fun endless runner game written in Go using Ebiten. Control a running cat, jump over obstacles, and survive as long as you can! The game features animated sprites, sound effects, rain, and randomly selected obstacles from your asset folder.

## Game Description

* You play as a cat running across the screen.
* Jump over various obstacles (random images from the assets folder).
* Rain and background music add atmosphere.
* The game gets harder as you survive longer!
* Press `Space` to jump, `R` to restart after losing.

## Tech Stack

* **Go** (Golang)
* **Ebiten** ([https://ebiten.org/](https://ebiten.org/)) for 2D game rendering and input
* **Ebiten/audio** for sound and music
* **Standard Go image/audio libraries** for asset loading

## Package Structure

```
assets/         # Game images and audio files
assetsgo/       # Asset loading and helpers (Go code)
core/game/      # Main game logic (Game struct, update/draw/logic)
gameobjects/    # Obstacle and raindrop types and logic
animation/      # Animator for sprite animation
utils/          # Utility functions (collision, drawing, input)
main.go         # Entry point, window setup, game loop
```

## How to Run

1. **Install Go**
   Download and install Go from the official site: [https://golang.org/dl/](https://golang.org/dl/)

2. **Clone the repository**
   Open a terminal and run:

   ```bash
   git clone https://github.com/your-username/your-repo-name.git
   cd your-repo-name
   ```

3. **Add assets**
   Put your game assets (images and audio) into the `assets/` folder:

   * `cat.gif` — your animated player sprite
   * `ground.png` — ground texture
   * `*.png`, `*.jpg`, or `*.jpeg` (excluding cat/ground) — obstacles
   * `bgm.mp3` — background music
   * `jump.mp3` — jump sound

4. **Download dependencies**
   If not done automatically, use:

   ```bash
   go mod tidy
   ```

5. **Run the game**
   Simply run:

   ```bash
   go run main.go
   ```

## Gameplay video

[https://github.com/user-attachments/assets/aa07a27e-69f0-4ad5-9a74-146beeaf7fcd](https://github.com/user-attachments/assets/aa07a27e-69f0-4ad5-9a74-146beeaf7fcd)
