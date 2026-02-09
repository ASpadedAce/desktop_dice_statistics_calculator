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
