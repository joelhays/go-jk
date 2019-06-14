package scene

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/joelhays/go-jk/camera"
	"github.com/joelhays/go-jk/jk"
	"github.com/joelhays/go-jk/opengl"
)

type Jk3doScene struct {
	jk3doName     string
	shaderProgram *opengl.ShaderProgram
	renderers     []opengl.Renderer
	cam           *camera.Camera
	window        *glfw.Window
	obj           *jk.Jk3doFile
	objRenderer   opengl.Renderer
}

func NewJk3doScene(jk3doName string, window *glfw.Window, cam *camera.Camera, shaderProgram *opengl.ShaderProgram) *Jk3doScene {
	return &Jk3doScene{jk3doName: jk3doName, window: window, cam: cam, shaderProgram: shaderProgram}
}

func (s *Jk3doScene) Load() {
	obj := jk.GetLoader().Load3DO(s.jk3doName)
	s.obj = &obj
	s.cam.Position = mgl32.Vec3{0, 1, 0}
	s.cam.Up = mgl32.Vec3{0, 0, 1}
	s.cam.Yaw = 90
	s.cam.Pitch = 0
	s.cam.UpdateCameraVectors()
}

func (s *Jk3doScene) Unload() {

}

func (s *Jk3doScene) Update() {
	if s.obj != nil && s.objRenderer == nil {
		s.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

		s.objRenderer = opengl.NewOpenGl3doRenderer(&jk.Thing{Position: mgl32.Vec3{float32(0), float32(0), float32(0)}, Yaw: 0, Pitch: 0, Roll: 0}, nil, s.obj, s.shaderProgram)
		s.renderers = append(s.renderers, s.objRenderer)
	}

	if len(s.renderers) > 0 {
		opengl.Draw(s.window, s.cam, s.renderers)
	}
}
