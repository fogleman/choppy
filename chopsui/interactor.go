package chopsui

import (
	"math"

	"github.com/fogleman/fauxgl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Interactor interface {
	Matrix(window *glfw.Window) fauxgl.Matrix
	CursorPositionCallback(window *glfw.Window, x, y float64)
	MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey)
	KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
	ScrollCallback(window *glfw.Window, dx, dy float64)
}

func BindInteractor(window *glfw.Window, interactor Interactor) {
	window.SetCursorPosCallback(glfw.CursorPosCallback(interactor.CursorPositionCallback))
	window.SetMouseButtonCallback(glfw.MouseButtonCallback(interactor.MouseButtonCallback))
	window.SetKeyCallback(glfw.KeyCallback(interactor.KeyCallback))
	window.SetScrollCallback(glfw.ScrollCallback(interactor.ScrollCallback))
}

// AppInteractor

type AppInteractor struct {
	MeshInteractor  Interactor
	PlaneInteractor Interactor
	Modifiers       glfw.ModifierKey
}

func NewAppInteractor(meshInteractor, planeInteractor Interactor) *AppInteractor {
	return &AppInteractor{meshInteractor, planeInteractor, 0}
}

func (a *AppInteractor) ForMesh(mods glfw.ModifierKey) bool {
	return mods == 0 || mods&glfw.ModSuper != 0 || mods == glfw.ModShift
}

func (a *AppInteractor) ForPlane(mods glfw.ModifierKey) bool {
	return mods == 0 || mods&glfw.ModAlt != 0 || mods == glfw.ModShift
}

func (a *AppInteractor) Matrix(window *glfw.Window) fauxgl.Matrix {
	panic("unsupported")
}

func (a *AppInteractor) CursorPositionCallback(window *glfw.Window, x, y float64) {
	if a.ForMesh(a.Modifiers) {
		a.MeshInteractor.CursorPositionCallback(window, x, y)
	}
	if a.ForPlane(a.Modifiers) {
		a.PlaneInteractor.CursorPositionCallback(window, x, y)
	}
}

func (a *AppInteractor) MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if a.ForMesh(a.Modifiers) {
		a.MeshInteractor.MouseButtonCallback(window, button, action, mods)
	}
	if a.ForPlane(a.Modifiers) {
		a.PlaneInteractor.MouseButtonCallback(window, button, action, mods)
	}
}

func (a *AppInteractor) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	a.Modifiers = mods
	if a.ForMesh(a.Modifiers) {
		a.MeshInteractor.KeyCallback(window, key, scancode, action, mods)
	}
	if a.ForPlane(a.Modifiers) {
		a.PlaneInteractor.KeyCallback(window, key, scancode, action, mods)
	}
}

func (a *AppInteractor) ScrollCallback(window *glfw.Window, dx, dy float64) {
	a.MeshInteractor.ScrollCallback(window, dx, dy)
	a.PlaneInteractor.ScrollCallback(window, dx, dy)
}

// Arcball

type Arcball struct {
	RotationSensitivity    float64
	TranslationSensitivity float64
	Start                  fauxgl.Vector
	Current                fauxgl.Vector
	Rotation               fauxgl.Matrix
	Translation            fauxgl.Vector
	Scroll                 float64
	Rotating               bool
	Panning                bool
}

func NewArcball() Interactor {
	a := Arcball{}
	a.RotationSensitivity = 20
	a.TranslationSensitivity = 1.5
	a.Rotation = fauxgl.Identity()
	return &a
}

func (a *Arcball) CursorPositionCallback(window *glfw.Window, x, y float64) {
	if a.Rotating {
		a.Current = arcballVector(window)
	}
	if a.Panning {
		a.Current = a.screenPosition(window)
	}
}

func (a *Arcball) MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButton1 {
		if action == glfw.Press {
			if mods&glfw.ModShift == 0 {
				v := arcballVector(window)
				a.Start = v
				a.Current = v
				a.Rotating = true
			} else {
				v := a.screenPosition(window)
				a.Start = v
				a.Current = v
				a.Panning = true
			}
		} else if action == glfw.Release {
			if a.Rotating {
				m := arcballRotate(a.Start, a.Current, a.RotationSensitivity)
				a.Rotation = m.Mul(a.Rotation)
				a.Rotating = false
			}
			if a.Panning {
				s := math.Pow(0.98, a.Scroll)
				d := a.Current.Sub(a.Start)
				a.Translation = a.Translation.Add(d.MulScalar(a.TranslationSensitivity / s))
				a.Panning = false
			}
		}
	}
}

