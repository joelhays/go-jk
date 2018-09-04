package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	// "github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 800
	height = 600
)

var jklData Jkl
var centroid [3]float32

var camera Camera
var previousTime float64

var lightPos = mgl32.Vec3{1.2, 1.0, 2.0}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	// jklData = ReadJKL("./jkl/01narshadda.jkl")
	jklData = ReadJKL("./jkl/m_boss15.jkl")
	// jklData = ReadJKL("./jkl/test.jkl")

	fmt.Println(jklData)

	var points []float32
	for _, surface := range jklData.Surfaces {
		for _, id := range surface.VertexIds {
			points = append(points, float32(jklData.Vertices[id].X))
			points = append(points, float32(jklData.Vertices[id].Y))
			points = append(points, float32(jklData.Vertices[id].Z))

			points = append(points, float32(surface.Normal.X))
			points = append(points, float32(surface.Normal.Y))
			points = append(points, float32(surface.Normal.Z))
		}
	}

	for _, vertex := range jklData.Vertices {
		centroid[0] += vertex.X
		centroid[1] += vertex.Y
		centroid[2] += vertex.Z
	}

	centroid[0] /= float32(len(jklData.Vertices))
	centroid[1] /= float32(len(jklData.Vertices))
	centroid[2] /= float32(len(jklData.Vertices))

	camera = NewCamera(mgl32.Vec3{5, 5, 0}, mgl32.Vec3{0, 0, 1}, 0, -90)

	// vao := makeVao(triangle)
	// vao := makeVao(cube)
	vao := makeVao(points)
	for !window.ShouldClose() {
		draw(vao, window, program)
	}
}

func draw(vao uint32, window *glfw.Window, program uint32) {
	deltaTime := glfw.GetTime() - previousTime
	previousTime = glfw.GetTime()

	doMovement(deltaTime)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(program)

	// vertex shader uniforms

	projection := mgl32.Perspective(mgl32.DegToRad(float32(camera.Zoom)), float32(width)/height, 0.1, 1000.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	cameraView := camera.GetViewMatrix()
	cameraUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &cameraView[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// fragment shader uniforms

	objectColor := mgl32.Vec3{1, 1, 1}
	objectColorUniform := gl.GetUniformLocation(program, gl.Str("objectColor\x00"))
	gl.Uniform3fv(objectColorUniform, 1, &objectColor[0])

	lightColor := mgl32.Vec3{1, 1, 1}
	lightColorUniform := gl.GetUniformLocation(program, gl.Str("lightColor\x00"))
	gl.Uniform3fv(lightColorUniform, 1, &lightColor[0])

	lightPosUniform := gl.GetUniformLocation(program, gl.Str("lightPos\x00"))
	// gl.Uniform3fv(lightPosUniform, 1, &lightPos[0])
	gl.Uniform3fv(lightPosUniform, 1, &camera.Position[0])

	viewPos := camera.Position
	viewPosUniform := gl.GetUniformLocation(program, gl.Str("viewPos\x00"))
	gl.Uniform3fv(viewPosUniform, 1, &viewPos[0])

	gl.BindVertexArray(vao)

	var offset int32

	for _, surface := range jklData.Surfaces {
		numVerts := int32(len(surface.VertexIds))

		if surface.Geo != 0 {
			gl.DrawArrays(gl.TRIANGLE_FAN, offset, int32(len(surface.VertexIds)))
			// gl.DrawArrays(gl.LINE_LOOP, offset, int32(len(surface.VertexIds)))
		}

		offset = offset + numVerts
	}
	// gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
	// gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cube)/6))

	glfw.PollEvents()
	window.SwapBuffers()
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "J", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetKeyCallback(KeyCallback)
	window.SetCursorPosCallback(MouseCallback)

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)

	vertexShaderSource := readShader("./shaders/vertex.glsl")
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShaderSource := readShader("./shaders/fragment.glsl")
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))

	return vao
}

func doMovement(deltaTime float64) {

	if keyW, keyUp := keys[glfw.KeyW], keys[glfw.KeyUp]; keyW || keyUp {
		camera.ProcessKeyboard(CAMERA_FORWARD, deltaTime)
	}

	if keyS, keyDown := keys[glfw.KeyS], keys[glfw.KeyDown]; keyS || keyDown {
		camera.ProcessKeyboard(CAMERA_BACKWARD, deltaTime)
	}

	if keyA, keyLeft := keys[glfw.KeyA], keys[glfw.KeyLeft]; keyA || keyLeft {
		camera.ProcessKeyboard(CAMERA_LEFT, deltaTime)
	}

	if keyD, keyRight := keys[glfw.KeyD], keys[glfw.KeyRight]; keyD || keyRight {
		camera.ProcessKeyboard(CAMERA_RIGHT, deltaTime)
	}
}
