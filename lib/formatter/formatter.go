package formatter

import (
	"fmt"
	"math"
)

func FormatWindDeg(deg float64) string {
	cds := []string{"N", "NNE", "E", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
	i := int(math.Round(deg / 22.5))
	i = i % 16
	return cds[i]
}

func FormatTemp(temp float64) string {
	return fmt.Sprintf("%v°F", math.Round(temp))
}

func FormatPercent(pct float64) string {
	return fmt.Sprintf("%v%%", math.Round(pct*100))
}

func FormatWindSpeed(windSpeed float64) string {
	return fmt.Sprintf("%vmph", math.Round(windSpeed))
}
