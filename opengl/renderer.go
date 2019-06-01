package opengl

type Renderer interface {
	Render()
	ShaderProgram() *ShaderProgram
}
