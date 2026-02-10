package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Dice Statistics Calculator")

	var calculations []*calculation
	var historyList *widget.List

	// Dice input bar
	diceInputEntry := widget.NewEntry()
	diceInputEntry.SetPlaceHolder("e.g., 2d20H, 3d6+5")

	historyList = widget.NewList(
		func() int {
			return len(calculations)
		},
		func() fyne.CanvasObject {
			return newHistoryItem(&calculation{})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*historyItem).SetCalculation(calculations[i])
		},
	)

	// Roll button
	rollButton := widget.NewButton("ROLL", func() {
		diceInput := strings.TrimSpace(diceInputEntry.Text)
		if diceInput == "" {
			return
		}

		result, diceRolls, err := CalculateDice(diceInput)
		if err != nil {
			// TODO: show error to user
		} else {
			id := 0
			if len(calculations) > 0 {
				id = calculations[len(calculations)-1].id + 1
			}
			c := &calculation{
				id:        id,
				equation:  diceInput,
				diceRolls: diceRolls,
				result:    fmt.Sprintf("= %s", strconv.FormatFloat(result, 'g', -1, 64)),
			}
			calculations = append([]*calculation{c}, calculations...)
			historyList.Refresh()
			diceInputEntry.SetText("")
		}
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

	buttonsContainer := container.New(newAspectRatioLayout(3.0/2.0, 7, 5))
	for _, button := range buttons {
		buttonsContainer.Add(button)
	}

	topContent := container.NewBorder(
		diceInputEntry,
		nil,
		nil,
		nil,
		historyList,
	)

	split := container.NewVSplit(topContent, buttonsContainer)
	split.Offset = 0.5 // Start with a 50/50 split

	myWindow.SetContent(split)
	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.ShowAndRun()
}

func createCalcButton(label string, entry *widget.Entry) *widget.Button {
	return widget.NewButton(label, func() {
		entry.SetText(entry.Text + label)
	})
}
