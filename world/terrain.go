package world

import (
	"github.com/adamcolton/vec2d"
	"github.com/adamcolton/vec2d/grid"
	"math/rand"
)

type Terrain int

// Impassible tiles
const (
	Wall Terrain = iota
)

// Passble tiles
const (
	Passible Terrain = iota + 2000 // can add more Impassible tiles
	Open
	StairUp
	StairDown
)

func (t Terrain) Passible() bool {
	return t > Passible
}

type Map struct {
	*grid.Grid
	OpenTiles []vec2d.I
}

func (m *Map) Get(pt vec2d.I) Terrain {
	return m.Grid.Get(pt).(Terrain)
}

func (m *Map) RandPassible() vec2d.I {
	rIdx := rand.Intn(len(m.OpenTiles))
	return m.OpenTiles[rIdx]
}

// RandOpenDist returns a random passible point at least dist away. After a given number
// of tries it will return the best point if it hasn't found a perfect match
func (m *Map) RandPassibleDist(dist float64, from vec2d.I, tries int) vec2d.I {
	var bestPt vec2d.I
	var bestDist float64
	for i := 0; i < tries && bestDist < dist; i++ {
		pt := m.RandPassible()
		ptDist := from.Distance(pt)
		if ptDist > bestDist {
			bestPt, bestDist = pt, ptDist
		}
	}
	return bestPt
}

func NewMap(size vec2d.I, density float64, smoothing int) *Map {
	g := grid.New(size, func(pt vec2d.I) interface{} {
		isBorder := pt.X == 0 || pt.Y == 0 || pt.X == size.X-1 || pt.Y == size.Y-1
		if isBorder || rand.Float64() < density {
			return int(Wall)
		}
		return int(Open)
	})

	conv := func(i interface{}) int { return i.(int) }
	smooth := grid.PluralProcessor(conv)
	for i := 0; i < smoothing; i++ {
		g = g.Process(smooth)
	}
	g = g.Process(func(pt vec2d.I, g *grid.Grid) interface{} {
		return Terrain(g.Get(pt).(int))
	})

	m := &Map{
		Grid: g,
	}
	m.setOpenTiles()

	return m
}

// find the largest continuous open area
func (m *Map) setOpenTiles() {
	seen := make(map[vec2d.I]bool)
	var t Terrain
	iter := m.IterAll(&t)
	for iter.Next() {
		if !t.Passible() || seen[iter.Pt()] {
			continue
		}
		area := m.Flood(iter.Pt(), dirs, isPassible)
		if len(area) > len(m.OpenTiles) {
			m.OpenTiles = area
		}
		for _, pt := range area {
			seen[pt] = true
		}
	}
}

func isPassible(pt vec2d.I, g *grid.Grid) bool {
	return g.Get(pt).(Terrain).Passible()
}
