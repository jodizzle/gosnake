package main

import (
	tl "github.com/JoelOtter/termloop"
)

var game *tl.Game
var firstPass bool
var snakeTime float64
var spawnTime float64

func GameOver() {
	end := tl.NewBaseLevel(tl.Cell {
				Bg: tl.ColorRed,
				Fg: tl.ColorBlack,
			})

	endText := LevelText {
		message: tl.NewText(0, 0, endMessage, tl.ColorGreen, tl.ColorBlack),
		instructions: tl.NewText(0, 0, endInstructions, tl.ColorGreen, tl.ColorBlack),
		instructions2: tl.NewText(0, 0, "", tl.ColorGreen, tl.ColorBlack),
	}

	end.AddEntity(&endText)

	firstPass = true
	game.Screen().SetLevel(end)
}

func GameStart() {
	start := tl.NewBaseLevel(
		tl.Cell{Bg: tl.ColorBlack, Fg: tl.ColorBlack, Ch: 'S'},
	)

	startText := LevelText{
		tl.NewText(0, 0, startMessage, tl.ColorGreen, tl.ColorBlack),
		tl.NewText(0, 0, instructions, tl.ColorGreen, tl.ColorBlack),
		tl.NewText(0, 0, instructions2, tl.ColorGreen, tl.ColorBlack),
	}
	
	start.AddEntity(&startText)

	firstPass = true
	game.Screen().SetLevel(start)
	game.Start()
}

func GamePlay() {
	level := tl.NewBaseLevel(tl.Cell {
				Bg: tl.ColorBlack,
				Fg: tl.ColorBlack,
			})

	player := Player{
		snake:	[]*tl.Rectangle{tl.NewRectangle(0, 0, 1, 1, tl.ColorRed)},
		level:	level,
	}

	level.AddEntity(&player)
	game.Screen().SetLevel(level)
}

type LevelText struct {
	message 		*tl.Text
	instructions	*tl.Text
	instructions2	*tl.Text
}

func (text *LevelText) Draw(screen *tl.Screen) {
	screenWidth, screenHeight := screen.Size()
	text.message.SetPosition(screenWidth/2, screenHeight/2)
	text.message.Draw(screen)
	text.instructions.SetPosition(screenWidth/2 - 5, screenHeight/2 + 1)
	text.instructions.Draw(screen)
	text.instructions2.SetPosition(screenWidth/2 - 5, screenHeight/2 + 2)
	text.instructions2.Draw(screen)
}

func (text *LevelText) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		if event.Key == tl.KeyEnter {
			GamePlay()
		}
	}
}

func NewGame() {
	game = tl.NewGame()
	game.SetDebugOn(false)
	GameStart()
}