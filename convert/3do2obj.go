package convert

import (
	"fmt"
	"os"

	"github.com/joelhays/go-vulkan/jk"
)

func From3do2obj(mesh *jk.JkMesh) {
	file, err := os.OpenFile("./test.obj", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("o Test\n")

	for _, v := range mesh.Vertices {
		line := fmt.Sprintf("v %f %f %f 1.0\n", v.X(), v.Y(), v.Z())
		file.WriteString(line)
	}

	for _, vt := range mesh.TextureVertices {
		line := fmt.Sprintf("vt %f %f\n", vt.X(), vt.Y())
		file.WriteString(line)
	}

	for _, vn := range mesh.VertexNormals {
		line := fmt.Sprintf("vn %f %f %f\n", vn.X(), vn.Y(), vn.Z())
		file.WriteString(line)
	}

	for _, f := range mesh.Surfaces {
		if f.Geo != 4 {
			continue
		}
		line := "f"

		for idx, v := range f.VertexIds {
			line = line + fmt.Sprintf(" %d/%d", v+1, f.TextureVertexIds[idx]+1)
		}

		line = line + "\n"

		file.WriteString(line)

		// f v1/vt1/vn1 v2/vt2/vn2 v3/vt3/vn3 ...
	}

	// f v1/vt1/vn1 v2/vt2/vn2 v3/vt3/vn3 ...
}
