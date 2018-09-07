package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type CameraDirection int

const (
	CAMERA_FORWARD  CameraDirection = 0
	CAMERA_BACKWARD CameraDirection = 1
	CAMERA_LEFT     CameraDirection = 2
	CAMERA_RIGHT    CameraDirection = 3
)

type Camera struct {
	Position mgl32.Vec3
	Front    mgl32.Vec3
	Up       mgl32.Vec3
	Right    mgl32.Vec3
	WorldUp  mgl32.Vec3

	Yaw   float64
	Pitch float64

	MovementSpeed    float64
	MouseSensitivity float64
	Zoom             float64
}

func NewCamera(position mgl32.Vec3, up mgl32.Vec3, yaw float64, pitch float64) Camera {
	c := Camera{
		Position:         position,
		Front:            mgl32.Vec3{0, 0, -1},
		Up:               up,
		WorldUp:          up,
		Yaw:              yaw,
		Pitch:            pitch,
		MovementSpeed:    6,
		MouseSensitivity: .25,
		Zoom:             45,
	}

	c.updateCameraVectors()

	return c
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}

func (c *Camera) ProcessKeyboard(direction CameraDirection, deltaTime float64) {
	velocity := c.MovementSpeed * deltaTime
	switch direction {
	case CAMERA_FORWARD:
		c.Position = c.Position.Add(c.Front.Mul(float32(velocity)))
	case CAMERA_BACKWARD:
		c.Position = c.Position.Sub(c.Front.Mul(float32(velocity)))
	case CAMERA_LEFT:
		c.Position = c.Position.Sub(c.Right.Mul(float32(velocity)))
	case CAMERA_RIGHT:
		c.Position = c.Position.Add(c.Right.Mul(float32(velocity)))
	}
}

func (c *Camera) ProcessMouseMovement(xOffset float64, yOffset float64, constrainPitch bool) {
	xOffset *= c.MouseSensitivity
	yOffset *= c.MouseSensitivity

	c.Yaw += xOffset
	c.Pitch += yOffset

	if constrainPitch {
		c.Pitch = float64(mgl32.Clamp(float32(c.Pitch), -89.0, 89.0))
	}

	c.updateCameraVectors()
}

func (c *Camera) updateCameraVectors() {
	x := math.Cos(float64(mgl32.DegToRad(float32(c.Yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(c.Pitch))))
	// y := math.Sin(float64(mgl32.DegToRad(float32(c.Pitch))))
	// z := math.Sin(float64(mgl32.DegToRad(float32(c.Yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(c.Pitch))))
	z := math.Sin(float64(mgl32.DegToRad(float32(c.Pitch))))
	y := -math.Sin(float64(mgl32.DegToRad(float32(c.Yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(c.Pitch))))

	c.Front = mgl32.Vec3{float32(x), float32(y), float32(z)}.Normalize()
	c.Right = c.Front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}
