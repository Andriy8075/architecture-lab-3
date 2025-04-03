package lang

import (
	"fmt"
	"github.com/roman-mazur/architecture-lab-3/painter"
	"image"
)

type Uistate struct {
	backgroundColor     painter.Operation
	backgroundRectangle *painter.BgRect
	figuresArray        []*painter.TFigure
	moveOperations      []painter.Operation
	updateOperation     painter.Operation
}

func (u *Uistate) Reset() {
	u.backgroundColor = painter.OperationFunc(painter.WhiteFill)
	u.backgroundRectangle = nil
	u.figuresArray = nil
	u.moveOperations = nil
	u.updateOperation = nil
}

func (u *Uistate) GetOperations() []painter.Operation {
	var ops []painter.Operation

	if u.backgroundColor != nil {
		ops = append(ops, u.backgroundColor)
	}
	if u.backgroundRectangle != nil {
		ops = append(ops, u.backgroundRectangle)
	}
	if len(u.moveOperations) != 0 {
		ops = append(ops, u.moveOperations...)
		u.moveOperations = nil
	}
	if len(u.figuresArray) != 0 {
		for _, figure := range u.figuresArray {
			ops = append(ops, figure)
		}
	}
	if u.updateOperation != nil {
		ops = append(ops, u.updateOperation)
	}

	return ops
}

func (u *Uistate) ResetOperations() {
	if u.updateOperation != nil {
		u.updateOperation = nil
	}
}

func (u *Uistate) GreenBackground() {
	u.backgroundColor = painter.OperationFunc(painter.GreenFill)
}

func (u *Uistate) WhiteBackground() {
	u.backgroundColor = painter.OperationFunc(painter.WhiteFill)
}

func (u *Uistate) BackgroundRectangle(firstPoint image.Point, secondPoint image.Point) {
	u.backgroundRectangle = &painter.BgRect{
		FirstPoint:  firstPoint,
		SecondPoint: secondPoint,
	}
}

func (u *Uistate) AddTFigure(centralPoint image.Point) {
	fmt.Println("Додано фігуру в:", centralPoint.X, centralPoint.Y) // було додано для перевірки
	figure := painter.TFigure{
		X: centralPoint.X,
		Y: centralPoint.Y,
	}
	u.figuresArray = append(u.figuresArray, &figure)
}

func (u *Uistate) AddMoveOperation(x int, y int) {
	moveOp := painter.Move{X: x, Y: y, FiguresArray: u.figuresArray}
	u.moveOperations = append(u.moveOperations, &moveOp)
}

func (u *Uistate) ResetStateAndBackground() {
	u.Reset()
}

func (u *Uistate) SetUpdateOperation() {
	u.updateOperation = painter.UpdateOp
}
