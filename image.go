package main

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
)

type Color struct {
	r int
	g int
	b int
}

type Image struct {
	img    [][]Color
	zBuf   [][]float64
	height int
	width  int
}

func MakeImage(height, width int) *Image {
	img := make([][]Color, height)
	zBuf := make([][]float64, height)
	for i := range img {
		img[i] = make([]Color, width)
		zBuf[i] = make([]float64, width)
		for j := range zBuf[i] {
			zBuf[i][j] = -1 * math.MaxFloat64
		}
	}
	image := &Image{
		img:    img,
		zBuf:   zBuf,
		height: height,
		width:  width,
	}
	image.Clear()
	return image
}

// Should switch x and y l o l
func (image Image) plot(c Color, x, y int, z float64) error {
	if x < 0 || x > image.height || y < 0 || y > image.width {
		return errors.New("Error: Coordinate invalid")
	}
	z = float64(int(z*1000)) / 1000.0
	if z > image.zBuf[x][y] {
		image.img[x][y] = c
		image.zBuf[x][y] = z
	}
	return nil
}

func (image Image) plotNoZ(c Color, x, y int) error {
	if x < 0 || x > image.height || y < 0 || y > image.width {
		return errors.New("Error: Coordinate invalid")
	}
	image.img[x][y] = c
	image.zBuf[x][y] = -1 * math.MaxFloat64
	return nil
}

func (image Image) fill(c Color) {
	for y := 0; y < image.width; y++ {
		for x := 0; x < image.height; x++ {
			image.plotNoZ(c, x, y)
		}
	}
}

func (image Image) Clear() {
	image.fill(Color{r: 255, g: 255, b: 255})
}

func (image Image) SavePPM(filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	// TODO: Take variant max color
	buffer.WriteString(fmt.Sprintf("P3 %d %d 255\n", image.height, image.width))

	for y := 0; y < image.width; y++ {
		for x := 0; x < image.height; x++ {
			newY := image.width - 1 - y
			buffer.WriteString(fmt.Sprintf("%d %d %d ", image.img[x][newY].r, image.img[x][newY].g, image.img[x][newY].b))
		}
		buffer.WriteString("\n")
	}

	f.WriteString(buffer.String())
	f.Close()
	return nil
}

func (image Image) Display() error {
	f := "temp"
	image.SavePPM(f)
	c := exec.Command("display", f)
	_, err := c.Output()
	os.Remove(f)
	return err
}

func (c *Color) Limit() {
	if c.r < 0 {
		c.r = 0
	} else if c.r > 255 {
		c.r = 255
	}
	if c.g < 0 {
		c.g = 0
	} else if c.g > 255 {
		c.g = 255
	}
	if c.b < 0 {
		c.b = 0
	} else if c.b > 255 {
		c.b = 255
	}
}
