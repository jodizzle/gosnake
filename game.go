package main

import (
	tl "github.com/JoelOtter/termloop"
)

// The game pointer is kept global for ease of access.
var game *tl.Game

// Makes a new termloop game and calls GameStart.
func NewGame() {
	game = tl.NewGame()
	game.SetDebugOn(false)
	GameStart()
}

// Sets the game level to the title screen.
func GameStart() {
	start := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
		Fg: tl.ColorBlack,
	})

	startText := LevelText{
		tl.NewText(0, 0, startMessage, tl.ColorGreen, tl.ColorBlack),
		tl.NewText(0, 0, instructions, tl.ColorGreen, tl.ColorBlack),
		tl.NewText(0, 0, instructions2, tl.ColorGreen, tl.ColorBlack),
	}

	start.AddEntity(&startText)

	game.Screen().SetLevel(start)
	game.Start()
}

// Sets the game level to the regular play screen and adds
// the player entity.
func GamePlay() {
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
		Fg: tl.ColorBlack,
	})

	player := Player{
		snake: []*tl.Rectangle{tl.NewRectangle(0, 0, 1, 1, tl.ColorRed)},
		level: level,
	}
	screenWidth, screenHeight := game.Screen().Size()
	player.SetPosition(screenWidth/2, screenHeight/2)

	level.AddEntity(&player)
	game.Screen().SetLevel(level)
}

// Sets the game level to the lose screen.
func GameOver() {
	end := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorRed,
		Fg: tl.ColorBlack,
	})

	endText := LevelText{
		tl.NewText(0, 0, endMessage, tl.ColorGreen, tl.ColorBlack),
		tl.NewText(0, 0, endInstructions, tl.ColorGreen, tl.ColorBlack),
		tl.NewText(0, 0, "", tl.ColorGreen, tl.ColorBlack),
	}

	end.AddEntity(&endText)

	game.Screen().SetLevel(end)
}

// Level text for the title and lose screens is currently made up
// of three termloop Text pointers.  It may be valuable to find a
// better way to handle Text messages.
type LevelText struct {
	message       *tl.Text
	instructions  *tl.Text
	instructions2 *tl.Text
}

// Draws the Level text starting at the middle of the screen down.
func (text *LevelText) Draw(screen *tl.Screen) {
	screenWidth, screenHeight := screen.Size()
	text.message.SetPosition(screenWidth/2, screenHeight/2)
	text.message.Draw(screen)
	text.instructions.SetPosition(screenWidth/2-5, screenHeight/2+1)
	text.instructions.Draw(screen)
	text.instructions2.SetPosition(screenWidth/2-5, screenHeight/2+2)
	text.instructions2.Draw(screen)
}

// Switches from the title or lose screen to the play screen when the
// user presses enter.
func (text *LevelText) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		if event.Key == tl.KeyEnter {
			GamePlay()
		}
	}
}
