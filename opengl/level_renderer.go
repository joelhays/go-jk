package opengl

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/joelhays/go-jk/jk/jktypes"

	"github.com/go-gl/gl/v3.2-core/gl"
)

type OpenGlLevelRenderer struct {
	thing    *jktypes.Thing
	template *jktypes.Template
	object   *jktypes.JkMesh
	program  *ShaderProgram
	vao      uint32
	textures []uint32
}

func NewOpenGlLevelRenderer(thing *jktypes.Thing, template *jktypes.Template, object *jktypes.JkMesh, program *ShaderProgram) Renderer {
	r := &OpenGlLevelRenderer{thing: thing, template: template, object: object, program: program}
	r.setupMesh()
	return r
}

func (r *OpenGlLevelRenderer) Render() {
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.BindVertexArray(r.vao)
	defer gl.BindVertexArray(0)

	var offset int32
	model := mgl32.Ident4()
	r.ShaderProgram().SetMatrixUniform("model", model)

	for _, surface := range r.object.Surfaces {
		numVerts := int32(len(surface.VertexIds))

		if surface.Geo != 0 {

			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, r.textures[surface.MaterialID])

			r.ShaderProgram().SetIntegerUniform("objectTexture", 0)

			gl.DrawArrays(gl.TRIANGLE_FAN, offset, int32(len(surface.VertexIds)))

			gl.BindTexture(gl.TEXTURE_2D, 0)
		}

		offset = offset + numVerts
	}
}

func (r *OpenGlLevelRenderer) ShaderProgram() *ShaderProgram {
	return r.program
}

func (r *OpenGlLevelRenderer) setupMesh() {
	points := r.makePoints()
	r.vao = loadToVAO(points)
	r.makeTextures()
}

func (r *OpenGlLevelRenderer) makePoints() []float32 {
	var points []float32
	for _, surface := range r.object.Surfaces {
		var mat jktypes.Material
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
				points = append(points, r.object.TextureVertices[textureVertexID][0]/float32(mat.SizeX))  // /mat.XTile)
				points = append(points, -r.object.TextureVertices[textureVertexID][1]/float32(mat.SizeY)) // /mat.YTile)
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

func (r *OpenGlLevelRenderer) makeTextures() {

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
				finalTexture[j*4] = r.object.ColorMaps[0].Palette[material.Texture[j]].R
				finalTexture[j*4+1] = r.object.ColorMaps[0].Palette[material.Texture[j]].G
				finalTexture[j*4+2] = r.object.ColorMaps[0].Palette[material.Texture[j]].B

				if material.Texture[j] == 0 {
					finalTexture[j*4+3] = 0
				} else {
					finalTexture[j*4+3] = 255
				}
			}
		} else {
			finalTexture = make([]byte, material.SizeX*material.SizeY*3)
			for j := 0; j < int(material.SizeX*material.SizeY); j++ {
				finalTexture[j*3] = r.object.ColorMaps[0].Palette[material.Texture[j]].R
				finalTexture[j*3+1] = r.object.ColorMaps[0].Palette[material.Texture[j]].G
				finalTexture[j*3+2] = r.object.ColorMaps[0].Palette[material.Texture[j]].B
			}
		}

		loadToTexture(textureID, material.SizeX, material.SizeY, &finalTexture, material.Transparent)
	}
}
