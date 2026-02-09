package main

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// init initializes the random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

// CalculateDice parses a dice expression and returns the result
// Supports formats like: 2d20, 3d6+5, 2d20H, 2d20L, 1d20+2d6, etc.
func CalculateDice(expression string) (int, error) {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return 0, fmt.Errorf("empty expression")
	}

	// Expand all dice notations to their rolled values
	expanded, err := expandDiceNotation(expression)
	if err != nil {
		return 0, err
	}

	// Evaluate the resulting mathematical expression
	result, err := evaluateMathExpression(expanded)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// expandDiceNotation finds all dice notation in the expression and replaces them with rolled values
func expandDiceNotation(expression string) (string, error) {
	result := expression

	// Pattern to match dice notation: [count]d[sides][H|L]
	// Examples: d20, 2d6, 3d6H, 4d8L, dx (where x is placeholder)
	dicePattern := regexp.MustCompile(`(\d+)?d(\d+|x)([HL])?`)

	// Process all dice matches
	matches := dicePattern.FindAllStringSubmatchIndex(result, -1)

	// Process matches in reverse to maintain string indices
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		start := match[0]
		end := match[1]

		// Extract components
		countStr := ""
		if match[2] != -1 {
			countStr = result[match[2]:match[3]]
		}
		sidesStr := result[match[4]:match[5]]
		modifier := ""
		if match[6] != -1 {
			modifier = result[match[6]:match[7]]
		}

		// Determine count (default is 1)
		count := 1
		if countStr != "" {
			var err error
			count, err = strconv.Atoi(countStr)
			if err != nil || count <= 0 {
				return "", fmt.Errorf("invalid dice count: %s", countStr)
			}
		}

		// Determine sides
		var sides int
		if sidesStr == "x" {
			return "", fmt.Errorf("dx requires a number (e.g., d20). Please use a specific die like d20 or d100")
		}
		var err error
		sides, err = strconv.Atoi(sidesStr)
		if err != nil || sides <= 0 {
			return "", fmt.Errorf("invalid dice sides: %s", sidesStr)
		}

		// Roll the dice
		rolls := rollDiceSet(count, sides)

		// Apply modifier (H for highest, L for lowest)
		var value int
		if modifier == "H" {
			if count == 1 {
				value = rolls[0]
			} else {
				sort.Ints(rolls)
				value = rolls[len(rolls)-1] // Highest
			}
		} else if modifier == "L" {
			if count == 1 {
				value = rolls[0]
			} else {
				sort.Ints(rolls)
				value = rolls[0] // Lowest
			}
		} else {
			// Sum all rolls
			value = 0
			for _, roll := range rolls {
				value += roll
			}
		}

		// Replace the dice notation with its value in the result string
		result = result[:start] + strconv.Itoa(value) + result[end:]
	}

	return result, nil
}

// rollDiceSet rolls count dice with the given number of sides
func rollDiceSet(count int, sides int) []int {
	rolls := make([]int, count)
	for i := 0; i < count; i++ {
		rolls[i] = rand.Intn(sides) + 1 // Results in 1 to sides inclusive
	}
	return rolls
}

// evaluateMathExpression evaluates a mathematical expression with +, -, *, /, and parentheses
// Uses a recursive descent parser to handle operator precedence
func evaluateMathExpression(expr string) (int, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return 0, fmt.Errorf("empty expression")
	}

	parser := &parser{expr: expr, pos: 0}
	result, err := parser.parseExpression()
	if err != nil {
		return 0, err
	}

	parser.skipWhitespace()
	if parser.pos < len(parser.expr) {
		return 0, fmt.Errorf("unexpected character at position %d: '%c'", parser.pos, parser.expr[parser.pos])
	}

	return result, nil
}

// parser is a simple recursive descent parser for mathematical expressions
type parser struct {
	expr string
	pos  int
}

// parseExpression handles addition and subtraction (lowest precedence)
func (p *parser) parseExpression() (int, error) {
	left, err := p.parseTerm()
	if err != nil {
		return 0, err
	}

	for {
		p.skipWhitespace()
		if p.pos >= len(p.expr) {
			break
		}

		if p.expr[p.pos] == '+' {
			p.pos++
			right, err := p.parseTerm()
			if err != nil {
				return 0, err
			}
			left = left + right
		} else if p.expr[p.pos] == '-' {
			p.pos++
			right, err := p.parseTerm()
			if err != nil {
				return 0, err
			}
			left = left - right
		} else {
			break
		}
	}

	return left, nil
}

// parseTerm handles multiplication and division (higher precedence)
func (p *parser) parseTerm() (int, error) {
	left, err := p.parseFactor()
	if err != nil {
		return 0, err
	}

	for {
		p.skipWhitespace()
		if p.pos >= len(p.expr) {
			break
		}

		if p.expr[p.pos] == '*' {
			p.pos++
			right, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			left = left * right
		} else if p.expr[p.pos] == '/' {
			p.pos++
			right, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			left = int(math.Floor(float64(left) / float64(right)))
		} else {
			break
		}
	}

	return left, nil
}

// parseFactor handles parentheses and unary operators (highest precedence)
func (p *parser) parseFactor() (int, error) {
	p.skipWhitespace()

	if p.pos >= len(p.expr) {
		return 0, fmt.Errorf("unexpected end of expression")
	}

	// Handle parentheses
	if p.expr[p.pos] == '(' {
		p.pos++
		result, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		p.skipWhitespace()
		if p.pos >= len(p.expr) || p.expr[p.pos] != ')' {
			return 0, fmt.Errorf("missing closing parenthesis")
		}
		p.pos++
		return result, nil
	}

	// Handle unary minus
	if p.expr[p.pos] == '-' {
		p.pos++
		value, err := p.parseFactor()
		if err != nil {
			return 0, err
		}
		return -value, nil
	}

	// Parse a number
	return p.parseNumber()
}

// parseNumber parses an integer from the expression
func (p *parser) parseNumber() (int, error) {
	p.skipWhitespace()
	start := p.pos

	for p.pos < len(p.expr) && isDigit(p.expr[p.pos]) {
		p.pos++
	}

	if start == p.pos {
		if p.pos < len(p.expr) {
			return 0, fmt.Errorf("expected number at position %d, got '%c'", p.pos, p.expr[p.pos])
		}
		return 0, fmt.Errorf("expected number at end of expression")
	}

	numStr := p.expr[start:p.pos]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", numStr)
	}

	return num, nil
}

// skipWhitespace skips over whitespace characters
func (p *parser) skipWhitespace() {
	for p.pos < len(p.expr) && (p.expr[p.pos] == ' ' || p.expr[p.pos] == '\t' || p.expr[p.pos] == '\n' || p.expr[p.pos] == '\r') {
		p.pos++
	}
}

// isDigit checks if a character is a digit
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
