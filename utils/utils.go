package utils

import (
	"go-cat/gameobjects"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func Clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func RectsOverlap(a, b [4]float64) bool {
	return a[0] < b[2] && a[2] > b[0] && a[1] < b[3] && a[3] > b[1]
}

func DrawRain(screen *ebiten.Image, rainDrops []gameobjects.Raindrop) {
	rainColor := color.RGBA{150, 150, 255, 180}
	for _, drop := range rainDrops {
		ebitenutil.DrawLine(screen, drop.X, drop.Y, drop.X, drop.Y+drop.Length, rainColor)
	}
}

func DebugPrintAt(screen *ebiten.Image, msg string, x, y int) {
	ebitenutil.DebugPrintAt(screen, msg, x, y)
}

func IsKeyPressedSpace() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}

func IsKeyPressedR() bool {
	return ebiten.IsKeyPressed(ebiten.KeyR)
}
