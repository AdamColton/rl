package world

import (
	"github.com/adamcolton/vec2d"
	"sync"
)

var dirs = []vec2d.I{
	{-1, 0},
	{1, 0},
	{0, -1},
	{0, 1},
}

type Vec3D struct {
	Depth int
	vec2d.I
}

type Level struct {
	Depth int
	World *World
	Map   *Map
	sync.RWMutex
	CreaturesMap map[vec2d.I]*CreatureTile
}

func (l *Level) CreaturesAt(at vec2d.I) *CreatureTile {
	l.RLock()
	ct, ok := l.CreaturesMap[at]
	l.RUnlock()
	if !ok {
		ct = NewCreatureTile()
		l.Lock()
		l.CreaturesMap[at] = ct
		l.Unlock()
	}
	return ct
}

type World struct {
	Levels      []*Level
	PlayerStart vec2d.I
	sync.RWMutex
	Creatures map[CreatureID]Creature
	size      vec2d.I
}

func New(size vec2d.I, levels int) *World {
	w := &World{
		Levels:    make([]*Level, levels),
		size:      size,
		Creatures: make(map[CreatureID]Creature),
	}
	var stair vec2d.I
	stairDist := size.Distance(vec2d.I{0, 0}) / 2
	for i := range w.Levels {
		l := w.addLevel(i)
		if i == 0 {
			stair = l.Map.RandPassible()
			w.PlayerStart = stair
		} else {
			prevLvl := w.Levels[i-1]
			newStair := l.Map.RandPassibleDist(stairDist, stair, 50)
			for !prevLvl.Map.Get(newStair).Passible() {
				newStair = l.Map.RandPassibleDist(stairDist, stair, 50)
			}
			l.Map.Set(newStair, StairUp)
			prevLvl.Map.Set(newStair, StairDown)
			stair = newStair
		}
	}
	return w
}

func (w *World) addLevel(depth int) *Level {
	l := &Level{
		Depth:        depth,
		World:        w,
		Map:          NewMap(w.size, 0.45, 3),
		CreaturesMap: make(map[vec2d.I]*CreatureTile),
	}
	w.Levels[depth] = l
	return l
}

func (w *World) AddPlayer(control <-chan rune) *Player {
	p := NewPlayer(control)
	p.Loc = w.PlayerStart
	p.lvl = w.Levels[0]
	w.AddCreature(p)
	return p
}

func (w *World) AddCreature(c Creature) {
	w.Levels[c.Depth()].CreaturesAt(c.Location()).Add(c.ID())
	w.Lock()
	w.Creatures[c.ID()] = c
	w.Unlock()
}

func (w *World) GetCreature(cid CreatureID) Creature {
	w.RLock()
	c := w.Creatures[cid]
	w.RUnlock()
	return c
}

type changed struct {
	depth int
	tiles []vec2d.I
}

func (w *World) Step() []Vec3D {
	wg := &sync.WaitGroup{}
	wg.Add(len(w.Creatures))

	ch := make(chan []Vec3D)
	var all []Vec3D
	go func() {
		for m := range ch {
			all = append(all, m...)
			wg.Done()
		}
	}()

	for _, c := range w.Creatures {
		go func(c Creature) {
			ch <- c.Step()
		}(c)
	}
	wg.Wait()
	close(ch)
	return all
}

func (w *World) MoveCreature(c CreatureID, move []Vec3D) {
	w.Levels[move[0].Depth].CreaturesAt(move[0].I).Remove(c)
	w.Levels[move[1].Depth].CreaturesAt(move[1].I).Add(c)
}
