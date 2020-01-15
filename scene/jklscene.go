package scene

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/joelhays/go-jk/camera"
	"github.com/joelhays/go-jk/jk"
	"github.com/joelhays/go-jk/jk/jktypes"
	"github.com/joelhays/go-jk/opengl"
)

type JklScene struct {
	jklName       string
	shaderProgram *opengl.ShaderProgram
	renderers     []opengl.Renderer
	cam           *camera.Camera
	window        *glfw.Window
	levelRenderer opengl.Renderer
	level         *jktypes.Jkl
}

func NewJklScene(jklName string, window *glfw.Window, cam *camera.Camera, shaderProgram *opengl.ShaderProgram) *JklScene {
	return &JklScene{jklName: jklName, window: window, cam: cam, shaderProgram: shaderProgram}
}

func (s *JklScene) Load() {
	if s.level == nil {
		level := jk.GetLoader().LoadJKL(s.jklName)
		s.level = &level
	}
}

func (s *JklScene) Unload() {
	s.renderers = make([]opengl.Renderer, 0)
	s.levelRenderer = nil
	s.level = nil
}

func (s *JklScene) Update() {
	if s.level != nil && s.levelRenderer == nil {
		s.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

		s.levelRenderer = opengl.NewOpenGlLevelRenderer(nil, nil, s.level.Model, s.shaderProgram)
		s.renderers = append(s.renderers, s.levelRenderer)

		var foundPlayer bool
		for i := 0; i < len(s.level.Things); i++ {
			thing := s.level.Things[i]
			if thing.TemplateName == "walkplayer" {
				if !foundPlayer {
					s.cam.Position = thing.Position
					foundPlayer = true
				}
				continue
			}

			template := s.level.Jk3doTemplates[thing.TemplateName]
			jk3do := s.level.Jk3dos[template.Jk3doName]

			if len(jk3do.GeoSets) > 0 {
				objRenderer := opengl.NewOpenGl3doRenderer(&thing, &template, &jk3do, s.shaderProgram)
				s.renderers = append(s.renderers, objRenderer)
			}
		}
	}

	if len(s.renderers) > 0 {
		opengl.Draw(s.window, s.cam, s.renderers)
	}
}
