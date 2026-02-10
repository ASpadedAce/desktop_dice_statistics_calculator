package main

import (
	"fmt"
	"image/color"
	"strconv"
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
			outputDisplay.Text = strconv.FormatFloat(result, 'g', -1, 64)
		}
		outputDisplay.Refresh()
	})
	rollButton.Importance = widget.HighImportance

	buttons := []fyne.CanvasObject{
		widget.NewButton("H", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "H")
		}),
		widget.NewButton("dX", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "dX")
		}),
		widget.NewButton("d4", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d4")
		}),
		widget.NewButton("d6", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d6")
		}),
		widget.NewButton("d8", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d8")
		}),
		widget.NewButton("L", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "L")
		}),
		widget.NewButton("d10", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d10")
		}),
		widget.NewButton("d12", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d12")
		}),
		widget.NewButton("d20", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d20")
		}),
		widget.NewButton("d100", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d100")
		}),
		createCalcButton("7", diceInputEntry),
		createCalcButton("8", diceInputEntry),
		createCalcButton("9", diceInputEntry),
		widget.NewButton("*", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "*")
		}),
		widget.NewButton("⌫", func() {
			text := diceInputEntry.Text
			if len(text) > 0 {
				diceInputEntry.SetText(text[:len(text)-1])
			}
		}),
		createCalcButton("4", diceInputEntry),
		createCalcButton("5", diceInputEntry),
		createCalcButton("6", diceInputEntry),
		widget.NewButton("/", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "/")
		}),
		widget.NewButton("C", func() {
			diceInputEntry.SetText("")
			outputDisplay.Text = "0"
			outputDisplay.Refresh()
		}),
		createCalcButton("1", diceInputEntry),
		createCalcButton("2", diceInputEntry),
		createCalcButton("3", diceInputEntry),
		widget.NewButton("+", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "+")
		}),
		widget.NewButton("📊", func() {
			// TODO: Implement statistics logic
		}),
		widget.NewButton(".", func() {
			diceInputEntry.SetText(diceInputEntry.Text + ".")
		}),
		createCalcButton("0", diceInputEntry),
		widget.NewButton("^", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "^")
		}),
		widget.NewButton("-", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "-")
		}),
		rollButton,
	}

	buttonsContainer := container.New(newAspectRatioLayout(2.0/3.0, 7, 5))
	for _, button := range buttons {
		buttonsContainer.Add(button)
	}

	// Output section
	outputSection := container.NewVBox(
		widget.NewLabel("Result:"),
		outputDisplay,
	)

	// Main layout: output at top, dice bar below, calculator buttons below that
	mainContent := container.NewBorder(
		outputSection,
		container.NewVBox(
			widget.NewLabel("Dice Expression:"),
			diceInputEntry,
		),
		nil,
		nil,
		buttonsContainer,
	)

	myWindow.SetContent(mainContent)
	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.ShowAndRun()
}

func createCalcButton(label string, entry *widget.Entry) *widget.Button {
	return widget.NewButton(label, func() {
		entry.SetText(entry.Text + label)
	})
}
