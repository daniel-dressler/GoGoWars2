package main

import (
	"math"
	"fmt"
	"os"
	"runtime"
	"github.com/go-gl/gl"
	//"github.com/go-gl/glh"
	"github.com/jackyb/go-sdl2/sdl"
	//"github.com/Jragonmiris/mathgl"
)


type Display struct {
	window *sdl.Window
	context sdl.GLContext
	planet *ds_sphere
	prog *ds_program
}

func new_display() Display {
	runtime.LockOSThread()
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

	sdl.GL_SetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	this.context = sdl.GL_CreateContext(this.window)
	if this.context == nil {
		panic(sdl.GetError())
	}

	gl.Init()

	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.2, 0.2, 0.3, 1.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
	gl.Viewport(0, 0, winWidth, winHeight)

	this.planet = gen_sphere(0.5, 6, 6);
	this.prog = gen_program(vShaderSrc, fShaderSrc);
	runtime.UnlockOSThread()
	return this
}

var vShaderSrc = `#version 120
// Input vertex data, different for all executions of this shader.
attribute vec3 vertexPosition_modelspace;
 
void main(){
 
gl_Position = vec4(vertexPosition_modelspace, 1.0);
 
}`

var fShaderSrc = `#version 120
 
void main()
{
 
// Output color = red
gl_FragColor = vec4(1,0,0,1);
 
}`

func (this Display) update() {
	runtime.LockOSThread()
	drawgl()
	this.planet.draw(0,0,0);
	sdl.GL_SwapWindow(this.window)
}

func (this Display) shutdown() {
	sdl.GL_DeleteContext(this.context)
	this.window.Destroy()
	sdl.Quit()
}

func drawgl() {
	runtime.LockOSThread()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	//gl.ShadeModel(gl.FLAT)

/*
	gl.Begin(gl.TRIANGLES)
	gl.Color3f(1.0, 0.0, 0.0)
	gl.Vertex2f(0.5, 0.0)
	gl.Color3f(0.0, 1.0, 0.0)
	gl.Vertex2f(-0.5, -0.5)
	gl.Color3f(0.0, 0.0, 1.0)
	gl.Vertex2f(-0.5, 0.5)
	gl.End()
	*/

}

type ds_program struct {
	program gl.Program
}

func compile_shader(shader *gl.Shader, src string) {
	shader.Source(src)
	shader.Compile()

	status := shader.Get(gl.COMPILE_STATUS)
	infoLog := shader.GetInfoLog()
	if status != 1 {
		fmt.Println("Shader err:", status, infoLog)
		fmt.Println("source ", src)
		os.Exit(1)
	}
}

func gen_program(vSrc, fSrc string) *ds_program {
	vShader := gl.CreateShader(gl.VERTEX_SHADER)
	fShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	compile_shader(&vShader, vSrc);
	compile_shader(&fShader, fSrc);

	program := gl.CreateProgram()
	program.AttachShader(vShader)
	program.AttachShader(fShader)
	program.Link()

	status := program.Get(gl.LINK_STATUS)
	infoLog := program.GetInfoLog()
	if status != 1 {
		fmt.Println("Program link err:", status, infoLog)
		os.Exit(1)
	}

	program.Use()

	var this ds_program
	this.program = program
	return &this
}

type ds_sphere struct {
	vertices []float64
	normals []float64
	texcoords []float64
	indices []float64
	vBuffer gl.Buffer
}

func gen_sphere(radius float64, rings uint, sectors uint) *ds_sphere {
	sp := new(ds_sphere)
	sp.vertices  = make([]float64, 0)
	sp.normals   = make([]float64, 0)
	sp.texcoords = make([]float64, 0)

	R := 1.0 / float64(rings - 1)
	S := 1.0 / float64(sectors - 1);
	for r := uint(0); r < rings; r++ {
		for s := uint(0); s < sectors; s++ {
			Pi := math.Pi
			fr := float64(r)
			fs := float64(s)
			y := math.Sin( -(Pi / 2) + (Pi * fr * R))
			x := math.Cos( 2*Pi * fs * S) * math.Sin( Pi * fr * R)
			z := math.Sin( 2*Pi * fs * S) * math.Cos( Pi * fr * R)
			fmt.Println("x: ", x, "y: ", y, "z: ", z)
			//append(sp.texcoords, fs*S, fr*R)
			sp.vertices = append(sp.vertices, x * radius, y * radius, z * radius)
			//append(sp.normals, x, y, z)
		}
	}

	sp.indices = make([]float64, 0)
	for r := uint(0); r < rings - 1; r++ {
		for s := uint(0); s < sectors - 1; s++ {
			/*
			append(sp.indices,
				float64(r * sectors + s),
				float64(r * sectors + (s + 1)),
				float64((r + 1) * sectors + (s + 1)),
				float64((r + 1) * sectors + s))
			*/
		}
	}

	sp.vBuffer = gl.GenBuffer()
	sp.vBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(sp.vertices) * 8, sp.vertices, gl.STATIC_DRAW)

	return sp
}

func (this ds_sphere) draw(x float64, y float64, z float64) {
	gl.PushMatrix()
	gl.MatrixMode(gl.MODELVIEW)
	gl.Translated(x, y, z)

	gl.DrawArrays(gl.TRIANGLES, 0, len(this.vertices));

/*
	gl.EnableClientState(gl.VERTEX_ARRAY);
	gl.EnableClientState(gl.NORMAL_ARRAY);
	gl.EnableClientState(gl.TEXTURE_COORD_ARRAY);

	gl.Begin(gl.TRIANGLES);
	for i := 0; i < len(this.vertices); i += 3 {
		gl.Color3f(0.0, 1.0, 0.0)
		gl.Vertex3d(this.vertices[i], this.vertices[i + 1], this.vertices[i + 2])
	}
	gl.End();
	*/
	gl.PopMatrix()
}
