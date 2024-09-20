package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCREEN_WIDTH          = 800
	SCREEN_HEIGHT         = 450
	CENTER_X              = SCREEN_WIDTH / 2
	CENTER_Y              = SCREEN_HEIGHT / 2
	SIDE_BUFFER_WIDTH     = 40
	INITIAL_PADDLE_HEIGHT = 80
	INITIAL_PADDLE_WIDTH  = 20
	INITIAL_PADDLE_SPEED  = 4
	INITIAL_BALL_SIZE     = 16
	INITIAL_BALL_SPEED    = 4
	RIGHT_BOUNDS          = SCREEN_WIDTH
	LEFT_BOUNDS           = 0
	TOP_BOUNDS            = 0
	BOTTOM_BOUNDS         = SCREEN_HEIGHT
)

type Paddle struct {
	X      int32
	Y      int32
	VX     int32
	VY     int32
	Width  int32
	Height int32
	Score  int32
	Speed  int32
	Color  color.RGBA
}

func (p *Paddle) Draw() {
	rl.DrawRectangle(p.X, p.Y, p.Width, p.Height, p.Color)
}

func (p *Paddle) IncrementScore() {
	p.Score += 1
}

func (p *Paddle) Update() {
	p.X += p.VX
	p.Y += p.VY
}

type Ball struct {
	X     int32
	Y     int32
	VX    int32
	VY    int32
	Size  int32
	Serve int32
	Color color.RGBA
}

func (b *Ball) Draw() {
	rl.DrawRectangle(b.X, b.Y, b.Size, b.Size, b.Color)
}

func (b *Ball) Update() {
	b.X += b.VX
	b.Y += b.VY
}

func (b *Ball) Reset() {
	b.Size = INITIAL_BALL_SIZE
	b.X = CENTER_X - (INITIAL_BALL_SIZE / 2)
	b.Y = CENTER_Y - (INITIAL_BALL_SIZE / 2)
	b.Color = rl.White
	b.VX = INITIAL_BALL_SPEED * b.serveModifier()
	b.VY = INITIAL_BALL_SPEED * pickRandomVariant()
	b.Serve += 1
}

func (b *Ball) BounceVertically() {
	b.VY = b.VY * -1
}

func (b *Ball) BounceHorizontally() {
	b.VX = b.VX * -1
}

func (b *Ball) serveModifier() int32 {
	if b.Serve%2 == 0 {
		return 1
	}

	return -1
}

func isBallOutToLeft(b *Ball) bool {
	return b.X < LEFT_BOUNDS
}

func isBallOutToRight(b *Ball) bool {
	return b.X-b.Size > RIGHT_BOUNDS
}

func ballAtBottom(b *Ball) bool {
	return b.Y > SCREEN_HEIGHT
}

func ballAtTop(b *Ball) bool {
	return b.Y+b.Size < 0
}

func paddleAtTop(p *Paddle) bool {
	return p.Y < TOP_BOUNDS
}

func paddleAtBottom(p *Paddle) bool {
	return p.Y+p.Height > BOTTOM_BOUNDS
}

func pickRandomVariant() int32 {
	return int32(rand.Intn(2)*2 - 1)
}

func getScoreText(leftPaddle *Paddle, rightPaddle *Paddle) string {
	return fmt.Sprintf("%d   %d", leftPaddle.Score, rightPaddle.Score)
}

func movePaddleToInterceptBall(p *Paddle, b *Ball) {
	paddleCenter := p.Y + p.Height/2
	if b.Y > paddleCenter {
		p.VY = p.Speed
	} else if b.Y < paddleCenter {
		p.VY = -p.Speed
	} else {
		p.VY = 0
	}
}

func movePaddleTowardsCenter(p *Paddle) {
	paddleCenter := p.Y + p.Height/2
	if paddleCenter < CENTER_Y {
		p.VY = p.Speed
	} else if paddleCenter > CENTER_Y {
		p.VY = -p.Speed
	} else {
		p.VY = 0
	}
}

func ballIntersectsPaddle(b *Ball, p *Paddle) bool {
	return b.X < p.X+p.Width &&
		b.X+b.Size > p.X &&
		b.Y+b.Size > p.Y &&
		b.Y < p.Y+p.Height
}

func main() {
	rand.Seed(time.Now().UnixNano())
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Pong")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	ball := Ball{}
	ball.Reset()

	leftPaddle := Paddle{
		Width:  INITIAL_PADDLE_WIDTH,
		Height: INITIAL_PADDLE_HEIGHT,
		Speed:  INITIAL_PADDLE_SPEED,
		Color:  rl.White,
		X:      SIDE_BUFFER_WIDTH,
		Y:      CENTER_Y - (INITIAL_PADDLE_HEIGHT / 2),
		Score:  0,
	}

	rightPaddle := Paddle{
		Width:  INITIAL_PADDLE_WIDTH,
		Height: INITIAL_PADDLE_HEIGHT,
		Speed:  INITIAL_PADDLE_SPEED,
		Color:  rl.White,
		X:      SCREEN_WIDTH - SIDE_BUFFER_WIDTH - INITIAL_PADDLE_WIDTH,
		Y:      CENTER_Y - (INITIAL_PADDLE_HEIGHT / 2),
		Score:  0,
	}

	for !rl.WindowShouldClose() {
		ball.Update()

		if isBallOutToLeft(&ball) {
			rightPaddle.IncrementScore()
			ball.Reset()
		}

		if isBallOutToRight(&ball) {
			leftPaddle.IncrementScore()
			ball.Reset()
		}

		if ballAtTop(&ball) || ballAtBottom(&ball) {
			ball.BounceVertically()
		}

		leftPaddle.Update()
		rightPaddle.Update()

		if ball.VX < 0 {
			movePaddleTowardsCenter(&rightPaddle)
		} else {
			movePaddleToInterceptBall(&rightPaddle, &ball)
		}

		if rl.IsKeyDown(rl.KeyUp) {
			leftPaddle.VY = -leftPaddle.Speed
		} else if rl.IsKeyDown(rl.KeyDown) {
			leftPaddle.VY = leftPaddle.Speed
		} else {
			leftPaddle.VY = 0
		}

		if paddleAtTop(&rightPaddle) {
			rightPaddle.Y = 0
		}

		if paddleAtBottom(&rightPaddle) {
			rightPaddle.Y = SCREEN_HEIGHT - rightPaddle.Height
		}

		if paddleAtTop(&leftPaddle) {
			leftPaddle.Y = 0
		}

		if paddleAtBottom(&leftPaddle) {
			leftPaddle.Y = SCREEN_HEIGHT - leftPaddle.Height
		}

		if ballIntersectsPaddle(&ball, &leftPaddle) ||
			ballIntersectsPaddle(&ball, &rightPaddle) {
			ball.BounceHorizontally()
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		rl.DrawLine(CENTER_X, TOP_BOUNDS, CENTER_X, BOTTOM_BOUNDS, rl.LightGray)

		scoreText := getScoreText(&leftPaddle, &rightPaddle)
		textSize := rl.MeasureText(scoreText, 20)

		rl.DrawText(scoreText, CENTER_X-(textSize/2), 20, 20, rl.White)

		leftPaddle.Draw()
		rightPaddle.Draw()
		ball.Draw()

		rl.EndDrawing()
	}
}
