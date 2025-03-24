package lang

import (
	"strings"
	"testing"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []painter.Operation
	}{
		{
			name:  "white command",
			input: "white",
			expected: []painter.Operation{
				painter.OperationFunc(painter.WhiteFill),
			},
		},
		{
			name:  "green command",
			input: "green",
			expected: []painter.Operation{
				painter.OperationFunc(painter.GreenFill),
			},
		},
		{
			name:  "bgrect command",
			input: "bgrect 0.1 0.2 0.3 0.4",
			expected: []painter.Operation{
				&painter.BgRect{X1: 0.1, Y1: 0.2, X2: 0.3, Y2: 0.4},
			},
		},
		{
			name:  "figure command",
			input: "figure 0.5 0.5",
			expected: []painter.Operation{
				&painter.TFigure{X: 0.5, Y: 0.5},
			},
		},
		{
			name:  "reset command",
			input: "reset",
			expected: []painter.Operation{
				&painter.Reset{},
			},
		},
		{
			name:  "update command",
			input: "update",
			expected: []painter.Operation{
				painter.UpdateOp,
			},
		},
	}

	p := &Parser{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			ops, err := p.Parse(r)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, ops)
		})
	}
}
