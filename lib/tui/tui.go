package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"weather/lib/locations"
	"weather/lib/weather"
)

type TUI struct {
	screen tcell.Screen
}

func (t *TUI) Init() error {
	var err error
	t.screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}
	err = t.screen.Init()
	return err
}

func (t *TUI) Start(weather weather.Weather, loc locations.Location) {
	tv := TextView{Text: fmt.Sprintf("%v, %v", loc.Name, loc.State), Style: tcell.StyleDefault}

	w, h := t.screen.Size()

	c := NewVBox(w, h, []Drawer{tv})

	c.Draw(t.screen, 0, 0)
	t.screen.Show()
}

// Drawer is a type that can draw to a Screen
type Drawer interface {
	Draw(screen tcell.Screen, x, y int)
	Size() (int, int)
}

type TextView struct {
	Text  string
	Style tcell.Style
}

func (t TextView) Draw(screen tcell.Screen, x, y int) {
	for i, c := range t.Text {
		screen.SetContent(x+i, y, c, nil, t.Style)
	}
}

func (t TextView) Size() (int, int) {
	return len(t.Text), 1
}

type VBox struct {
	children []Drawer
	width    int
	height   int
}

func NewVBox(width, height int, children []Drawer) VBox {
	return VBox{
		width:    width,
		height:   height,
		children: children,
	}
}

func (v VBox) Draw(screen tcell.Screen, x, y int) {
	var row int
	for _, c := range v.children {
		cw, ch := c.Size()
		c.Draw(screen, ((v.width-cw)/2)+x, row+y)
		row += ch
	}
}

func (v VBox) Size() (int, int) {
	return v.width, v.height
}

type HBox struct {
	children []Drawer
	width    int
	height   int
}

func (h HBox) Draw(screen tcell.Screen, x, y int) {

}

func (h HBox) Size() (int, int) {
	return h.width, h.height
}

type Title struct {
	Text string
}

func (t Title) Draw(screen tcell.Screen, x, y int) {
	tv := TextView{Style: tcell.StyleDefault, Text: t.Text}
	b := NewVBox()
}
