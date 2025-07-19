# Go Cat Game

A simple, fun endless runner game written in Go using Ebiten. Control a running cat, jump over obstacles, and survive as long as you can! The game features animated sprites, sound effects, rain, and randomly selected obstacles from your asset folder.

## Game Description
- You play as a cat running across the screen.
- Jump over various obstacles (random images from the assets folder).
- Rain and background music add atmosphere.
- The game gets harder as you survive longer!
- Press `Space` to jump, `R` to restart after losing.

## Tech Stack
- **Go** (Golang)
- **Ebiten** (https://ebiten.org/) for 2D game rendering and input
- **Ebiten/audio** for sound and music
- **Standard Go image/audio libraries** for asset loading

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
1. **Install Go** (https://golang.org/dl/)
2. **Install Ebiten** (if not already):
   ```
go get github.com/hajimehoshi/ebiten/v2
   ```
3. **Clone this repository** and `cd` into the project folder.
4. **Add your assets** (images and audio) to the `assets/` folder. The game will automatically use all PNG/JPEG images (except backgrounds/cat/ground) as obstacles.
5. **Run the game:**
   ```
go run .
   

## Gameplay video

https://github.com/user-attachments/assets/aa07a27e-69f0-4ad5-9a74-146beeaf7fcd

