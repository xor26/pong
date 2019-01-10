package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

const winWidth, winHeight int = 1200, 800
const ballBasicVelocityX, ballBasicVelocityY int = 10,10

type color struct {
	r, g, b byte
}

type paddle struct {
	x, y          int
	width, height int
	color         color
}

type aiPaddle struct {
	paddle
}

func (padddle *paddle) draw(pixels []byte) {
	for x := -padddle.width; x < padddle.width; x++ {
		for y := -padddle.height; y < padddle.height; y++ {
			setPixelColor(padddle.x+x, padddle.y+y, padddle.color, pixels)
		}
	}
}

func (padddle *aiPaddle) draw(pixels []byte) {
	for x := -padddle.width; x < padddle.width; x++ {
		for y := -padddle.height; y < padddle.height; y++ {
			setPixelColor(padddle.x+x, padddle.y+y, padddle.color, pixels)
		}
	}
}

func (padddle *aiPaddle) update(ball *ball) {
	if ball.y >= 10+padddle.height &&  ball.y < winHeight - padddle.height {
		padddle.y = ball.y
	}
}

func (padddle *paddle) update(keyState []byte) {
	if keyState[sdl.SCANCODE_W] != 0 && padddle.y >= 10+padddle.height {
		padddle.y -= 10
	}
	if keyState[sdl.SCANCODE_S] != 0 && padddle.y < winHeight - padddle.height{
		padddle.y += 10
	}
}

type ball struct {
	x, y       int
	xVel, yVel int
	radius     int
	color      color
}

func (ball *ball) draw(pixels []byte) {
	for x := -ball.radius; x < ball.radius; x++ {
		for y := -ball.radius; y < ball.radius; y++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixelColor(ball.x+x, ball.y+y, ball.color, pixels)
			}
		}
	}
}

func (ball *ball) update(leftPaddle *paddle, rightPaddle *aiPaddle) {
	ball.x += ball.xVel
	ball.y += ball.yVel

	if ball.x-ball.radius <= 0 {
		ball.x = 800
		ball.y = 600
		ball.xVel = ballBasicVelocityX
		ball.yVel = -ballBasicVelocityY
	}

	if ball.x+ball.radius >= winWidth {
		ball.xVel = -ball.xVel
	}

	if ball.y-ball.radius <= 0 || ball.y+ball.radius >= winHeight {
		ball.yVel = -ball.yVel
	}

	isBallWithinLeftPaddleBorders := ball.y-ball.radius >= leftPaddle.y-leftPaddle.height && ball.y+ball.radius <= leftPaddle.y+leftPaddle.height
	isBallTouchLeftPaddle := ball.x-ball.radius <= leftPaddle.x + leftPaddle.width
	if isBallTouchLeftPaddle && isBallWithinLeftPaddleBorders {
		ball.xVel = int(math.Abs(float64(ball.xVel)))
	}

	isBallWithinRightPaddleBorders := ball.y-ball.radius >= rightPaddle.y-rightPaddle.height && ball.y+ball.radius <= rightPaddle.y+rightPaddle.height
	isBallTouchRightPaddle := ball.x+ball.radius >= rightPaddle.x - rightPaddle.width
	if isBallTouchRightPaddle && isBallWithinRightPaddleBorders {
		if ball.xVel > 0 {
			ball.xVel = - ball.xVel
		}
	}
}

func setPixelColor(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index > len(pixels)-4 {
		return
	}

	pixels[index] = c.r
	pixels[index+1] = c.g
	pixels[index+2] = c.b
}

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Pong", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	ballColor := color{0, 155, 0}
	ball := ball{600, 400, ballBasicVelocityX, ballBasicVelocityY, 20, ballColor}
	player1 := paddle{60, 400, 10, 60, color{0, 0, 100}}
	player2 := aiPaddle{paddle{winWidth-60, 400, 10, 60, color{100, 0, 0}}}

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		pixels := make([]byte, winWidth*winHeight*4)

		keyState := sdl.GetKeyboardState()
		player1.update(keyState)
		player1.draw(pixels)

		player2.update(&ball)
		player2.draw(pixels)

		ball.update(&player1, &player2)
		ball.draw(pixels)

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		sdl.Delay(10)

	}
}
