package main

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
