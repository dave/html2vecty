package main

import (
	"github.com/dave/splitter"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

type Page struct {
	vecty.Core
	app *App

	split *splitter.Split
}

func NewPage(app *App) *Page {
	v := &Page{
		app: app,
	}
	return v
}

func (v *Page) Mount() {
	v.app.Watch(v, func(done chan struct{}) {
		defer close(done)
		// Only top-level page should fire vecty.Rerender
		vecty.Rerender(v)
	})

	v.split = splitter.New("split")
	v.split.Init(
		js.S{"#left", "#right"},
		js.M{"sizes": []float64{50, 50}},
	)
}

func (v *Page) Unmount() {
	v.app.Delete(v)
}

func (v *Page) Render() vecty.ComponentOrHTML {
	return elem.Body(
		elem.Div(
			vecty.Markup(
				vecty.Class("container-fluid", "p-0", "split", "split-horizontal"),
			),
			v.renderLeft(),
			v.renderRight(),
		),
	)
}

func (v *Page) renderLeft() *vecty.HTML {
	return elem.Div(
		vecty.Markup(
			prop.ID("left"),
			vecty.Class("split"),
		),
		NewEditor(v.app, "html-editor", "html", v.app.Editor.Html(), func(value string) {
			v.app.Dispatch(&UserChangedTextAction{
				Text: value,
			})
		}),
	)
}

func (v *Page) renderRight() *vecty.HTML {
	return elem.Div(
		vecty.Markup(
			prop.ID("right"),
			vecty.Class("split"),
		),
		NewEditor(v.app, "code-editor", "golang", v.app.Editor.Code(), nil),
	)
}
