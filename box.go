package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/yobert/vector"
)

type Box struct {
	Origin   vector.V3
	HalfSize vector.V3
	Blue     bool
}

func (box Box) Draw() {
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	gl.Begin(gl.QUADS)

	r := 0.5
	g := 0.5
	b := 0.5

	if box.Blue {
		r = 0
		g = 0
		b = 1
	}

	// -Z
	v := 0.1
	gl.Color3d(r*v, g*v, b*v)
	gl.Normal3d(0, 0, -1)
	gl.Vertex3d((-box.HalfSize.X), (-box.HalfSize.Y), (-box.HalfSize.Z))
	gl.Vertex3d((-box.HalfSize.X), (box.HalfSize.Y), (-box.HalfSize.Z))
	gl.Vertex3d((box.HalfSize.X), (box.HalfSize.Y), (-box.HalfSize.Z))
	gl.Vertex3d((box.HalfSize.X), (-box.HalfSize.Y), (-box.HalfSize.Z))

	// -X
	v = 0.2
	gl.Color3d(r*v, g*v, b*v)
	gl.Normal3d(-1, 0, 0)
	gl.Vertex3d((-box.HalfSize.X), (-box.HalfSize.Y), (-box.HalfSize.Z))
	gl.Vertex3d((-box.HalfSize.X), (-box.HalfSize.Y), (box.HalfSize.Z))
	gl.Vertex3d((-box.HalfSize.X), (box.HalfSize.Y), (box.HalfSize.Z))
	gl.Vertex3d((-box.HalfSize.X), (box.HalfSize.Y), (-box.HalfSize.Z))

	// -Y
	v = 0.3
	gl.Color3d(r*v, g*v, b*v)
	gl.Normal3d(0, -1, 0)
	gl.Vertex3d((-box.HalfSize.X), (-box.HalfSize.Y), (-box.HalfSize.Z))
	gl.Vertex3d((box.HalfSize.X), (-box.HalfSize.Y), (-box.HalfSize.Z))
	gl.Vertex3d((box.HalfSize.X), (-box.HalfSize.Y), (box.HalfSize.Z))
	gl.Vertex3d((-box.HalfSize.X), (-box.HalfSize.Y), (box.HalfSize.Z))

	// +Z
	v = 0.4
	gl.Color3d(r*v, g*v, b*v)
	gl.Normal3d(0, 0, 1)
	gl.Vertex3d((box.HalfSize.X), (box.HalfSize.Y), (box.HalfSize.Z))
	gl.Vertex3d((-box.HalfSize.X), (box.HalfSize.Y), (box.HalfSize.Z))
	gl.Vertex3d((-box.HalfSize.X), (-box.HalfSize.Y), (box.HalfSize.Z))
	gl.Vertex3d((box.HalfSize.X), (-box.HalfSize.Y), (box.HalfSize.Z))

	// +X
	v = 0.5
	gl.Color3d(r*v, g*v, b*v)
	gl.Normal3d(1, 0, 0)
	gl.Vertex3d((box.HalfSize.X), (box.HalfSize.Y), (box.HalfSize.Z))
	gl.Vertex3d((box.HalfSize.X), (-box.HalfSize.Y), (box.HalfSize.Z))
	gl.Vertex3d((box.HalfSize.X), (-box.HalfSize.Y), (-box.HalfSize.Z))
	gl.Vertex3d((box.HalfSize.X), (box.HalfSize.Y), (-box.HalfSize.Z))

	// +Y
	v = 0.6
	gl.Color3d(r*v, g*v, b*v)
	gl.Normal3d(0, 1, 0)
	gl.Vertex3d(box.HalfSize.X, box.HalfSize.Y, box.HalfSize.Z)
	gl.Vertex3d(box.HalfSize.X, box.HalfSize.Y, -box.HalfSize.Z)
	gl.Vertex3d(-box.HalfSize.X, box.HalfSize.Y, -box.HalfSize.Z)
	gl.Vertex3d(-box.HalfSize.X, box.HalfSize.Y, box.HalfSize.Z)

	gl.End()
}
