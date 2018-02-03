package world

import (
	"github.com/adamcolton/vec2d"
	"unicode"
)

type Player struct {
	control <-chan rune
	*BaseCreature
}

func (p *Player) Type() CreatureType { return PlayerType }

func NewPlayer(control <-chan rune) *Player {
	return &Player{
		control:      control,
		BaseCreature: NewCreature(),
	}
}

var moveKeys = map[rune]vec2d.I{
	'w': {0, -1},
	'a': {-1, 0},
	's': {0, 1},
	'd': {1, 0},
}

func (p *Player) Step() []Vec3D {
	for {
		r := unicode.ToLower(<-p.control)
		if d, ok := moveKeys[r]; ok {
			if !p.checkMove(d) {
				continue
			}
			return p.Move(d)
		}
		switch r {
		case ' ':
			return nil // wait a turn
		}
	}
}
