package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)


const winWidth, winHeight int = 1200, 800

type color struct {
	r, g, b byte
}

type ball struct {
	x,y int
	xVel,yVel int
	radius int
	color color
}

func (ball *ball) draw(pixels []byte)  {
	for x := -ball.radius; x < ball.radius; x++ {
		for y := -ball.radius; y < ball.radius; y++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixelColor(ball.x+x, ball.y+y, ball.color, pixels)

			}
		}
	}
}

func (ball *ball) update()  {
	ball.x += ball.xVel
	ball.y += ball.yVel

	if ball.x-ball.radius <= 0 || ball.x+ball.radius >= winWidth {
		ball.xVel = -ball.xVel
	}

	if ball.y-ball.radius <= 0 || ball.y+ball.radius >= winHeight {
		ball.yVel = -ball.yVel
	}

}


func setPixelColor(x,y int, c color, pixels []byte){
	index := (y*winWidth + x) *4

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

	ballColor := color{255,0,0}
	ball := ball{600,400, 10,10, 20, ballColor}

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		pixels := make([]byte, winWidth*winHeight*4)

		ball.update()
		ball.draw(pixels)

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		sdl.Delay(10)

	}
}