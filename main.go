package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/joelhays/go-jk/camera"
	"github.com/joelhays/go-jk/opengl"
	"github.com/joelhays/go-jk/scene"

	// "github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
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
	defer runtime.UnlockOSThread()

	sceneManager := scene.NewSceneManager()
	defer sceneManager.Unload()
	inputManager := NewInputManager(sceneManager)

	window := opengl.InitGlfw(1024, 768, inputManager.KeyCallback, inputManager.MouseCallback)
	defer glfw.Terminate()

	opengl.InitOpenGL()

	shaderProgram := opengl.NewShaderProgram("./shaders/vertex.glsl", "./shaders/fragment.glsl")
	defer shaderProgram.Cleanup()

	guiShaderProgram := opengl.NewShaderProgram("./shaders/gui_vertex.glsl", "./shaders/gui_fragment.glsl")
	defer guiShaderProgram.Cleanup()

	cam = camera.NewCamera(mgl32.Vec3{0, 0, 1}, mgl32.Vec3{0, 0, 1}, 0, -90)
	cam.MovementSpeed = 2

	for _, level := range spLevels {
		sceneManager.Add(level, scene.NewJklScene(level, window, &cam, shaderProgram))
	}
	for _, level := range mpLevels {
		sceneManager.Add(level, scene.NewJklScene(level, window, &cam, shaderProgram))
	}
	for _, level := range ctfLevels {
		sceneManager.Add(level, scene.NewJklScene(level, window, &cam, shaderProgram))
	}
	sceneManager.Add("spLevel", scene.NewJklScene(spLevels[0], window, &cam, shaderProgram))
	sceneManager.Add("mpLevel", scene.NewJklScene(mpLevels[0], window, &cam, shaderProgram))
	sceneManager.Add("ctfLevel", scene.NewJklScene(ctfLevels[0], window, &cam, shaderProgram))
	sceneManager.Add("3do", scene.NewJk3doScene("rystr.3do", window, &cam, shaderProgram))
	sceneManager.Add("menu", scene.NewMainMenuScene(window, sceneManager))
	sceneManager.Add("sft", scene.NewSFTSceneScene("creditlarge.sft", window, &cam, guiShaderProgram))
	sceneManager.LoadScene("menu")

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		deltaTime := glfw.GetTime() - previousTime
		previousTime = glfw.GetTime()

		doMovement(deltaTime)

		sceneManager.Update()

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
