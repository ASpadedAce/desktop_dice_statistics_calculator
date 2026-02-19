package main

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// DiceStatistics holds the theoretical statistics for a dice roll
type DiceStatistics struct {
	MinValue    int
	MaxValue    int
	Results     map[int]int     // outcome -> count of ways to achieve it
	Total       int             // total number of possible outcomes
	Percentages map[int]float64 // outcome -> percentage
	Average     float64         // average/mean value
	MostCommon  int             // most common (median) value
}

// Distribution represents the frequency distribution of outcomes
type Distribution map[int]int

// Regex patterns for parsing
var (
	diceTokenPattern = regexp.MustCompile(`^([HL])?(\d*)d(\d+)([HL])?`)
	// Updated numberTokenPattern to include optional decimal part
	numberTokenPattern = regexp.MustCompile(`^(\d+(\.\d+)?)`)
)

// CalculateDiceStatistics calculates the theoretical distribution of possible outcomes for a dice expression
func CalculateDiceStatistics(expression string) (*DiceStatistics, error) {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return nil, fmt.Errorf("empty expression")
	}

	parser := &statParser{expr: expression, pos: 0}
	outcomes, err := parser.parseExpression()
	if err != nil {
		return nil, err
	}

	parser.skipWhitespace()
	if parser.pos < len(parser.expr) {
		return nil, fmt.Errorf("unexpected character at position %d: '%c'", parser.pos, parser.expr[parser.pos])
	}

	if len(outcomes) == 0 {
		return nil, fmt.Errorf("no valid outcomes for expression")
	}

	// Find min and max
	minVal := 0
	maxVal := 0
	first := true
	totalCount := 0

	for value, count := range outcomes {
		totalCount += count
		if first {
			minVal = value
			maxVal = value
			first = false
		} else {
			if value < minVal {
				minVal = value
			}
			if value > maxVal {
				maxVal = value
			}
		}
	}

	// Calculate percentages
	percentages := make(map[int]float64)
	for value, count := range outcomes {
		percentages[value] = (float64(count) / float64(totalCount)) * 100
	}

	stats := &DiceStatistics{
		MinValue:    minVal,
		MaxValue:    maxVal,
		Results:     outcomes,
		Total:       totalCount,
		Percentages: percentages,
	}

	// Calculate average and most common value
	stats.calculateAverageAndMedian()

	return stats, nil
}

// statParser implementation
type statParser struct {
	expr string
	pos  int
}

func (p *statParser) skipWhitespace() {
	for p.pos < len(p.expr) && (p.expr[p.pos] == ' ' || p.expr[p.pos] == '\t') {
		p.pos++
	}
}

// parseExpression handles addition and subtraction
func (p *statParser) parseExpression() (Distribution, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
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
				return nil, err
			}
			left = addDist(left, right)
		} else if p.expr[p.pos] == '-' {
			p.pos++
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			left = subDist(left, right)
		} else {
			break
		}
	}

	return left, nil
}

// parseTerm handles multiplication, division and implicit multiplication
func (p *statParser) parseTerm() (Distribution, error) {
	left, err := p.parsePower()
	if err != nil {
		return nil, err
	}

	for {
		p.skipWhitespace()
		if p.pos >= len(p.expr) {
			break
		}

		c := p.expr[p.pos]
		if c == '*' {
			p.pos++
			right, err := p.parsePower()
			if err != nil {
				return nil, err
			}
			left = multDist(left, right)
		} else if c == '/' {
			p.pos++
			right, err := p.parsePower()
			if err != nil {
				return nil, err
			}
			left = divDist(left, right)
		} else if c == '(' || (c >= '0' && c <= '9') || c == 'd' || c == 'H' || c == 'L' {
			// Implicit multiplication for things that look like factors
			right, err := p.parsePower()
			if err != nil {
				return nil, err
			}
			left = multDist(left, right)
		} else {
			break
		}
	}

	return left, nil
}

// parsePower handles exponentiation
func (p *statParser) parsePower() (Distribution, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for {
		p.skipWhitespace()
		if p.pos >= len(p.expr) {
			break
		}

		if p.expr[p.pos] == '^' {
			p.pos++
			right, err := p.parseFactor() // Left-associative to match calculator
			if err != nil {
				return nil, err
			}
			left = powDist(left, right)
		} else {
			break
		}
	}

	return left, nil
}

// parseFactor handles parentheses, dice, and numbers
func (p *statParser) parseFactor() (Distribution, error) {
	p.skipWhitespace()
	if p.pos >= len(p.expr) {
		return nil, fmt.Errorf("unexpected end of expression")
	}

	// Parentheses
	if p.expr[p.pos] == '(' {
		p.pos++
		dist, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		p.skipWhitespace()
		if p.pos >= len(p.expr) || p.expr[p.pos] != ')' {
			return nil, fmt.Errorf("missing closing parenthesis")
		}
		p.pos++
		return dist, nil
	}

	// Try Dice Pattern
	remaining := p.expr[p.pos:]
	if loc := diceTokenPattern.FindStringIndex(remaining); loc != nil {
		token := remaining[loc[0]:loc[1]]
		p.pos += loc[1]
		return parseDiceToken(token)
	}

	// Try Number Pattern
	if loc := numberTokenPattern.FindStringIndex(remaining); loc != nil {
		token := remaining[loc[0]:loc[1]]
		p.pos += loc[1]
		// Parse as float then cast to int (truncate/floor) to handle buttons like "."
		valFloat, err := strconv.ParseFloat(token, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", token)
		}
		return Distribution{int(valFloat): 1}, nil
	}

	return nil, fmt.Errorf("unexpected character: %c", p.expr[p.pos])
}

