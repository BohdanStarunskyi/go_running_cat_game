package animation

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animator struct {
	frames     []*ebiten.Image
	frameDelay time.Duration
	elapsed    time.Duration
	index      int
}

func NewAnimator(frames []*ebiten.Image, frameDelay time.Duration) *Animator {
	return &Animator{
		frames:     frames,
		frameDelay: frameDelay,
		elapsed:    0,
		index:      0,
	}
}

func (a *Animator) Update(delta time.Duration) {
	a.elapsed += delta
	for a.elapsed >= a.frameDelay {
		a.elapsed -= a.frameDelay
		a.index = (a.index + 1) % len(a.frames)
	}
}

func (a *Animator) CurrentFrame() *ebiten.Image {
	return a.frames[a.index]
}
