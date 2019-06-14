package scene

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
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
	bmRenderer    opengl.Renderer
	bm            *jk.BMFile
}

func NewBMScene(bmName string, window *glfw.Window, cam *camera.Camera, shaderProgram *opengl.ShaderProgram) *BMScene {
	return &BMScene{bmName: bmName, window: window, cam: cam, shaderProgram: shaderProgram}
}

func (s *BMScene) Load() {
	bm := jk.GetLoader().LoadBM(s.bmName)
	s.bm = &bm
}

func (s *BMScene) Unload() {

}

func (s *BMScene) Update() {
	if s.bm != nil && s.bmRenderer == nil {
		w, h := s.window.GetSize()
		windowAspect := float32(w) / float32(h)

		var scale mgl32.Vec2
		imageRatio := float32(s.bm.Images[0].SizeX) / float32(s.bm.Images[0].SizeY) / windowAspect
		scale = mgl32.Vec2{imageRatio, 1}
		if imageRatio > 1 {
			imageRatio = float32(s.bm.Images[0].SizeY) / float32(s.bm.Images[0].SizeX) / windowAspect
			scale = mgl32.Vec2{1, imageRatio}
		}
		//fmt.Printf("%d, %d, %+v\n", s.bm.Images[0].SizeX, s.bm.Images[0].SizeY, scale)

		s.bmRenderer = opengl.NewOpenGlBmRenderer(s.bm, scale, s.shaderProgram)
		s.renderers = append(s.renderers, s.bmRenderer)
	}

	if len(s.renderers) > 0 {
		opengl.Draw(s.window, s.cam, s.renderers)
	}
}
