package main

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/Knetic/govaluate"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Calculate returns the result of a mathematical expression
func (a *App) Calculate(expression string) string {
	// Define custom functions
	functions := map[string]govaluate.ExpressionFunction{
		"sqrt": func(args ...interface{}) (interface{}, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("sqrt: expected 1 argument, got %d", len(args))
			}
			val, ok := args[0].(float64)
			if !ok {
				return nil, fmt.Errorf("sqrt: argument must be a number")
			}
			if val < 0 {
				return nil, fmt.Errorf("sqrt: cannot take square root of negative number")
			}
			return math.Sqrt(val), nil
		},
		"PI": func(args ...interface{}) (interface{}, error) {
			if len(args) != 0 {
				return nil, fmt.Errorf("PI: expected 0 arguments, got %d", len(args))
			}
			return math.Pi, nil
		},
		// govaluate supports ^ for power and % for modulo, so custom pow2 and mod might be redundant
		// but implementing as requested.
		"pow2": func(args ...interface{}) (interface{}, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("pow2: expected 1 argument, got %d", len(args))
			}
			val, ok := args[0].(float64)
			if !ok {
				return nil, fmt.Errorf("pow2: argument must be a number")
			}
			return math.Pow(val, 2), nil
		},
		"mod": func(args ...interface{}) (interface{}, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("mod: expected 2 arguments, got %d", len(args))
			}
			arg1, ok1 := args[0].(float64)
			arg2, ok2 := args[1].(float64)
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("mod: arguments must be numbers")
			}
			if arg2 == 0 {
				return nil, fmt.Errorf("mod: division by zero")
			}
			return math.Mod(arg1, arg2), nil
		},
	}

	// Pre-process the expression string
	processedExpression := strings.ReplaceAll(expression, "π", "PI()")
	processedExpression = strings.ReplaceAll(processedExpression, "√", "sqrt")
	// For x², replace "²" with "^2" to use govaluate's built-in power operator.
	// This is more robust than trying to match "x" for a "pow2(x)" function.
	// If a pow2 function is strictly required, it should be invoked as pow2(argument) by user.
	// The user is currently inputting "expression²" from frontend.
	// So we should replace "²" with "^2".
	processedExpression = strings.ReplaceAll(processedExpression, "²", "^2")
	// Frontend sends "%" for modulo, govaluate supports it directly.
	// If "mod(a,b)" was sent from frontend, it would work with the custom "mod" func.
	// Since frontend sends "a%b", no specific replacement for mod is needed here for govaluate.

	evaluableExpression, err := govaluate.NewEvaluableExpressionWithFunctions(processedExpression, functions)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}

	result, err := evaluableExpression.Evaluate(nil)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}

	return fmt.Sprintf("%v", result)
}
