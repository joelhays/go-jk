package main

import (
	"runtime"

	// "github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/joelhays/go-vulkan/jk"
)

type levelFile struct {
	Gob  string
	Name string
}

var mpLevels = []string{"jkl\\m10.jkl", "jkl\\m2.jkl", "jkl\\m4.jkl", "jkl\\m5.jkl", "jkl\\m_boss15.jkl", "jkl\\m_boss17.jkl"}
var ctfLevels = []string{"jkl\\c1.jkl", "jkl\\c2.jkl", "jkl\\c3.jkl"}
var spLevels = []string{
	"jkl\\01narshadda.jkl", "jkl\\02narshadda.jkl", "jkl\\03katarn.jkl", "jkl\\04escapehouse.jkl", "jkl\\06abarons.jkl",
	"jkl\\06bbarons.jkl", "jkl\\07yun.jkl", "jkl\\08escape88.jkl", "jkl\\09fuelstation.jkl", "jkl\\10cargo.jkl", "jkl\\11gorc.jkl",
	"jkl\\12escape.jkl", "jkl\\14tower.jkl", "jkl\\15maw.jkl", "jkl\\16aescapeship.jkl", "jkl\\16bescapeship.jkl", "jkl\\17asarris.jkl",
	"jkl\\17bsarris.jkl", "jkl\\18ascend.jkl", "jkl\\19a.jkl", "jkl\\19b.jkl", "jkl\\20aboc.jkl", "jkl\\20bboc.jkl", "jkl\\21ajarec.jkl",
	"jkl\\21bjarec.jkl",
}

const (
	width  = 1024
	height = 768
)

var camera Camera
var previousTime float64

var lightPos = mgl32.Vec3{1.2, 1.0, 2.0}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	camera = NewCamera(mgl32.Vec3{0, 0, 1}, mgl32.Vec3{0, 0, 1}, 0, -90)

	// matBytes := jk.LoadFileFromGOB("J:\\Resource\\Res2.gob", "00cglyf3.mat")
	// data := jk.ParseMatFile(matBytes)
	// _ = data
	// return

	// cmpBytes := jk.LoadFileFromGOB("J:\\Resource\\Res2.gob", "dflt.cmp")
	// data := jk.ParseCmpFile(cmpBytes)
	// _ = data
	// return

	// jk3doBytes := jk.LoadFileFromGOB("J:\\Resource\\Res2.gob", "00crte6x6.3do") // "landpad.3do")
	// jklData := jk.Parse3doFromString(string(jk3doBytes))

	// fmt.Println(jklData.Surfaces)
	// return

	// vao := makeVao(triangle)
	// vao := makeVao(cube)

	// jklBytes := jk.LoadFileFromGOB("J:\\Episode\\JK1CTF.GOB", ctfLevels[2])
	// jklBytes := jk.LoadFileFromGOB("J:\\Episode\\JK1.GOB", spLevels[0])
	jklBytes := jk.LoadFileFromGOB("J:\\Episode\\JK1MP.GOB", mpLevels[4])
	jklLevel := jk.ReadJKLFromString(string(jklBytes))
	level := NewOpenGlLevelRenderer(nil, nil, jklLevel.Model, program)

	models := make([]*OpenGl3doRenderer, len(jklLevel.Things))

	for i := 0; i < len(jklLevel.Things); i++ {
		thing := jklLevel.Things[i]
		if thing.TemplateName == "walkplayer" {
			models[i] = nil
			continue
		}

		template := jklLevel.Jk3doTemplates[thing.TemplateName]
		jk3do := jklLevel.Jk3dos[template.Jk3doName]
		jk3do.ColorMap = jklLevel.Model.ColorMaps[0]

		if len(jk3do.GeoSets) > 0 {
			models[i] = NewOpenGl3doRenderer(&thing, &template, &jk3do, program)
		} else {
			models[i] = nil
		}
	}

	// jk3doBytes := jk.LoadFileFromGOB("J:\\Resource\\Res2.gob", "landpad.3do")
	// jklModel := jk.Parse3doFile(string(jk3doBytes))
	// jklModel.ColorMap = jklLevel.Model.Mesh.ColorMaps[0]
	// models := make([]*OpenGl3doRenderer, 1)
	// thing := &jk.Thing{Position: mgl32.Vec3{float32(0), float32(0), float32(0)}}
	// models[0] = NewOpenGl3doRenderer(thing, nil, &jklModel, program)

	for !window.ShouldClose() {
		drawRenderer(window, level, models)
	}
}

func doMovement(deltaTime float64) {

	if keyMinus := keys[glfw.KeyKPSubtract]; keyMinus {
		camera.MovementSpeed = .75
	}

	if keyDecimal := keys[glfw.KeyKPDecimal]; keyDecimal {
		camera.MovementSpeed = 6
	}

	if keyPlus := keys[glfw.KeyKPAdd]; keyPlus {
		camera.MovementSpeed = 12
	}

	if keyW, keyUp := keys[glfw.KeyW], keys[glfw.KeyUp]; keyW || keyUp {
		camera.ProcessKeyboard(CAMERA_FORWARD, deltaTime)
	}

	if keyS, keyDown := keys[glfw.KeyS], keys[glfw.KeyDown]; keyS || keyDown {
		camera.ProcessKeyboard(CAMERA_BACKWARD, deltaTime)
	}

	if keyA, keyLeft := keys[glfw.KeyA], keys[glfw.KeyLeft]; keyA || keyLeft {
		camera.ProcessKeyboard(CAMERA_LEFT, deltaTime)
	}

	if keyD, keyRight := keys[glfw.KeyD], keys[glfw.KeyRight]; keyD || keyRight {
		camera.ProcessKeyboard(CAMERA_RIGHT, deltaTime)
	}
}
