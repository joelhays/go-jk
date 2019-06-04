package scene

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/joelhays/go-jk/camera"
	"github.com/joelhays/go-jk/jk"
	"github.com/joelhays/go-jk/opengl"
)

type MenuScene struct {
	backgroundImage string
	shaderProgram   *opengl.ShaderProgram
	renderers       []opengl.Renderer
	cam             *camera.Camera
	window          *glfw.Window
}

func NewMenuScene(backgroundImage string, window *glfw.Window, cam *camera.Camera, shaderProgram *opengl.ShaderProgram) *MenuScene {
	return &MenuScene{backgroundImage: backgroundImage, window: window, cam: cam, shaderProgram: shaderProgram}
}

func (s *MenuScene) Load() {
	bmFile := jk.GetLoader().LoadBM(s.backgroundImage)
	bmRenderer := opengl.NewOpenGlBmRenderer(&bmFile, s.shaderProgram)
	s.renderers = append(s.renderers, bmRenderer)
}

func (s *MenuScene) Unload() {

}

func (s *MenuScene) Update() {
	opengl.Draw(s.window, s.cam, s.renderers)
}
