package jk

import (
	"github.com/go-gl/mathgl/mgl32"
)

type JklParser interface {
	ParseJKLFromFile(filePath string) Jkl
	ParseJKLFromString(jklString string) Jkl
}

// Jkl contains the information extracted from the Jedi Knight Level (.jkl) file
type Jkl struct {
	Model          *JkMesh
	Jk3dos         map[string]Jk3doFile
	Jk3doTemplates map[string]Template
	Things         []Thing
}

type surface struct {
	VertexIds        []int64
	TextureVertexIds []int64
	LightIntensities []float64
	Normal           mgl32.Vec3
	Geo              int64
	MaterialID       int64
}

type Template struct {
	Name      string
	Jk3doName string
	Size      float64
}

type Thing struct {
	TemplateName string
	Position     mgl32.Vec3
	Pitch        float64
	Yaw          float64
	Roll         float64
}

type JkMesh struct {
	Name            string
	Vertices        []mgl32.Vec3
	TextureVertices []mgl32.Vec2
	VertexNormals   []mgl32.Vec3
	Surfaces        []surface
	Materials       []Material
	ColorMaps       []ColorMap
}
