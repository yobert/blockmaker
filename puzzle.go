package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/yobert/vector"
)

type Puzzle struct {
	Segments [27]Segment
}

type Segment struct {
	Blue   bool
	Kind   Kind
	Rotate int

	// For animation
	LastRotate int
}

type Kind int

const (
	End Kind = iota
	Corner
	Straight
)

func clamp(i int) int {
	if i == -1 {
		return 3
	}
	if i == 4 {
		return 0
	}
	return i
}

func MakePuzzle() Puzzle {
	p := Puzzle{}

	blue := true
	for i := range p.Segments {
		p.Segments[i].Blue = blue
		blue = !blue
	}

	p.Segments[0].Kind = End
	p.Segments[26].Kind = End

	corners := []int{
		2, 3, 4, 6, 7,
		9, 10, 11, 13, 15,
		16, 17, 18, 20, 22,
		24,
	}

	rotate := 1
	for _, v := range corners {
		p.Segments[v].Kind = Corner
		p.Segments[v].Rotate = rotate
		p.Segments[v].LastRotate = clamp(rotate - 1)
		rotate = clamp(rotate + 1)
	}

	return p
}

func draw_puzzle(p Puzzle, animate float64) {
	pos := vector.V3{}
	rot := vector.IdentityQ()

	for _, segment := range p.Segments {
		if segment.Kind == Corner {
			rot = rot.Mult(vector.AxisAngleQ(vector.V3{0, 0, 1}, vector.Degree(90).Radian()))

			anglefrom := 90.0 * float64(segment.LastRotate)
			angleto := 90.0 * float64(segment.Rotate)

			angle := ((angleto - anglefrom) * animate) + anglefrom

			rot = rot.Mult(vector.AxisAngleQ(vector.V3{1, 0, 0}, vector.Degree(angle).Radian()))
		}

		b := Box{
			Blue:     segment.Blue,
			Origin:   pos,
			HalfSize: vector.V3{0.5, 0.5, 0.5},
		}

		gl.PushMatrix()
		gl.Translated(
			b.Origin.X,
			b.Origin.Y,
			b.Origin.Z)

		rotmat := rot.M33().M44()
		gl.MultMatrixd(&rotmat[0])
		b.Draw()
		gl.PopMatrix()

		dir := vector.V3{0, 1, 0}

		dir = rot.M33().MultV3(dir)

		//		v := vector.V3{0, 1, 0}
		//		if dir {
		//			v = vector.V3{1, 0, 0}
		//		}

		pos = pos.Add(dir)
	}
}
