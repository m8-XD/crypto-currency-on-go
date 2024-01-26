// SPDX-License-Identifier: Unlicense OR MIT

package ui

// A simple Gio program. See https://gioui.org for more information.

import (
	"blockchainCentralServer/pkg/entity"
	"fmt"
	"image/color"
	"sync"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

func Start(si *entity.ServerInfo, outerWG *sync.WaitGroup) {
	w := app.NewWindow()
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
OuterLoop:
	for {
		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			break OuterLoop
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			l := material.H1(th, "Hello, Gio")
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			l.Color = maroon
			l.Alignment = text.Middle
			l.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}

	fmt.Println("ui stopped")

	outerWG.Done()
	return
}
