package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Dice Statistics Calculator")

	var calculations []*calculation
	var historyList *widget.List

	// Dice input bar
	diceInputEntry := newCustomEntry(fyne.CurrentApp().Settings().Theme().Size(theme.SizeNameText) * 2)
	diceInputEntry.SetPlaceHolder("e.g., 2d20H, 3d6+5")

	historyList = widget.NewList(
		func() int {
			return len(calculations)
		},
		func() fyne.CanvasObject {
			return newHistoryItemWithCallback(&calculation{}, func(equation string) {
				diceInputEntry.SetText(equation)
			})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*historyItem).SetCalculation(calculations[i])
		},
	)

	roll := func() {
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
		}
	}

	// Roll button
	rollButton := newCustomButton2WithImportance("ROLL", widget.HighImportance, roll)

	diceInputEntry.OnSubmitted = func(s string) {
		roll()
	}

	buttons := []fyne.CanvasObject{
		newCustomButton2("H", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "H")
		}),
		newCustomButton2("dX", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d")
		}),
		newCustomButton2("d4", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d4")
		}),
		newCustomButton2("d6", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d6")
		}),
		newCustomButton2("d8", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d8")
		}),
		newCustomButton2("L", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "L")
		}),
		newCustomButton2("d10", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d10")
		}),
		newCustomButton2("d12", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d12")
		}),
		newCustomButton2("d20", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d20")
		}),
		newCustomButton2("d100", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "d100")
		}),
		createCalcButton("7", diceInputEntry),
		createCalcButton("8", diceInputEntry),
		createCalcButton("9", diceInputEntry),
		newCustomButton2("*", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "*")
		}),
		newCustomButton2("BKSP", func() {
			text := diceInputEntry.Text
			if len(text) > 0 {
				diceInputEntry.SetText(text[:len(text)-1])
			}
		}),
		createCalcButton("4", diceInputEntry),
		createCalcButton("5", diceInputEntry),
		createCalcButton("6", diceInputEntry),
		newCustomButton2("/", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "/")
		}),
		newCustomButton2("C", func() {
			diceInputEntry.SetText("")
		}),
		createCalcButton("1", diceInputEntry),
		createCalcButton("2", diceInputEntry),
		createCalcButton("3", diceInputEntry),
		newCustomButton2("+", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "+")
		}),
		newCustomButton2("📊", func() {
			diceInput := strings.TrimSpace(diceInputEntry.Text)
			if diceInput == "" {
				return
			}
			ShowStatisticsWindow(diceInput)
			roll()
		}),
		newCustomButton2(".", func() {
			diceInputEntry.SetText(diceInputEntry.Text + ".")
		}),
		createCalcButton("0", diceInputEntry),
		newCustomButton2("^", func() {
			diceInputEntry.SetText(diceInputEntry.Text + "^")
		}),
		newCustomButton2("-", func() {
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

func createCalcButton(label string, entry *customEntry) *customButton2 {
	return newCustomButton2(label, func() {
		entry.SetText(entry.Text + label)
	})
}
