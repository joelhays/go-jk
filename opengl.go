package main

import (
	"fmt"
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
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 9*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 9*4, gl.PtrOffset(3*4))

	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 9*4, gl.PtrOffset(6*4))

	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointer(3, 1, gl.FLOAT, false, 9*4, gl.PtrOffset(7*4))

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

		var finalTexture []byte
		if false { //} material.Transparent {
			finalTexture = make([]byte, material.SizeX*material.SizeY*4)
			for j := 0; j < int(material.SizeX*material.SizeY); j++ {
				finalTexture[j*4] = jklData.ColorMaps[0].Pallette[material.Texture[j]].R
				finalTexture[j*4+1] = jklData.ColorMaps[0].Pallette[material.Texture[j]].G
				finalTexture[j*4+2] = jklData.ColorMaps[0].Pallette[material.Texture[j]].B

				if material.Texture[j] == 0 {
					finalTexture[j*4+3] = 0
				} else {
					finalTexture[j*4+3] = 255
				}
				gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, material.SizeX, material.SizeY, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(finalTexture))
			}
		} else {
			finalTexture = make([]byte, material.SizeX*material.SizeY*3)
			for j := 0; j < int(material.SizeX*material.SizeY); j++ {
				finalTexture[j*3] = jklData.ColorMaps[0].Pallette[material.Texture[j]].R
				finalTexture[j*3+1] = jklData.ColorMaps[0].Pallette[material.Texture[j]].G
				finalTexture[j*3+2] = jklData.ColorMaps[0].Pallette[material.Texture[j]].B
			}
			gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, material.SizeX, material.SizeY, 0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(finalTexture))
		}

		gl.GenerateMipmap(gl.TEXTURE_2D)

		gl.BindTexture(gl.TEXTURE_2D, 0)
	}

	return textures
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
			// gl.DrawArrays(gl.POINTS, offset, int32(len(surface.VertexIds)))
			// gl.PointSize(10)

			gl.BindTexture(gl.TEXTURE_2D, 0)
		}

		offset = offset + numVerts
	}

	// gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))

	// gl.ActiveTexture(gl.TEXTURE0)
	// gl.BindTexture(gl.TEXTURE_2D, (*textures)[185])
	// textureUniform := gl.GetUniformLocation(program, gl.Str("objectTexture\x00"))
	// gl.Uniform1i(textureUniform, 0)
	// gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cube)/6))

	glfw.PollEvents()
	window.SwapBuffers()
}
