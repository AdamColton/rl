package world

import (
	"github.com/adamcolton/vec2d"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestWorld(t *testing.T) {
	levels := 20

	w := New(vec2d.I{100, 100}, levels)
	ch := make(chan rune, 4)
	w.AddPlayer(ch)

	for i := 0; i < 100; i++ {
		d := rand.Intn(levels)
		NewSnake(w.Levels[d])
	}

	assert.NotNil(t, w)
	populateMoveChan(ch)
	allTiles := w.Step()
	assert.Len(t, allTiles, 202)
}

func populateMoveChan(ch chan rune) {
	for len(ch) > 0 {
		<-ch
	}
	ch <- 'w'
	ch <- 'a'
	ch <- 's'
	ch <- 'd'
}
