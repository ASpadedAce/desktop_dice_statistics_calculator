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
	container      fyne.CanvasObject
	objects        []fyne.CanvasObject
}

func (r *historyItemRenderer) MinSize() fyne.Size {
	minSize := r.container.MinSize()
	// Ensure a minimum height so items don't overlap
	if minSize.Height < 90 {
		minSize = fyne.NewSize(minSize.Width, 90)
	}
	return minSize
}

func (r *historyItemRenderer) Layout(size fyne.Size) {
	r.container.Resize(size)
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
	r.container.Refresh()
}

func (r *historyItemRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *historyItemRenderer) Destroy() {
}

type historyItem struct {
	widget.BaseWidget
	calc     *calculation
	onTapped func(equation string)
}

func (h *historyItem) Tapped(*fyne.PointEvent) {
	if h.onTapped != nil && h.calc != nil {
		h.onTapped(h.calc.equation)
	}
}

func (h *historyItem) TappedSecondary(*fyne.PointEvent) {
}

func (h *historyItem) CreateRenderer() fyne.WidgetRenderer {
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
		container:      layout,
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

func newHistoryItemWithCallback(c *calculation, onTapped func(equation string)) *historyItem {
	item := &historyItem{calc: c, onTapped: onTapped}
	item.ExtendBaseWidget(item)
	return item
}
