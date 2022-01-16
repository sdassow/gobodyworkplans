package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TapLabel struct {
	widget.Label
	OnTap       func(string)
	OnTapString string
}

func NewTapLabel(t string, s string, fn func(string)) *TapLabel {
	w := &TapLabel{
		OnTap:       fn,
		OnTapString: s,
	}
	w.ExtendBaseWidget(w)
	w.SetText(t)

	return w
}

func (w *TapLabel) Tapped(_ *fyne.PointEvent) {
	w.OnTap(w.OnTapString)
}
