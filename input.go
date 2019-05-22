package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/joelhays/go-jk/camera"
)

var (
	keys  = make(map[glfw.Key]bool)
	lastX float64
	lastY float64
)

func KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}

	if action == glfw.Press {
		keys[key] = true
	} else if action == glfw.Release {
		delete(keys, key)
	}
}

func MouseCallback(window *glfw.Window, xpos float64, ypos float64) {
	xOffset := xpos - lastX
	yOffset := lastY - ypos
	lastX = xpos
	lastY = ypos

	cam.ProcessMouseMovement(xOffset, yOffset, true)
}

func doMovement(deltaTime float64) {

	if keyMinus := keys[glfw.KeyKPSubtract]; keyMinus {
		cam.MovementSpeed = .75
	}

	if keyDecimal := keys[glfw.KeyKPDecimal]; keyDecimal {
		cam.MovementSpeed = 6
	}

	if keyPlus := keys[glfw.KeyKPAdd]; keyPlus {
		cam.MovementSpeed = 12
	}

	if key := keys[glfw.KeyKP1]; key {
		cam.MovementSpeed = 1
	}
	if key := keys[glfw.KeyKP2]; key {
		cam.MovementSpeed = 2
	}
	if key := keys[glfw.KeyKP3]; key {
		cam.MovementSpeed = 3
	}
	if key := keys[glfw.KeyKP4]; key {
		cam.MovementSpeed = 4
	}

	if keyW, keyUp := keys[glfw.KeyW], keys[glfw.KeyUp]; keyW || keyUp {
		cam.ProcessKeyboard(camera.CAMERA_FORWARD, deltaTime)
	}

	if keyS, keyDown := keys[glfw.KeyS], keys[glfw.KeyDown]; keyS || keyDown {
		cam.ProcessKeyboard(camera.CAMERA_BACKWARD, deltaTime)
	}

	if keyA, keyLeft := keys[glfw.KeyA], keys[glfw.KeyLeft]; keyA || keyLeft {
		cam.ProcessKeyboard(camera.CAMERA_LEFT, deltaTime)
	}

	if keyD, keyRight := keys[glfw.KeyD], keys[glfw.KeyRight]; keyD || keyRight {
		cam.ProcessKeyboard(camera.CAMERA_RIGHT, deltaTime)
	}
}
