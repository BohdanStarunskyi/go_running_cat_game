package assetsgo

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const sampleRate = 44100

func LoadImage(path string) (*ebiten.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

func LoadImageOrPanic(path string) *ebiten.Image {
	img, err := LoadImage(path)
	if err != nil {
		log.Fatalf("load image %s: %v", path, err)
	}
	return img
}

func LoadGifFrames(path string) ([]*ebiten.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	gifImg, err := gif.DecodeAll(f)
	if err != nil {
		return nil, err
	}
	frames := make([]*ebiten.Image, len(gifImg.Image))
	for i, frame := range gifImg.Image {
		frames[i] = ebiten.NewImageFromImage(frame)
	}
	return frames, nil
}

func LoadGifFramesOrPanic(path string) []*ebiten.Image {
	frames, err := LoadGifFrames(path)
	if err != nil {
		log.Fatalf("load gif %s: %v", path, err)
	}
	return frames
}

func GenerateGroundImage(screenWidth int) *ebiten.Image {
	width := screenWidth
	height := 80
	img := ebiten.NewImage(width, height)
	baseR, baseG, baseB := 76, 76, 76
	maxNoise := 20
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8(clampInt(baseR+rand.Intn(maxNoise*2+1)-maxNoise, 0, 255))
			g := uint8(clampInt(baseG+rand.Intn(maxNoise*2+1)-maxNoise, 0, 255))
			b := uint8(clampInt(baseB+rand.Intn(maxNoise*2+1)-maxNoise, 0, 255))
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return img
}

func clampInt(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

type audioBytes struct {
	data []byte
	pos  int64
}

func (ab *audioBytes) Read(p []byte) (n int, err error) {
	if ab.pos >= int64(len(ab.data)) {
		return 0, io.EOF
	}
	n = copy(p, ab.data[ab.pos:])
	ab.pos += int64(n)
	return n, nil
}

func (ab *audioBytes) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		ab.pos = offset
	case io.SeekCurrent:
		ab.pos += offset
	case io.SeekEnd:
		ab.pos = int64(len(ab.data)) + offset
	}
	if ab.pos < 0 {
		ab.pos = 0
	}
	if ab.pos > int64(len(ab.data)) {
		ab.pos = int64(len(ab.data))
	}
	return ab.pos, nil
}

func loadAudio(audioContext *audio.Context, path string) (*audio.Player, interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open audio file %s: %v", path, err)
	}
	defer f.Close()
	stream, err := mp3.DecodeWithSampleRate(sampleRate, f)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode MP3 %s: %v", path, err)
	}
	data, err := io.ReadAll(stream)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read audio data from %s: %v", path, err)
	}
	audioReader := &audioBytes{data: data}
	infiniteLoop := audio.NewInfiniteLoop(audioReader, int64(len(data)))
	player, err := audioContext.NewPlayer(infiniteLoop)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create player for %s: %v", path, err)
	}
	return player, infiniteLoop, nil
}

func InitAudioPlayers(audioContext *audio.Context) (*audio.Player, *audio.Player, *audio.Player, interface{}, interface{}, interface{}) {
	var bgPlayer, runPlayer, rainPlayer *audio.Player
	var bgStream, runStream, rainStream interface{}
	if player, stream, err := loadAudio(audioContext, "assets/bg.mp3"); err == nil {
		bgPlayer = player
		bgStream = stream
		bgPlayer.SetVolume(0.3)
		bgPlayer.Play()
	}
	if player, stream, err := loadAudio(audioContext, "assets/rain.mp3"); err == nil {
		rainPlayer = player
		rainStream = stream
		rainPlayer.SetVolume(0.15)
		rainPlayer.Play()
	}
	if player, stream, err := loadAudio(audioContext, "assets/running.mp3"); err == nil {
		runPlayer = player
		runStream = stream
		runPlayer.SetVolume(0.15)
		runPlayer.Play()
	}
	return bgPlayer, runPlayer, rainPlayer, bgStream, runStream, rainStream
}

func LoadObstacleImages() ([]*ebiten.Image, []string) {
	files, err := os.ReadDir("assets")
	if err != nil {
		log.Fatalf("failed to read assets dir: %v", err)
	}
	var images []*ebiten.Image
	var names []string
	exclude := map[string]bool{
		"bg.png":  true,
		"cat.gif": true,
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		ext := filepath.Ext(name)
		if (ext == ".png" || ext == ".jpeg" || ext == ".jpg") && !exclude[name] {
			img, err := LoadImage("assets/" + name)
			if err == nil {
				images = append(images, img)
				names = append(names, name)
			}
		}
	}
	return images, names
}
