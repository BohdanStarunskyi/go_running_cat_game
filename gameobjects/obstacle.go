package gameobjects

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Obstacle struct {
	X, Y  float64
	Scale float64
	Img   *ebiten.Image
}

func (o *Obstacle) GetWidth() float64 {
	return float64(o.Img.Bounds().Dx()) * o.Scale
}

func (o *Obstacle) GetHeight() float64 {
	return float64(o.Img.Bounds().Dy()) * o.Scale
}

func (o *Obstacle) GetRect() [4]float64 {
	return [4]float64{o.X, o.Y, o.X + o.GetWidth(), o.Y + o.GetHeight()}
}

func CreateObstacle(x, y, scale float64, img *ebiten.Image) *Obstacle {
	return &Obstacle{
		X:     x,
		Y:     y,
		Scale: scale,
		Img:   img,
	}
}

func SpawnRandomObstacle(img *ebiten.Image, groundY float64, screenWidth int) *Obstacle {
	scale := 5.0
	scaledHeight := float64(img.Bounds().Dy()) * scale
	obstacleY := groundY - scaledHeight
	return CreateObstacle(float64(screenWidth)+200, obstacleY, scale, img)
}
