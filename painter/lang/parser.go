package lang

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"strconv"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

type Parser struct {
	uistate Uistate
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	p.uistate.ResetOperations()

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		cmdLine := scanner.Text()

		err := p.parse(cmdLine)
		if err != nil {
			return nil, err
		}
	}

	res := p.uistate.GetOperations()

	return res, nil
}

func (p *Parser) parse(cmdl string) error {
	words := strings.Split(cmdl, " ")
	command := words[0]

	switch command {
	case "white":
		_, err := checkForErrorsInParameters(words)
		if err != nil {
			return err
		}
		p.uistate.WhiteBackground()
	case "green":
		_, err := checkForErrorsInParameters(words)
		if err != nil {
			return err
		}
		p.uistate.GreenBackground()
	case "bgrect":
		parameters, err := checkForErrorsInParameters(words)
		if err != nil {
			return err
		}
		p.uistate.BackgroundRectangle(image.Point{X: parameters[0], Y: parameters[1]}, image.Point{X: parameters[2], Y: parameters[3]})
	case "figure":
		parameters, err := checkForErrorsInParameters(words)
		if err != nil {
			return err
		}
		p.uistate.AddTFigure(image.Point{X: parameters[0], Y: parameters[1]})
	case "move":
		parameters, err := checkForErrorsInParameters(words)
		if err != nil {
			return err
		}
		p.uistate.AddMoveOperation(parameters[0], parameters[1])
	case "reset":
		_, err := checkForErrorsInParameters(words)
		if err != nil {
			return err
		}
		p.uistate.ResetStateAndBackground()
	case "update":
		_, err := checkForErrorsInParameters(words)
		if err != nil {
			return err
		}
		p.uistate.SetUpdateOperation()
	default:
		return fmt.Errorf("invalid command %v", words[0])
	}
	return nil
}

func checkForErrorsInParameters(words []string) ([]int, error) {
	if len(words) == 0 {
		return nil, fmt.Errorf("got empty string as command")
	}
	var command = words[0]
	fmt.Println(command)
	fmt.Println(countsOfArguments[command])
	fmt.Println(len(words))
	if len(words) != (countsOfArguments[command] + 1) {
		return nil, fmt.Errorf("wrong number of arguments for '%v' command", words[0])
	}
	var params []int
	for _, param := range words[1:] {
		p, err := parseInt(param)
		if err != nil {
			return nil, fmt.Errorf("invalid parameter for '%s' command: '%s' is not a number", command, param)
		}
		params = append(params, p)
	}
	return params, nil
}

var countsOfArguments = map[string]int{
	"white":  0,
	"green":  0,
	"bgrect": 4,
	"figure": 2,
	"move":   2,
	"reset":  0,
	"update": 0,
}

func parseInt(s string) (int, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse float: %s", s)
	}
	return int(f * 800), nil
}
