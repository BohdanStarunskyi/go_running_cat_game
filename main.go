package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"

	"go-cat/core/game"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

type GameWithAudio struct {
	game         *game.Game
	audioContext *audio.Context
}

func (g *GameWithAudio) Update() error {
	return g.game.Update(g.audioContext)
}
func (g *GameWithAudio) Draw(screen *ebiten.Image) {
	g.game.Draw(screen)
}
func (g *GameWithAudio) Layout(w, h int) (int, int) {
	return g.game.Layout(w, h)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	imageAssets := []string{"assets/bg.png", "assets/cat.gif", "assets/box.png"}
	audioAssets := []string{"assets/bg.mp3", "assets/running.mp3"}

	for _, asset := range imageAssets {
		if _, err := os.Stat(asset); err != nil {
			log.Fatalf("missing required asset %s: %v", asset, err)
		}
	}

	for _, asset := range audioAssets {
		if _, err := os.Stat(asset); err != nil {
			log.Printf("Warning: missing audio asset %s: %v", asset, err)
		} else {
			log.Printf("Found audio asset: %s", asset)
		}
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Go Jumpy Cat with Rain and Sound")

	audioContext := audio.NewContext(44100)
	g := &GameWithAudio{
		game:         game.NewGame(audioContext),
		audioContext: audioContext,
	}

	log.Printf("Starting game...")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
