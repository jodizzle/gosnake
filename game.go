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
	snake 		[]*tl.Rectangle
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
		player.SetPosition(screenWidth/2, screenHeight/2)
		firstPass = false
	} /*else {
		player.r.SetPosition(50,50)
		firstPass = true
	}*/

	snakeTime += screen.TimeDelta()
	if snakeTime > 0.1 {
		snakeTime -= 0.1

		player.prevX, player.prevY = player.Position()
		switch player.direction {
		case "right":
			player.SetPosition(player.prevX+1, player.prevY)
			//player.SnakeMovement(player.prevX+1, player.prevY)
		case "left":
			player.SetPosition(player.prevX-1, player.prevY)
			//player.SnakeMovement(player.prevX-1, player.prevY)
		case "up":
			player.SetPosition(player.prevX, player.prevY-1)
			//player.SnakeMovement(player.prevX, player.prevY-1)
		case "down":
			player.SetPosition(player.prevX, player.prevY+1)
			//player.SnakeMovement(player.prevX, player.prevY+1)
		}

		player.SnakeMovement()
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

	//Check box boundaries
	playerX, playerY := player.Position()
	screenWidth, screenHeight := game.Screen().Size()

	//<= is used on the upper-boundaries to prevent the player from disappearing offscreen
	//by one square
	if playerX < 0 || playerX >= screenWidth {
		player.SetPosition(player.prevX, player.prevY)
	}
	if playerY < 0 || playerY >= screenHeight {
		player.SetPosition(player.prevX, player.prevY)
	}
}

// Old approach
// func (player *Player) SnakeMovement(x, y int) {
// 	// Appends a copy of the snake tail
// 	tailX, tailY := player.snake[len(player.snake)-1].Position()
// 	rect := tl.NewRectangle(tailX, tailY, 1, 1, tl.ColorGreen)
// 	game.Screen().Level().AddEntity(rect)
// 	player.snake = append(player.snake, rect)
// 	// Set new position of head
// 	player.SetPosition(x, y)
// 	// Deletes old tail
// 	// oldTail := player.snake[len(player.snake)-2]
// 	// game.Screen().Level().RemoveEntity(oldTail)
// 	// player.snake = append(player.snake[:len(player.snake)-2], player.snake[len(player.snake)-1:]...)
// 	// Updates positions

	// 	for i := 0; i < len(player.snake)-1; i++ {
// 		prevX, prevY := player.snake[i].Position()
// 		player.snake[i+1].SetPosition(prevX, prevY)
// 	}
// }


func (player *Player) SnakeMovement() {
	//Don't do anything if it's currently just the head
	if len(player.snake) > 1 {

		refSnake := make([]*tl.Rectangle, len(player.snake))
		//Doesn't work because it 'copy' copies the pointers
		// copy(refSnake, player.snake)
		//Need Deep Copy:
		for i := 0; i < len(player.snake); i++ {
			prevX, prevY := player.snake[i].Position()
			refSnake[i] = tl.NewRectangle(prevX, prevY, 1, 1, tl.ColorGreen)
		}

		for i := 1; i < len(player.snake)-1; i++ {
			prevX, prevY := refSnake[i].Position()
			player.snake[i+1].SetPosition(prevX, prevY)
		}

		player.snake[1].SetPosition(player.prevX, player.prevY)
	}
}

func (player *Player) Draw(screen *tl.Screen) {
	// screenWidth, screenHeight := screen.Size()
	// x, y := player.entity.Position()
	// player.level.SetOffset(screenWidth/2-x, screenHeight/2-y)
	
	player.Update(screen)

	for _,s := range player.snake {
		s.Draw(screen)
	}
}

// Remove rectangle and add to snake length
func (player *Player) Eat(rect *tl.Rectangle) {

	// Have to delete from level entities, not screen
	//game.Screen().Level().RemoveEntity(rect)
	game.Log("Entity Removed")

	px, py := player.snake[len(player.snake)-1].Position()
	rect.SetColor(tl.ColorWhite)
	switch player.direction {
	case "right":
		rect.SetPosition(px-1, py)
	case "left":
		rect.SetPosition(px+1, py)
	case "up":
		rect.SetPosition(px, py+1)
	case "down":
		rect.SetPosition(px, py-1)
	}

	player.snake = append(player.snake, rect)
}

//Order seems to be Tick then Draw, but only if there is an event to activate Tick
func (player *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		player.prevX, player.prevY = player.Position()
		//x, y := player.entity.Position()
		switch event.Key {
		case tl.KeyArrowRight:
			player.direction = "right"
		case tl.KeyArrowLeft:
			player.direction = "left"
		case tl.KeyArrowUp:
			player.direction = "up"
		case tl.KeyArrowDown:
			player.direction = "down"
		}
	}
}

func (player *Player) Size() (int, int) {
	return player.snake[0].Size()
}

func (player *Player) Position() (int, int) {
	return player.snake[0].Position()
}

func (player *Player) SetPosition(x, y int) {
	player.snake[0].SetPosition(x, y)
}

func (player *Player) Collide(collision tl.Physical) {
	//Check if it's a rectangle we're colliding with
	if rect, ok := collision.(*tl.Rectangle); ok {
		if (rect.Color() == tl.ColorGreen) {
			player.Eat(rect)
			game.Log("Rectangle Collision: %d", rect)
		}
	}
}

func main() {
	game = tl.NewGame()
	game.SetDebugOn(true)

	level := tl.NewBaseLevel(tl.Cell {
		Bg: tl.ColorBlack,
		Fg: tl.ColorBlack,
	})

	player := Player{
		snake:	[]*tl.Rectangle{tl.NewRectangle(50, 50, 1, 1, tl.ColorRed)},
		level:	level,
	}

	//player.entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorRed, Ch: '☺'})
	level.AddEntity(&player)

	game.Screen().SetLevel(level)
	firstPass = true

	game.Start()
}