package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"golang.org/x/term"
)

var speed = flag.Int("speed", 25, "frame speed ms")
var density = flag.Float64("density", 0.7, "rain density 0.1-1.0")

var chars = []rune("アイウエオカキクケコサシスセソワヲンABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type Drop struct {
	head   float64
	length int
	vel    float64
	active bool
}

var buffer [][]int
var width, height int
var drops []Drop

func main() {

	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	initScreen()

	hideCursor()
	defer showCursor()
	defer reset()
	defer clear()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	ticker := time.NewTicker(time.Duration(*speed) * time.Millisecond)

	for {

		select {

		case <-sig:
			return

		case <-ticker.C:

			checkResize()
			update()
			render()
		}
	}
}

func initScreen() {

	width, height, _ = term.GetSize(int(os.Stdout.Fd()))

	buffer = make([][]int, height)

	for y := range buffer {
		buffer[y] = make([]int, width)
	}

	drops = make([]Drop, width)

	for i := range drops {

		if rand.Float64() < *density {

			drops[i] = newDrop()
		}
	}
}

func newDrop() Drop {

	return Drop{
		head:   rand.Float64() * float64(height),
		length: rand.Intn(height/2) + 6,
		vel:    rand.Float64()*0.5 + 0.5,
		active: true,
	}
}

func update() {

	// fade buffer

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			if buffer[y][x] > 0 {
				buffer[y][x]--
			}
		}
	}

	// update drops

	for x := range drops {

		d := &drops[x]

		if !d.active {

			if rand.Float64() < 0.002 {

				drops[x] = newDrop()
			}

			continue
		}

		d.head += d.vel

		for i := 0; i < d.length; i++ {

			y := int(d.head) - i

			if y >= 0 && y < height {

				buffer[y][x] = 255 - (i * 20)

				if buffer[y][x] < 20 {
					buffer[y][x] = 20
				}
			}
		}

		if int(d.head)-d.length > height {

			d.active = false
		}
	}
}

func render() {

	fmt.Print("\033[H")

	for y := 0; y < height; y++ {

		for x := 0; x < width; x++ {

			v := buffer[y][x]

			if v == 0 {

				fmt.Print(" ")

				continue
			}

			color(v)

			fmt.Print(string(chars[rand.Intn(len(chars))]))
		}

		fmt.Print("\n")
	}
}

func color(v int) {

	switch {

	case v > 220:
		fmt.Print("\033[97m")

	case v > 180:
		fmt.Print("\033[92m")

	case v > 120:
		fmt.Print("\033[32m")

	case v > 60:
		fmt.Print("\033[32;2m")

	default:
		fmt.Print("\033[30;1m")
	}
}

func checkResize() {

	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	if w != width || h != height {

		clear()
		initScreen()
	}
}

func clear() {
	fmt.Print("\033[2J")
}

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}

func reset() {
	fmt.Print("\033[0m")
}
