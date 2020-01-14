package jk

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Jk3doParser interface {
	Parse3doFromFile(data string) Jk3doFile
	Parse3doFromString(data string) Jk3doFile
}

type Jk3doFile struct {
	Materials []Material
	GeoSets   []GeoSet
	Hierarchy []HierarchyDef
	ColorMap  ColorMap
}

type GeoSet struct {
	Meshes []Mesh
}

type Mesh struct {
	GeometryMode    int64
	Vertices        []mgl32.Vec3
	TextureVertices []mgl32.Vec2
	VertexNormals   []mgl32.Vec3
	Faces           []Face
	FaceNormals     []mgl32.Vec3
}

type Face struct {
	VertexIds        []int64
	TextureVertexIds []int64
	LightIntensities []float64
	GeometryMode     int64
	MaterialID       int64
}

type HierarchyDef struct {
	MeshID      int64
	ParentID    int64
	ChildID     int64
	SiblingID   int64
	NumChildren int64
	Position    mgl32.Vec3
	Pitch       float64
	Yaw         float64
	Roll        float64
	Pivot       mgl32.Vec3
	NodeName    string
}
