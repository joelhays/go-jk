package main

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

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

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

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

func drawRenderer(window *glfw.Window, modelRenderers []*OpenGlModelRenderer) {
	deltaTime := glfw.GetTime() - previousTime
	previousTime = glfw.GetTime()

	doMovement(deltaTime)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for idx, modelRenderer := range modelRenderers {
		if modelRenderer == nil {
			_ = idx
			// fmt.Println("nil renderer at", idx)
			continue
		}

		gl.UseProgram(modelRenderer.Program)

		// vertex shader uniforms

		projection := mgl32.Perspective(mgl32.DegToRad(float32(camera.Zoom)), float32(width)/height, 0.1, 1000.0)
		projectionUniform := gl.GetUniformLocation(modelRenderer.Program, gl.Str("projection\x00"))
		gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

		cameraView := camera.GetViewMatrix()
		cameraUniform := gl.GetUniformLocation(modelRenderer.Program, gl.Str("view\x00"))
		gl.UniformMatrix4fv(cameraUniform, 1, false, &cameraView[0])

		// fragment shader uniforms

		objectColor := mgl32.Vec3{1, 1, 1}
		objectColorUniform := gl.GetUniformLocation(modelRenderer.Program, gl.Str("objectColor\x00"))
		gl.Uniform3fv(objectColorUniform, 1, &objectColor[0])

		lightColor := mgl32.Vec3{1, 1, 1}
		lightColorUniform := gl.GetUniformLocation(modelRenderer.Program, gl.Str("lightColor\x00"))
		gl.Uniform3fv(lightColorUniform, 1, &lightColor[0])

		lightPosUniform := gl.GetUniformLocation(modelRenderer.Program, gl.Str("lightPos\x00"))
		// gl.Uniform3fv(lightPosUniform, 1, &lightPos[0])
		gl.Uniform3fv(lightPosUniform, 1, &camera.Position[0])

		viewPos := camera.Position
		viewPosUniform := gl.GetUniformLocation(modelRenderer.Program, gl.Str("viewPos\x00"))
		gl.Uniform3fv(viewPosUniform, 1, &viewPos[0])

		modelRenderer.Render()
	}

	glfw.PollEvents()
	window.SwapBuffers()
}
