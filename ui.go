package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// customLabel is a custom label widget with configurable text styling
type customLabel struct {
	widget.BaseWidget
	Text      string
	TextSize  float32
	Alignment fyne.TextAlign
	TextStyle fyne.TextStyle
	Color     color.Color
}

func (l *customLabel) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(l.Text, l.Color)
	text.TextSize = l.TextSize
	text.Alignment = l.Alignment
	text.TextStyle = l.TextStyle
	return &customLabelRenderer{
		text:    text,
		label:   l,
		objects: []fyne.CanvasObject{text},
	}
}

type customLabelRenderer struct {
	text    *canvas.Text
	label   *customLabel
	objects []fyne.CanvasObject
}

func (r *customLabelRenderer) Layout(size fyne.Size) {
	r.text.Resize(size)
}

func (r *customLabelRenderer) MinSize() fyne.Size {
	return r.text.MinSize()
}

func (r *customLabelRenderer) Refresh() {
	r.text.Text = r.label.Text
	r.text.TextSize = r.label.TextSize
	r.text.Alignment = r.label.Alignment
	r.text.TextStyle = r.label.TextStyle
	r.text.Color = r.label.Color
	r.text.Refresh()
}

func (r *customLabelRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *customLabelRenderer) Destroy() {
}

func newCustomLabel(text string, size float32, style fyne.TextStyle, align fyne.TextAlign, textColor color.Color) *customLabel {
	l := &customLabel{
		Text:      text,
		TextSize:  size,
		TextStyle: style,
		Alignment: align,
		Color:     textColor,
	}
	l.ExtendBaseWidget(l)
	return l
}

// customButton2 is a custom button widget with theme-aware styling
type customButton2 struct {
	widget.BaseWidget
	Text       string
	OnTapped   func()
	Importance widget.Importance
}

func (b *customButton2) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(b.Text, b.textColor())
	text.Alignment = fyne.TextAlignCenter
	background := canvas.NewRectangle(color.Transparent)
	return &customButton2Renderer{
		text:       text,
		background: background,
		button:     b,
		objects:    []fyne.CanvasObject{background, text},
	}
}

func (b *customButton2) Tapped(*fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *customButton2) TappedSecondary(*fyne.PointEvent) {
}

func (b *customButton2) MouseIn(*desktop.MouseEvent) {
}

func (b *customButton2) MouseOut() {
}

func (b *customButton2) MouseMoved(*desktop.MouseEvent) {
}

type customButton2Renderer struct {
	text       *canvas.Text
	background *canvas.Rectangle
	button     *customButton2
	objects    []fyne.CanvasObject
}

func (r *customButton2Renderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.text.Resize(size)
	newTextSize := size.Height * 0.4
	if size.Width < newTextSize {
		newTextSize = size.Width * 0.4
	}
	r.text.TextSize = newTextSize
}

func (r *customButton2Renderer) MinSize() fyne.Size {
	return r.text.MinSize()
}

func (r *customButton2Renderer) Refresh() {
	r.text.Text = r.button.Text
	r.text.Color = r.button.textColor()
	r.text.Refresh()
	r.background.FillColor = r.button.backgroundColor()
	r.background.Refresh()
}

func (r *customButton2Renderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *customButton2Renderer) Destroy() {
}

func (b *customButton2) backgroundColor() color.Color {
	variant := fyne.CurrentApp().Settings().ThemeVariant()
	switch b.Importance {
	case widget.HighImportance:
		return fyne.CurrentApp().Settings().Theme().Color(theme.ColorNamePrimary, variant)
	default:
		return fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameButton, variant)
	}
}

func (b *customButton2) textColor() color.Color {
	variant := fyne.CurrentApp().Settings().ThemeVariant()
	switch b.Importance {
	case widget.HighImportance:
		return fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameBackground, variant)
	default:
		return fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameForeground, variant)
	}
}

