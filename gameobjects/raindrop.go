package gameobjects

import "math/rand"

type Raindrop struct {
	X, Y   float64
	Length float64
	Speed  float64
}

func NewRaindrop(screenWidth, screenHeight int) Raindrop {
	return Raindrop{
		X:      rand.Float64() * float64(screenWidth),
		Y:      rand.Float64() * float64(screenHeight),
		Length: 10 + rand.Float64()*10,
		Speed:  400 + rand.Float64()*200,
	}
}
