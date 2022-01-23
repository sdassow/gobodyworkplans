package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TapLabel struct {
	widget.Label
	OnTapped       func(*TapLabel)
}

func NewTapLabel(t string, fn func(*TapLabel)) *TapLabel {
	w := &TapLabel{
		OnTapped:       fn,
	}
	w.ExtendBaseWidget(w)
	w.SetText(t)

	return w
}

func (w *TapLabel) Tapped(_ *fyne.PointEvent) {
	w.OnTapped(w)
}
