package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

type Operation interface {
	Do(t screen.Texture) (ready bool)
}

type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

type BgRect struct {
	X1, Y1, X2, Y2 float32
}

func (op *BgRect) Do(t screen.Texture) bool {
	// Нормалізація координат
	x1 := int(op.X1 * float32(t.Bounds().Dx()))
	y1 := int(op.Y1 * float32(t.Bounds().Dy()))
	x2 := int(op.X2 * float32(t.Bounds().Dx()))
	y2 := int(op.Y2 * float32(t.Bounds().Dy()))

	rect := image.Rect(x1, y1, x2, y2)
	t.Fill(rect, color.Black, screen.Src)
	return false
}

type TFigure struct {
	X, Y float32
}

func (op *TFigure) Do(t screen.Texture) bool {
	size := 200
	thickness := 50
	x := int(op.X * float32(t.Bounds().Dx()))
	y := int(op.Y * float32(t.Bounds().Dy()))

	// Горизонтальна частина "Т"
	rect1 := image.Rect(x-size/2, y-thickness/2, x+size/2, y+thickness/2)
	t.Fill(rect1, color.RGBA{B: 0xff, A: 0xff}, screen.Src)

	// Вертикальна частина "Т"
	rect2 := image.Rect(x-thickness/2, y-thickness/2, x+thickness/2, y+size/2)
	t.Fill(rect2, color.RGBA{B: 0xff, A: 0xff}, screen.Src)
	return false
}

type Move struct {
	X, Y float32
}

func (op *Move) Do(t screen.Texture) bool {
	// Для реалізації move потрібно буде зберігати стан фігур
	// Це буде реалізовано в State
	return false
}

type Reset struct{}

func (op *Reset) Do(t screen.Texture) bool {
	t.Fill(t.Bounds(), color.Black, screen.Src)
	return true
}
