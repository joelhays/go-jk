package opengl

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type ShaderProgram struct {
	ProgramID        uint32
	vertexShaderID   uint32
	fragmentShaderID uint32
}

func NewShaderProgram(vertexFile string, fragmentFile string) *ShaderProgram {
	program := &ShaderProgram{}

	var err error

	vertexShaderSource := readShader(vertexFile)
	program.vertexShaderID, err = compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShaderSource := readShader(fragmentFile)
	program.fragmentShaderID, err = compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program.ProgramID = gl.CreateProgram()
	gl.AttachShader(program.ProgramID, program.vertexShaderID)
	gl.AttachShader(program.ProgramID, program.fragmentShaderID)
	gl.LinkProgram(program.ProgramID)
	gl.ValidateProgram(program.ProgramID)

	return program
}

func (p *ShaderProgram) Start() {
	gl.UseProgram(p.ProgramID)
}

func (p *ShaderProgram) Stop() {
	gl.UseProgram(0)
}

func (p *ShaderProgram) SetMatrixUniform(uniformName string, mat mgl32.Mat4) {
	uniform := gl.GetUniformLocation(p.ProgramID, gl.Str(uniformName+"\x00"))
	gl.UniformMatrix4fv(uniform, 1, false, &mat[0])
}

func (p *ShaderProgram) SetVectorUniform(uniformName string, vec mgl32.Vec3) {
	uniform := gl.GetUniformLocation(p.ProgramID, gl.Str(uniformName+"\x00"))
	gl.Uniform3fv(uniform, 1, &vec[0])
}

func (p *ShaderProgram) SetIntegerUniform(uniformName string, value int32) {
	uniform := gl.GetUniformLocation(p.ProgramID, gl.Str(uniformName+"\x00"))
	gl.Uniform1i(uniform, value)
}

func (p *ShaderProgram) Cleanup() {
	gl.DetachShader(p.ProgramID, p.vertexShaderID)
	gl.DetachShader(p.ProgramID, p.fragmentShaderID)
	gl.DeleteShader(p.vertexShaderID)
	gl.DeleteShader(p.fragmentShaderID)
	gl.DeleteProgram(p.ProgramID)
}

func readShader(filePath string) string {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	shaderSource := string(bytes) + "\x00"

	return shaderSource
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

		shaderLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(shaderLog))

		return 0, fmt.Errorf("failed to compile %v: %v", source, shaderLog)
	}

	return shader, nil
}
