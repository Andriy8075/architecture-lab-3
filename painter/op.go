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
	FirstPoint  image.Point
	SecondPoint image.Point
}

func (op *BgRect) Do(t screen.Texture) bool {
	c := color.Black
	t.Fill(image.Rect(op.FirstPoint.X, op.FirstPoint.Y, op.SecondPoint.X, op.SecondPoint.Y), c, screen.Src)
	return false
}

type TFigure struct {
	X, Y int
}

func (op *TFigure) Do(t screen.Texture) bool {
	size := 200
	thickness := 50
	x := op.X
	y := op.Y

	// Горизонтальна частина "Т"
	rect1 := image.Rect(x-size/2, y-thickness/2, x+size/2, y+thickness/2)
	t.Fill(rect1, color.RGBA{B: 0xff, A: 0xff}, screen.Src)

	// Вертикальна частина "Т"
	rect2 := image.Rect(x-thickness/2, y-thickness/2, x+thickness/2, y+size/2)
	t.Fill(rect2, color.RGBA{B: 0xff, A: 0xff}, screen.Src)
	return false
}

type Move struct {
	X            int
	Y            int
	FiguresArray []*TFigure
}

func (op *Move) Do(t screen.Texture) bool {
	for i := range op.FiguresArray {
		op.FiguresArray[i].X += op.X
		op.FiguresArray[i].Y += op.Y
	}
	return false
}

func Reset(t screen.Texture) {
	t.Fill(t.Bounds(), color.Black, screen.Src)
}
