package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"

	"go-cat/animation"
	"go-cat/assetsgo"
	"go-cat/gameobjects"
	"go-cat/utils"
)

const (
	screenWidth   = 1280
	screenHeight  = 720
	catScale      = 0.2
	gravity       = 2.0
	jumpVelocity  = -35.0
	groundSpeed   = 20.0
	obstacleSpeed = 20.0
)

type Game struct {
	bg, ground     *ebiten.Image
	animator       *animation.Animator
	bgScale        float64
	groundOffset   float64
	catY, catVY    float64
	isJumping      bool
	catBaseY       float64
	obstacles      []*gameobjects.Obstacle
	gameOver       bool
	lastUpdateTime time.Time
	rainDrops      []gameobjects.Raindrop
	numRainDrops   int
	bgPlayer       *audio.Player
	runPlayer      *audio.Player
	rainPlayer     *audio.Player
	bgStream       interface{}
	runStream      interface{}
	rainStream     interface{}
	startTime      time.Time
	finalScore     int
	obstacleImages []*ebiten.Image
	obstacleNames  []string
	highScore      int
}

func NewGame(audioContext *audio.Context) *Game {
	bg := assetsgo.LoadImageOrPanic("assets/bg.png")
	catFrames := assetsgo.LoadGifFramesOrPanic("assets/cat.gif")
	ground := assetsgo.GenerateGroundImage(screenWidth)

	bgW, bgH := bg.Size()
	bgScaleX := float64(screenWidth) / float64(bgW)
	bgScaleY := float64(screenHeight) / float64(bgH)
	bgScale := bgScaleX
	if bgScaleY > bgScaleX {
		bgScale = bgScaleY
	}

	groundY := float64(screenHeight) - float64(ground.Bounds().Dy())
	_, ch := catFrames[0].Size()
	catBaseY := groundY - float64(ch)*catScale

	obstacleImages, obstacleNames := assetsgo.LoadObstacleImages()

	obstacles := make([]*gameobjects.Obstacle, 2)
	gap := float64(600 + rand.Intn(301))
	startX := float64(screenWidth)
	for i := 0; i < 2; i++ {
		img := obstacleImages[rand.Intn(len(obstacleImages))]
		obs := gameobjects.SpawnRandomObstacle(img, groundY, screenWidth)
		obs.X = startX
		obstacles[i] = obs
		startX += gap
		gap = float64(600 + rand.Intn(301))
	}

	numRainDrops := 200
	rainDrops := make([]gameobjects.Raindrop, numRainDrops)
	for i := range rainDrops {
		rainDrops[i] = gameobjects.NewRaindrop(screenWidth, screenHeight)
	}

	bgPlayer, runPlayer, rainPlayer, bgStream, runStream, rainStream := assetsgo.InitAudioPlayers(audioContext)

	return &Game{
		bg:             bg,
		ground:         ground,
		animator:       animation.NewAnimator(catFrames, 100*time.Millisecond),
		bgScale:        bgScale,
		obstacles:      obstacles,
		catBaseY:       catBaseY,
		catY:           0,
		catVY:          0,
		isJumping:      false,
		groundOffset:   0,
		lastUpdateTime: time.Now(),
		rainDrops:      rainDrops,
		numRainDrops:   numRainDrops,
		bgPlayer:       bgPlayer,
		runPlayer:      runPlayer,
		bgStream:       bgStream,
		runStream:      runStream,
		rainPlayer:     rainPlayer,
		rainStream:     rainStream,
		startTime:      time.Now(),
		obstacleImages: obstacleImages,
		obstacleNames:  obstacleNames,
		highScore:      0,
	}
}

func (g *Game) Close() {
	if g.bgPlayer != nil {
		g.bgPlayer.Close()
	}
	if g.runPlayer != nil {
		g.runPlayer.Close()
	}
	if g.rainPlayer != nil {
		g.rainPlayer.Close()
	}
}

