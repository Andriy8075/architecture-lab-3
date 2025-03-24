package lang

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

type Parser struct{}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		cmdLine := scanner.Text()
		if cmdLine == "" {
			continue
		}

		fields := strings.Fields(cmdLine)
		if len(fields) == 0 {
			continue
		}

		cmd := fields[0]
		args := fields[1:]

		switch cmd {
		case "white":
			res = append(res, painter.OperationFunc(painter.WhiteFill))
		case "green":
			res = append(res, painter.OperationFunc(painter.GreenFill))
		case "update":
			res = append(res, painter.UpdateOp)
		case "bgrect":
			if len(args) != 4 {
				println("less than 4 args")
				continue
			}
			println("enough args")
			x1, err1 := strconv.ParseFloat(args[0], 32)
			y1, err2 := strconv.ParseFloat(args[1], 32)
			x2, err3 := strconv.ParseFloat(args[2], 32)
			y2, err4 := strconv.ParseFloat(args[3], 32)
			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				println("errors")
				continue
			}
			res = append(res, &painter.BgRect{
				X1: float32(x1),
				Y1: float32(y1),
				X2: float32(x2),
				Y2: float32(y2),
			})
		case "figure":
			if len(args) != 2 {
				continue
			}
			x, err1 := strconv.ParseFloat(args[0], 32)
			y, err2 := strconv.ParseFloat(args[1], 32)
			if err1 != nil || err2 != nil {
				continue
			}
			res = append(res, &painter.TFigure{
				X: float32(x),
				Y: float32(y),
			})
		case "move":
			if len(args) != 2 {
				continue
			}
			x, err1 := strconv.ParseFloat(args[0], 32)
			y, err2 := strconv.ParseFloat(args[1], 32)
			if err1 != nil || err2 != nil {
				continue
			}
			res = append(res, &painter.Move{
				X: float32(x),
				Y: float32(y),
			})
		case "reset":
			res = append(res, &painter.Reset{})
		}
	}

	return res, nil
}
