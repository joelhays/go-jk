package scene

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/joelhays/go-jk/camera"
	"github.com/joelhays/go-jk/jk"
	"github.com/joelhays/go-jk/opengl"
)

type BMScene struct {
	bmName        string
	shaderProgram *opengl.ShaderProgram
	renderers     []opengl.Renderer
	cam           *camera.Camera
	window        *glfw.Window
}

func NewBMScene(bmName string, window *glfw.Window, cam *camera.Camera, shaderProgram *opengl.ShaderProgram) *BMScene {
	return &BMScene{bmName: bmName, window: window, cam: cam, shaderProgram: shaderProgram}
}

func (s *BMScene) Load() {
	bm := jk.GetLoader().LoadBM(s.bmName)
	sftRenderer := opengl.NewOpenGlBmRenderer(&bm, s.shaderProgram)
	s.renderers = append(s.renderers, sftRenderer)
}

func (s *BMScene) Unload() {

}

func (s *BMScene) Update() {
	opengl.Draw(s.window, s.cam, s.renderers)
}