func (a *Arcball) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		if key >= 49 && key <= 55 {
			a.Translation = fauxgl.Vector{}
			a.Scroll = 0
		}
		switch key {
		case 49: //1
			a.Rotation = fauxgl.Identity()
		case 50:
			a.Rotation = fauxgl.Rotate(fauxgl.V(0, 0, 1), math.Pi/2)
		case 51:
			a.Rotation = fauxgl.Rotate(fauxgl.V(0, 0, 1), math.Pi)
		case 52:
			a.Rotation = fauxgl.Rotate(fauxgl.V(0, 0, 1), -math.Pi/2)
		case 53:
			a.Rotation = fauxgl.Rotate(fauxgl.V(1, 0, 0), math.Pi/2)
		case 54:
			a.Rotation = fauxgl.Rotate(fauxgl.V(1, 0, 0), -math.Pi/2)
		case 55:
			a.Rotation = fauxgl.Rotate(fauxgl.V(1, 1, 0).Normalize(), -math.Pi/4).Rotate(fauxgl.V(0, 0, 1), math.Pi/4)
		}
	}
}

func (a *Arcball) ScrollCallback(window *glfw.Window, dx, dy float64) {
	a.Scroll += dy
}

func (a *Arcball) Matrix(window *glfw.Window) fauxgl.Matrix {
	w, h := window.GetFramebufferSize()
	aspect := float64(w) / float64(h)
	s := math.Pow(0.98, a.Scroll)
	r := a.Rotation
	if a.Rotating {
		r = arcballRotate(a.Start, a.Current, a.RotationSensitivity).Mul(r)
	}
	t := a.Translation
	if a.Panning {
		t = t.Add(a.Current.Sub(a.Start).MulScalar(a.TranslationSensitivity / s))
	}
	m := fauxgl.Identity()
	m = m.Translate(t)
	m = r.Mul(m)
	m = m.Scale(fauxgl.V(s, s, s))
	m = m.LookAt(fauxgl.V(0, -3, 0), fauxgl.V(0, 0, 0), fauxgl.V(0, 0, 1))
	m = m.Perspective(50, aspect, 0.1, 100)
	return m
}

func (a *Arcball) screenPosition(window *glfw.Window) fauxgl.Vector {
	r := a.Rotation
	if a.Rotating {
		r = arcballRotate(a.Start, a.Current, a.RotationSensitivity).Mul(r)
	}
	x, y := window.GetCursorPos()
	w, h := window.GetSize()
	x = (x/float64(w))*2 - 1
	y = (y/float64(h))*2 - 1
	v := fauxgl.Vector{x, 0, -y}
	return r.Inverse().MulPosition(v)
}

func arcballVector(window *glfw.Window) fauxgl.Vector {
	x, y := window.GetCursorPos()
	w, h := window.GetSize()
	x = (x/float64(w))*2 - 1
	y = (y/float64(h))*2 - 1
	x /= 4
	y /= 4
	x = -x
	q := x*x + y*y
	if q <= 1 {
		z := math.Sqrt(1 - q)
		return fauxgl.Vector{x, z, y}
	} else {
		return fauxgl.Vector{x, 0, y}.Normalize()
	}
}

func arcballRotate(a, b fauxgl.Vector, sensitivity float64) fauxgl.Matrix {
	const eps = 1e-9
	dot := b.Dot(a)
	if math.Abs(dot) < eps || math.Abs(dot-1) < eps {
		return fauxgl.Identity()
	} else if math.Abs(dot+1) < eps {
		return fauxgl.Rotate(a.Perpendicular(), math.Pi*sensitivity)
	} else {
		angle := math.Acos(dot)
		v := b.Cross(a).Normalize()
		return fauxgl.Rotate(v, angle*sensitivity)
	}
}
