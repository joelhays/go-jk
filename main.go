package main

import (
	"github.com/joelhays/go-jk/camera"
	"github.com/joelhays/go-jk/opengl"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	// "github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/joelhays/go-jk/jk"
)

var (
	mpLevels  = []string{"m10.jkl", "m2.jkl", "m4.jkl", "m5.jkl", "m_boss15.jkl", "m_boss17.jkl"}
	ctfLevels = []string{"c1.jkl", "c2.jkl", "c3.jkl"}
	spLevels  = []string{"01narshadda.jkl", "02narshadda.jkl", "03katarn.jkl", "04escapehouse.jkl", "06abarons.jkl",
		"06bbarons.jkl", "07yun.jkl", "08escape88.jkl", "09fuelstation.jkl", "10cargo.jkl", "11gorc.jkl", "12escape.jkl",
		"14tower.jkl", "15maw.jkl", "16aescapeship.jkl", "16bescapeship.jkl", "17asarris.jkl", "17bsarris.jkl",
		"18ascend.jkl", "19a.jkl", "19b.jkl", "20aboc.jkl", "20bboc.jkl", "21ajarec.jkl", "21bjarec.jkl"}
)

var (
	cam          camera.Camera
	previousTime float64
	cpuprofile   = "go-jk.prof"
)

func main() {
	f, err := os.Create(cpuprofile)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	runtime.LockOSThread()

	window := opengl.InitGlfw(1024, 768, KeyCallback, MouseCallback)
	defer glfw.Terminate()

	opengl.InitOpenGL()

	shaderProgram := opengl.NewShaderProgram("./shaders/vertex.glsl", "./shaders/fragment.glsl")
	defer shaderProgram.Cleanup()

	guiShaderProgram := opengl.NewShaderProgram("./shaders/gui_vertex.glsl", "./shaders/gui_fragment.glsl")
	defer guiShaderProgram.Cleanup()

	cam = camera.NewCamera(mgl32.Vec3{0, 0, 1}, mgl32.Vec3{0, 0, 1}, 0, -90)
	cam.MovementSpeed = 2

	var renderers []opengl.Renderer

	//debug_createBmRenderer(&renderers, guiShaderProgram)
	//debug_create3doRenderer(&renderers, shaderProgram)

	createJklRenderer(&renderers, shaderProgram)

	for !window.ShouldClose() {
		deltaTime := glfw.GetTime() - previousTime
		previousTime = glfw.GetTime()

		doMovement(deltaTime)

		opengl.Draw(window, &cam, renderers)
	}
}

func createJklRenderer(renderers *[]opengl.Renderer, shaderProgram *opengl.ShaderProgram) {
	jklLevel := jk.GetLoader().LoadJKL(spLevels[0])
	level := opengl.NewOpenGlLevelRenderer(nil, nil, jklLevel.Model, shaderProgram)
	*renderers = append(*renderers, level)

	var foundPlayer bool
	for i := 0; i < len(jklLevel.Things); i++ {
		thing := jklLevel.Things[i]
		if thing.TemplateName == "walkplayer" {
			if !foundPlayer {
				cam.Position = thing.Position
				foundPlayer = true
			}
			continue
		}

		template := jklLevel.Jk3doTemplates[thing.TemplateName]
		jk3do := jklLevel.Jk3dos[template.Jk3doName]

		if len(jk3do.GeoSets) > 0 {
			objRenderer := opengl.NewOpenGl3doRenderer(&thing, &template, &jk3do, shaderProgram)
			*renderers = append(*renderers, objRenderer)
		}
	}
}

func debug_createBmRenderer(renderers *[]opengl.Renderer, shaderProgram *opengl.ShaderProgram) {
	bmFile := jk.GetLoader().LoadBM("bkmain.bm")
	bmRenderer := opengl.NewOpenGlBmRenderer(&bmFile, shaderProgram)
	*renderers = append(*renderers, bmRenderer)
}
func debug_create3doRenderer(renderers *[]opengl.Renderer, shaderProgram *opengl.ShaderProgram) {
	/* RENDER 3DO AT ORIGIN */
	obj := jk.GetLoader().Load3DO("rystr.3do")
	objRenderer := opengl.NewOpenGl3doRenderer(&jk.Thing{Position: mgl32.Vec3{float32(0), float32(0), float32(0)}, Yaw: 0, Pitch: 0, Roll: 0}, nil, &obj, shaderProgram)
	*renderers = append(*renderers, objRenderer)
}
