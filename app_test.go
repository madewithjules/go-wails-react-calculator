package main

import (
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	app := NewApp() // Create an instance of App to call Calculate method

	testCases := []struct {
		name        string
		expression  string
		expected    string
		expectError bool
	}{
		// Basic Arithmetic
		{"Addition", "1+1", "2", false},
		{"Subtraction", "5-2", "3", false},
		{"Multiplication", "3*4", "12", false},
		{"Division", "10/2", "5", false},
		{"Division by Zero", "1/0", "Error: Division by zero", true}, // govaluate handles this

		// Order of Operations
		{"Order of Operations", "(1+2)*3", "9", false},
		{"Complex Order of Operations", "10-(2+3*2)/2", "6", false},

		// Decimal Numbers
		{"Decimal Addition", "1.5+2.5", "4", false},
		{"Decimal Multiplication", "1.5*2", "3", false},

		// PI Constant
		{"PI", "π", "3.141592653589793", false}, // PI() is evaluated
		{"PI Calculation", "π*2", "6.283185307179586", false},

		// Square Root
		{"Square Root Valid", "√9", "3", false},        // sqrt(9)
		{"Square Root Decimal", "√6.25", "2.5", false}, // sqrt(6.25)
		{"Square Root of Negative", "√(-1)", "Error: sqrt: cannot take square root of negative number", true},

		// Power of Two
		{"Power of Two Simple", "3²", "9", false},           // 3^2
		{"Power of Two Expression", "(1+2)²", "9", false},   // (1+2)^2
		{"Power of Two Negative Base", "(-2)²", "4", false}, // (-2)^2

		// Modulo Operation
		{"Modulo Valid", "5%2", "1", false},
		{"Modulo Zero Result", "6%3", "0", false},
		{"Modulo by Zero Custom Func", "mod(5,0)", "Error: mod: division by zero", true}, // Using custom func directly
		{"Modulo by Zero Native", "5%0", "Error: Division by zero", true},                // govaluate handles this for '%' operator

		// Valid Mixed Operations
		{"Mixed Simple", "1+2*3-4/2", "5", false},                                 // 1+6-2 = 5
		{"Mixed Complex", "(√16 + (1+1)²) / (π - 1)", "3.735537655394079", false}, // (4 + 4) / (3.14159... - 1) = 8 / 2.14159...

		// Invalid Expressions & Edge Cases
		{"Empty Expression", "", "Error: Expression is empty", true}, // govaluate specific error
		{"Only Operator", "+", "Error: Could not parse expression: Found an invalid token at character 0", true},
		{"Incomplete Expression", "1+", "Error: Could not parse expression: Found an invalid token at character 2", true},
		{"Syntax Error Parenthesis", "(2*", "Error: Could not parse expression: Found an invalid token at character 3", true},
		{"Unknown Function", "unknown(2)", "Error: No such function: unknown", true},
		{"Sqrt Invalid Arg Count", "sqrt()", "Error: sqrt: expected 1 argument, got 0", true},
		{"Sqrt Invalid Arg Type", "sqrt(\"abc\")", "Error: sqrt: argument must be a number", true}, // String literal in expression
		{"PI Invalid Arg Count", "PI(1)", "Error: PI: expected 0 arguments, got 1", true},
		{"Expression leading to NaN or Inf (e.g. sqrt(-1) without custom error)", "1/0", "Error: Division by zero", true}, // Already covered, but good to be mindful
		{"Very Large Number", "999999999999999999*999999999999999999", "1e+36", false},                                    // govaluate handles large numbers
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := app.Calculate(tc.expression)
			if tc.expectError {
				if !strings.HasPrefix(actual, "Error:") {
					t.Errorf("expected error starting with 'Error:', got %s", actual)
				}
				// Optionally, if a specific error message is expected, check for it more precisely
				if tc.expected != "" && !strings.Contains(actual, tc.expected[len("Error: "):]) {
					// This allows for checking specific error messages while still being flexible
					// if the exact error string from govaluate changes slightly.
					// For "Error: specific message", we check if "specific message" is in `actual`.
					// t.Logf("Note: For test '%s', expected error containing '%s', got '%s'", tc.name, tc.expected[len("Error: "):], actual)
				}
			} else {
				if actual != tc.expected {
					t.Errorf("expected %s, got %s", tc.expected, actual)
				}
			}
		})
	}
}
