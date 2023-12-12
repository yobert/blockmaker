package main

import (
	"fmt"
	"runtime"
	"os"
	"math/rand"

	"github.com/yobert/vector"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/gl/v2.1/gl"
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
	if key == glfw.KeyEscape || key == glfw.KeyQ {
		w.SetShouldClose(true)
	}
}

func run() error {

	err := glfw.Init()
	if err != nil {
		return err
	}
	defer glfw.Terminate()

	//glfw.WindowHint(glfw.Samples, 0)

	window, err := glfw.CreateWindow(640, 480, "blockmaker", nil, nil)
	if err != nil {
		return err
	}

	window.SetKeyCallback(window_key)
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return err
	}

	glfw.SwapInterval(1)

	cam := &vector.Camera{}

	lastTime := glfw.GetTime()
	tt := 0.0
	ti := 0
	oldmsg := ""

	for !window.ShouldClose() {
		glfw.PollEvents()

		window.MakeContextCurrent()

		wx, wy := window.GetSize()

		cam.Width = float64(wx)
		cam.Height = float64(wy)
		cam.YFov = 70
		cam.Near = 1
		cam.Far = 1000

		cam.Position = vector.V3{0, 0, 10}

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

		rand.Seed(666)
		gl.Begin(gl.POINTS)
		for i := 0; i < 1000; i++ {
			gl.Color3f(1, 0, 0)
			x := r(-10, 10)
			y := r(-10, 10)
			z := r(-10, 10)
			gl.Vertex3d(x, y, z)
		}
		gl.End()

		if ge := gl.GetError(); ge != gl.NO_ERROR {
			fmt.Println(ge)
		}

		window.SwapBuffers()

		newTime := glfw.GetTime()
		t := newTime - lastTime
		lastTime = newTime

		tt += t
		ti++
		if tt > 1 {
			//window.SetTitle()
			newmsg := fmt.Sprintf("%.0fhz", float64(ti) / tt)
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
