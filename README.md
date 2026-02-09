# Desktop Dice Statistics Calculator

A native desktop application for calculating dice statistics for tabletop RPGs. Built with Go and the Fyne GUI library to support both Windows and Linux.

## Features

- **Dice Rolling**: Roll any combination of dice with notation like `2d20`, `3d6+5`, etc.
- **Highest/Lowest Selection**: Use `H` and `L` modifiers to take the highest or lowest result from multiple dice
  - Example: `2d20H` rolls two d20s and takes the highest
  - Example: `4d6L` rolls four d6s and takes the lowest
- **Calculator Functionality**: Perform arithmetic operations alongside dice rolls
  - Supports: `+`, `-`, `*`, `/`, and parentheses
  - Example: `2d6 + 5 * 3`
- **Standard Dice Support**: d4, d6, d8, d10, d12, d20, d100
- **Custom Dice**: Use `dx` to define custom dice (e.g., `d24`, `d30`)
- **Calculator-Style Interface**: Familiar button layout resembling a traditional calculator

## Building

### Prerequisites

- Go 1.25.7 or later
- [Fyne GUI library](https://fyne.io/) (automatically fetched by `go mod`)

### Build Instructions

```bash
cd desktop_dice_statistics_calculator
go build -o dice_calculator
```

This creates an executable named `dice_calculator` in the current directory.

### Running

```bash
./dice_calculator
```

On Windows, use:
```bash
dice_calculator.exe
```

## Application Structure

### main.go
The main application file containing all UI logic:
- Window setup and layout management
- Dice control bar with quick-access buttons (d4, d6, d8, d10, d12, d20, d100, dx, H, L)
- Calculator-style number pad and operation buttons
- Input field for dice expressions
- Output display for results
- Event handlers for button clicks

### calculations.go
The calculation engine containing all dice logic:
- **CalculateDice()**: Main entry point for evaluating dice expressions
- **expandDiceNotations()**: Parses and rolls dice notation (e.g., `2d20H`)
- **rollDiceSet()**: Rolls a specified number of dice with a given number of sides
- **evaluateMathExpression()**: Evaluates mathematical expressions with proper operator precedence
- **ExpressionParser**: Recursive descent parser handling +, -, *, /, and parentheses

## Usage Examples

### Basic Dice Rolls
- `d20` - Roll a single d20
- `2d6` - Roll two d6s and sum them
- `3d4` - Roll three d4s and sum them

### With Modifiers
- `2d20H` - Roll two d20s, take the highest (advantage in D&D 5e)
- `2d20L` - Roll two d20s, take the lowest (disadvantage in D&D 5e)
- `4d6L` - Roll four d6s, take the lowest (typical stat rolling method)

### With Arithmetic
- `2d6 + 5` - Roll 2d6 and add 5
- `3d6 + 2d4 + 3` - Multiple dice and modifiers
- `2d6 * 2` - Roll 2d6 and multiply by 2
- `(2d6 + 1) * 3` - Using parentheses for complex calculations

### Custom Dice
- `d24` - Roll a 24-sided die
- `3d30` - Roll three 30-sided dice
- `d100` - Roll a percentile die

## Interface Layout

```
┌─────────────────────────────────────┐
│         Result: [output]            │
├─────────────────────────────────────┤
│ Dice Options                        │
│ [d4] [d6] [d8] [d10] [d12]         │
│ [d20] [d100] [dx] [H] [L]          │
├─────────────────────────────────────┤
│ Calculator                          │
│ [7] [8] [9] [+] [-]                │
│ [4] [5] [6] [*] [/]                │
│ [1] [2] [3] [CLR] [BACKSPACE]      │
│ [0] [.] [ ] [ ] [ ]                │
├─────────────────────────────────────┤
│ Dice Expression:                    │
│ [input field showing: 2d20H]        │
│              [ROLL]                 │
└─────────────────────────────────────┘
```

## Supported Operations

### Dice Notation
- `[count]d[sides]` - Standard dice notation (count defaults to 1)
- `H` - Take highest result (when count > 1)
- `L` - Take lowest result (when count > 1)

### Arithmetic Operations
- `+` Addition
- `-` Subtraction
- `*` Multiplication
- `/` Division (integer division)
- `()` Parentheses for grouping

### Operator Precedence
1. Parentheses
2. Unary minus (e.g., `-5`)
3. Multiplication and Division (left-to-right)
4. Addition and Subtraction (left-to-right)

## Error Handling

The calculator provides error messages for:
- Invalid dice notation (e.g., `d0`, `0d6`)
- Division by zero
- Malformed expressions
- Unexpected characters

When an error occurs, it displays in the result field with an "Error: " prefix.

## Dependencies

- `fyne.io/fyne/v2` - Cross-platform GUI library
- Go standard library (math/rand, regexp, strconv, strings, time, etc.)

## Cross-Platform Support

This application is built with Fyne, which supports:
- **Linux** - All major distributions
- **Windows** - Windows 7 and later
- **macOS** - (Can be built, but not tested in this project scope)

To build for a different platform, use Go's cross-compilation flags:
```bash
GOOS=windows GOARCH=amd64 go build
GOOS=linux GOARCH=amd64 go build
```

## Future Enhancement Ideas

- History of recent rolls
- Statistics display (average, min, max for dice rolls)
- Save/load custom dice definitions
- Keyboard shortcuts for common operations
- Dark/light theme support
- Sound effects for dice rolls

## License

See LICENSE file for details.