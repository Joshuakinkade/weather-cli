package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Timestamp represents the time for a weather forecast.
type Timestamp time.Time

// UnmarshalJSON converts a UNIX timestamp from JSON to a Timestamp.
func (t *Timestamp) UnmarshalJSON(d []byte) error {
	i, err := strconv.Atoi(string(d))
	if err != nil {
		return err
	}
	pt := time.Unix(int64(i), 0)
	*t = Timestamp(pt)
	return nil
}

// IsZero returns true if the timestamp is at it's zero value.
func (t Timestamp) IsZero() bool {
	dt := time.Time(t)
	return dt.IsZero()
}

// Condition contains description information about the weather.
type Condition struct {
	ID          int
	Main        string
	Description string
	Icon        string
}

// Current contains data for the current weather.
type Current struct {
	DT         Timestamp
	Sunrise    Timestamp
	Sunset     Timestamp
	Temp       float64
	FeelsLike  float64 `json:"feels_like"`
	Pressure   float64
	Humidity   float64
	DewPoint   float64
	UVI        float64
	Clouds     float64
	Visibility float64
	WindSpeed  float64 `json:"wind_speed"`
	WindDeg    float64 `json:"wind_deg"`
	Weather    []Condition
	Rain       struct {
		OneHour float64 `json:"1hr"`
	}
}

// Hourly contains data for the forecast at an hour.
type Hourly struct {
	DT         Timestamp
	Temp       float64
	FeelsLike  float64 `json:"feels_like"`
	Pressure   float64
	UVI        float64
	Clouds     float64
	Visibility float64
	WindSpeed  float64 `json:"wind_speed"`
	WindGust   float64 `json:"wind_gust"`
	WindDeg    float64 `json:"wind_deg"`
	Weather    []Condition
	Pop        float64
}

// Daily contains forecast data for a day.
type Daily struct {
	DT        Timestamp
	Sunrise   Timestamp
	Sunset    Timestamp
	Moonrise  Timestamp
	Moonset   Timestamp
	MoonPhase float64 `json:"moon_phase"`
	Temp      struct {
		Day     float64
		Min     float64
		Max     float64
		Night   float64
		Eve     float64
		Morning float64 `json:"morn"`
	}
	FeelsLike struct {
		Day     float64
		Night   float64
		Eve     float64
		Morning float64 `json:"morn"`
	} `json:"feels_like"`
	Pressure  float64
	Humidity  float64
	DewPoint  float64 `json:"dew_point"`
	WindSpeed float64 `json:"wind_speed"`
	WindDeg   float64 `json:"wind_deg"`
	Weather   []Condition
	Clouds    float64
	Pop       float64
	Rain      float64
	UVI       float64
}

// Weather contains weather data returned by the API.
type Weather struct {
	Lat      float64
	Lng      float64 `json:"lon"`
	Timezone string
	Current  Current
	Hourly   []Hourly
	Daily    []Daily
}

// Client connects to OpenWeatherAPI.
type Client struct {
	client http.Client
	apiURL string
	apiKey string
}

// NewClient returns an initialized Client.
func NewClient(apiURL, apiKey string) Client {
	return Client{
		client: http.Client{},
		apiURL: apiURL,
		apiKey: apiKey,
	}
}

// Get requests the weather for a location given by latitude and longitude.
func (c Client) Get(lat, lng float64) (Weather, error) {
	url := fmt.Sprintf("%v/data/2.5/onecall?lat=%v&lon=%v&appid=%v&units=imperial", c.apiURL, lat, lng, c.apiKey)

	fmt.Println(url)

	r, err := c.client.Get(url)
	if err != nil {
		return Weather{}, err
	}

	var w Weather
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&w)
	return w, err
}
