package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/yobert/vector"
)

var (
	cam  *vector.Camera
	keys [1024]bool
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func r(min, max float64) float64 {
	return (rand.Float64() * (max - min)) + min
}

func window_key(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

	if action == glfw.Press {
		keys[key] = true
	} else if action == glfw.Release {
		keys[key] = false
	}

	if key == glfw.KeyEscape || key == glfw.KeyQ {
		w.SetShouldClose(true)
	}
}

func run() error {
	cam = &vector.Camera{}
	cam.Position = vector.V3{0, 0, 10}

	err := glfw.Init()
	if err != nil {
		return err
	}
	defer glfw.Terminate()

	//glfw.WindowHint(glfw.Samples, 0)

	window, err := glfw.CreateWindow(1920, 1080, "blockmaker", nil, nil)
	if err != nil {
		return err
	}

	window.SetKeyCallback(window_key)
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return err
	}

	glfw.SwapInterval(1)

	lastTime := glfw.GetTime()
	tt := 0.0
	ti := 0
	oldmsg := ""

	p := MakePuzzle()

	for !window.ShouldClose() {
		glfw.PollEvents()

		window.MakeContextCurrent()

		wx, wy := window.GetSize()

		cam.Width = float64(wx)
		cam.Height = float64(wy)
		cam.YFov = 70
		cam.Near = 0.01
		cam.Far = 100

		cam.SetupViewProjection()
		cam.SetupModelView()

		gl.Viewport(0, 0, int32(wx), int32(wy))

		gl.MatrixMode(gl.PROJECTION)
		gl.LoadMatrixd(&cam.Projection[0])

		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadMatrixd(&cam.ModelView[0])

		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LESS)

		//gl.Enable(gl.CULL_FACE)

		//gl.Enable(gl.LIGHTING)

		rand.Seed(666)
		gl.Begin(gl.POINTS)
		for i := 0; i < 1000; i++ {
			gl.Color3f(1, 1, 1)
			x := r(-10, 10)
			y := r(-10, 10)
			z := r(-10, 10)
			gl.Vertex3d(x, y, z)
		}
		gl.End()

		draw_puzzle(p)

		if ge := gl.GetError(); ge != gl.NO_ERROR {
			fmt.Println(ge)
		}

		move_speed := 0.05
		move := vector.V3{}

		if keys[glfw.KeyA] {
			move.Y = 1
		}
		if keys[glfw.KeyZ] {
			move.Y = -1
		}
		if keys[glfw.KeyS] {
			move.X = -1
		}
		if keys[glfw.KeyF] {
			move.X = 1
		}
		if keys[glfw.KeyE] {
			move.Z = -1
		}
		if keys[glfw.KeyD] {
			move.Z = 1
		}

		look_speed := 0.02

		e := cam.RotAxis
		e.Z = 0

		if keys[glfw.KeyUp] {
			e.X += (vector.Radian)(look_speed)
		}
		if keys[glfw.KeyDown] {
			e.X -= (vector.Radian)(look_speed)
		}
		if keys[glfw.KeyLeft] {
			e.Y += (vector.Radian)(look_speed)
		}
		if keys[glfw.KeyRight] {
			e.Y -= (vector.Radian)(look_speed)
		}
		if e.X > math.Pi/2 {
			e.X = math.Pi / 2
		}
		if e.X < -math.Pi/2 {
			e.X = -math.Pi / 2
		}

		cam.RotAxis = e

		dir := cam.RotAxis.M33().MultV3(move).Scale(move_speed)
		cam.Position = cam.Position.Add(dir)

		window.SwapBuffers()

		newTime := glfw.GetTime()
		t := newTime - lastTime
		lastTime = newTime

		tt += t
		ti++
		if tt > 1 {
			newmsg := fmt.Sprintf("%.0fhz", float64(ti)/tt)
			if newmsg != oldmsg {
				for _ = range oldmsg {
					fmt.Print("\b")
				}
				for _ = range oldmsg {
					fmt.Print(" ")
				}
				for _ = range oldmsg {
					fmt.Print("\b")
				}
				fmt.Print(newmsg)
				oldmsg = newmsg
			}

			tt = 0
			ti = 0
		}
	}
	fmt.Println()

	return nil
}

func draw_puzzle(p Puzzle) {
	dir := false
	pos := vector.V3{}
	for _, segment := range p.Segments {
		b := Box{
			Blue:     segment.Blue,
			Origin:   pos,
			HalfSize: vector.V3{0.5, 0.5, 0.5},
		}
		b.Draw()

		if segment.Kind == Corner {
			dir = !dir
		}

		v := vector.V3{0, 1, 0}
		if dir {
			v = vector.V3{1, 0, 0}
		}

		pos = pos.Add(v)
	}
}
