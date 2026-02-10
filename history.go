package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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
	equationLabel  *customLabel
	diceRollsLabel *customLabel
	resultLabel    *customLabel
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
	r.equationLabel.Text = r.item.calc.equation
	r.diceRollsLabel.Text = r.item.calc.diceRolls
	r.resultLabel.Text = r.item.calc.result
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
	equationLabel := newCustomLabel(h.calc.equation, fyne.CurrentApp().Settings().Theme().Size(theme.SizeNameText)*2, fyne.TextStyle{Bold: true}, fyne.TextAlignLeading, theme.ForegroundColor())
	diceRollsLabel := newCustomLabel(h.calc.diceRolls, fyne.CurrentApp().Settings().Theme().Size(theme.SizeNameText)*1.5, fyne.TextStyle{Italic: true}, fyne.TextAlignLeading, theme.ForegroundColor())
	resultLabel := newCustomLabel(h.calc.result, fyne.CurrentApp().Settings().Theme().Size(theme.SizeNameText)*2, fyne.TextStyle{}, fyne.TextAlignTrailing, theme.ForegroundColor())

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
