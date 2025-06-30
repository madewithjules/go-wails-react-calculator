package main

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	app := NewApp() // Create an instance of App to call Calculate method

	// Helper for PI for exact comparison
	piStr := fmt.Sprintf("%.15f", math.Pi) // Match precision of typical math.Pi string output
	// For strconv.FormatFloat(result, 'f', -1, 64), -1 precision means minimum digits.
	// Let's use a helper for float comparison or ensure expected strings are precise.
	// For simplicity, we'll use string comparison and be careful with expected values.
	// strconv.FormatFloat(math.Pi, 'f', -1, 64) is the reference for PI output.
	piCalcStr := strconv.FormatFloat(math.Pi, 'f', -1, 64)


	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic Arithmetic
		{"Addition", "1+1", "2"},
		{"Subtraction", "5-2", "3"},
		{"Multiplication", "3*4", "12"},
		{"Division", "10/2", "5"},
		{"Floating Point Addition", "1.5 + 2.5", "4"},
		{"Floating Point Subtraction", "3.7 - 1.2", "2.5"},
		{"Floating Point Multiplication", "1.5 * 2", "3"},
		{"Floating Point Division", "5.0 / 2.0", "2.5"},
		{"Zero Result", "1-1", "0"},

		// Operator Precedence
		{"Precedence 1", "2+3*4", "14"},    // 2 + 12 = 14
		{"Precedence 2", "10-4/2", "8"},    // 10 - 2 = 8
		{"Precedence 3", "2*3+4*5", "26"},  // 6 + 20 = 26
		{"Precedence 4", "10/2*5", "25"},   // 5 * 5 = 25 (left-to-right for same precedence)
		{"Precedence 5", "10-2+3", "11"},   // 8 + 3 = 11 (left-to-right for same precedence)

		// Parentheses
		{"Parentheses 1", "(2+3)*4", "20"},
		{"Parentheses 2", "10/(2+3)", "2"},
		{"Parentheses 3", "2*(3+4)/2", "7"},
		{"Nested Parentheses", "((1+1)*2+1)*2", "10"}, // ((2)*2+1)*2 = (4+1)*2 = 5*2 = 10

		// Unary Minus
		{"Unary Minus 1", "-5", "0"}, // Tokenizer makes this 0-5, but RPN results in -5. Let's check app.go logic.
		// The current Tokenizer for "-5" produces: Number(0), Operator(-), Number(5)
		// RPN: Number(0), Number(5), Operator(-) -> 0-5 = -5. Correct.
		{"Unary Minus Test Actual -5", "-5", "-5"},
		{"Unary Minus 2", "-5+10", "5"},
		{"Unary Minus 3", "10+-2", "8"}, // 10 + (0-2)
		{"Unary Minus 4", "10*-2", "-20"}, // 10 * (0-2)
		{"Unary Minus Parentheses", "-(2+3)", "-5"},
		{"Unary Minus Complex", "-sqrt(4)+1", "-1"}, // -(2)+1 = -1

		// Custom Functions & Preprocessing
		// SQRT and √
		{"sqrt Valid", "sqrt(9)", "3"},
		{"sqrt Decimal", "sqrt(6.25)", "2.5"},
		{"sqrt Zero", "sqrt(0)", "0"},
		{"sqrt Preprocessed", "√9", "3"}, // Preprocessor: √9 -> sqrt(9)
		{"sqrt Preprocessed Decimal", "√6.25", "2.5"},
		{"sqrt Preprocessed Complex Arg", "√(1+3)", "2"}, // Preprocessor: √(1+3) -> sqrt(1+3)
		{"sqrt Error Negative", "sqrt(-1)", "Error: sqrt: cannot take square root of negative number (-1)"},
		{"sqrt Preprocessed Error Negative", "√-1", "Error: sqrt: cannot take square root of negative number (-1)"}, // Preprocessor makes sqrt(-1)
		{"Expression with sqrt", "1+sqrt(16)", "5"},
		{"Expression with √", "1+√16", "5"},
		{"Nested sqrt", "sqrt(sqrt(81))", "3"},

		// PI and π
		{"PI Function", "PI()", piCalcStr},
		{"PI Preprocessed", "π", piCalcStr}, // Preprocessor: π -> PI()
		{"PI Calculation", "2*PI()", strconv.FormatFloat(2*math.Pi, 'f', -1, 64)},
		{"PI Preprocessed Calculation", "2*π", strconv.FormatFloat(2*math.Pi, 'f', -1, 64)},
		{"PI with other ops", "PI()/2", strconv.FormatFloat(math.Pi/2, 'f', -1, 64)},

		// POW2 and ²
		{"pow2 Valid", "pow2(3)", "9"},
		{"pow2 Negative Base", "pow2(-2)", "4"},
		{"pow2 Decimal Base", "pow2(1.5)", "2.25"},
		{"pow2 Zero Base", "pow2(0)", "0"},
		{"pow2 Preprocessed", "3²", "9"},       // Preprocessor: 3² -> pow2(3)
		{"pow2 Preprocessed Negative", "(-2)²", "4"}, // Preprocessor: (-2)² -> pow2((-2))
		{"pow2 Preprocessed Decimal", "1.5²", "2.25"},
		// Behavior of -3²: My preprocessor for `X²` tries to capture X.
		// If X is "3", then "-3²" becomes "-pow2(3)". Result: -9.
		// If X is "-3", then "(-3)²" becomes "pow2(-3)". Result: 9.
		// Test the "-3²" case based on current preprocessor (likely "-pow2(3)")
		{"pow2 Negative vs Square", "-3²", "-9"}, // Expects -(pow2(3)) due to preprocessing of `3²` part first.
		{"Expression with pow2", "1+pow2(4)", "17"},
		{"Expression with ²", "1+4²", "17"},

		// MOD
		{"mod Valid", "mod(5,2)", "1"},
		{"mod Valid Negative", "mod(-5,2)", "-1"}, // math.Mod behavior
		{"mod Valid Negative Divisor", "mod(5,-2)", "1"}, // math.Mod behavior
		{"mod Decimal", "mod(7.5, 2.2)", strconv.FormatFloat(math.Mod(7.5, 2.2), 'f', -1, 64)}, // 0.9
		{"mod Zero Result", "mod(6,3)", "0"},
		{"mod Error Div By Zero", "mod(5,0)", "Error: mod: division by zero (5 % 0)"},

		// Complex Expressions
		{"Complex 1", "sqrt(pow2(3)+7)", "4"},                      // sqrt(9+7) = sqrt(16) = 4
		{"Complex 2", "2*π+1", strconv.FormatFloat(2*math.Pi+1, 'f', -1, 64)}, // Using π
		{"Complex 3", "(sqrt(9)+pow2(2))/mod(7,4)", "2.3333333333333335"}, // (3+4)/3 = 7/3. Careful with float precision.
		{"Complex 4", "-π + (10/2)*3", strconv.FormatFloat(-math.Pi+(10/2)*3, 'f', -1, 64)}, // -pi + 15
		{"Complex 5", "100/pow2(sqrt(4)+sqrt(1))", "11.11111111111111"}, // 100 / pow2(2+1) = 100 / pow2(3) = 100/9

		// Error Handling
		{"Error Empty Expression", "", ""}, // Behavior for empty string
		{"Error Only Spaces", "   ", ""},   // Behavior for only spaces
		{"Error Invalid Character", "1a+2", "Error: unknown identifier or function: a"}, // Tokenizer error
		{"Error Invalid Character 2", "$+1", "Error: invalid character: $"},
		{"Error Mismatched Parentheses 1", "(2+3", "Error: mismatched parentheses (during shunting yard)"},
		{"Error Mismatched Parentheses 2", "2+3)", "Error: mismatched parentheses (during shunting yard)"},
		{"Error Mismatched Parentheses 3", "((1+2)", "Error: mismatched parentheses (during shunting yard)"},
		{"Error Syntax Invalid Operator Chain", "1+*2", "Error: syntax error (not enough operands for *)"}, // RPN eval error
		{"Error Syntax Leading Operator", "*2+1", "Error: syntax error (not enough operands for *)"},
		{"Error Syntax Trailing Operator", "1+2*", "Error: syntax error (not enough operands for *)"},
		{"Error Insufficient Operands", "1 2", "Error: invalid expression (final stack size 2, expected 1)"}, // RPN eval error
		{"Error Unknown Function", "unknown(1)", "Error: unknown identifier or function: unknown"},
		{"Error sqrt Arity 0", "sqrt()", "Error: not enough arguments for function sqrt (expected 1, got 0)"},
		{"Error sqrt Arity 2", "sqrt(1,2)", "Error: invalid RPN expression (stack not empty at end)"}, // This will parse as sqrt, (, 1, ,, 2, ). SY: 1, 2, sqrt. Eval: depends on comma handling.
		// Let's refine sqrt(1,2). Tokenizer: sqrt, (, 1, ,, 2, ). SY: 1, 2, sqrt. Eval: sqrt takes 1 arg (2), leaves 1 on stack. Error.
		{"Error mod Arity 1", "mod(1)", "Error: not enough arguments for function mod (expected 2, got 1)"},
		{"Error Division by Zero", "1/0", "Error: division by zero"},
		{"Error Invalid Number Format", "1.2.3+4", "Error: invalid number: 1.2.3"},
		{"Error Number Then Paren", "5(1+1)", "Error: invalid character: ("}, // Needs explicit operator
		{"Error Paren Then Number", "(1+1)5", "Error: invalid character: 5"}, // Needs explicit operator
		{"Error Func No Parens", "sqrt 9", "Error: unknown identifier or function: sqrt9"}, // Assuming 'sqrt9' if no space, or 'sqrt' then '9'
		// If "sqrt 9": Tokenizer (sqrt, 9). ShuntingYard (sqrt, 9) -> RPN (9, sqrt). Error in Eval (sqrt needs 1 from stack, gets 9, ok).
		// This depends on how Tokenizer handles "sqrt 9". If it's TokenFunc("sqrt"), TokenNum("9"), then RPN eval is fine.
		// My tokenizer based on `(char >= 'a' && char <= 'z')` and then flushing buffer would make `sqrt` a token, then `9` a token.
		// This would evaluate `sqrt(9)` correctly.
		{"Test Func Space Number", "sqrt 9", "3"}, // This should work if tokenizer separates function name and number
		{"Test Func Space Number Mod", "mod 10 3", "1"}, // mod(10,3)

		// Edge cases from previous tests
		{"Edge Complex Div Zero", "10/(2-2)", "Error: division by zero"},
		{"Edge PI Div Zero", "PI()/0", "Error: division by zero"},
		{"Edge Sqrt Div By Zero In Arg", "sqrt(10/0)", "Error: division by zero"}, // Division by zero happens before sqrt
		{"Edge Very Large Number (Overflow not explicitly handled, relies on float64)", "999999999*999999999", "999999998000000000"}, // Relies on float64 precision/range
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := app.Calculate(tc.input)
			// For floating point results, direct string comparison can be tricky.
			// However, strconv.FormatFloat with 'f' and -1 precision is quite deterministic.
			// If we hit issues, we might need a custom comparison for floats.
			if strings.HasPrefix(tc.expected, "Error:") {
				if !strings.HasPrefix(actual, "Error:") {
					t.Errorf("input: \"%s\": expected error string \"%s\", got \"%s\"", tc.input, tc.expected, actual)
				} else if actual != tc.expected {
					// Check if the core message is there, useful if error details (like specific numbers) vary slightly but are acceptable.
					// For now, demand exact match for errors too.
					t.Errorf("input: \"%s\": expected error string \"%s\", got \"%s\"", tc.input, tc.expected, actual)
				}
			} else if actual != tc.expected {
				t.Errorf("input: \"%s\": expected \"%s\", got \"%s\"", tc.input, tc.expected, actual)
			}
		})
	}
}
