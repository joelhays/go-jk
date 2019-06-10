package opengl

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/joelhays/go-jk/jk"
)

type OpenGlBmRenderer struct {
	bm       *jk.BMFile
	program  *ShaderProgram
	vao      uint32
	textures []uint32
}

func NewOpenGlBmRenderer(bm *jk.BMFile, program *ShaderProgram) Renderer {
	r := &OpenGlBmRenderer{bm: bm, program: program}

	r.setupMesh()
	return r
}

func (r *OpenGlBmRenderer) Render() {
	gl.BindVertexArray(r.vao)
	defer gl.BindVertexArray(0)

	var offset int32 = 0
	model := mgl32.Ident4()
	// model = mgl32.Scale3D(0.5, 0.5, 0.5)
	r.ShaderProgram().SetMatrixUniform("model", model)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, r.textures[0])

	r.ShaderProgram().SetIntegerUniform("objectTexture", 0)

	gl.DrawArrays(gl.TRIANGLE_FAN, offset, 6)

	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (r *OpenGlBmRenderer) ShaderProgram() *ShaderProgram {
	return r.program
}

func (r *OpenGlBmRenderer) setupMesh() {
	points := r.makePoints()
	r.vao = loadToVAO(points)
	r.makeTextures()
}

func (r *OpenGlBmRenderer) makePoints() []float32 {
	// VERTICES (3), NORMALS (3), UV (2), LIGHT (1)
	var points = []float32{
		/*pos bl*/ -1, -1, 0 /*norm*/, 0.0, 1.0, 0.0 /*tex*/, 0.0, 0.0 /*light*/, 1,
		/*pos br*/ +1, -1, 0 /*norm*/, 0.0, 1.0, 0.0 /*tex*/, 1.0, 0.0 /*light*/, 1,
		/*pos tr*/ +1, +1, 0 /*norm*/, 0.0, 1.0, 0.0 /*tex*/, 1.0, 1.0 /*light*/, 1,
		/*pos tr*/ +1, +1, 0 /*norm*/, 0.0, 1.0, 0.0 /*tex*/, 1.0, 1.0 /*light*/, 1,
		/*pos tl*/ -1, +1, 0 /*norm*/, 0.0, 1.0, 0.0 /*tex*/, 0.0, 1.0 /*light*/, 1,
		/*pos bl*/ -1, -1, 0 /*norm*/, 0.0, 1.0, 0.0 /*tex*/, 0.0, 1.0 /*light*/, 1,
	}
	return points
}

func (r *OpenGlBmRenderer) makeTextures() {
	numTextures := int32(len(r.bm.Images))

	if numTextures == 0 {
		fmt.Println("bm contains no images")
		return
	}

	r.textures = make([]uint32, numTextures)

	gl.GenTextures(numTextures, &r.textures[0])

	for i := int32(0); i < numTextures; i++ {
		textureID := r.textures[i]
		material := r.bm.Images[i]

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
		loadToTexture(textureID, material.SizeX, material.SizeY, &finalTexture, false)
	}
}

func (r *OpenGlBmRenderer) GetTextureID() uint32 {
	return r.textures[0]
}
