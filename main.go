package main

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Dice Statistics Calculator")

	// Output display
	outputDisplay := canvas.NewText("0", color.White)
	outputDisplay.TextSize = 32

	// Dice input bar
	diceInputEntry := widget.NewEntry()
	diceInputEntry.SetPlaceHolder("e.g., 2d20H, 3d6+5")

	// Dice buttons for quick input
	diceButtonsContainer := container.NewVBox(
		container.NewHBox(
			widget.NewButton("d4", func() {
				diceInputEntry.SetText(diceInputEntry.Text + "d4")
			}),
			widget.NewButton("d6", func() {
				diceInputEntry.SetText(diceInputEntry.Text + "d6")
			}),
			widget.NewButton("d8", func() {
				diceInputEntry.SetText(diceInputEntry.Text + "d8")
			}),
			widget.NewButton("d10", func() {
				diceInputEntry.SetText(diceInputEntry.Text + "d10")
			}),
			widget.NewButton("d12", func() {
				diceInputEntry.SetText(diceInputEntry.Text + "d12")
			}),
		),
		container.NewHBox(
			widget.NewButton("d20", func() {
				diceInputEntry.SetText(diceInputEntry.Text + "d20")
			}),
			widget.NewButton("d100", func() {
				diceInputEntry.SetText(diceInputEntry.Text + "d100")
			}),
			widget.NewButton("dx", func() {
				diceInputEntry.SetText(diceInputEntry.Text + "dx")
			}),
			widget.NewButton("H", func() {
				diceInputEntry.SetText(diceInputEntry.Text + "H")
			}),
			widget.NewButton("L", func() {
				diceInputEntry.SetText(diceInputEntry.Text + "L")
			}),
		),
	)

	// Roll button
	rollButton := widget.NewButton("ROLL", func() {
		diceInput := strings.TrimSpace(diceInputEntry.Text)
		if diceInput == "" {
			outputDisplay.Text = "Error: Empty input"
			outputDisplay.Refresh()
			return
		}

		result, err := CalculateDice(diceInput)
		if err != nil {
			outputDisplay.Text = fmt.Sprintf("Error: %v", err)
		} else {
			outputDisplay.Text = fmt.Sprintf("%d", result)
		}
		outputDisplay.Refresh()
	})
	rollButton.Importance = widget.HighImportance

	// Clear button
	clearButton := widget.NewButton("CLEAR", func() {
		diceInputEntry.SetText("")
		outputDisplay.Text = "0"
		outputDisplay.Refresh()
	})

	// Calculator-style number buttons
	numberButtonsContainer := container.NewVBox(
		container.NewHBox(
			createCalcButton("7", diceInputEntry),
			createCalcButton("8", diceInputEntry),
			createCalcButton("9", diceInputEntry),
			widget.NewButton("+", func() {
				diceInputEntry.SetText(diceInputEntry.Text + "+")
			}),
			widget.NewButton("-", func() {
				diceInputEntry.SetText(diceInputEntry.Text + "-")
			}),
		),
		container.NewHBox(
			createCalcButton("4", diceInputEntry),
			createCalcButton("5", diceInputEntry),
			createCalcButton("6", diceInputEntry),
			widget.NewButton("*", func() {
				diceInputEntry.SetText(diceInputEntry.Text + "*")
			}),
			widget.NewButton("/", func() {
				diceInputEntry.SetText(diceInputEntry.Text + "/")
			}),
		),
		container.NewHBox(
			createCalcButton("1", diceInputEntry),
			createCalcButton("2", diceInputEntry),
			createCalcButton("3", diceInputEntry),
			clearButton,
			widget.NewButton("Backspace", func() {
				text := diceInputEntry.Text
				if len(text) > 0 {
					diceInputEntry.SetText(text[:len(text)-1])
				}
			}),
		),
		container.NewHBox(
			createCalcButton("0", diceInputEntry),
			widget.NewButton(".", func() {
				diceInputEntry.SetText(diceInputEntry.Text + ".")
			}),
			widget.NewLabel(""),
			widget.NewLabel(""),
			widget.NewLabel(""),
		),
	)

	// Output section
	outputSection := container.NewVBox(
		widget.NewLabel("Result:"),
		outputDisplay,
	)

	// Main layout: output at top, dice bar below, calculator buttons below that
	mainContent := container.NewVBox(
		outputSection,
		widget.NewCard("Dice Options", "", diceButtonsContainer),
		widget.NewCard("Calculator", "", numberButtonsContainer),
		widget.NewLabel("Dice Expression:"),
		diceInputEntry,
		rollButton,
	)

	scrollContainer := container.NewScroll(mainContent)
	myWindow.SetContent(scrollContainer)
	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.ShowAndRun()
}

func createCalcButton(label string, entry *widget.Entry) *widget.Button {
	return widget.NewButton(label, func() {
		entry.SetText(entry.Text + label)
	})
}
