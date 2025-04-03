package lang

import (
	"fmt"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedOps int
		expectError bool
	}{
		{
			name:        "white command",
			input:       "white\n",
			expectedOps: 1,
		},
		{
			name:        "green command",
			input:       "green\n",
			expectedOps: 1,
		},
		{
			name:        "bgrect command",
			input:       "bgrect 0.1 0.2 0.3 0.4\n",
			expectedOps: 1,
		},
		{
			name:        "figure command",
			input:       "figure 0.5 0.5\n",
			expectedOps: 1,
		},
		{
			name:        "move command",
			input:       "move 0.1 0.1\n",
			expectedOps: 1,
		},
		{
			name:        "reset command",
			input:       "reset\n",
			expectedOps: 1,
		},
		{
			name:        "update command",
			input:       "update\n",
			expectedOps: 1,
		},
		{
			name:        "multiple commands",
			input:       "white\nfigure 0.5 0.5\nupdate\n",
			expectedOps: 3,
		},
		{
			name:        "invalid command",
			input:       "invalid\n",
			expectError: true,
		},
		{
			name:        "bgrect with wrong params",
			input:       "bgrect 0.1 0.2\n",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &Parser{}
			reader := strings.NewReader(tc.input)
			ops, err := parser.Parse(reader)

			if tc.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(ops) != tc.expectedOps {
				t.Errorf("Expected %d operations, got %d\ntest name: %s", tc.expectedOps, len(ops), tc.name)
			}
		})
	}
}

func TestCheckForErrorsInParameters(t *testing.T) {
	tests := []struct {
		name        string
		input       []string
		expectedLen int
		expectError bool
	}{
		{
			name:        "without parameters",
			input:       []string{"green"},
			expectedLen: 0,
		},
		{
			name:        "valid parameters",
			input:       []string{"bgrect", "0.1", "0.2", "0.3", "0.4"},
			expectedLen: 4,
		},
		{
			name:        "wrong number of parameters",
			input:       []string{"bgrect", "0.1", "0.2"},
			expectError: true,
		},
		{
			name:        "invalid parameter",
			input:       []string{"bgrect", "0.1", "abc", "0.3", "0.4"},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("name: %s\n", tc.name)
			params, err := checkForErrorsInParameters(tc.input)
			fmt.Print("params: ", params)
			fmt.Print("\nerror: ", err)
			fmt.Print("\n")
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none. Test name: %s", tc.name)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(params) != tc.expectedLen {
				t.Errorf("Expected %d parameters, got %d", tc.expectedLen, len(params))
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{
		{"0.1", false},
		{"0.5", false},
		{"1.0", false},
		{"abc", true},
		{"", true},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			_, err := parseInt(tc.input)

			if tc.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
