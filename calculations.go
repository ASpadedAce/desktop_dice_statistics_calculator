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
func CalculateDice(expression string) (float64, string, error) {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return 0, "", fmt.Errorf("empty expression")
	}

	// Expand all dice notations to their rolled values
	expanded, diceRolls, err := expandDiceNotation(expression)
	if err != nil {
		return 0, "", err
	}

	// Evaluate the resulting mathematical expression
	result, err := evaluateMathExpression(expanded)
	if err != nil {
		return 0, "", err
	}

	return result, diceRolls, nil
}

// expandDiceNotation finds all dice notation in the expression and replaces them with rolled values
func expandDiceNotation(expression string) (string, string, error) {
	result := expression
	diceRollsStr := expression

	// Pattern to match dice notation: [H|L]?[count]d[sides][H|L]?
	// Examples: d20, 2d6, 3d6H, 4d8L, h2d20, l3d6, dx (where x is placeholder)
	dicePattern := regexp.MustCompile(`([HL])?(\d+)?d(\d+|x)([HL])?`)

	// Process all dice matches
	matches := dicePattern.FindAllStringSubmatchIndex(result, -1)

	// Process matches in reverse to maintain string indices
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		start := match[0]
		end := match[1]

		// Extract components
		prefixModifier := ""
		if match[2] != -1 {
			prefixModifier = result[match[2]:match[3]]
		}
		countStr := ""
		if match[4] != -1 {
			countStr = result[match[4]:match[5]]
		}
		sidesStr := result[match[6]:match[7]]
		suffixModifier := ""
		if match[8] != -1 {
			suffixModifier = result[match[8]:match[9]]
		}

		// Determine which modifier to use (priority: suffix > prefix)
		modifier := ""
		if suffixModifier != "" {
			modifier = suffixModifier
		} else if prefixModifier != "" {
			modifier = prefixModifier
		}

		// Determine count (default is 1)
		count := 1
		if countStr != "" {
			var err error
			count, err = strconv.Atoi(countStr)
			if err != nil || count <= 0 {
				return "", "", fmt.Errorf("invalid dice count: %s", countStr)
			}
		}

		// Determine sides
		var sides int
		if sidesStr == "x" {
			return "", "", fmt.Errorf("dx requires a number (e.g., d20). Please use a specific die like d20 or d100")
		}
		var err error
		sides, err = strconv.Atoi(sidesStr)
		if err != nil || sides <= 0 {
			return "", "", fmt.Errorf("invalid dice sides: %s", sidesStr)
		}

		// Roll the dice
		rolls := rollDiceSet(count, sides)
		var rollsStr []string
		for _, r := range rolls {
			rollsStr = append(rollsStr, strconv.Itoa(r))
		}

		// Apply modifier (H for highest, L for lowest)
		var value float64
		if modifier == "H" {
			if count == 1 {
				value = float64(rolls[0])
			} else {
				sort.Ints(rolls)
				value = float64(rolls[len(rolls)-1]) // Highest
			}
		} else if modifier == "L" {
			if count == 1 {
				value = float64(rolls[0])
			} else {
				sort.Ints(rolls)
				value = float64(rolls[0]) // Lowest
			}
		} else {
			// Sum all rolls
			value = 0
			for _, roll := range rolls {
				value += float64(roll)
			}
		}

		// Replace the dice notation with its value in the result string
		result = result[:start] + strconv.FormatFloat(value, 'f', -1, 64) + result[end:]
		diceRollsStr = diceRollsStr[:start] + fmt.Sprintf("(%dd%d: %s)", count, sides, strings.Join(rollsStr, ", ")) + diceRollsStr[end:]
	}

	return result, diceRollsStr, nil
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
func evaluateMathExpression(expr string) (float64, error) {
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
func (p *parser) parseExpression() (float64, error) {
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
func (p *parser) parseTerm() (float64, error) {
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
			left = left / right
		} else {
			break
		}
	}

	return left, nil
}

// parseFactor handles parentheses and unary operators (highest precedence)
func (p *parser) parseFactor() (float64, error) {
	left, err := p.parsePower()
	if err != nil {
		return 0, err
	}

	for {
		p.skipWhitespace()
		if p.pos >= len(p.expr) {
			break
		}

		if p.expr[p.pos] == '^' {
			p.pos++
			right, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			left = math.Pow(left, right)
		} else {
			break
		}
	}

	return left, nil
}

// parsePower handles parentheses and unary operators
func (p *parser) parsePower() (float64, error) {
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

// parseNumber parses a number (integer or float) from the expression
func (p *parser) parseNumber() (float64, error) {
	p.skipWhitespace()
	start := p.pos

	for p.pos < len(p.expr) && (isDigit(p.expr[p.pos]) || p.expr[p.pos] == '.') {
		p.pos++
	}

	if start == p.pos {
		if p.pos < len(p.expr) {
			return 0, fmt.Errorf("expected number at position %d, got '%c'", p.pos, p.expr[p.pos])
		}
		return 0, fmt.Errorf("expected number at end of expression")
	}

	numStr := p.expr[start:p.pos]
	num, err := strconv.ParseFloat(numStr, 64)
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
