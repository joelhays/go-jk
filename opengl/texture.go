package opengl

import (
	"github.com/go-gl/gl/v3.2-core/gl"
)

func loadToTexture(textureID uint32, sizeX int32, sizeY int32, data *[]byte, useAlpha bool) {
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	if useAlpha {
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, sizeX, sizeY, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(*data))
	} else {
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, sizeX, sizeY, 0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(*data))
	}

	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)
}
