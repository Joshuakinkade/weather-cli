package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"
	"time"
	"weather/lib/formatter"
	"weather/lib/locations"
	"weather/lib/weather"
)

const (
	hl = rune('\u2500')
	tl = rune('\u250c')
	tr = rune('\u2510')
	vl = rune('\u2502')
	bl = rune('\u2514')
	br = rune('\u2518')
)

type Frame struct {
	X      int
	Y      int
	Height int
	Width  int
	Title  string
}

func (f Frame) Draw(sc tcell.Screen) {
	tw := len(f.Title)
	to := (f.Width - tw) / 2

	// draw border
	sc.SetContent(f.X, f.Y, tl, nil, tcell.StyleDefault)
	for i := 1; i < to-1; i++ {
		sc.SetContent(i+f.X, f.Y, hl, nil, tcell.StyleDefault)
	}
	for i, c := range f.Title {
		sc.SetContent(to+i+f.X, f.Y, c, nil, tcell.StyleDefault)
	}
	for i := to + tw + 1 + f.X; i < f.Width-1+f.X; i++ {
		sc.SetContent(i, f.Y, hl, nil, tcell.StyleDefault)
	}
	sc.SetContent(f.Width-1+f.X, f.Y, tr, nil, tcell.StyleDefault)

	for i := 1 + f.Y; i < f.Height+f.Y; i++ {
		sc.SetContent(f.X, i, vl, nil, tcell.StyleDefault)
		sc.SetContent(f.Width-1+f.X, i, vl, nil, tcell.StyleDefault)
	}

	sc.SetContent(f.X, f.Height+f.Y, bl, nil, tcell.StyleDefault)
	for i := f.X + 1; i < f.Width-1+f.X; i++ {
		sc.SetContent(i, f.Height+f.Y, hl, nil, tcell.StyleDefault)
	}
	sc.SetContent(f.Width-1+f.X, f.Height+f.Y, br, nil, tcell.StyleDefault)
}

type TUI struct {
	screen  tcell.Screen
	weather weather.Client
}

func (t *TUI) Init(apiURL, apiKey, units string) error {
	t.weather = weather.NewClient(apiURL, apiKey, units)

	var err error
	t.screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}

	return t.screen.Init()
}

func (t *TUI) Start(loc locations.Location) {
	var w weather.Weather
	w, _ = t.weather.Get(loc.Lat, loc.Lng)

	t.Draw(w, loc)
	for {
		event := t.screen.PollEvent()

		switch event := event.(type) {
		case *tcell.EventKey:
			if event.Rune() == 'q' {
				t.screen.Clear()
				os.Exit(0)
			} else if event.Rune() == 'r' {
				t.screen.Clear()
				w, _ = t.weather.Get(loc.Lat, loc.Lng)
				t.Draw(w, loc)
			}
		case *tcell.EventResize:
			t.screen.Clear()
			t.Draw(w, loc)
		}
	}
}

func (t *TUI) Draw(w weather.Weather, loc locations.Location) {
	width, _ := t.screen.Size()
	if width > 64 {
		width = 64
	}
	h := 22

	f := Frame{
		X:      0,
		Y:      0,
		Height: h,
		Width:  width,
		Title:  loc.Format(),
	}

	f.Draw(t.screen)

	tf := Frame{
		X:      2,
		Y:      2,
		Height: 5,
		Width:  width/2 - 2,
		Title:  "Today",
	}

	tf.Draw(t.screen)

	today := w.Daily[0]
	t.DrawString(today.Weather[0].Main, 4, 3)
	t.DrawString(fmt.Sprintf("High: %v", formatter.FormatTemp(today.Temp.Max)), 4, 4)
	t.DrawString(fmt.Sprintf("Low: %v", formatter.FormatTemp(today.Temp.Min)), 4, 5)
	t.DrawString(fmt.Sprintf("%% Precip: %v", formatter.FormatPercent(today.Pop)), 4, 6)

	cf := Frame{
		X:      width/2 + 1,
		Y:      2,
		Height: 5,
		Width:  width/2 - 3,
		Title:  "Current",
	}

	cf.Draw(t.screen)

	cw := w.Current
	t.DrawString(cw.Weather[0].Main, width/2+3, 3)
	t.DrawString(fmt.Sprintf("Temperature: %v", formatter.FormatTemp(cw.Temp)), width/2+3, 4)
	t.DrawString(fmt.Sprintf("Feels Like: %v", formatter.FormatTemp(cw.FeelsLike)), width/2+3, 5)
	t.DrawString(fmt.Sprintf("Wind: %v from the %v", formatter.FormatWindSpeed(cw.WindSpeed), formatter.FormatWindDeg(cw.WindDeg)), width/2+3, 6)

	hf := Frame{
		X:      2,
		Y:      9,
		Height: 4,
		Width:  width - 4,
		Title:  "Hourly",
	}

	hf.Draw(t.screen)

	var hs []weather.Hourly
	for i, h := range w.Hourly[1:11] {
		if i%2 == 0 {
			hs = append(hs, h)
		}
	}

	hw := (width - 4) / 5

	for i, hr := range hs {
		dt := time.Time(hr.DT)
		xOffset := 4 + (i * hw)
		t.DrawString(dt.Format(time.Kitchen), xOffset, 10)
		t.DrawString(hr.Weather[0].Main, xOffset, 11)
		t.DrawString(formatter.FormatTemp(hr.Temp), xOffset, 12)
	}

	df := Frame{
		X:      2,
		Y:      15,
		Height: 6,
		Width:  width - 4,
		Title:  "Daily",
	}

	df.Draw(t.screen)

	ds := w.Daily[1:5]

	dw := (width - 4) / 4

	for i, d := range ds {
		day := time.Time(d.DT)
		xOffset := 4 + (i * dw)
		t.DrawString(day.Format("Mon"), xOffset, 16)
		t.DrawString(d.Weather[0].Main, xOffset, 17)
		t.DrawString(fmt.Sprintf("H: %v", formatter.FormatTemp(d.Temp.Max)), xOffset, 18)
		t.DrawString(fmt.Sprintf("L: %v", formatter.FormatTemp(d.Temp.Min)), xOffset, 19)
		t.DrawString(fmt.Sprintf("Precip: %v", formatter.FormatPercent(d.Pop)), xOffset, 20)
	}

	t.screen.Show()
}

func (t TUI) DrawString(text string, x, y int) {
	for i, c := range text {
		t.screen.SetContent(x+i, y, c, nil, tcell.StyleDefault)
	}
}