func parseDiceToken(token string) (Distribution, error) {
	matches := diceTokenPattern.FindStringSubmatch(token)

	if matches != nil {
		// It is a dice expression
		prefixModifier := matches[1]
		countStr := matches[2]
		sidesStr := matches[3]
		suffixModifier := matches[4]

		count := 1
		if countStr != "" {
			c, err := strconv.Atoi(countStr)
			if err != nil {
				return nil, err
			}
			count = c
		}

		sides, err := strconv.Atoi(sidesStr)
		if err != nil {
			return nil, err
		}

		modifier := ""
		if suffixModifier != "" {
			modifier = suffixModifier
		} else if prefixModifier != "" {
			modifier = prefixModifier
		}

		return getDiceOutcomes(count, sides, modifier), nil
	}

	return nil, fmt.Errorf("invalid dice term: %s", token)
}

// Operations on Distributions

func addDist(a, b Distribution) Distribution {
	res := make(Distribution)
	for valA, countA := range a {
		for valB, countB := range b {
			res[valA+valB] += countA * countB
		}
	}
	return res
}

func subDist(a, b Distribution) Distribution {
	res := make(Distribution)
	for valA, countA := range a {
		for valB, countB := range b {
			res[valA-valB] += countA * countB
		}
	}
	return res
}

func multDist(a, b Distribution) Distribution {
	res := make(Distribution)
	for valA, countA := range a {
		for valB, countB := range b {
			res[valA*valB] += countA * countB
		}
	}
	return res
}

func divDist(a, b Distribution) Distribution {
	res := make(Distribution)
	for valA, countA := range a {
		for valB, countB := range b {
			if valB == 0 {
				continue // Division by zero yields no outcome
			}
			res[valA/valB] += countA * countB
		}
	}
	return res
}

func powDist(a, b Distribution) Distribution {
	res := make(Distribution)
	for valA, countA := range a {
		for valB, countB := range b {
			// Integer exponentiation
			// Standard behavior for non-negative exponents
			// Negative exponents with int base result in 0 (unless -1, 1).
			val := 0
			if valB >= 0 {
				val = int(math.Pow(float64(valA), float64(valB)))
			} else {
				// Integer division for 1/(a^-b) usually 0
				val = int(math.Pow(float64(valA), float64(valB)))
			}
			res[val] += countA * countB
		}
	}
	return res
}

// getDiceOutcomes returns a map of all possible outcomes for a dice roll and their frequencies
func getDiceOutcomes(count int, sides int, modifier string) map[int]int {
	outcomes := make(map[int]int)

	if modifier == "H" {
		// Keep only the highest die
		generateHighestOutcomes(count, sides, []int{}, outcomes)
	} else if modifier == "L" {
		// Keep only the lowest die
		generateLowestOutcomes(count, sides, []int{}, outcomes)
	} else {
		// Sum all dice
		generateSumOutcomes(count, sides, []int{}, outcomes)
	}

	return outcomes
}

// generateSumOutcomes recursively generates all sums
func generateSumOutcomes(remaining int, sides int, current []int, outcomes map[int]int) {
	if remaining == 0 {
		sum := 0
		for _, val := range current {
			sum += val
		}
		outcomes[sum]++
		return
	}

	for die := 1; die <= sides; die++ {
		generateSumOutcomes(remaining-1, sides, append(current, die), outcomes)
	}
}

// generateHighestOutcomes recursively generates all highest-die outcomes
func generateHighestOutcomes(remaining int, sides int, current []int, outcomes map[int]int) {
	if remaining == 0 {
		highest := 0
		for _, val := range current {
			if val > highest {
				highest = val
			}
		}
		outcomes[highest]++
		return
	}

	for die := 1; die <= sides; die++ {
		generateHighestOutcomes(remaining-1, sides, append(current, die), outcomes)
	}
}

// generateLowestOutcomes recursively generates all lowest-die outcomes
func generateLowestOutcomes(remaining int, sides int, current []int, outcomes map[int]int) {
	if remaining == 0 {
		lowest := sides + 1
		for _, val := range current {
			if val < lowest {
				lowest = val
			}
		}
		outcomes[lowest]++
		return
	}

	for die := 1; die <= sides; die++ {
		generateLowestOutcomes(remaining-1, sides, append(current, die), outcomes)
	}
}

// GetSortedOutcomes returns sorted unique outcomes
func (s *DiceStatistics) GetSortedOutcomes() []int {
	var outcomes []int
	for value := range s.Results {
		outcomes = append(outcomes, value)
	}
	sort.Ints(outcomes)
	return outcomes
}

// GetMaxPercentage returns the maximum percentage value
func (s *DiceStatistics) GetMaxPercentage() float64 {
	maxPercentage := 0.0
	for _, percentage := range s.Percentages {
		if percentage > maxPercentage {
			maxPercentage = percentage
		}
	}
	return maxPercentage
}

// calculateAverageAndMedian calculates the average and most common value
func (s *DiceStatistics) calculateAverageAndMedian() {
	if len(s.Results) == 0 {
		s.Average = 0
		s.MostCommon = 0
		return
	}

	// Calculate average (mean)
	sum := 0
	totalCount := 0
	for value, count := range s.Results {
		sum += value * count
		totalCount += count
	}
	s.Average = float64(sum) / float64(totalCount)

	// Find most common (mode) - the value with highest count
	maxCount := 0
	for value, count := range s.Results {
		if count > maxCount {
			maxCount = count
			s.MostCommon = value
		}
	}

	// If there are tied values, choose the smallest one
	if maxCount > 0 {
		for value, count := range s.Results {
			if count == maxCount && value < s.MostCommon {
				s.MostCommon = value
			}
		}
	}
}
