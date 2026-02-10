package main

import (
	"fyne.io/fyne/v2"
)

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
