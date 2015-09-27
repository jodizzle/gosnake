package main

import (
	tl "github.com/JoelOtter/termloop"
	"math/rand"
	"time"
)

var game *tl.Game
var firstPass bool
var snakeTime float64
var spawnTime float64

type Player struct {
	r 			*tl.Rectangle
	direction	string
	prevX 		int
	prevY 		int
	level 		*tl.BaseLevel
}

//Generates a random number in a given range
// func Random(min, max, int) int {
// 	rando := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	return rando.Int
// }

//Handles auto events
func (player *Player) Update(screen *tl.Screen) {
	//tl.Screen.size() parameters are evidently zero until the game.Start(),
	//So this is a crude solution intended to center the player after the game has begun
	if firstPass {
		screenWidth, screenHeight := screen.Size()
		player.r.SetPosition(screenWidth/2, screenHeight/2)
		firstPass = false
	} /*else {
		player.r.SetPosition(50,50)
		firstPass = true
	}*/

	snakeTime += screen.TimeDelta()
	if snakeTime > 0.1 {
		snakeTime -= 0.1

		player.prevX, player.prevY = player.r.Position()
		switch player.direction {
		case "right":
			player.r.SetPosition(player.prevX+2, player.prevY)
		case "left":
			player.r.SetPosition(player.prevX-2, player.prevY)
		case "up":
			player.r.SetPosition(player.prevX, player.prevY-1)
		case "down":
			player.r.SetPosition(player.prevX, player.prevY+1)
		}
	}

	spawnTime += screen.TimeDelta()
	if spawnTime > 1 {
		spawnTime -= 1

		screenWidth, screenHeight := screen.Size()
		rando := rand.New(rand.NewSource(time.Now().UnixNano()))
		spawnX, spawnY := rando.Intn(screenWidth), rando.Intn(screenHeight)
		screen.Level().AddEntity(tl.NewRectangle(spawnX, spawnY, 1, 1, tl.ColorGreen))
		
		game.Log("Spawn at (%d,%d)", spawnX, spawnY)
	}
}

func (player *Player) Draw(screen *tl.Screen) {
	// screenWidth, screenHeight := screen.Size()
	// x, y := player.entity.Position()
	// player.level.SetOffset(screenWidth/2-x, screenHeight/2-y)
	
	player.Update(screen)

	player.r.Draw(screen)
}

// Remove rectangle add to snake length
func (player *Player) Eat(rect *tl.Rectangle) {
	game.Screen().Level().RemoveEntity(rect)
	game.Log("Entity Removed")
}

//Order seems to be Tick then Draw, but only if there is an event to activate Tick
func (player *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		player.prevX, player.prevY = player.r.Position()
		//x, y := player.entity.Position()
		switch event.Key {
		case tl.KeyArrowRight:
			player.direction = "right"
			//player.r.SetPosition(player.prevX+1, player.prevY)
		case tl.KeyArrowLeft:
			player.direction = "left"
			//player.r.SetPosition(player.prevX-1, player.prevY)
		case tl.KeyArrowUp:
			player.direction = "up"
			//player.r.SetPosition(player.prevX, player.prevY-1)
		case tl.KeyArrowDown:
			player.direction = "down"
			//player.r.SetPosition(player.prevX, player.prevY+1)
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
	if rect, ok := collision.(*tl.Rectangle); ok {
		//player.r.SetPosition(player.prevX, player.prevY)
		player.Eat(rect)
		game.Log("Rectangle Collision: %d", rect)
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