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
func InitOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	vertexShaderSource := ReadShader("./shaders/vertex.glsl")
	vertexShader, err := CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShaderSource := ReadShader("./shaders/fragment.glsl")
	fragmentShader, err := CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func DrawRenderer(window *glfw.Window, camera *camera.Camera, levelRenderer *OpenGlLevelRenderer, modelRenderers []*OpenGl3doRenderer, bmRenderer *OpenGlBmRenderer) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	width, height := window.GetSize()

	if levelRenderer != nil {
		configureProgram(levelRenderer.Program, camera, width, height)
		levelRenderer.Render()
	}

	if modelRenderers != nil {
		for idx, modelRenderer := range modelRenderers {
			if modelRenderer == nil {
				_ = idx
				// fmt.Println("nil renderer at", idx)
				continue
			}

			configureProgram(modelRenderer.Program, camera, width, height)
			modelRenderer.Render()
		}
	}

	if bmRenderer != nil {
		configureProgram(bmRenderer.Program, camera, width, height)
		bmRenderer.Render()
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

func configureProgram(program uint32, camera *camera.Camera, width int, height int) {
	gl.UseProgram(program)

	// vertex shader uniforms

	projection := mgl32.Perspective(mgl32.DegToRad(float32(camera.Zoom)), float32(width)/float32(height), 0.1, 1000.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	cameraView := camera.GetViewMatrix()
	cameraUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &cameraView[0])

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
}
