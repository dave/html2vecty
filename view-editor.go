package main

import (
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	"github.com/tulir/gopher-ace"
)

type Editor struct {
	vecty.Core
	app *App

	Text string `vecty:"prop"`

	editor   ace.Editor
	id, lang string
	change   func(string)
}

func NewEditor(app *App, id, lang, text string, change func(string)) *Editor {
	v := &Editor{
		app:    app,
		lang:   lang,
		id:     id,
		change: change,
		Text:   text,
	}
	return v
}

func (v *Editor) Mount() {
	v.editor = ace.Edit(v.id)
	v.editor.SetOptions(map[string]interface{}{
		"mode": "ace/mode/" + v.lang,
	})
	if v.Text != "" {
		v.editor.SetValue(v.Text)
		v.editor.ClearSelection()
		v.editor.MoveCursorTo(0, 0)
	}
	if v.change != nil {
		var changes int
		v.editor.OnChange(func(ev *js.Object) {
			changes++
			before := changes
			go func() {
				<-time.After(time.Millisecond * 250)
				if before == changes {
					v.change(v.editor.GetValue())
				}
			}()
		})
	}
}

func (v *Editor) Render() vecty.ComponentOrHTML {
	if v.editor.Object != nil && v.Text != v.editor.GetValue() {
		// only update the editor if the text is changed
		v.editor.SetValue(v.Text)
		v.editor.ClearSelection()
		v.editor.MoveCursorTo(0, 0)
	}

	return elem.Div(
		vecty.Markup(
			prop.ID(v.id),
			vecty.Class("editor"),
		),
		vecty.Text(v.Text),
	)
}
