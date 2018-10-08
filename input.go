package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
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

	camera.ProcessMouseMovement(xOffset, yOffset, true)
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

	if key := keys[glfw.KeyKP1]; key {
		camera.MovementSpeed = 1
	}
	if key := keys[glfw.KeyKP2]; key {
		camera.MovementSpeed = 2
	}
	if key := keys[glfw.KeyKP3]; key {
		camera.MovementSpeed = 3
	}
	if key := keys[glfw.KeyKP4]; key {
		camera.MovementSpeed = 4
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
