package main

import tl "github.com/JoelOtter/termloop"

var game *tl.Game
var firstPass bool

type Player struct {
	r *tl.Rectangle
	prevX int
	prevY int
	level *tl.BaseLevel
}

func (player *Player) Draw(screen *tl.Screen) {
	// screenWidth, screenHeight := screen.Size()
	// x, y := player.entity.Position()
	// player.level.SetOffset(screenWidth/2-x, screenHeight/2-y)

	//tl.Screen.size() parameters are evidently zero until the game.Start(),
	//So this is a crude solution intended to center the player after the game has begun
	if firstPass {
		screenWidth, screenHeight := screen.Size()
		player.r.SetPosition(screenWidth/2, screenHeight/2)
		firstPass = false
	}

	player.r.Draw(screen)
}

//Order seems to be Tick then Draw, but only if there is an event to activate Tick
func (player *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		player.prevX, player.prevY = player.r.Position()
		//x, y := player.entity.Position()
		switch event.Key {
		case tl.KeyArrowRight:
			player.r.SetPosition(player.prevX+1, player.prevY)
		case tl.KeyArrowLeft:
			player.r.SetPosition(player.prevX-1, player.prevY)
		case tl.KeyArrowUp:
			player.r.SetPosition(player.prevX, player.prevY-1)
		case tl.KeyArrowDown:
			player.r.SetPosition(player.prevX, player.prevY+1)
		}
	}

	//Check box boundaries
	playerX, playerY := player.r.Position()
	screenWidth, screenHeight := game.Screen().Size()

	//<= is used on the upper-boundaries to prevent the player from disappearing offscreen
	//by one square
	if playerX < 0 || playerX >= screenWidth {
		player.r.SetPosition(player.prevX, player.prevY)
	}
	if playerY < 0 || playerY >= screenHeight {
		player.r.SetPosition(player.prevX, player.prevY)
	}
}

func (player *Player) Size() (int, int) {
	return player.r.Size()
}

func (player *Player) Position() (int, int) {
	return player.r.Position()
}

func (player *Player) Collide(collision tl.Physical) {
	//Check if it's a rectangle we're colliding with
	if _, ok := collision.(*tl.Rectangle); ok {
		player.r.SetPosition(player.prevX, player.prevY)
	}
}

func main() {
	game = tl.NewGame()
	game.SetDebugOn(true)

	level := tl.NewBaseLevel(tl.Cell {
		Bg: tl.ColorBlack,
		Fg: tl.ColorBlack,
	})

	level.AddEntity(tl.NewRectangle(10, 10, 50, 20, tl.ColorBlue))

	player := Player{
		r:		tl.NewRectangle(50, 50, 1, 1, tl.ColorRed),
		level:	level,
	}

	//player.entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorRed, Ch: '☺'})
	level.AddEntity(&player)

	game.Screen().SetLevel(level)
	firstPass = true

	game.Start()
}