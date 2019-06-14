package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/joelhays/go-jk/camera"
	"github.com/joelhays/go-jk/jk"
	"github.com/joelhays/go-jk/opengl"
	"github.com/joelhays/go-jk/scene"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
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

	for _, gobFileName := range jk.GetLoader().LoadJKLManifest() {
		sceneManager.Add(gobFileName, scene.NewJklScene(gobFileName, window, &cam, shaderProgram))
	}
	for _, gobFileName := range jk.GetLoader().Load3DOManifest() {
		sceneManager.Add(gobFileName, scene.NewJk3doScene(gobFileName, window, &cam, shaderProgram))
	}
	for _, gobFileName := range jk.GetLoader().LoadBMManifest() {
		sceneManager.Add(gobFileName, scene.NewBMScene(gobFileName, window, &cam, guiShaderProgram))
	}
	sceneManager.Add("3do", scene.NewJk3doScene("rystr.3do", window, &cam, shaderProgram))
	sceneManager.Add("menu", scene.NewMainMenuScene(window, sceneManager))
	sceneManager.Add("sft", scene.NewSFTScene("large0.sft", window, &cam, guiShaderProgram))
	sceneManager.Add("bm", scene.NewBMScene("bkdialog.bm", window, &cam, guiShaderProgram))
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
