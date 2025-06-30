package main

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
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

// TokenType defines the type of a token
type TokenType int

const (
	TokenNumber TokenType = iota
	TokenOperator
	TokenLeftParen
	TokenRightParen
	TokenFunction
	TokenComma
	TokenIdentifier // For function names before they are confirmed
)

// Token represents a token in the expression
type Token struct {
	Type  TokenType
	Value string
	Prec  int // Precedence for operators
	Arity int // Arity for functions
}

var knownFunctions = map[string]int{
	"sqrt": 1,
	"PI":   0,
	"pow2": 1,
	"mod":  2,
}

// preprocess handles textual replacements for special characters
func preprocess(expression string) string {
	s := expression
	s = strings.ReplaceAll(s, "π", "PI()")
	// Handle '√' followed by a number or parenthesized expression
	// Example: √9 -> sqrt(9), √(1+2) -> sqrt(1+2)
	// This regex looks for '√' followed by a number or a parenthesized expression
	reSqrt := regexp.MustCompile(`√\s*(\d+(\.\d+)?|\([^)]+\))`)
	s = reSqrt.ReplaceAllStringFunc(s, func(match string) string {
		arg := strings.TrimSpace(strings.TrimPrefix(match, "√"))
		return fmt.Sprintf("sqrt(%s)", arg)
	})
	// Fallback for simple sqrt if regex doesn't match (e.g. √ followed by variable name, not handled yet)
	// For now, this simpler replacement ensures "√" becomes "sqrt" if the regex didn't catch it.
	// The tokenizer will then expect "sqrt" followed by "(".
	if !strings.Contains(s,"sqrt("){
		s = strings.ReplaceAll(s, "√", "sqrt")
	}


	// Handle '²' (squared) applied to a number or parenthesized expression
	// Example: 3² -> pow2(3), (1+2)² -> pow2((1+2))
	// This regex looks for a number or parenthesized expression followed by '²'
	reSq := regexp.MustCompile(`(\d+(\.\d+)?|\([^)]+\))\s*²`)
	s = reSq.ReplaceAllStringFunc(s, func(match string) string {
		base := strings.TrimSpace(strings.TrimSuffix(match, "²"))
		return fmt.Sprintf("pow2(%s)", base)
	})

	return s
}

// Tokenizer converts an infix expression string into a list of tokens
func Tokenizer(expression string) ([]Token, error) {
	var tokens []Token
	var buffer strings.Builder // Used for numbers or identifiers (function names)

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])
		sChar := string(char)

		if (char >= '0' && char <= '9') || char == '.' {
			buffer.WriteString(sChar)
		} else if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			buffer.WriteString(sChar) // Start of an identifier
		} else {
			// End of number or identifier
			if buffer.Len() > 0 {
				val := buffer.String()
				buffer.Reset()
				if arity, isFunc := knownFunctions[val]; isFunc {
					tokens = append(tokens, Token{Type: TokenFunction, Value: val, Arity: arity})
				} else if _, err := strconv.ParseFloat(val, 64); err == nil {
					if strings.Count(val, ".") > 1 {
						return nil, fmt.Errorf("invalid number: %s", val)
					}
					tokens = append(tokens, Token{Type: TokenNumber, Value: val})
				} else {
					// Was not a known function and not a number. Could be an unknown identifier.
					// For this calculator, unknown identifiers are errors.
					return nil, fmt.Errorf("unknown identifier or function: %s", val)
				}
			}

			// Handle current character
			if strings.ContainsRune("+-*/", char) {
				prec := 1 // Default for +,-
				if char == '*' || char == '/' {
					prec = 2
				}
				if char == '-' {
					isUnary := (len(tokens) == 0) || (tokens[len(tokens)-1].Type == TokenOperator || tokens[len(tokens)-1].Type == TokenLeftParen || tokens[len(tokens)-1].Type == TokenComma)
					if isUnary {
						tokens = append(tokens, Token{Type: TokenNumber, Value: "0"})
					}
				}
				tokens = append(tokens, Token{Type: TokenOperator, Value: sChar, Prec: prec})
			} else if char == '(' {
				tokens = append(tokens, Token{Type: TokenLeftParen, Value: sChar})
			} else if char == ')' {
				tokens = append(tokens, Token{Type: TokenRightParen, Value: sChar})
			} else if char == ',' {
				tokens = append(tokens, Token{Type: TokenComma, Value: sChar})
			} else if char == ' ' {
				// Ignore whitespace
			} else {
				return nil, fmt.Errorf("invalid character: %s", sChar)
			}
		}
	}

	// After loop, check if buffer has anything left (number or identifier)
	if buffer.Len() > 0 {
		val := buffer.String()
		if arity, isFunc := knownFunctions[val]; isFunc {
			tokens = append(tokens, Token{Type: TokenFunction, Value: val, Arity: arity})
		} else if _, err := strconv.ParseFloat(val, 64); err == nil {
			if strings.Count(val, ".") > 1 {
				return nil, fmt.Errorf("invalid number: %s", val)
			}
			tokens = append(tokens, Token{Type: TokenNumber, Value: val})
		} else {
			return nil, fmt.Errorf("unknown identifier or function at end of expression: %s", val)
		}
	}

	return tokens, nil
}

