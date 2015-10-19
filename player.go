package main

import (
	tl "github.com/JoelOtter/termloop"
	"math/rand"
	"time"
)

// The player is a direction and a slice of rectangles that represent
// the snakebody, with snake[0] being the head. The direction and the
// length of the slice influence how the full snake is drawn on the screen.
type Player struct {
	snake     []*tl.Rectangle
	direction string
	prevX     int
	prevY     int
	level     *tl.BaseLevel
	snakeTime float64
	spawnTime float64
}

// Returns the size of the player (width, height).
func (player *Player) Size() (int, int) {
	return player.snake[0].Size()
}

// Returns the position of the player (x, y).
func (player *Player) Position() (int, int) {
	return player.snake[0].Position()
}

// Sets the position of the player to (x, y).
func (player *Player) SetPosition(x, y int) {
	player.snake[0].SetPosition(x, y)
}

// Handles player movement and food spawn rate.
func (player *Player) Update(screen *tl.Screen) {
	player.snakeTime += screen.TimeDelta()
	if player.snakeTime > snakeRate {
		player.snakeTime -= snakeRate

		player.prevX, player.prevY = player.Position()
		switch player.direction {
		case "right":
			player.SetPosition(player.prevX+1, player.prevY)
		case "left":
			player.SetPosition(player.prevX-1, player.prevY)
		case "up":
			player.SetPosition(player.prevX, player.prevY-1)
		case "down":
			player.SetPosition(player.prevX, player.prevY+1)
		}

		player.SnakeMovement()
	}

	player.spawnTime += screen.TimeDelta()
	if player.spawnTime > spawnRate {
		player.spawnTime -= spawnRate

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
	}
	if playerY < 0 || playerY >= screenHeight {
		GameOver()
	}
}

// Draws the player and the snake body.  Also calls player.Update.
func (player *Player) Draw(screen *tl.Screen) {
	player.Update(screen)

	for _, s := range player.snake {
		s.Draw(screen)
	}
}

// Handles arrow key inputs for movement.
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

// Handles collision with food and snake body.
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

// Controls snake-like movement of the snake body.
func (player *Player) SnakeMovement() {
	//Don't do anything if it's currently just the head
	if len(player.snake) > 1 {
		// Change color to white here because of difficulties putting it in player.Eat()
		player.snake[len(player.snake)-1].SetColor(tl.ColorWhite)

		// Need Deep Copy:
		refSnake := make([]*tl.Rectangle, len(player.snake))
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

// Remove rect and add to snake length.
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

// Returns true if rect is part of the snake body; false otherwise.
func (player *Player) InSnake(rect *tl.Rectangle) bool {
	for i, s := range player.snake {
		if rect == s {
			game.Log("Death by %d at %d", s, i)
			return true
		}
	}

	return false
}
