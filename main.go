package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"sync"
)

var width = 2048
var height = 2048

var xFrom, xTo, yFrom, yTo float64 = -2, 1, -1.5, 1.5

var maxIter uint8 = 30

func main() {
	fmt.Println("Starting")

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	var w sync.WaitGroup
	w.Add(width)

	for x := 0; x < width; x++ {
		go func(x int) {
			for y := 0; y < height; y++ {
				col := color.RGBA{}
				mandelX := (float64(x)/float64(width))*(xTo-xFrom) + xFrom
				mandelY := (float64(y)/float64(height))*(yTo-yFrom) + yFrom
				iterationen := Mandelbrot(mandelX, mandelY, mandelX, mandelY, int(maxIter))

				col.R = 255 / maxIter * iterationen
				col.G = 255 / maxIter * iterationen
				col.B = 255 / maxIter * iterationen

				col.A = 255 / maxIter * iterationen

				img.Set(x, y, col)
			}
			w.Done()
		}(x)

	}

	w.Wait()

	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(f, img)

	fmt.Println("Finished")
}

//Mandelbrot function from Wikipedia
func Mandelbrot(x, y, xAdd, yAdd float64, maxIter int) uint8 {
	maxBetrag2 := 4.0
	remainIter := maxIter
	xx := x * x
	yy := y * y
	xy := x * y
	betrag2 := xx + yy

	for betrag2 <= maxBetrag2 && remainIter > 0 {
		remainIter = remainIter - 1
		x = xx - yy + xAdd
		y = xy + xy + yAdd
		xx = x * x
		yy = y * y
		xy = x * y
		betrag2 = xx + yy
	}
	return uint8(maxIter - remainIter)
}
