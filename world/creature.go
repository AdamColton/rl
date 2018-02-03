package world

import (
	"github.com/adamcolton/vec2d"
	"sync"
)

type CreatureType uint16

const (
	_ CreatureType = iota
	PlayerType
	SnakeType
)

type BaseCreature struct {
	CID CreatureID
	Loc vec2d.I
	lvl *Level
}

func NewCreature() *BaseCreature {
	return &BaseCreature{
		CID: <-getCreatureID,
	}
}

func (bc *BaseCreature) ID() CreatureID { return bc.CID }
func (bc *BaseCreature) Depth() int     { return bc.lvl.Depth }

func (bc *BaseCreature) Location() vec2d.I { return bc.Loc }
func (bc *BaseCreature) Vec3D() Vec3D {
	return Vec3D{
		Depth: bc.lvl.Depth,
		I:     bc.Loc,
	}
}

func (bc *BaseCreature) Move(d vec2d.I) []Vec3D {
	start := bc.Vec3D()
	bc.Loc = bc.Loc.Add(d)

	tile := bc.lvl.Map.Get(bc.Loc)
	if tile == StairDown {
		d := bc.Depth() + 1
		bc.lvl = bc.lvl.World.Levels[d]
	}
	if tile == StairUp {
		d := bc.Depth() - 1
		bc.lvl = bc.lvl.World.Levels[d]
	}

	end := bc.Vec3D()
	move := []Vec3D{start, end}
	bc.lvl.World.MoveCreature(bc.CID, move)
	return move
}

func (bc *BaseCreature) checkMove(d vec2d.I) bool {
	return bc.lvl.Map.Get(bc.Loc.Add(d)).Passible()
}

type CreatureID uint32

// generator to guarentee unique ids
var getCreatureID = func() <-chan CreatureID {
	ch := make(chan CreatureID)
	var cid CreatureID
	go func() {
		for {
			cid++
			ch <- cid
		}
	}()
	return ch
}()

type Creature interface {
	Location() vec2d.I
	Type() CreatureType
	ID() CreatureID
	Step() []Vec3D
	Depth() int
}

type CreatureTile struct {
	Creatures map[CreatureID]bool
	sync.RWMutex
}

func NewCreatureTile() *CreatureTile {
	return &CreatureTile{
		Creatures: make(map[CreatureID]bool),
	}
}

func (cs *CreatureTile) Add(id CreatureID) {
	cs.Lock()
	cs.Creatures[id] = true
	cs.Unlock()
}

func (cs *CreatureTile) Remove(id CreatureID) {
	cs.Lock()
	delete(cs.Creatures, id)
	cs.Unlock()
}

func (cs *CreatureTile) First() CreatureID {
	cs.RLock()
	defer cs.RUnlock()
	for id := range cs.Creatures {
		return id
	}
	return 0
}
