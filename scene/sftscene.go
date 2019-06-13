package scene

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/joelhays/go-jk/camera"
	"github.com/joelhays/go-jk/jk"
	"github.com/joelhays/go-jk/opengl"
)

type SFTScene struct {
	sftName       string
	shaderProgram *opengl.ShaderProgram
	renderers     []opengl.Renderer
	cam           *camera.Camera
	window        *glfw.Window
}

func NewSFTScene(sftName string, window *glfw.Window, cam *camera.Camera, shaderProgram *opengl.ShaderProgram) *SFTScene {
	return &SFTScene{sftName: sftName, window: window, cam: cam, shaderProgram: shaderProgram}
}

func (s *SFTScene) Load() {
	sft := jk.GetLoader().LoadSFT(s.sftName)
	fmt.Printf("%+v\n", sft)
	sftRenderer := opengl.NewOpenGlBmRenderer(&sft.BMFile, s.shaderProgram)
	s.renderers = append(s.renderers, sftRenderer)
}

func (s *SFTScene) Unload() {

}

func (s *SFTScene) Update() {
	opengl.Draw(s.window, s.cam, s.renderers)
}