func (g *Game) Update(audioContext *audio.Context) error {
	if g.gameOver {
		if utils.IsKeyPressedR() {
			high := g.highScore
			if g.finalScore > high {
				high = g.finalScore
			}
			g.Close()
			newGame := NewGame(audioContext)
			newGame.highScore = high
			*g = *newGame
		}
		return nil
	}

	now := time.Now()
	delta := now.Sub(g.lastUpdateTime)
	g.lastUpdateTime = now

	g.groundOffset -= groundSpeed
	if g.groundOffset < -float64(screenWidth) {
		g.groundOffset += float64(screenWidth)
	}

	if !g.isJumping && utils.IsKeyPressedSpace() {
		g.catVY = jumpVelocity
		g.isJumping = true
		if g.runPlayer != nil && g.runPlayer.IsPlaying() {
			g.runPlayer.Pause()
		}
	}

	if g.isJumping {
		g.catVY += gravity
		g.catY += g.catVY
		if g.catY > 0 {
			g.catY = 0
			g.catVY = 0
			g.isJumping = false
			if g.runPlayer != nil && !g.runPlayer.IsPlaying() {
				g.runPlayer.Play()
			}
		}
	}

	groundY := float64(screenHeight) - float64(g.ground.Bounds().Dy())
	for _, obs := range g.obstacles {
		obs.X -= obstacleSpeed
		if obs.X+obs.GetWidth() < 0 {
			maxX := float64(screenWidth)
			for _, o := range g.obstacles {
				if o.X > maxX {
					maxX = o.X
				}
			}
			img := g.obstacleImages[rand.Intn(len(g.obstacleImages))]
			*obs = *gameobjects.SpawnRandomObstacle(img, groundY, screenWidth)
			gap := float64(600 + rand.Intn(301))
			obs.X = maxX + gap
		}
	}

	g.animator.Update(delta)
	dt := delta.Seconds()
	for i := range g.rainDrops {
		g.rainDrops[i].Y += g.rainDrops[i].Speed * dt
		if g.rainDrops[i].Y > screenHeight {
			g.rainDrops[i].Y = -g.rainDrops[i].Length
			g.rainDrops[i].X = rand.Float64() * screenWidth
		}
	}

	catX := 50.0
	catW := float64(g.animator.CurrentFrame().Bounds().Dx()) * (catScale * 0.8)
	_, ch := g.animator.CurrentFrame().Size()
	catY := g.catBaseY + g.catY
	catRect := [4]float64{catX, catY, catX + catW, catY + float64(ch)*catScale}

	for _, obs := range g.obstacles {
		obsRect := obs.GetRect()
		if utils.RectsOverlap(catRect, obsRect) {
			g.gameOver = true
			if g.runPlayer != nil && g.runPlayer.IsPlaying() {
				g.runPlayer.Pause()
				g.rainPlayer.Pause()
			}
		}
	}

	if g.gameOver && g.finalScore == 0 {
		g.finalScore = int(time.Since(g.startTime).Milliseconds()) / 10
		if g.finalScore > g.highScore {
			g.highScore = g.finalScore
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.bgScale, g.bgScale)
	screen.DrawImage(g.bg, op)

	groundY := float64(screenHeight) - float64(g.ground.Bounds().Dy())
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.groundOffset, groundY)
	screen.DrawImage(g.ground, op)

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(g.groundOffset+float64(screenWidth), groundY)
	screen.DrawImage(g.ground, op2)

	catX := 50.0
	catY := g.catBaseY + g.catY
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(catScale, catScale)
	op.GeoM.Translate(catX, catY)
	screen.DrawImage(g.animator.CurrentFrame(), op)

	for _, obs := range g.obstacles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(obs.Scale, obs.Scale)
		op.GeoM.Translate(obs.X, obs.Y)
		screen.DrawImage(obs.Img, op)
	}

	utils.DrawRain(screen, g.rainDrops)

	if g.gameOver {
		utils.DebugPrintAt(screen, "GAME OVER! Press R to restart", screenWidth/2-150, screenHeight/2)
	}

	var score int
	if g.gameOver {
		score = g.finalScore
	} else {
		score = int(time.Since(g.startTime).Milliseconds()) / 10
	}
	utils.DebugPrintAt(screen, fmt.Sprintf("Score: %v", score), screenWidth-100, 10)
	utils.DebugPrintAt(screen, fmt.Sprintf("High Score: %v", g.highScore), screenWidth-100, 25)
}

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}
