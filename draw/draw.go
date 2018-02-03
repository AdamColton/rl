package draw

import (
	"github.com/adamcolton/rl/world"
	"github.com/adamcolton/vec2d"
	"github.com/gdamore/tcell"
)

type tile struct {
	symbol rune
	style  tcell.Style
}

var dftl = tcell.StyleDefault.
	Foreground(tcell.ColorWhite).
	Background(tcell.ColorBlack)
var red = tcell.StyleDefault.Foreground(tcell.NewHexColor(0xff0000))
var green = tcell.StyleDefault.Foreground(tcell.NewHexColor(0x00ff00))
var blue = tcell.StyleDefault.Foreground(tcell.NewHexColor(0xccccff))
var darkGrey = tcell.StyleDefault.Foreground(tcell.NewHexColor(0x333333))
var lightGrey = tcell.StyleDefault.Foreground(tcell.NewHexColor(0x999999))

var terrainTiles = map[world.Terrain]tile{
	world.Wall:      {'#', lightGrey},
	world.Open:      {'.', darkGrey},
	world.StairUp:   {'<', blue},
	world.StairDown: {'>', blue},
}

var creatureTiles = map[world.CreatureType]tile{
	world.PlayerType: {'@', green},
	world.SnakeType:  {'s', red},
}

type Draw struct {
	screen tcell.Screen
	world  *world.World
	player *world.Player
	depth  int
}

func New() (*Draw, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	err = screen.Init()
	if err != nil {
		return nil, err
	}

	screen.SetStyle(dftl)
	screen.Clear()

	return &Draw{
		screen: screen,
	}, nil
}

func (d *Draw) AddWorldAndPlayer(w *world.World, p *world.Player) {
	d.world = w
	d.player = p
}

func (d *Draw) Size() vec2d.I {
	sx, sy := d.screen.Size()
	return vec2d.I{sx, sy}
}

func (d *Draw) Fini() {
	d.screen.Fini()
}

func (d *Draw) DrawLevel() {
	d.depth = d.player.Depth()
	lvl := d.world.Levels[d.depth]
	var t world.Terrain
	iter := lvl.Map.IterAll(&t)
	for iter.Next() {
		if r, ok := terrainTiles[t]; ok {
			pt := iter.Pt()
			d.screen.SetContent(pt.X, pt.Y, r.symbol, nil, r.style)
		}
	}

	lvl.RLock()
	for pt, cs := range lvl.CreaturesMap {
		if cID := cs.First(); cID > 0 {
			c := lvl.World.GetCreature(cID)
			r := creatureTiles[c.Type()]
			d.screen.SetContent(pt.X, pt.Y, r.symbol, nil, r.style)
		}
	}
	lvl.RUnlock()

	d.screen.Show()
}

var KeyMap = map[tcell.Key]rune{}

func (d *Draw) Poll() rune {
	for {
		switch ev := d.screen.PollEvent().(type) {
		case *tcell.EventResize:
			d.screen.Sync()
		case *tcell.EventKey:
			k := ev.Key()
			if r, ok := KeyMap[k]; ok {
				return r
			}
			return ev.Rune()
		}
	}
}

func (d *Draw) Update(tiles []world.Vec3D) {
	if d.depth != d.player.Depth() {
		d.DrawLevel()
		return
	}
	lvl := d.world.Levels[d.depth]
	for _, pt3d := range tiles {
		if pt3d.Depth != d.depth {
			continue
		}

		pt := pt3d.I
		if cid := lvl.CreaturesAt(pt).First(); cid != 0 {
			c := d.world.GetCreature(cid)
			r := creatureTiles[c.Type()]
			d.screen.SetContent(pt.X, pt.Y, r.symbol, nil, r.style)
			continue
		}
		r := terrainTiles[lvl.Map.Get(pt)]
		d.screen.SetContent(pt.X, pt.Y, r.symbol, nil, r.style)
	}
	d.screen.Show()
}
