package world

import (
	"math/rand"
)

type Snake struct {
	*BaseCreature
}

func (s *Snake) Type() CreatureType { return SnakeType }

func NewSnake(lvl *Level) *Snake {
	s := &Snake{
		BaseCreature: NewCreature(),
	}
	s.lvl = lvl
	s.Loc = lvl.Map.RandPassible()
	lvl.World.AddCreature(s)
	return s
}

func (s *Snake) Step() []Vec3D {
	d := dirs[rand.Intn(4)]
	for !s.checkMove(d) {
		d = dirs[rand.Intn(4)]
	}
	return s.Move(d)
}
