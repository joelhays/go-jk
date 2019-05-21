package main

import (
	"runtime"

	// "github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/joelhays/go-jk/jk"
)

type levelFile struct {
	Gob  string
	Name string
}

var mpLevels = []string{"m10.jkl", "m2.jkl", "m4.jkl", "m5.jkl", "m_boss15.jkl", "m_boss17.jkl"}
var ctfLevels = []string{"c1.jkl", "c2.jkl", "c3.jkl"}
var spLevels = []string{"01narshadda.jkl", "02narshadda.jkl", "03katarn.jkl", "04escapehouse.jkl", "06abarons.jkl",
	"06bbarons.jkl", "07yun.jkl", "08escape88.jkl", "09fuelstation.jkl", "10cargo.jkl", "11gorc.jkl", "12escape.jkl",
	"14tower.jkl", "15maw.jkl", "16aescapeship.jkl", "16bescapeship.jkl", "17asarris.jkl", "17bsarris.jkl",
	"18ascend.jkl", "19a.jkl", "19b.jkl", "20aboc.jkl", "20bboc.jkl", "21ajarec.jkl", "21bjarec.jkl",
}

const (
	width  = 1024
	height = 768
)

var camera Camera
var previousTime float64
var foundPlayer bool

var lightPos = mgl32.Vec3{1.2, 1.0, 2.0}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	camera = NewCamera(mgl32.Vec3{0, 0, 1}, mgl32.Vec3{0, 0, 1}, 0, -90)
	camera.MovementSpeed = 2

	bmBytes := jk.LoadFileFromGOB("J:\\Resource\\Res1hi.gob", "bkmain.bm")
	bmFile := jk.ParseBmFile(bmBytes)

	bmRenderer := NewOpenGlBmRenderer(&bmFile, program)

	//for !window.ShouldClose() {
	//	drawRenderer(window, nil, nil, bmRenderer)
	//}
	//
	//return

	//jklBytes := jk.LoadFileFromGOB("J:\\Episode\\JK1CTF.GOB", ctfLevels[2])
	//jklBytes := jk.LoadFileFromGOB("J:\\Episode\\JK1MP.GOB", mpLevels[4])
	jklBytes := jk.LoadFileFromGOB("J:\\Episode\\JK1.GOB", spLevels[0])
	jklLevel := jk.ReadJKLFromString(string(jklBytes))
	level := NewOpenGlLevelRenderer(nil, nil, jklLevel.Model, program)

	models := make([]*OpenGl3doRenderer, len(jklLevel.Things))

	for i := 0; i < len(jklLevel.Things); i++ {
		thing := jklLevel.Things[i]
		if thing.TemplateName == "walkplayer" {
			if !foundPlayer {
				camera.Position = thing.Position
				foundPlayer = true
			}
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

	/* RENDER 3DO AT ORIGIN */
	//jk3doBytes := jk.LoadFileFromGOB("J:\\Resource\\Res2.gob", "rystr.3do")
	//jklModel := jk.Parse3doFile(string(jk3doBytes))
	//jklModel.ColorMap = jklLevel.Model.ColorMaps[0]
	//thing := &jk.Thing{Position: mgl32.Vec3{float32(0), float32(0), float32(0)}, Yaw: 45, Pitch: 45, Roll: 45}
	//models = append(models, NewOpenGl3doRenderer(thing, nil, &jklModel, program))

	for !window.ShouldClose() {
		drawRenderer(window, level, models, bmRenderer)
	}
}
