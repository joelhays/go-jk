package main

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/joelhays/go-vulkan/jk"
)

type ModelRenderer interface {
	Render()
}

type OpenGlModelRenderer struct {
	thing    *jk.Thing
	template *jk.Template
	object   *jk.Jk3do
	Program  uint32
	vao      uint32
	vbo      uint32
	textures []uint32
}

func NewOpenGlModelRenderer(thing *jk.Thing, template *jk.Template, object *jk.Jk3do, program uint32) *OpenGlModelRenderer {
	r := &OpenGlModelRenderer{thing: thing, template: template, object: object, Program: program}
	r.setupMesh()
	return r
}

func (r *OpenGlModelRenderer) Render() {
	model := mgl32.Ident4()
	if r.thing != nil {
		// scale := mgl32.Scale3D(float32(r.template.Size), float32(r.template.Size), float32(r.template.Size))
		rotateX := mgl32.HomogRotate3DX(mgl32.DegToRad(float32(r.thing.Pitch)))
		rotateY := mgl32.HomogRotate3DY(mgl32.DegToRad(float32(r.thing.Roll)))
		rotateZ := mgl32.HomogRotate3DZ(mgl32.DegToRad(float32(r.thing.Yaw)))
		rotation := rotateX.Mul4(rotateY.Mul4(rotateZ))
		translation := mgl32.Translate3D(r.thing.Position.X(), r.thing.Position.Y(), r.thing.Position.Z())
		model = translation.Mul4(rotation)
	}
	modelUniform := gl.GetUniformLocation(r.Program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.BindVertexArray(r.vao)

	var offset int32
	for _, surface := range r.object.Surfaces {
		numVerts := int32(len(surface.VertexIds))

		if surface.Geo != 0 {

			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, r.textures[surface.MaterialID])
			textureUniform := gl.GetUniformLocation(r.Program, gl.Str("objectTexture\x00"))
			gl.Uniform1i(textureUniform, 0)

			gl.DrawArrays(gl.TRIANGLE_FAN, offset, int32(len(surface.VertexIds)))

			gl.BindTexture(gl.TEXTURE_2D, 0)
		}

		offset = offset + numVerts
	}
}

func (r *OpenGlModelRenderer) setupMesh() {
	points := r.makePoints()
	r.makeVao(points)
	r.makeTextures()
}

func (r *OpenGlModelRenderer) makePoints() []float32 {
	var points []float32
	for _, surface := range r.object.Surfaces {
		var mat jk.Material
		if surface.MaterialID != -1 {
			mat = r.object.Materials[surface.MaterialID]
		}

		for idx, id := range surface.VertexIds {
			points = append(points, float32(r.object.Vertices[id][0]))
			points = append(points, float32(r.object.Vertices[id][1]))
			points = append(points, float32(r.object.Vertices[id][2]))

			points = append(points, float32(surface.Normal[0]))
			points = append(points, float32(surface.Normal[1]))
			points = append(points, float32(surface.Normal[2]))

			textureVertexID := surface.TextureVertexIds[idx]
			if textureVertexID != -1 {
				points = append(points, r.object.TextureVertices[textureVertexID][0]/float32(mat.SizeX)) // /mat.XTile)
				points = append(points, r.object.TextureVertices[textureVertexID][1]/float32(mat.SizeY)) // /mat.YTile)
			} else {
				points = append(points, 0)
				points = append(points, 0)
			}

			lightIntensity := surface.LightIntensities[idx]
			points = append(points, float32(lightIntensity))
		}
	}
	return points
}

func (r *OpenGlModelRenderer) makeVao(points []float32) {
	gl.GenBuffers(1, &r.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &r.vao)
	gl.BindVertexArray(r.vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 9*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 9*4, gl.PtrOffset(3*4))

	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 9*4, gl.PtrOffset(6*4))

	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointer(3, 1, gl.FLOAT, false, 9*4, gl.PtrOffset(7*4))
}

func (r *OpenGlModelRenderer) makeTextures() {

	numTextures := int32(len(r.object.Materials))

	r.textures = make([]uint32, numTextures)

	gl.GenTextures(numTextures, &r.textures[0])

	for i := int32(0); i < numTextures; i++ {
		textureID := r.textures[i]
		material := r.object.Materials[i]

		gl.BindTexture(gl.TEXTURE_2D, textureID)

		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

		if len(r.object.Materials[i].Texture) == 0 {
			fmt.Println("empty material")
			continue
		}

		var finalTexture []byte
		if material.Transparent {
			finalTexture = make([]byte, material.SizeX*material.SizeY*4)
			for j := 0; j < int(material.SizeX*material.SizeY); j++ {
				finalTexture[j*4] = r.object.ColorMaps[0].Pallette[material.Texture[j]].R
				finalTexture[j*4+1] = r.object.ColorMaps[0].Pallette[material.Texture[j]].G
				finalTexture[j*4+2] = r.object.ColorMaps[0].Pallette[material.Texture[j]].B

				if material.Texture[j] == 0 {
					finalTexture[j*4+3] = 0
				} else {
					finalTexture[j*4+3] = 255
				}
				gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, material.SizeX, material.SizeY, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(finalTexture))
			}
		} else {
			finalTexture = make([]byte, material.SizeX*material.SizeY*3)
			for j := 0; j < int(material.SizeX*material.SizeY); j++ {
				finalTexture[j*3] = r.object.ColorMaps[0].Pallette[material.Texture[j]].R
				finalTexture[j*3+1] = r.object.ColorMaps[0].Pallette[material.Texture[j]].G
				finalTexture[j*3+2] = r.object.ColorMaps[0].Pallette[material.Texture[j]].B
			}
			gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, material.SizeX, material.SizeY, 0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(finalTexture))
		}

		gl.GenerateMipmap(gl.TEXTURE_2D)

		gl.BindTexture(gl.TEXTURE_2D, 0)
	}
}
