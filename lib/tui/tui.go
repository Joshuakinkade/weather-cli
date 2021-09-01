package tui

import (
	"github.com/gdamore/tcell/v2"
	"weather/lib/formatter"
	"weather/lib/locations"
	"weather/lib/weather"
)

// TextCell represents a single, one-line, styled piece of text.
type TextCell struct {
	Text  string
	Width int
	X     int
	Style tcell.Style
	Align string
}

// Draw writes the text cell to the Screen on the given row.
func (t TextCell) Draw(sc tcell.Screen, row int) {
	var offset int
	if t.Align == "center" {
		offset = (t.Width - len(t.Text)) / 2
	} else if t.Align == "left" {
		offset = 0
	} else if t.Align == "right" {
		offset = t.Width - len(t.Text)
	}
	for i, c := range t.Text {
		sc.SetContent(i+offset, row, c, nil, t.Style)
	}
}

// Box represents a rectangular area on the screen.
type Box struct {
	X         int
	Y         int
	Height    int
	Width     int
	textCells []TextCell
	row       int
}

// AddText adds a text cell to the Box
func (b *Box) AddText(t TextCell) {
	b.textCells = append(b.textCells, t)
}

// AddSpace adds a blank line to the Box
func (b *Box) AddSpace() {
	b.row++
}

// Draw writes its textCells to the Screen
func (b *Box) Draw(sc tcell.Screen) {
	for _, t := range b.textCells {
		t.Draw(sc, b.row+b.Y)
		b.row++
	}
	b.row = 0
}

type TUI struct {
	screen tcell.Screen
	row    int
}

func (t *TUI) Init() error {
	var err error
	t.screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}

	err = t.screen.Init()
	if err != nil {
		return err
	}

	return nil
}

func (t *TUI) PrintWeather(wtr weather.Weather, loc locations.Location) error {
	var bs []Box

	bs = append(bs, t.title(loc.String()))
	bs = append(bs, t.current(wtr.Current))

	for _, b := range bs {
		b.Draw(t.screen)
	}

	t.screen.Show()
	return nil
}

func (t *TUI) title(title string) Box {
	w, _ := t.screen.Size()
	box := Box{Width: w, Height: 3, X: 0, Y: 0}
	box.AddText(t.hr(w, 0))
	box.AddText(TextCell{
		X:     0,
		Width: w,
		Text:  title,
		Align: "center",
		Style: tcell.StyleDefault,
	})
	box.AddText(t.hr(w, 0))

	return box
}

func (t *TUI) current(cur weather.Current) Box {
	w, _ := t.screen.Size()
	box := Box{Width: w / 2, Height: 0, X: 0, Y: 3}
	box.AddText(TextCell{
		Text:  "Current Conditions",
		Width: w / 2,
		Align: "center",
		Style: tcell.StyleDefault,
	})
	box.AddText(TextCell{
		Text:  "Temperature: " + formatter.FormatTemp(cur.Temp),
		Width: w / 2,
		Align: "left",
		Style: tcell.StyleDefault,
	})
	return box
}

func (t *TUI) hr(w, offset int) TextCell {
	return TextCell{
		X:     offset,
		Text:  "--------------------------------------------------------------------------------",
		Width: w,
		Style: tcell.StyleDefault,
		Align: "left",
	}
}
