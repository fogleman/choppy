package chopsui

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fogleman/choppy"
	"github.com/fogleman/fauxgl"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var vertexShader = `
#version 120

uniform mat4 matrix;

attribute vec4 position;

varying vec3 ec_pos;

void main() {
	gl_Position = matrix * position;
	ec_pos = vec3(gl_Position);
}
`

var fragmentShader = `
#version 120

varying vec3 ec_pos;

const vec3 light_direction = normalize(vec3(1, -1.5, 1));
const vec3 object_color = vec3(0x5b / 255.0, 0xac / 255.0, 0xe3 / 255.0);

void main() {
	vec3 ec_normal = normalize(cross(dFdx(ec_pos), dFdy(ec_pos)));
	float diffuse = max(0, dot(ec_normal, light_direction)) * 0.85 + 0.2;
	vec3 color = object_color * diffuse;
	gl_FragColor = vec4(color, 1);
}
`

var planeVertexShader = `
#version 120

uniform mat4 matrix;

attribute vec4 position;

void main() {
	gl_Position = matrix * position;
}
`

var planeFragmentShader = `
#version 120

void main() {
	gl_FragColor = vec4(1, 0, 1, 0.5);
}
`

func init() {
	runtime.LockOSThread()
}

func loadMesh(path string, ch chan *MeshData) {
	go func() {
		start := time.Now()
		data, err := LoadMesh(path)
		if err != nil {
			return // TODO: display an error
		}
		fmt.Printf(
			"loaded %d triangles in %.3f seconds\n",
			len(data.Buffer)/9, time.Since(start).Seconds())
		ch <- data
	}()
}

func Run(path string) {
	start := time.Now()

	// load mesh in the background
	ch := make(chan *MeshData)
	loadMesh(path, ch)

	// initialize glfw
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	// create the window
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(640, 640, path, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	fmt.Printf("window shown at %.3f seconds\n", time.Since(start).Seconds())

	// initialize gl
	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.DEPTH_TEST)
	// gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.ClearColor(float32(0xd4)/255, float32(0xd9)/255, float32(0xde)/255, 1)

	// compile shaders
	program, err := compileProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	// gl.UseProgram(program)

	matrixUniform := uniformLocation(program, "matrix")
	positionAttrib := attribLocation(program, "position")

	planeProgram, err := compileProgram(planeVertexShader, planeFragmentShader)
	if err != nil {
		panic(err)
	}
	planeMatrixUniform := uniformLocation(planeProgram, "matrix")
	planePositionAttrib := attribLocation(planeProgram, "position")

	var mesh *Mesh
	planeMesh := NewPlaneMesh()

	// chop function
	chop := func(a *AppInteractor) {
		if mesh == nil {
			return
		}
		start := time.Now()
		fm := mesh.ToFauxgl()
		m1 := a.MeshInteractor.matrix().Mul(mesh.Transform)
		m2 := a.PlaneInteractor.matrix().Mul(planeMesh.Transform)
		fm.Transform(m1)
		point := m2.MulPosition(fauxgl.Vector{})
		normal := m2.MulDirection(fauxgl.Vector{0, 0, 1})
		fm1 := choppy.Chop(fm, point, normal)
		fm2 := choppy.Chop(fm, point, normal.Negate())
		fm1.Transform(m1.Inverse())
		fm2.Transform(m1.Inverse())
		fmt.Printf(
			"chopped mesh in %.3f seconds\n", time.Since(start).Seconds())
		fm1.SaveSTL("out1.stl")
		fm2.SaveSTL("out2.stl")
	}

	// create interactor
	interactor := NewAppInteractor(chop)
	BindInteractor(window, interactor)

	// render function
	render := func() {
		gl.Clear(gl.DEPTH_BUFFER_BIT | gl.COLOR_BUFFER_BIT)
		if mesh != nil {
			gl.UseProgram(program)
			gl.Enable(gl.CULL_FACE)
			matrix := getMatrix(window, interactor.MeshInteractor, mesh)
			setMatrix(matrixUniform, matrix)
			mesh.Draw(positionAttrib)
		}
		{
			gl.UseProgram(planeProgram)
			gl.Disable(gl.CULL_FACE)
			matrix := getMatrix(window, interactor.PlaneInteractor, planeMesh)
			setMatrix(planeMatrixUniform, matrix)
			planeMesh.Draw(planePositionAttrib)
		}
		window.SwapBuffers()
	}

	// render during resize
	window.SetFramebufferSizeCallback(func(window *glfw.Window, w, h int) {
		render()
	})

	// handle drop events
	window.SetDropCallback(func(window *glfw.Window, filenames []string) {
		loadMesh(filenames[0], ch)
		window.SetTitle(filenames[0])
	})

	// main loop
	for !window.ShouldClose() {
		select {
		case data := <-ch:
			if mesh != nil {
				mesh.Destroy()
			}
			mesh = NewMesh(data)
			fmt.Printf("first frame at %.3f seconds\n", time.Since(start).Seconds())
		default:
		}
		render()
		glfw.PollEvents()
	}
}

func getMatrix(window *glfw.Window, interactor Interactor, mesh *Mesh) fauxgl.Matrix {
	matrix := mesh.Transform
	matrix = interactor.Matrix(window).Mul(matrix)
	return matrix
}
