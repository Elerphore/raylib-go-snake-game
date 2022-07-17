package main

import (
	"fmt"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type food struct {
	eated bool
	pos   point
}

type GameState int

const (
	Pause GameState = 0
	Game  GameState = 1
	Died  GameState = 2
)

type SnakeVector int

const (
	Up    SnakeVector = 0
	Right SnakeVector = 1
	Down  SnakeVector = 2
	Left  SnakeVector = 3
)

type point struct {
	x int32
	y int32
}

var score = 0
var lastScore = 0

var fontsize int32 = 30

var descriptionTitle = "Press any key to START"

var height int32 = 1000
var width int32 = 1000
var cellSize int32 = 50

var rowsCount = height / cellSize
var columnsCount = width / cellSize

var snake = make([]point, 1)
var state = Pause

var snakeVector = Right

func main() {
	x := func() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		switch state {
		case 0:
			pause()
		case 1:
			game()
		case 2:
			died()
		}

		rl.EndDrawing()
	}

	render(x)
}

func pause() {
	title := "Yep, Another Snake Game..."

	titleLength := rl.MeasureText(title, fontsize)
	descriptionTitleLength := rl.MeasureText(descriptionTitle, fontsize)

	rl.DrawText(title, width/2-int32(titleLength/2), height/2-fontsize-30, fontsize, rl.LightGray)
	rl.DrawText(descriptionTitle, width/2-int32(descriptionTitleLength/2), height/2-fontsize, fontsize, rl.LightGray)

	handleGameStarting()
}

func died() {

	title := fmt.Sprintf("You've died... Your score: %d", lastScore)

	titleLength := rl.MeasureText(title, fontsize)
	descriptionTitleLength := rl.MeasureText(descriptionTitle, fontsize)

	rl.DrawText(title, width/2-int32(titleLength/2), height/2-fontsize-30, fontsize, rl.LightGray)
	rl.DrawText(descriptionTitle, width/2-int32(descriptionTitleLength/2), height/2-fontsize, fontsize, rl.LightGray)

	handleGameStarting()
}

func handleGameStarting() {
	if rl.GetKeyPressed() != 0 {
		state = Game
	}
}

func game() {
	snakePositionHandling()
	renderSnake()
	renderGrid()
	renderScore()
	foodManager()

	collisions()
}

func collisions() {
	handleSnakeFoodCollision()
	handleSnakeHeadBodyCollision()
}

func handleSnakeHeadBodyCollision() {
	head := rl.Rectangle{X: float32(snake[0].x), Y: float32(snake[0].y), Width: 1, Height: 1}

	for index, bodyPart := range snake {
		if index != 0 {
			body := rl.Rectangle{X: float32(bodyPart.x), Y: float32(bodyPart.y), Width: 1, Height: 1}
			if rl.CheckCollisionRecs(body, head) {
				snake = make([]point, 1)
				lastScore = score
				score = 0
				state = Died
				break
			}
		}
	}
}

func handleSnakeFoodCollision() {
	fdRect := rl.Vector2{
		X: float32(apple.pos.x),
		Y: float32(apple.pos.y),
	}

	snakeRect := rl.Rectangle{
		X:      float32(snake[0].x),
		Y:      float32(snake[0].y),
		Width:  float32(0.5),
		Height: float32(0.5),
	}

	if rl.CheckCollisionPointRec(fdRect, snakeRect) && !apple.eated {
		fmt.Println(
			" fdRect: ", fdRect,
			" snakeRect: ", snakeRect,
		)

		onSnakeFoodCollision()
	}
}

func onSnakeFoodCollision() {
	lastPartPosition := snake[len(snake)-1]

	pt := getSnakeNextPosition()
	if pt != nil {
		moveSnakeBody(pt, 0)
	}
	frameCounter = 0

	snake =
		append(
			snake,
			lastPartPosition,
		)

	apple.eated = true
	apple = spawnFood()
	score++
}

func spawnFood() food {
	var fd food
	for {
		fd = food{eated: false, pos: point{x: rl.GetRandomValue(1, width/cellSize-1), y: rl.GetRandomValue(1, height/cellSize-1)}}
		var collision bool

		for _, bodyPart := range snake {
			collision = rl.CheckCollisionPointRec(
				rl.Vector2{X: float32(fd.pos.x), Y: float32(fd.pos.y)},
				rl.Rectangle{X: float32(bodyPart.x), Y: float32(bodyPart.y), Width: 0.5, Height: 0.5},
			)

			if collision {
				break
			}
		}

		if !collision {
			break
		}
	}

	return fd
}

var apple food = spawnFood()

func foodManager() {
	renderFood()
}

func renderFood() {
	rl.DrawCircle(
		apple.pos.x*cellSize+int32(float32(cellSize)/2),
		apple.pos.y*cellSize+int32(float32(cellSize)/2),
		float32(cellSize)/2, rl.Brown)
}

var frameCounter = 0

func snakePositionHandling() {
	if snake[0].x > (width/cellSize)-1 {
		snake[0].x = 0
	} else if snake[0].x < 0 {
		snake[0].x = (width / cellSize) - 1
	}

	if snake[0].y < 0 {
		snake[0].y = (height / cellSize) - 1
	} else if snake[0].y > (height/cellSize)-1 {
		snake[0].y = 0
	}

	changeSnakeVector()

	frameCounter++

	if frameCounter > 10 {
		pt := getSnakeNextPosition()
		if pt != nil {
			moveSnakeBody(pt, 0)
		}
		frameCounter = 0
	}
}

func changeSnakeVector() {
	key := rl.GetKeyPressed()

	if key == rl.KeyRight && snakeVector != Left {
		snakeVector = Right
	}

	if key == rl.KeyLeft && snakeVector != Right {
		snakeVector = Left
	}

	if key == rl.KeyUp && snakeVector != Down {
		snakeVector = Up
	}

	if key == rl.KeyDown && snakeVector != Up {
		snakeVector = Down
	}
}

func getSnakeNextPosition() *point {
	var pt *point = nil

	switch snakeVector {
	case Right:
		pt := point{x: snake[0].x + 1.0, y: snake[0].y}
		return &pt
	case Left:
		pt := point{x: snake[0].x - 1.0, y: snake[0].y}
		return &pt
	case Up:
		pt := point{x: snake[0].x, y: snake[0].y - 1.0}
		return &pt
	case Down:
		pt := point{x: snake[0].x, y: snake[0].y + 1.0}
		return &pt
	}

	return pt
}

func moveSnakeBody(point *point, index int) {
	oldPosition := snake[index]
	snake[index] = *point

	if index < (len(snake) - 1) {
		moveSnakeBody(&oldPosition, index+1)
	}
}

func renderSnake() {
	for _, item := range snake {
		rect := rl.Rectangle{X: float32(item.x * cellSize), Y: float32(item.y * cellSize), Width: float32(cellSize), Height: float32(cellSize)}
		rl.DrawRectangleRec(rect, rl.Blue)
	}
}

func renderScore() {
	rl.DrawText(strconv.Itoa(score), 10, 10, 50, rl.LightGray)
}

func renderGrid() {
	for i := int32(0); i < rowsCount; i++ {
		rl.DrawLine(
			int32(0),
			i*cellSize,
			width,
			i*cellSize,
			rl.Blue,
		)
	}

	for i := int32(0); i < columnsCount; i++ {
		rl.DrawLine(
			i*cellSize,
			int32(0),
			i*cellSize,
			height,
			rl.Blue,
		)
	}
}

func render(frame func()) {
	rl.InitWindow(width, height, "Yep, Another Snake Game")
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		frame()
	}

	rl.CloseWindow()
}
