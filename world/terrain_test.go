package world

import (
	"github.com/adamcolton/vec2d"
	"github.com/adamcolton/vec2d/grid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tiles = map[Terrain]string{
	Wall: "#",
	Open: ".",
}

var format = grid.Formatter{
	Stringer: func(t interface{}) string {
		if s, ok := tiles[t.(Terrain)]; ok {
			return s
		}
		return "?"
	},
}

func TestMap(t *testing.T) {
	m := NewMap(vec2d.I{100, 100}, .45, 3)
	assert.NotNil(t, m)
	t.Log(format.Format(m.Grid))
}

func TestCheckOpenness(t *testing.T) {
	gen := func(pt vec2d.I) interface{} {
		if pt.In(vec2d.I{1, 1}, vec2d.I{4, 4}) && (vec2d.I{2, 2} != pt) {
			return Wall
		}
		return Open
	}
	m := &Map{
		Grid: grid.New(vec2d.I{5, 5}, gen),
	}
	m.setOpenTiles()
	assert.Len(t, m.OpenTiles, 16)
}
