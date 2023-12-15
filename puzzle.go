package main

import (
	"math"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/yobert/vector"
)

type Puzzle struct {
	Segments    [27]Segment
	CornerCount int
	Step        int
	MaxStep     int
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

func MakePuzzle() *Puzzle {
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

	//	rotate := 1
	for _, v := range corners {
		p.Segments[v].Kind = Corner
		//		p.Segments[v].Rotate = rotate
		//		p.Segments[v].LastRotate = clamp(rotate - 1)
		//		rotate = clamp(rotate + 1)
		p.CornerCount = p.CornerCount + 1
	}

	p.MaxStep = int(math.Pow(4, float64(p.CornerCount)))

	return &p
}

func (p *Puzzle) Advance() *Puzzle {
	np := *p

	np.Step++
	if np.Step >= np.MaxStep {
		return nil
	}

	v := np.Step
	for i := range np.Segments {
		np.Segments[i].Rotate = v % 4
		v = v >> 2
	}

	return &np
}

const (
	m = 10
)

var (
	grid  [m * 2][m * 2][m * 2]bool
	blank [m * 2][m * 2][m * 2]bool
)

func (p *Puzzle) Analyze() (sx, sy, sz int, valid bool) {

	// reset to zero
	grid = blank

	var (
		min_sx int
		min_sy int
		min_sz int
		max_sx int
		max_sy int
		max_sz int
	)

	min_sx = m
	min_sy = m
	min_sz = m

	max_sx = -m
	max_sy = -m
	max_sz = -m

	pos := vector.V3{}
	rot := vector.IdentityQ()

	// "Render" into our integer grid
	for _, segment := range p.Segments {

		ix := int(math.Round(pos.X))
		iy := int(math.Round(pos.Y))
		iz := int(math.Round(pos.Z))

		ix += m
		iy += m
		iz += m

		if grid[ix][iy][iz] {
			return
		}
		grid[ix][iy][iz] = true

		if ix < min_sx {
			min_sx = ix
		}
		if ix > max_sx {
			max_sx = ix
		}
		if iy < min_sy {
			min_sy = iy
		}
		if iy > max_sy {
			max_sy = iy
		}
		if iz < min_sz {
			min_sz = iz
		}
		if iz > max_sz {
			max_sz = iz
		}

		if segment.Kind == Corner {
			rot = rot.Mult(vector.AxisAngleQ(vector.V3{0, 0, 1}, vector.Degree(90).Radian()))

			angle := 90.0 * float64(segment.Rotate)

			rot = rot.Mult(vector.AxisAngleQ(vector.V3{1, 0, 0}, vector.Degree(angle).Radian()))
		}
		dir := vector.V3{0, 1, 0}
		dir = rot.M33().MultV3(dir)
		pos = pos.Add(dir)
	}

	valid = true
	sx = max_sx - min_sx
	sy = max_sy - min_sy
	sz = max_sz - min_sz
	return
}

func draw_puzzle(p *Puzzle, animate float64) {

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
		if segment.Blue {
			debug_material_blue()
		} else {
			debug_material()
		}
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
