package struktury_i_funkcje

import (
	"time"
	"image"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"
)

// Struktura opisujaca lzy
type Tear struct {
	Position pixel.Vec
	Velocity pixel.Vec
	StartTime time.Time
	Lifetime time.Duration
	Active bool
	Damage int
}

// Struktura opisujaca lzy bossa
type BossTear struct {
	Position pixel.Vec
	Velocity pixel.Vec
	Active bool
}

// Struktura opisujaca Gurdiego
type Gurdy struct {
	Position pixel.Vec
	Active bool
	Health int
}

// Struktura opisujaca Husha
type Hush struct {
	Position pixel.Vec
	Active bool
	Health int
}

func Abs(n float64) float64 {
	if n < 0 {
		return -n
	} else {
		return n
	}
}

func LoadPicture(path string) (pixel.Picture) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	return pixel.PictureDataFromImage(img)
}
