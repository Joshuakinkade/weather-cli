package formatter

import (
	"fmt"
	"time"
	"weather/lib/locations"
	"weather/lib/weather"
)

type Today struct {
	High      float64
	Low       float64
	Condition weather.Condition
	Sunrise   time.Time
	Sunset    time.Time
}

func PrettyPrint(w weather.Weather, loc locations.Location) {
	PrintLocation(loc)

	if !w.Current.DT.IsZero() {
		PrintCurrent(w.Current)
	}

	t := Today{
		High:      w.Daily[0].Temp.Max,
		Low:       w.Daily[0].Temp.Min,
		Condition: w.Daily[0].Weather[0],
		Sunrise:   time.Time(w.Current.Sunrise),
		Sunset:    time.Time(w.Current.Sunset),
	}

	PrintToday(t)

	fmt.Println("Hourly:")
	for _, h := range w.Hourly[1:9] {
		PrintHour(h)
	}
	fmt.Print("\n")

	fmt.Println("Daily:")
	for _, d := range w.Daily[1:] {
		PrintDay(d)
	}
	fmt.Print("\n")
}

func PrintLocation(loc locations.Location) {
	fmt.Printf(`
---------------------------------------
  Weather in %v
---------------------------------------

`, loc)
}

func PrintCurrent(c weather.Current) {
	fmt.Printf(`Currently:
  %v
  Temperature: %v, Feels Like: %v
  Wind from the %v at %v
  Relative Humidity: %v, Dew Point, %v

`,
		c.Weather[0].Main, FormatTemp(c.Temp), FormatTemp(c.FeelsLike), FormatWindDeg(c.WindDeg), FormatWindSpeed(c.WindSpeed), FormatPercent(c.Humidity/100), FormatTemp(c.DewPoint))
}

func PrintToday(t Today) {
	fmt.Printf(`Today:
  %v
  High: %v, Low %v
  Sunrise: %v, Sunset: %v

`, t.Condition.Main, t.High, t.Low, t.Sunrise.Format(time.Kitchen), t.Sunset.Format(time.Kitchen))
}

func PrintHour(h weather.Hourly) {
	t := time.Time(h.DT)
	fmt.Printf("  %v: %v, Temperature: %v, Precipitation: %v\n", t.Format(time.Kitchen),
		h.Weather[0].Main, FormatTemp(h.Temp), FormatPercent(h.Pop))
}

func PrintDay(d weather.Daily) {
	t := time.Time(d.DT)
	fmt.Printf("  %v: High: %v, Low: %v, Precipitation: %v\n", t.Format("Mon, 1/2"),
		FormatTemp(d.Temp.Max), FormatTemp(d.Temp.Min), FormatPercent(d.Pop))

}
