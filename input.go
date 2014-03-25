package main

import (
	"fmt"
	"github.com/jackyb/go-sdl2/sdl"
)

type Input struct {

}

func new_input() Input {
	var this Input
	return this
}

func (this Input) update() {
	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.MouseMotionEvent:
			fmt.Printf("[%d ms] MouseMotion\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n", t.Timestamp, t.Which, t.X, t.Y, t.XRel, t.YRel)
		}
	}
}

func (this Input) shutdown() {
}
