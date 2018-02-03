package world

import (
	"github.com/adamcolton/vec2d"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayer(t *testing.T) {
	m := NewMap(vec2d.I{20, 20}, 0, 0)

	ch := make(chan rune, 2)
	p := NewPlayer(ch)
	p.Loc = vec2d.I{10, 2}
	p.lvl = &Level{
		CreaturesMap: make(map[vec2d.I]*CreatureTile),
		Map:          m,
		World: &World{
			Levels: []*Level{
				&Level{
					Map:          m,
					CreaturesMap: make(map[vec2d.I]*CreatureTile),
				},
			},
		},
	}
	assert.NotNil(t, p)

	// Normal move works correctly
	ch <- 'w'
	expected := []Vec3D{
		{0, vec2d.I{10, 2}},
		{0, vec2d.I{10, 1}},
	}
	assert.Equal(t, expected, p.Step())
	assert.Equal(t, vec2d.I{10, 1}, p.Location())

	// walk into wall
	// the 'w' (up) should be consumed and the ' ' is actually used
	ch <- 'w'
	ch <- ' '
	assert.Nil(t, p.Step())
	assert.Len(t, ch, 0)
}
