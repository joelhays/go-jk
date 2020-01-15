package scene

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/joelhays/go-jk/camera"
	"github.com/joelhays/go-jk/jk"
	"github.com/joelhays/go-jk/jk/jkparsers"
	"github.com/joelhays/go-jk/jk/jktypes"
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
	var sft jktypes.SFTFile
	fileBytes := jk.GetLoader().LoadResource(s.sftName)
	if fileBytes != nil {
		sft = jkparsers.NewSftParser().ParseFromBytes(fileBytes)
	}
	//fmt.Printf("%+v\n", sft)

	w, h := s.window.GetSize()
	windowAspect := float32(w) / float32(h)

	var scale mgl32.Vec2
	imageRatio := float32(sft.BMFile.Images[0].SizeX) / float32(sft.BMFile.Images[0].SizeY) / windowAspect
	scale = mgl32.Vec2{imageRatio, 1}
	if imageRatio > 1 {
		imageRatio = float32(sft.BMFile.Images[0].SizeY) / float32(sft.BMFile.Images[0].SizeX) / windowAspect
		scale = mgl32.Vec2{1, imageRatio}
	}
	//fmt.Printf("%d, %d, %+v\n", sft.BMFile.Images[0].SizeX, sft.BMFile.Images[0].SizeY, scale)

	sftRenderer := opengl.NewOpenGlBmRenderer(&sft.BMFile, scale, s.shaderProgram)
	s.renderers = append(s.renderers, sftRenderer)
}

func (s *SFTScene) Unload() {

}

func (s *SFTScene) Update() {
	opengl.Draw(s.window, s.cam, s.renderers)
}