// ShuntingYard converts infix tokens to RPN (postfix)
func ShuntingYard(tokens []Token) ([]Token, error) {
	var outputQueue []Token
	var operatorStack []Token

	for _, token := range tokens {
		switch token.Type {
		case TokenNumber:
			outputQueue = append(outputQueue, token)
		case TokenFunction:
			operatorStack = append(operatorStack, token)
		case TokenOperator:
			for len(operatorStack) > 0 {
				topOp := operatorStack[len(operatorStack)-1]
				if topOp.Type == TokenLeftParen || topOp.Type == TokenComma { // Comma check might be redundant here based on overall logic
					break
				}
				if topOp.Type == TokenFunction || (topOp.Type == TokenOperator && topOp.Prec >= token.Prec) {
					outputQueue = append(outputQueue, topOp)
					operatorStack = operatorStack[:len(operatorStack)-1]
				} else {
					break
				}
			}
			operatorStack = append(operatorStack, token)
		case TokenComma:
			foundLeftParenOrFunc := false
			for len(operatorStack) > 0 {
				topOp := operatorStack[len(operatorStack)-1]
				if topOp.Type == TokenLeftParen {
					foundLeftParenOrFunc = true
					break
				}
				outputQueue = append(outputQueue, topOp)
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if !foundLeftParenOrFunc {
				return nil, fmt.Errorf("misplaced comma or mismatched parentheses for function arguments")
			}
		case TokenLeftParen:
			operatorStack = append(operatorStack, token)
		case TokenRightParen:
			foundLeftParen := false
			for len(operatorStack) > 0 {
				topOp := operatorStack[len(operatorStack)-1]
				operatorStack = operatorStack[:len(operatorStack)-1]
				if topOp.Type == TokenLeftParen {
					foundLeftParen = true
					// If the token before the left parenthesis is a function, pop it to output.
					if len(operatorStack) > 0 && operatorStack[len(operatorStack)-1].Type == TokenFunction {
						outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
						operatorStack = operatorStack[:len(operatorStack)-1]
					}
					break
				}
				outputQueue = append(outputQueue, topOp)
			}
			if !foundLeftParen {
				return nil, fmt.Errorf("mismatched parentheses")
			}
		}
	}

	for len(operatorStack) > 0 {
		topOp := operatorStack[len(operatorStack)-1]
		if topOp.Type == TokenLeftParen || topOp.Type == TokenComma { // Comma should not be on stack here
			return nil, fmt.Errorf("mismatched parentheses or comma issue")
		}
		outputQueue = append(outputQueue, topOp)
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	return outputQueue, nil
}

// EvaluateRPN evaluates an RPN token queue
func EvaluateRPN(rpnTokens []Token) (float64, error) {
	var operandStack []float64

	for _, token := range rpnTokens {
		switch token.Type {
		case TokenNumber:
			num, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number in RPN: %s", token.Value)
			}
			operandStack = append(operandStack, num)
		case TokenOperator:
			if len(operandStack) < 2 {
				return 0, fmt.Errorf("syntax error (not enough operands for %s)", token.Value)
			}
			op2 := operandStack[len(operandStack)-1]
			op1 := operandStack[len(operandStack)-2]
			operandStack = operandStack[:len(operandStack)-2]
			var result float64
			switch token.Value {
			case "+": result = op1 + op2
			case "-": result = op1 - op2
			case "*": result = op1 * op2
			case "/":
				if op2 == 0 { return 0, fmt.Errorf("division by zero") }
				result = op1 / op2
			default: return 0, fmt.Errorf("unknown operator %s", token.Value)
			}
			operandStack = append(operandStack, result)
		case TokenFunction:
			if len(operandStack) < token.Arity {
				return 0, fmt.Errorf("not enough arguments for function %s (expected %d, got %d)", token.Value, token.Arity, len(operandStack))
			}
			args := make([]float64, token.Arity)
			for i := token.Arity - 1; i >= 0; i-- {
				args[i] = operandStack[len(operandStack)-1]
				operandStack = operandStack[:len(operandStack)-1]
			}
			var result float64
			var err error
			switch token.Value {
			case "sqrt":
				if args[0] < 0 { err = fmt.Errorf("sqrt: cannot take square root of negative number (%v)", args[0]) } else { result = math.Sqrt(args[0]) }
			case "PI":
				result = math.Pi
			case "pow2":
				result = math.Pow(args[0], 2)
			case "mod":
				if args[1] == 0 { err = fmt.Errorf("mod: division by zero (%v %% %v)", args[0], args[1]) } else { result = math.Mod(args[0], args[1]) }
			default:
				err = fmt.Errorf("unknown function %s", token.Value)
			}
			if err != nil { return 0, err }
			operandStack = append(operandStack, result)
		}
	}

	if len(operandStack) != 1 {
		// This could be due to mismatched parens not caught, or functions not eating args correctly, or just a number.
		// If input is just "5", RPN is ["5"], stack has [5]. This is valid.
		// If input is "5 5", RPN is ["5", "5"], stack has [5,5]. Invalid.
		// If input is "sin" (unknown func), tokenizer errors.
		// If input is "1 +", tokenizer is fine, SY might be fine, RPN eval fails on operator.
		// If input is "(1+2" (mismatched paren), SY should error.
		// If input is "1 2 3", tokenizer gives three numbers. SY gives [1,2,3]. Eval RPN leads to stack [1,2,3] -> error.
		if len(rpnTokens) > 1 && (len(rpnTokens) != 2*len(operandStack) -1 && tokenIsOperatorOrFunction(rpnTokens[len(rpnTokens)-1])) {
			// A more specific check or rely on the generic one.
		}
		return 0, fmt.Errorf("invalid expression (final stack size %d, expected 1)", len(operandStack))
	}
	return operandStack[0], nil
}
func tokenIsOperatorOrFunction(token Token) bool {
    return token.Type == TokenOperator || token.Type == TokenFunction
}

// Calculate returns the result of a mathematical expression
func (a *App) Calculate(expression string) string {
	if strings.TrimSpace(expression) == "" {
		return ""
	}

	preprocessedExpression := preprocess(expression)

	tokens, err := Tokenizer(preprocessedExpression)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}

	// Handle case like "PI" which becomes "PI()" -> tokens: PI, (, )
	// If only a number is entered, e.g. "5", tokens: 5. RPN: 5. Result: 5. This is fine.
	// If "PI()" is input, RPN: PI. Result: 3.14...
	if len(tokens) == 0 && strings.TrimSpace(preprocessedExpression) != "" {
		return "Error: Invalid expression (empty token list after preprocessing)"
	}
     if len(tokens) == 0 && strings.TrimSpace(preprocessedExpression) == ""{
		return ""
	}


	rpnTokens, err := ShuntingYard(tokens)
	if err != nil {
		return fmt.Sprintf("Error: %s (during shunting yard)", err.Error())
	}

    if len(rpnTokens) == 0 && strings.TrimSpace(preprocessedExpression) != "" {
		// This can happen for input "()" -> preprocessed "()" -> tokens LPAREN, RPAREN -> shunting yard empty.
		// EvaluateRPN will then correctly error "invalid RPN expression (stack not empty at end)" or similar.
		// Or if input is just "PI", preprocess makes "PI()", Shunting yard: PI, LPAR, RPAR -> RPN: PI. EvalRPN: PI -> result.
    }


	result, err := EvaluateRPN(rpnTokens)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}

	return strconv.FormatFloat(result, 'f', -1, 64)
}
