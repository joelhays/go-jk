package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	// "github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 500
	height = 500

	vertexShaderSource = `
		#version 410
		layout (location = 0) in vec3 position;
		layout (location = 1) in vec3 normal;

		out vec3 Normal;
		out vec3 FragPos;

		uniform mat4 model;
		uniform mat4 view;
		uniform mat4 projection;

		void main() {
			gl_Position = projection * view * model * vec4(position, 1.0);
			FragPos = vec3(model * vec4(position, 1.0));
			Normal = mat3(transpose(inverse(model))) * normal;
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410

		in vec3 FragPos;
		in vec3 Normal;

		uniform vec3 lightPos;
		uniform vec3 viewPos;
		uniform vec3 objectColor;
		uniform vec3 lightColor;

		out vec4 frag_color;
		void main() {
			// frag_color = vec4(1, 1, 1, 1.0);

			// ambient
			float ambientStrength = .1f;
			vec3 ambient = ambientStrength * lightColor;

			// diffuse
			vec3 norm = normalize(Normal);
			vec3 lightDirection = normalize(lightPos - FragPos);
			float diff = max(dot(norm, lightDirection), 0.0);
			vec3 diffuse = diff * lightColor;

			// specular
			float specularStrength = 0.5f;
			vec3 viewDirection = normalize(viewPos - FragPos);
			vec3 reflectDirection = reflect(-lightDirection, norm);
			float spec = pow(max(dot(viewDirection, reflectDirection), 0.0), 32);
			vec3 specular = specularStrength * spec* lightColor;

			vec3 result = (ambient + diffuse + specular) * objectColor;
			frag_color = vec4(result, 1.0f);

			// frag_color = vec4(objectColor, 1);
		}
	` + "\x00"
)

var (
	triangle = []float32{
		0, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}
)

var jklData jkl
var centroid [3]float32

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	// vao := makeVao(triangle)

	// jklFile, _ := filepath.Abs("./01narshadda.jkl")
	jklFile, _ := filepath.Abs("./m_boss15.jkl")
	jklData = ReadJKL(jklFile)
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

	vao := makeVao(points)
	for !window.ShouldClose() {
		draw(vao, window, program)
	}
}

func draw(vao uint32, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(program)

	// vertex shader uniforms

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.1, 1000.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{300, 300, centroid[2]}, mgl32.Vec3{centroid[0], centroid[1], centroid[2]}, mgl32.Vec3{0, 0, 1})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

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

	lightPos := mgl32.Vec3{300, 300, centroid[2]}
	lightPosUniform := gl.GetUniformLocation(program, gl.Str("lightPos\x00"))
	gl.Uniform3fv(lightPosUniform, 1, &lightPos[0])

	viewPos := mgl32.Vec3{300, 300, centroid[2]}
	viewPosUniform := gl.GetUniformLocation(program, gl.Str("viewPos\x00"))
	gl.Uniform3fv(viewPosUniform, 1, &viewPos[0])

	gl.BindVertexArray(vao)

	var offset int32
	for _, surface := range jklData.Surfaces {
		numVerts := int32(len(surface.VertexIds))
		gl.DrawArrays(gl.TRIANGLE_FAN, offset, int32(len(surface.VertexIds)))
		// gl.DrawArrays(gl.LINE_LOOP, offset, int32(len(surface.VertexIds)))
		offset = offset + numVerts
	}
	// gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))

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
	gl.CullFace(gl.BACK)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

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

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))

	return vao
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
