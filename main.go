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
	width  = 1024
	height = 768
)

var jklData Jkl

var camera Camera
var previousTime float64

var lightPos = mgl32.Vec3{1.2, 1.0, 2.0}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	// matBytes := LoadFileFromGOB("J:\\Resource\\Res2.gob", "00cglyf3.mat")
	// data := parseMatFile(matBytes)
	// _ = data
	// return

	jklBytes := LoadFileFromGOB("J:\\Episode\\JK1MP.GOB", "jkl\\m_boss15.jkl")
	jklData = ReadJKLFromString(string(jklBytes))

	// jklData = ReadJKLFromFile("./jkl/01narshadda.jkl")
	// jklData = ReadJKL("./jkl/test.jkl")

	// fmt.Println(jklData)

	var points []float32
	for _, surface := range jklData.Surfaces {
		var mat material
		if surface.MaterialID != -1 {
			mat = jklData.Materials[surface.MaterialID]
		}

		for _, id := range surface.VertexIds {
			points = append(points, float32(jklData.Vertices[id][0]))
			points = append(points, float32(jklData.Vertices[id][1]))
			points = append(points, float32(jklData.Vertices[id][2]))

			points = append(points, float32(surface.Normal[0]))
			points = append(points, float32(surface.Normal[1]))
			points = append(points, float32(surface.Normal[2]))

			points = append(points, jklData.TextureVertices[id][0]/float32(mat.SizeX))
			points = append(points, jklData.TextureVertices[id][1]/float32(mat.SizeY))
		}
	}

	camera = NewCamera(mgl32.Vec3{5, 5, 0}, mgl32.Vec3{0, 0, 1}, 0, -90)

	// vao := makeVao(triangle)
	// vao := makeVao(cube)
	vao := makeVao(points)

	textures := makeTextures()
	for !window.ShouldClose() {
		draw(vao, &textures, window, program)
	}
}

func draw(vao uint32, textures *[]uint32, window *glfw.Window, program uint32) {
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

			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, (*textures)[surface.MaterialID])
			textureUniform := gl.GetUniformLocation(program, gl.Str("objectTexture\x00"))
			gl.Uniform1i(textureUniform, 0)

			gl.DrawArrays(gl.TRIANGLE_FAN, offset, int32(len(surface.VertexIds)))
			// gl.DrawArrays(gl.LINE_LOOP, offset, int32(len(surface.VertexIds)))

			gl.BindTexture(gl.TEXTURE_2D, 0)
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
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))

	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))

	return vao
}

func makeTextures() []uint32 {

	numTextures := int32(len(jklData.Materials))

	textures := make([]uint32, numTextures)

	gl.GenTextures(numTextures, &textures[0])

	for i := int32(0); i < numTextures; i++ {
		textureID := textures[i]
		material := jklData.Materials[i]

		gl.BindTexture(gl.TEXTURE_2D, textureID)

		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

		if len(jklData.Materials[i].Texture) == 0 {
			fmt.Println("empty material")
			continue
		}

		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, material.SizeX, material.SizeY, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(material.Texture))
		gl.GenerateMipmap(gl.TEXTURE_2D)

		gl.BindTexture(gl.TEXTURE_2D, 0)
	}

	return textures
}

func doMovement(deltaTime float64) {

	if keyMinus := keys[glfw.KeyKPSubtract]; keyMinus {
		camera.MovementSpeed = .75
	}

	if keyDecimal := keys[glfw.KeyKPDecimal]; keyDecimal {
		camera.MovementSpeed = 6
	}

	if keyPlus := keys[glfw.KeyKPAdd]; keyPlus {
		camera.MovementSpeed = 12
	}

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
