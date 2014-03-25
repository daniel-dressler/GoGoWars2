package main

import (
	"github.com/jackyb/go-gl/gl"
	"github.com/jackyb/go-sdl2/sdl"
)


type Display struct {
	window *sdl.Window
	context sdl.GLContext
}

func new_display() Display {
	var winTitle string = "Go-SDL2 + Go-GL"
	var winWidth, winHeight int = 800, 600
	var this Display

	if 0 != sdl.Init(sdl.INIT_EVERYTHING) {
		panic(sdl.GetError())
	}
	this.window = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_OPENGL)
	if this.window == nil {
		panic(sdl.GetError())
	}
	this.context = sdl.GL_CreateContext(this.window)
	if this.context == nil {
		panic(sdl.GetError())
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.2, 0.2, 0.3, 1.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
	gl.Viewport(0, 0, int32(winWidth), int32(winHeight))

	return this
}

func (this Display) update() {
	drawgl()
	sdl.GL_SwapWindow(this.window)
}

func (this Display) shutdown() {
	sdl.GL_DeleteContext(this.context)
	this.window.Destroy()
	sdl.Quit()
}

func drawgl() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.Begin(gl.TRIANGLES)
	gl.Color3f(1.0, 0.0, 0.0)
	gl.Vertex2f(0.5, 0.0)
	gl.Color3f(0.0, 1.0, 0.0)
	gl.Vertex2f(-0.5, -0.5)
	gl.Color3f(0.0, 0.0, 1.0)
	gl.Vertex2f(-0.5, 0.5)
	gl.End()
}
