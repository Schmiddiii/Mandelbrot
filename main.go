package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"sync"
)

var xMiddle, xDelta, yMiddle, yDelta float64 //= -2, 3, -1.5, 3
var width, height int
var fileName string

func main() {

	flag.Float64Var(&xMiddle, "x", -.5, "Starting x-value")
	flag.Float64Var(&yMiddle, "y", 0, "Starting y-value")

	flag.Float64Var(&xDelta, "dx", 3, "Delta x-value")
	flag.Float64Var(&yDelta, "dy", 3, "Delta y-value")

	flag.IntVar(&width, "w", 2048, "Width of Image")
	flag.IntVar(&height, "h", 2048, "Delta y-value")

	var maxIter64 uint64
	flag.Uint64Var(&maxIter64, "i", 30, "Maximal Iterations")
	maxIter := uint8(maxIter64)

	flag.StringVar(&fileName, "file", "image", "Name of the outputfile")

	flag.Parse()

	fmt.Printf("Starting with x: %f, y: %f, dx: %f, dy: %f, i: %d\n", xMiddle, xDelta, yMiddle, yDelta, maxIter)

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	var w sync.WaitGroup
	w.Add(width)

	for x := 0; x < width; x++ {
		go func(x int) {
			for y := 0; y < height; y++ {
				col := color.RGBA{}
				mandelX := (float64(x)/float64(width))*xDelta + (xMiddle - .5*xDelta)
				mandelY := (float64(y)/float64(height))*yDelta + (yMiddle - .5*yDelta)
				iterationen := mandelbrot(mandelX, mandelY, mandelX, mandelY, int(maxIter))

				if iterationen == maxIter {

				} else {
					col = hsv2rgb(1080*float64(iterationen)/float64(maxIter), 1, 1, 255)
				}

				col.A = 255 / maxIter * iterationen

				img.Set(x, y, col)
			}
			w.Done()
		}(x)

	}

	w.Wait()

	f, err := os.Create(fileName + ".png")
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(f, img)

	fmt.Println("Finished")
}

//Mandelbrot function from Wikipedia
func mandelbrot(x, y, xAdd, yAdd float64, maxIter int) uint8 {
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

func hsv2rgb(h, s, v, a float64) color.RGBA {
	h = modFloat(h, 360)
	c := v * s
	x := c * (1 - math.Abs(float64(modFloat(h/60, 2.0)-1)))
	m := v - c
	if h < 60 {
		return color.RGBA{R: uint8((c + m) * 255), G: uint8((x + m) * 255), B: uint8(((0 + m) * 255)), A: uint8(a)}
	} else if h < 120 {
		return color.RGBA{R: uint8((x + m) * 255), G: uint8((c + m) * 255), B: uint8(((0 + m) * 255)), A: uint8(a)}
	} else if h < 180 {
		return color.RGBA{R: uint8((0 + m) * 255), G: uint8((c + m) * 255), B: uint8(((x + m) * 255)), A: uint8(a)}
	} else if h < 240 {
		return color.RGBA{R: uint8((0 + m) * 255), G: uint8((x + m) * 255), B: uint8(((c + m) * 255)), A: uint8(a)}
	} else if h < 300 {
		return color.RGBA{R: uint8((x + m) * 255), G: uint8((0 + m) * 255), B: uint8(((c + m) * 255)), A: uint8(a)}
	} else {
		return color.RGBA{R: uint8((x + m) * 255), G: uint8((0 + m) * 255), B: uint8(((c + m) * 255)), A: uint8(a)}
	}

}

func modFloat(m, f float64) float64 {
	num := int(m / f)
	return m - float64(num)*f
}
