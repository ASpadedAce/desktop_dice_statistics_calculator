package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type calculation struct {
	id        int
	equation  string
	diceRolls string
	result    string
}

type historyItemRenderer struct {
	item           *historyItem
	equationLabel  *widget.Label
	diceRollsLabel *widget.Label
	resultLabel    *widget.Label
	objects        []fyne.CanvasObject
}

func (r *historyItemRenderer) MinSize() fyne.Size {
	return r.objects[0].MinSize()
}

func (r *historyItemRenderer) Layout(size fyne.Size) {
	r.objects[0].Resize(size)
}

func (r *historyItemRenderer) ApplyTheme() {
}

func (r *historyItemRenderer) Refresh() {
	r.equationLabel.SetText(r.item.calc.equation)
	r.diceRollsLabel.SetText(r.item.calc.diceRolls)
	r.resultLabel.SetText(r.item.calc.result)
	r.equationLabel.Refresh()
	r.diceRollsLabel.Refresh()
	r.resultLabel.Refresh()
}

func (r *historyItemRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *historyItemRenderer) Destroy() {
}

type historyItem struct {
	widget.BaseWidget
	calc *calculation
}

func (h *historyItem) CreateRenderer() fyne.WidgetRenderer {
	h.ExtendBaseWidget(h)
	equationLabel := widget.NewLabel(h.calc.equation)
	equationLabel.TextStyle.Bold = true

	diceRollsLabel := widget.NewLabel(h.calc.diceRolls)
	diceRollsLabel.TextStyle.Italic = true

	resultLabel := widget.NewLabel(h.calc.result)
	resultLabel.Alignment = fyne.TextAlignTrailing

	layout := container.NewBorder(
		nil,
		nil,
		container.NewVBox(
			equationLabel,
			diceRollsLabel,
		),
		resultLabel,
	)

	return &historyItemRenderer{
		item:           h,
		equationLabel:  equationLabel,
		diceRollsLabel: diceRollsLabel,
		resultLabel:    resultLabel,
		objects:        []fyne.CanvasObject{layout},
	}
}

func (h *historyItem) SetCalculation(c *calculation) {
	h.calc = c
	h.Refresh()
}

func newHistoryItem(c *calculation) *historyItem {
	item := &historyItem{calc: c}
	item.ExtendBaseWidget(item)
	return item
}
