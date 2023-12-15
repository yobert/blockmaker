package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/yobert/progress"
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

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0, 0, 0, 1)

	lastTime := glfw.GetTime()
	tt := 0.0
	ti := 0
	oldmsg := ""

	p := MakePuzzle()

	// Solve the puzzle: Find the start shape (widest possible) and the solution shape (slimmest possible)

	var (
		max_p *Puzzle
		max_s int

		min_p *Puzzle
	)

	bar := progress.NewBar(p.MaxStep, "processing...")

	for p != nil {
		sx, sy, sz, valid := p.Analyze()
		if valid {
			sany := 0
			if sx > sany {
				sany = sx
			}
			if sy > sany {
				sany = sy
			}
			if sz > sany {
				sany = sz
			}
			if max_p == nil || sany > max_s {
				max_p = p
			}

			if sx == 3 && sy == 3 && sz == 3 {
				min_p = p
				fmt.Println("found a solution!")
			}
		}
		p = p.Advance()
		bar.Next()
	}
	bar.Done()

	if min_p == nil || max_p == nil {
		return fmt.Errorf("Awwww shucks! No solution.")
	}

	// Animate from the maximum shape to the minimum.
	for i := range min_p.Segments {
		min_p.Segments[i].LastRotate = max_p.Segments[i].Rotate
	}

	animate := 0.0

	for !window.ShouldClose() {
		glfw.PollEvents()

		window.MakeContextCurrent()

		setup_camera(window, cam)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		draw_debug()

		debug_lighting()
		draw_puzzle(p, animate)

		print_gl_errors()
		window.SwapBuffers()

		// Calculate how long it's been since the last frame. We'll use that in an FPS printout on the console,
		// as well as a multiplier for camera movement speed.
		newTime := glfw.GetTime()
		t := newTime - lastTime
		lastTime = newTime

		// Advance animation
		animate += t * 0.1 * 0.5
		if animate > 1 {
			animate = 0
		}

		// Print the FPS
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

		// Chose a movement vector
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

		// Look around?
		look_speed := t * 2
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

		// Move camera position
		move_speed := t * 5
		if keys[glfw.KeyRightShift] {
			move_speed *= 10.0
		}
		dir := cam.RotAxis.M33().MultV3(move).Scale(move_speed)
		cam.Position = cam.Position.Add(dir)
	}
	fmt.Println()

	return nil
}

func setup_camera(window *glfw.Window, cam *vector.Camera) {
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
}

func draw_debug() {
	//gl.Enable(gl.CULL_FACE)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.CULL_FACE)

	// Draw axis
	gl.Begin(gl.LINES)
	gl.Color3f(0, 0, 0)
	gl.Vertex3f(0, 0, 0)
	gl.Color3f(1, 0, 0)
	gl.Vertex3f(10, 0, 0)

	gl.Color3f(0, 0, 0)
	gl.Vertex3f(0, 0, 0)
	gl.Color3f(0, 1, 0)
	gl.Vertex3f(0, 10, 0)

	gl.Color3f(0, 0, 0)
	gl.Vertex3f(0, 0, 0)
	gl.Color3f(0, 0, 1)
	gl.Vertex3f(0, 0, 10)
	gl.End()

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
}

func debug_lighting() {
	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.LIGHT0)

	pos := []float32{10, 20, 30, 0}
	amb := []float32{0.5, 0.5, 0.5, 1}
	dif := []float32{0.7, 0.7, 0.7, 1}
	spe := []float32{1, 1, 1, 1}
	g_amb := []float32{0, 0, 0, 1}

	gl.Lightfv(gl.LIGHT0, gl.POSITION, &pos[0])
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &amb[0])
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, &dif[0])
	gl.Lightfv(gl.LIGHT0, gl.SPECULAR, &spe[0])

	gl.LightModelfv(gl.LIGHT_MODEL_AMBIENT, &g_amb[0])
}
func debug_material() {
	gl.Disable(gl.TEXTURE_2D)

	amb := []float32{0.1, 0.1, 0.1, 1}
	dif := []float32{0.4, 0.4, 0.4, 1}
	spe := []float32{0.5, 0.5, 0.5, 1}

	gl.Materialf(gl.FRONT, gl.SHININESS, 120) // specular exponent, range of 0..128
	gl.Materialfv(gl.FRONT, gl.AMBIENT, &amb[0])
	gl.Materialfv(gl.FRONT, gl.DIFFUSE, &dif[0])
	gl.Materialfv(gl.FRONT, gl.SPECULAR, &spe[0])
}
func debug_material_blue() {
	gl.Disable(gl.TEXTURE_2D)

	amb := []float32{0.0, 0.0, 0.1, 1}
	dif := []float32{0.0, 0.0, 0.4, 1}
	spe := []float32{0.5, 0.5, 0.5, 1}

	gl.Materialf(gl.FRONT, gl.SHININESS, 120) // specular exponent, range of 0..128
	gl.Materialfv(gl.FRONT, gl.AMBIENT, &amb[0])
	gl.Materialfv(gl.FRONT, gl.DIFFUSE, &dif[0])
	gl.Materialfv(gl.FRONT, gl.SPECULAR, &spe[0])
}

func print_gl_errors() {
	if ge := gl.GetError(); ge != gl.NO_ERROR {
		fmt.Println(ge)
	}
}
