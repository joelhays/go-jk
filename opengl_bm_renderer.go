package main

import (
	"fmt"
	"github.com/joelhays/go-jk/opengl"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/joelhays/go-jk/jk"
)

type OpenGlBmRenderer struct {
	bm       *jk.BMFile
	Program  uint32
	vao      uint32
	textures []uint32
}

func NewOpenGlBmRenderer(bm *jk.BMFile, program uint32) *OpenGlBmRenderer {
	r := &OpenGlBmRenderer{bm: bm, Program: program}

	r.setupMesh()
	return r
}

func (r *OpenGlBmRenderer) Render() {

	gl.BindVertexArray(r.vao)

	var offset int32 = 0
	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(r.Program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, r.textures[0])
	textureUniform := gl.GetUniformLocation(r.Program, gl.Str("objectTexture\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.DrawArrays(gl.TRIANGLE_FAN, offset, 6)

	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (r *OpenGlBmRenderer) setupMesh() {
	points := r.makePoints()
	r.vao = opengl.LoadToVAO(points)
	r.makeTextures()
}

func (r *OpenGlBmRenderer) makePoints() []float32 {
	// VERTICES (3), NORMALS (3), UV (2), LIGHT (1)
	var points = []float32{
		/*pos bl*/ -0.5, -0.5, -0.5 /*norm*/, 0.0, 1.0, 0.0 /*tex*/, 1.0, 0 /*light*/, 1,
		/*pos br*/ 0.5, -0.5, -0.5 /*norm*/, 0.0, 1.0, 0.0 /*tex*/, 1.0, 1.0 /*light*/, 1,
		/*pos tr*/ 0.5, 0.5, -0.5 /*norm*/, 0.0, 1.0, 0.0 /*tex*/, 0.0, 1.0 /*light*/, 1,
		/*pos tr*/ 0.5, 0.5, -0.5 /*norm*/, 0.0, 1.0, 0.0 /*tex*/, 0.0, 1.0 /*light*/, 1,
		/*pos tl*/ -0.5, 0.5, -0.5 /*norm*/, 0.0, 1.0, 0.0 /*tex*/, 0.0, 0.0 /*light*/, 1,
		/*pos bl*/ -0.5, -0.5, -0.5 /*norm*/, 0.0, 1.0, 0.0 /*tex*/, 1.0, 0.0 /*light*/, 1,
	}
	return points
}

func (r *OpenGlBmRenderer) makeTextures() {
	numTextures := int32(len(r.bm.Images))

	r.textures = make([]uint32, numTextures)

	gl.GenTextures(numTextures, &r.textures[0])

	for i := int32(0); i < numTextures; i++ {
		textureID := r.textures[i]
		material := r.bm.Images[i]

		gl.BindTexture(gl.TEXTURE_2D, textureID)

		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

		if len(r.bm.Images) == 0 {
			fmt.Println("empty material")
			continue
		}

		var finalTexture []byte
		finalTexture = make([]byte, material.SizeX*material.SizeY*3)
		for j := 0; j < int(material.SizeX*material.SizeY); j++ {
			finalTexture[j*3] = r.bm.Palette.Palette[material.Data[j]].R
			finalTexture[j*3+1] = r.bm.Palette.Palette[material.Data[j]].G
			finalTexture[j*3+2] = r.bm.Palette.Palette[material.Data[j]].B
		}
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, material.SizeX, material.SizeY, 0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(finalTexture))

		gl.GenerateMipmap(gl.TEXTURE_2D)

		gl.BindTexture(gl.TEXTURE_2D, 0)
	}
}
