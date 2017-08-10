package chopsui

import (
	"github.com/fogleman/fauxgl"
	"github.com/go-gl/gl/v2.1/gl"
)

type MeshData struct {
	Buffer []float32
	Box    fauxgl.Box
}

type Mesh struct {
	Data         *MeshData
	Transform    fauxgl.Matrix
	VertexBuffer uint32
	VertexCount  int32
}

func NewMesh(data *MeshData) *Mesh {
	// compute transform to scale and center mesh
	scale := fauxgl.V(2, 2, 2).Div(data.Box.Size()).MinComponent()
	transform := fauxgl.Identity()
	transform = transform.Translate(data.Box.Center().Negate())
	transform = transform.Scale(fauxgl.V(scale, scale, scale))

	// generate vbo
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data.Buffer)*4, gl.Ptr(data.Buffer), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// compute number of vertices
	count := int32(len(data.Buffer) / 3)

	return &Mesh{data, transform, vbo, count}
}

func NewPlaneMesh() *Mesh {
	const n = 1.5
	buffer := []float32{
		-n, -n, 0,
		n, -n, 0,
		n, n, 0,
		-n, -n, 0,
		n, n, 0,
		-n, n, 0,
	}
	box := fauxgl.Box{fauxgl.Vector{-n, -n, 0}, fauxgl.Vector{n, n, 0}}
	mesh := NewMesh(&MeshData{buffer, box})
	mesh.Transform = fauxgl.Identity()
	return mesh
}

func (mesh *Mesh) Draw(positionAttrib uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.VertexBuffer)
	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 3, gl.FLOAT, false, 12, gl.PtrOffset(0))
	gl.DrawArrays(gl.TRIANGLES, 0, mesh.VertexCount)
	gl.DisableVertexAttribArray(positionAttrib)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (mesh *Mesh) Destroy() {
	gl.DeleteBuffers(1, &mesh.VertexBuffer)
}

func (mesh *Mesh) ToFauxgl() *fauxgl.Mesh {
	b := mesh.Data.Buffer
	nv := len(b) / 3
	nt := nv / 3
	triangles := make([]*fauxgl.Triangle, nt)
	for i := 0; i < nt; i++ {
		j := i * 9
		p1 := fauxgl.Vector{float64(b[j+0]), float64(b[j+1]), float64(b[j+2])}
		p2 := fauxgl.Vector{float64(b[j+3]), float64(b[j+4]), float64(b[j+5])}
		p3 := fauxgl.Vector{float64(b[j+6]), float64(b[j+7]), float64(b[j+8])}
		t := fauxgl.NewTriangleForPoints(p1, p2, p3)
		triangles[i] = t
	}
	return fauxgl.NewTriangleMesh(triangles)
}