func newCustomButton2(label string, onTap func()) *customButton2 {
	b := &customButton2{
		Text:     label,
		OnTapped: onTap,
	}
	b.ExtendBaseWidget(b)
	return b
}

func newCustomButton2WithImportance(label string, importance widget.Importance, onTap func()) *customButton2 {
	b := &customButton2{
		Text:       label,
		OnTapped:   onTap,
		Importance: importance,
	}
	b.ExtendBaseWidget(b)
	return b
}

// customEntry is a custom entry widget with configurable text size
type customEntry struct {
	widget.Entry
	TextSize float32
}

func (e *customEntry) CreateRenderer() fyne.WidgetRenderer {
	// Call the parent's CreateRenderer to ensure proper initialization
	renderer := e.Entry.CreateRenderer()

	// Wrap the parent renderer to add our custom TextSize
	return &customEntryRenderer{
		entry:          e,
		parentRenderer: renderer,
	}
}

type customEntryRenderer struct {
	entry          *customEntry
	parentRenderer fyne.WidgetRenderer
}

func (r *customEntryRenderer) Layout(size fyne.Size) {
	// Update TextSize on all canvas text objects
	for _, obj := range r.parentRenderer.Objects() {
		if text, ok := obj.(*canvas.Text); ok {
			text.TextSize = r.entry.TextSize
		}
	}
	r.parentRenderer.Layout(size)
}

func (r *customEntryRenderer) MinSize() fyne.Size {
	return r.parentRenderer.MinSize()
}

func (r *customEntryRenderer) Refresh() {
	// Update TextSize on all canvas text objects
	for _, obj := range r.parentRenderer.Objects() {
		if text, ok := obj.(*canvas.Text); ok {
			text.TextSize = r.entry.TextSize
		}
	}
	r.parentRenderer.Refresh()
}

func (r *customEntryRenderer) Objects() []fyne.CanvasObject {
	return r.parentRenderer.Objects()
}

func (r *customEntryRenderer) Destroy() {
	r.parentRenderer.Destroy()
}

func newCustomEntry(size float32) *customEntry {
	e := &customEntry{
		TextSize: size,
	}
	e.ExtendBaseWidget(e)
	return e
}

// aspectRatioLayout is a custom layout that maintains aspect ratio for grid items
type aspectRatioLayout struct {
	aspectRatio float32
	rows        int
	cols        int
}

func newAspectRatioLayout(aspectRatio float32, rows, cols int) fyne.Layout {
	return &aspectRatioLayout{aspectRatio: aspectRatio, rows: rows, cols: cols}
}

func (a *aspectRatioLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	// Calculate cell size based on width
	cellWidthBasedOnWidth := size.Width / float32(a.cols)
	cellHeightFromWidth := cellWidthBasedOnWidth / a.aspectRatio

	// Calculate cell size based on height
	cellHeightBasedOnHeight := size.Height / float32(a.rows)
	cellWidthFromHeight := cellHeightBasedOnHeight * a.aspectRatio

	var cellWidth, cellHeight float32

	// Choose the smaller of the two to fit in the container
	if cellWidthBasedOnWidth*float32(a.cols) <= size.Width && cellHeightFromWidth*float32(a.rows) <= size.Height {
		cellWidth = cellWidthBasedOnWidth
		cellHeight = cellHeightFromWidth
	} else {
		cellWidth = cellWidthFromHeight
		cellHeight = cellHeightBasedOnHeight
	}

	if cellHeight*float32(a.rows) > size.Height {
		cellHeight = size.Height / float32(a.rows)
		cellWidth = cellHeight * a.aspectRatio
	}

	x, y := float32(0), float32(0)
	for i, o := range objects {
		o.Resize(fyne.NewSize(cellWidth, cellHeight))
		o.Move(fyne.NewPos(x, y))

		if (i+1)%a.cols == 0 {
			x = 0
			y += cellHeight
		} else {
			x += cellWidth
		}
	}
}

func (a *aspectRatioLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(10, 10)
}
