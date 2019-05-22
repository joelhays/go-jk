package opengl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

func LoadToVAO(data []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(data), gl.Ptr(data), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	/* position */
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 9*4, gl.PtrOffset(0))

	/* normal */
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 9*4, gl.PtrOffset(3*4))

	/* uv */
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 9*4, gl.PtrOffset(6*4))

	/* light intensity */
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointer(3, 1, gl.FLOAT, false, 9*4, gl.PtrOffset(7*4))

	return vao
}
