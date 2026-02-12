package main

import (
	"fmt"
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
}

// CalculateDiceStatistics calculates the theoretical distribution of possible outcomes for a dice expression
func CalculateDiceStatistics(expression string) (*DiceStatistics, error) {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return nil, fmt.Errorf("empty expression")
	}

	// Parse the expression to extract terms
	terms, err := parseTerms(expression)
	if err != nil {
		return nil, err
	}

	// Calculate all possible outcomes and their frequencies
	outcomes := calculateOutcomeDistribution(terms)

	if len(outcomes) == 0 {
		return nil, fmt.Errorf("no valid outcomes for expression")
	}

	// Find min and max
	minVal := -1
	maxVal := -1
	totalCount := 0

	for value, count := range outcomes {
		totalCount += count
		if minVal == -1 || value < minVal {
			minVal = value
		}
		if maxVal == -1 || value > maxVal {
			maxVal = value
		}
	}

	// Calculate percentages
	percentages := make(map[int]float64)
	for value, count := range outcomes {
		percentages[value] = (float64(count) / float64(totalCount)) * 100
	}

	return &DiceStatistics{
		MinValue:    minVal,
		MaxValue:    maxVal,
		Results:     outcomes,
		Total:       totalCount,
		Percentages: percentages,
	}, nil
}

// Term represents a single term in the expression (dice roll or constant)
type Term struct {
	isDice   bool
	count    int    // number of dice
	sides    int    // sides per die
	modifier string // "" for sum, "H" for highest, "L" for lowest
	value    int    // constant value if not dice
	op       string // operation before this term: "+", "-"
}

// parseTerms parses a dice expression into terms
func parseTerms(expression string) ([]Term, error) {
	var terms []Term

	// Split by + and -, keeping the operators
	parts := regexp.MustCompile(`([+\-])`).Split(expression, -1)

	currentOp := "+"
	dicePattern := regexp.MustCompile(`^(\d*)d(\d+)([HL])?$`)

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if part == "" {
			continue
		}

		// Check if this is an operator
		if part == "+" || part == "-" {
			currentOp = part
			continue
		}

		// Try to match dice notation
		matches := dicePattern.FindStringSubmatch(part)
		if matches != nil {
			count := 1
			if matches[1] != "" {
				c, err := strconv.Atoi(matches[1])
				if err != nil {
					return nil, err
				}
				count = c
			}

			sides, err := strconv.Atoi(matches[2])
			if err != nil {
				return nil, err
			}

			if sides <= 0 || count <= 0 {
				return nil, fmt.Errorf("invalid dice: %dd%d", count, sides)
			}

			modifier := matches[3]

			terms = append(terms, Term{
				isDice:   true,
				count:    count,
				sides:    sides,
				modifier: modifier,
				op:       currentOp,
			})
			currentOp = "+"
		} else {
			// Try to parse as constant
			val, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid term: %s", part)
			}

			terms = append(terms, Term{
				isDice: false,
				value:  val,
				op:     currentOp,
			})
			currentOp = "+"
		}
	}

	return terms, nil
}

// calculateOutcomeDistribution calculates all possible outcomes and their frequencies
func calculateOutcomeDistribution(terms []Term) map[int]int {
	// Start with base case: single outcome of 0 with 1 way to achieve it
	outcomes := map[int]int{0: 1}

	for _, term := range terms {
		outcomes = applyTerm(outcomes, term)
	}

	return outcomes
}

// applyTerm applies a term to the current outcomes distribution
func applyTerm(currentOutcomes map[int]int, term Term) map[int]int {
	newOutcomes := make(map[int]int)

	if term.isDice {
		// Get all possible values for this dice roll
		diceOutcomes := getDiceOutcomes(term.count, term.sides, term.modifier)

		// Combine with current outcomes
		for currentVal, currentCount := range currentOutcomes {
			for diceVal, diceCount := range diceOutcomes {
				var resultVal int
				if term.op == "-" {
					resultVal = currentVal - diceVal
				} else {
					resultVal = currentVal + diceVal
				}

				newOutcomes[resultVal] += currentCount * diceCount
			}
		}
	} else {
		// Constant value
		for currentVal, currentCount := range currentOutcomes {
			var resultVal int
			if term.op == "-" {
				resultVal = currentVal - term.value
			} else {
				resultVal = currentVal + term.value
			}

			newOutcomes[resultVal] += currentCount
		}
	}

	return newOutcomes
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
