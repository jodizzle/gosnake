package main

import (
	tl "github.com/JoelOtter/termloop"
	"math/rand"
	"time"
)

type Player struct {
	snake 		[]*tl.Rectangle
	direction	string
	prevX 		int
	prevY 		int
	level 		*tl.BaseLevel
}

//Handles auto events
func (player *Player) Update(screen *tl.Screen) {
	//tl.Screen.size() parameters are evidently zero until the game.Start(),
	//So this is a crude solution intended to center the player after the game has begun
	if firstPass {
		screenWidth, screenHeight := screen.Size()
		player.SetPosition(screenWidth/2, screenHeight/2)
		firstPass = false
	}

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
	//(Funnily enough, when player.snake is more than one unit long, just stopping the player at 
	//the boundaries also causes a game over state because the tail slides into the head)
	if playerX < 0 || playerX >= screenWidth {
		GameOver()
		//player.SetPosition(player.prevX, player.prevY)
	}
	if playerY < 0 || playerY >= screenHeight {
		GameOver()
		//player.SetPosition(player.prevX, player.prevY)
	}
}

func (player *Player) SnakeMovement() {
	//Don't do anything if it's currently just the head
	if len(player.snake) > 1 {
		//Change color to white up here because of difficulties putting it in player.Eat()
		player.snake[len(player.snake)-1].SetColor(tl.ColorWhite)

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
	player.Update(screen)

	for _,s := range player.snake {
		s.Draw(screen)
	}
}

// Remove rectangle and add to snake length
func (player *Player) Eat(rect *tl.Rectangle) {

	//Removing entities (at least the way that I approached it) doesn't seem to have much of an affect.
	//Have to delete from level entities, not screen
	// game.Screen().Level().RemoveEntity(rect)
	// game.Screen().RemoveEntity(rect)
	// game.Log("Entity Removed")

	px, py := player.snake[len(player.snake)-1].Position()
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
	//Trying to change Color here doesn't seem to work, possibly because of the order that collisions are
	//handled in.  Color is instead changed in player.Update()
	// rect.SetColor(tl.ColorWhite)
	// game.Screen().Level().AddEntity(rect)
}

//Order seems to be Tick then Draw, but only if there is an event to activate Tick
func (player *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		player.prevX, player.prevY = player.Position()
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

func (player *Player) InSnake(rect *tl.Rectangle) bool {
	for i,s := range player.snake {
		if rect == s {
			game.Log("Death by %d at %d", s, i)
			return true
		}
	}

	return false
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
		if rect.Color() == tl.ColorGreen {
			player.Eat(rect)
			game.Log("Rectangle Collision: %d", rect)
		}
		if rect.Color() == tl.ColorWhite {
			if player.InSnake(rect) {
				GameOver()
			}
		}
	}
}