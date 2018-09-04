package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	keys      = make(map[glfw.Key]bool)
	modifiers = make(map[glfw.ModifierKey]bool)
	lastX     float64
	lastY     float64
)

func KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}

	if action == glfw.Press {
		keys[key] = true
		modifiers[mods] = true
	} else if action == glfw.Release {
		delete(keys, key)
		delete(modifiers, mods)
	}
}

func MouseCallback(window *glfw.Window, xpos float64, ypos float64) {
	xOffset := xpos - lastX
	yOffset := lastY - ypos
	lastX = xpos
	lastY = ypos

	camera.ProcessMouseMovement(xOffset, yOffset, true)
}
