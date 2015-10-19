package main

const (
	startMessage = "SNAKE"
	instructions = "ARROW KEYS TO MOVE"
	instructions2 = "PRESS ENTER TO START"
	endMessage = "GAME OVER"
	endInstructions = "PRESS ENTER TO RESTART"
	snakeRate = 0.1 // Speed of snake movement (1 cell/0.1 second)
	spawnRate = 1 // Rate of spawning food (1 spawn/1 second)
)

func main() {
	NewGame()
}