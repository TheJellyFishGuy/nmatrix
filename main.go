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

var speed = flag.Int("speed", 35, "speed in ms (lower = faster)")

var charset = []rune("アイウエオカキクケコサシスセソABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type Drop struct {
	y      int
	length int
	speed  int
}

func main() {

	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	drops := make([]Drop, width)

	for i := range drops {

		drops[i] = Drop{
			y:      rand.Intn(height),
			length: rand.Intn(height/2) + 5,
			speed:  rand.Intn(3) + 1,
		}
	}

	hideCursor()
	clear()

	defer showCursor()
	defer reset()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	ticker := time.NewTicker(time.Duration(*speed) * time.Millisecond)

	for {

		select {

		case <-sig:
			return

		case <-ticker.C:

			draw(width, height, drops)

			update(height, drops)

		}
	}
}

func draw(width, height int, drops []Drop) {

	fmt.Print("\033[H")

	for x := 0; x < width; x++ {

		d := drops[x]

		for i := 0; i < d.length; i++ {

			y := d.y - i

			if y < 0 || y >= height {
				continue
			}

			move(x+1, y+1)

			char := string(charset[rand.Intn(len(charset))])

			if i == 0 {

				fmt.Print("\033[97m", char) // bright white head

			} else if i < 3 {

				fmt.Print("\033[92m", char)

			} else {

				fmt.Print("\033[32m", char)
			}
		}
	}

	reset()
}

func update(height int, drops []Drop) {

	for i := range drops {

		drops[i].y += drops[i].speed

		if drops[i].y-drops[i].length > height {

			drops[i].y = 0
			drops[i].length = rand.Intn(height/2) + 5
			drops[i].speed = rand.Intn(3) + 1
		}
	}
}

func move(x, y int) {

	fmt.Printf("\033[%d;%dH", y, x)
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
