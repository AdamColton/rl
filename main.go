package main

import (
	"github.com/adamcolton/rl/draw"
	"github.com/adamcolton/rl/world"
	"github.com/dist-ribut-us/log"
	"github.com/gdamore/tcell"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Just in case
	log.Contents = log.Truncate
	log.ToFile("rl.log")
	log.SetTrace(1, 2, 2)

	var esc rune = 27
	draw.KeyMap[tcell.KeyEscape] = esc
	draw.KeyMap[tcell.KeyCtrlC] = esc

	d, err := draw.New()
	if err != nil {
		panic(err)
	}

	levels := 20
	w := world.New(d.Size(), levels)
	playerControls := make(chan rune)
	p := w.AddPlayer(playerControls)
	d.AddWorldAndPlayer(w, p)

	for i := 0; i < 100; i++ {
		lvl := w.Levels[rand.Intn(levels)]
		world.NewSnake(lvl)
	}

	go func() {
		d.DrawLevel()
		for {
			d.Update(w.Step())
		}
	}()

OUTER:
	for {
		switch r := d.Poll(); r {
		case esc:
			break OUTER
		default:
			playerControls <- r
		}
	}

	d.Fini()
}
