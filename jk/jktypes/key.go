package jktypes

import "github.com/go-gl/mathgl/mgl32"

const (
	KEY_FLAG_LOOPING                        = 0x00
	KEY_FLAT_STOP_AT_LAST_FRAME_UNTIL_EVENT = 0x14
	KEY_FLAG_STOP_AFTER_LAST_FRAME          = 0x2c
)

type Key struct {
	Header        KeyHeader
	KeyframeNodes []KeyframeNode
}

type KeyHeader struct {
	Flags  byte
	Type   int32
	Frames int32
	FPS    float32
	Joints int32
}

type KeyframeNode struct {
	MeshName string
	Entries  []KeyframeNodeEntry
}

type KeyframeNodeEntry struct {
	Frame            int32
	Flags            byte
	Offset           mgl32.Vec3
	Orientation      mgl32.Vec3
	DeltaOffset      mgl32.Vec3
	DeltaOrientation mgl32.Vec3
}
