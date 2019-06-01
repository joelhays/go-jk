package opengl

import (
	"github.com/joelhays/go-jk/camera"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// InitGlfw initializes glfw and returns a Window to use.
func InitGlfw(windowWidth int, windowHeight int, keyCallback func(*glfw.Window, glfw.Key, int, glfw.Action, glfw.ModifierKey),
	mouseCallback func(*glfw.Window, float64, float64)) *glfw.Window {

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "JK Viewer", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetKeyCallback(keyCallback)
	window.SetCursorPosCallback(mouseCallback)

	return window
}

// initOpenGL initializes OpenGL and returns an initialized program.
func InitOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

func Draw(window *glfw.Window, camera *camera.Camera, renderers []Renderer) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	width, height := window.GetSize()

	for _, renderer := range renderers {
		program := renderer.ShaderProgram()
		program.Start()
		configureProgram(program, camera, width, height)
		renderer.Render()
		program.Stop()
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

func configureProgram(program *ShaderProgram, camera *camera.Camera, width int, height int) {
	// vertex shader uniforms
	projection := mgl32.Perspective(mgl32.DegToRad(float32(camera.Zoom)), float32(width)/float32(height), 0.1, 1000.0)
	program.SetMatrixUniform("projection", projection)
	program.SetMatrixUniform("view", camera.GetViewMatrix())

	// fragment shader uniforms
	program.SetVectorUniform("objectColor", mgl32.Vec3{1, 1, 1})
	program.SetVectorUniform("lightColor", mgl32.Vec3{1, 1, 1})
	program.SetVectorUniform("lightPos", camera.Position)
	program.SetVectorUniform("viewPos", camera.Position)
}
