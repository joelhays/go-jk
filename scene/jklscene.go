package scene

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/joelhays/go-jk/camera"
	"github.com/joelhays/go-jk/jk"
	"github.com/joelhays/go-jk/opengl"
)

type JklScene struct {
	jklName       string
	shaderProgram *opengl.ShaderProgram
	renderers     []opengl.Renderer
	cam           *camera.Camera
	window        *glfw.Window
}

func NewJklScene(jklName string, window *glfw.Window, cam *camera.Camera, shaderProgram *opengl.ShaderProgram) *JklScene {
	return &JklScene{jklName: jklName, window: window, cam: cam, shaderProgram: shaderProgram}
}

func (s *JklScene) Load() {
	jklLevel := jk.GetLoader().LoadJKL(s.jklName)
	level := opengl.NewOpenGlLevelRenderer(nil, nil, jklLevel.Model, s.shaderProgram)
	s.renderers = append(s.renderers, level)

	var foundPlayer bool
	for i := 0; i < len(jklLevel.Things); i++ {
		thing := jklLevel.Things[i]
		if thing.TemplateName == "walkplayer" {
			if !foundPlayer {
				s.cam.Position = thing.Position
				foundPlayer = true
			}
			continue
		}

		template := jklLevel.Jk3doTemplates[thing.TemplateName]
		jk3do := jklLevel.Jk3dos[template.Jk3doName]

		if len(jk3do.GeoSets) > 0 {
			objRenderer := opengl.NewOpenGl3doRenderer(&thing, &template, &jk3do, s.shaderProgram)
			s.renderers = append(s.renderers, objRenderer)
		}
	}
}

func (s *JklScene) Unload() {

}

func (s *JklScene) Update() {
	opengl.Draw(s.window, s.cam, s.renderers)
}
