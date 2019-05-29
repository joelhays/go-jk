package opengl

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/joelhays/go-jk/jk"
)

type OpenGl3doRenderer struct {
	thing    *jk.Thing
	template *jk.Template
	object   *jk.Jk3doFile
	Program  uint32
	vao      uint32
	textures []uint32
	lod      int32
}

func NewOpenGl3doRenderer(thing *jk.Thing, template *jk.Template, object *jk.Jk3doFile, program uint32) *OpenGl3doRenderer {
	if thing == nil {
		panic("Thing is nil!")
	}
	r := &OpenGl3doRenderer{thing: thing, template: template, object: object, Program: program}

	r.lod = 0

	numGeoSets := len(r.object.GeoSets)

	if r.lod >= int32(numGeoSets) {
		r.lod = int32(numGeoSets - 1)
	}

	r.setupMesh()
	return r
}

func (r *OpenGl3doRenderer) Render() {

	gl.BindVertexArray(r.vao)

	var offset int32
	// render the main mesh if it has vertices
	// render all child meshes with parent transform

	modelUniform := gl.GetUniformLocation(r.Program, gl.Str("model\x00"))
	textureUniform := gl.GetUniformLocation(r.Program, gl.Str("objectTexture\x00"))

	for meshIdx, mesh := range r.object.GeoSets[r.lod].Meshes {

		if len(mesh.Vertices) == 0 {
			continue
		}

		_ = meshIdx
		model := mgl32.Ident4()

		var hierarchy jk.HierarchyDef
		for i := 0; i < len(r.object.Hierarchy); i++ {
			h := r.object.Hierarchy[i]
			if h.MeshID == int64(meshIdx) {
				hierarchy = h
				break
			}
		}

		meshRotateValues := mgl32.Vec3{float32(hierarchy.Pitch), float32(hierarchy.Roll), float32(hierarchy.Yaw)}
		meshTranslateValues := hierarchy.Position

		parentID := hierarchy.ParentID
		for parentID != -1 {
			parent := r.object.Hierarchy[parentID]

			parentRotateValues := mgl32.Vec3{float32(parent.Pitch), float32(parent.Roll), float32(parent.Yaw)}
			meshRotateValues = meshRotateValues.Add(parentRotateValues)

			meshTranslateValues = meshTranslateValues.Add(parent.Position)

			parentID = parent.ParentID
		}

		meshRotateX := mgl32.HomogRotate3DX(mgl32.DegToRad(meshRotateValues.X()))
		meshRotateY := mgl32.HomogRotate3DY(mgl32.DegToRad(meshRotateValues.Y()))
		meshRotateZ := mgl32.HomogRotate3DZ(mgl32.DegToRad(meshRotateValues.Z()))
		meshRotation := meshRotateX.Mul4(meshRotateY.Mul4(meshRotateZ))
		meshTranslation := mgl32.Translate3D(meshTranslateValues.X(), meshTranslateValues.Y(), meshTranslateValues.Z())
		meshPivot := mgl32.Translate3D(hierarchy.Pivot.X(), hierarchy.Pivot.Y(), hierarchy.Pivot.Z())

		thingTranslate := mgl32.Translate3D(r.thing.Position.X(), r.thing.Position.Y(), r.thing.Position.Z())
		thingRotateX := mgl32.HomogRotate3DX(mgl32.DegToRad(float32(r.thing.Pitch)))
		thingRotateY := mgl32.HomogRotate3DY(mgl32.DegToRad(float32(r.thing.Roll)))
		thingRotateZ := mgl32.HomogRotate3DZ(mgl32.DegToRad(float32(r.thing.Yaw)))
		thingRotation := thingRotateX.Mul4(thingRotateY.Mul4(thingRotateZ))

		model = thingTranslate.Mul4(thingRotation).Mul4(meshTranslation).Mul4(meshRotation).Mul4(meshPivot)

		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		for _, surface := range mesh.Faces {
			numVerts := int32(len(surface.VertexIds))

			if surface.GeometryMode != 0 {

				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, r.textures[surface.MaterialID])
				gl.Uniform1i(textureUniform, 0)

				gl.DrawArrays(gl.TRIANGLE_FAN, offset, int32(len(surface.VertexIds)))

				gl.BindTexture(gl.TEXTURE_2D, 0)
			}

			offset = offset + numVerts
		}
	}
}

func (r *OpenGl3doRenderer) setupMesh() {
	points := r.makePoints()
	r.vao = loadToVAO(points)
	r.makeTextures()
}

func (r *OpenGl3doRenderer) makePoints() []float32 {
	var points []float32

	for _, mesh := range r.object.GeoSets[r.lod].Meshes {
		for surfaceIdx, surface := range mesh.Faces {
			var mat jk.Material
			if surface.MaterialID != -1 {
				mat = r.object.Materials[surface.MaterialID]
			}

			for idx, id := range surface.VertexIds {
				points = append(points, float32(mesh.Vertices[id][0]))
				points = append(points, float32(mesh.Vertices[id][1]))
				points = append(points, float32(mesh.Vertices[id][2]))

				points = append(points, float32(mesh.FaceNormals[surfaceIdx][0]))
				points = append(points, float32(mesh.FaceNormals[surfaceIdx][1]))
				points = append(points, float32(mesh.FaceNormals[surfaceIdx][2]))

				textureVertexID := surface.TextureVertexIds[idx]
				if len(mesh.TextureVertices) > 0 && textureVertexID != -1 {
					points = append(points, mesh.TextureVertices[textureVertexID][0]/float32(mat.SizeX))
					points = append(points, -mesh.TextureVertices[textureVertexID][1]/float32(mat.SizeY))
				} else {
					points = append(points, 0)
					points = append(points, 0)
				}

				lightIntensity := surface.LightIntensities[idx]
				points = append(points, float32(lightIntensity))
			}
		}
	}
	return points
}

func (r *OpenGl3doRenderer) makeTextures() {
	numTextures := int32(len(r.object.Materials))

	r.textures = make([]uint32, numTextures)

	gl.GenTextures(numTextures, &r.textures[0])

	for i := int32(0); i < numTextures; i++ {
		textureID := r.textures[i]
		material := r.object.Materials[i]

		if len(r.object.Materials[i].Texture) == 0 {
			fmt.Println("empty material")
			continue
		}

		var finalTexture []byte
		if material.Transparent {
			finalTexture = make([]byte, material.SizeX*material.SizeY*4)
			for j := 0; j < int(material.SizeX*material.SizeY); j++ {
				finalTexture[j*4] = r.object.ColorMap.Palette[material.Texture[j]].R
				finalTexture[j*4+1] = r.object.ColorMap.Palette[material.Texture[j]].G
				finalTexture[j*4+2] = r.object.ColorMap.Palette[material.Texture[j]].B

				if material.Texture[j] == 0 {
					finalTexture[j*4+3] = 0
				} else {
					finalTexture[j*4+3] = 255
				}
			}
		} else {
			finalTexture = make([]byte, material.SizeX*material.SizeY*3)
			for j := 0; j < int(material.SizeX*material.SizeY); j++ {
				finalTexture[j*3] = r.object.ColorMap.Palette[material.Texture[j]].R
				finalTexture[j*3+1] = r.object.ColorMap.Palette[material.Texture[j]].G
				finalTexture[j*3+2] = r.object.ColorMap.Palette[material.Texture[j]].B
			}
		}

		loadToTexture(textureID, material.SizeX, material.SizeY, &finalTexture, material.Transparent)
	}
}
