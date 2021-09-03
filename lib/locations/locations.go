package locations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

// LocationParseError represents an error with parsing a location string.
type LocationParseError struct {
	RawText string
}

// Error returns a string representation of the error.
func (lpe LocationParseError) Error() string {
	return fmt.Sprintf("%v is not a valid location.", lpe.RawText)
}

// Location represents a geographical location.
type Location struct {
	Name    string
	State   string
	Country string
	Lat     float64
	Lng     float64 `json:"lon"`
}

// Parse reads a location string and returns a Location.
func Parse(raw string) (Location, error) {
	if raw == "" {
		return Location{}, LocationParseError{RawText: raw}
	}
	var loc Location
	rps := strings.Split(raw, ",")

	var ps []string
	for _, p := range rps {
		ps = append(ps, strings.Trim(p, " "))
	}

	if len(ps) == 1 {
		loc.Name = raw
	}
	if len(ps) == 2 {
		loc.Name = ps[0]
		if isState(ps[1]) {
			loc.State = ps[1]
			loc.Country = "US"
		} else {
			loc.Country = ps[1]
		}
	}
	if len(ps) == 3 {
		loc.Name = ps[0]
		loc.State = ps[1]
		loc.Country = ps[2]
	}
	if len(ps) > 3 {
		err := LocationParseError{RawText: raw}
		return Location{}, err
	}
	return loc, nil
}

// isState returns true if the passed in string is a US state abbreviation, or false otherwise.
func isState(str string) bool {
	usStates := []string{"AK", "AL", "AR", "AZ", "CA", "CO", "CT", "DE", "FL", "GA", "HI", "IA", "ID", "IL", "IN", "KS", "KY", "LA", "MA", "MD", "ME", "MI", "MN", "MO", "MS", "MT", "NC", "ND", "NE", "NH", "NJ", "NM", "NV", "NY", "OH", "OK", "OR", "PA", "RI", "SC", "SD", "TN", "TX", "UT", "VA", "VT", "WA", "WI", "WV", "WY"}
	i := sort.SearchStrings(usStates, str)
	if usStates[i] == str {
		return true
	}
	return false
}

// Format returns a formatted representation to pass to the API.
func (loc Location) Format() string {
	str := loc.Name
	if loc.State != "" {
		str += "," + loc.State
	}
	str += "," + loc.Country
	return str
}

// String formats the location for debugging.
func (loc Location) String() string {
	ps := []string{loc.Name}
	if loc.State != "" {
		ps = append(ps, loc.State)
	}
	if loc.Country != "" {
		ps = append(ps, loc.Country)
	}

	txt := strings.Join(ps, ", ")

	if loc.Lat != 0 && loc.Lng != 0 {
		txt += fmt.Sprintf(" (%.4f, %.4f)", loc.Lat, loc.Lng)
	}

	return txt
}

// Geocoder manages a connection to OpenWeatherAPI's geocoding API endpoints.
type Geocoder struct {
	client http.Client
	apiURL string
	apiKey string
}

// NewGeocoder returns an initialized Geocoder.
func NewGeocoder(apiURL, apiKey string) Geocoder {
	return Geocoder{
		client: http.Client{},
		apiURL: apiURL,
		apiKey: apiKey,
	}
}

// GetCoords requests the coordinates for a location.
func (g Geocoder) GetCoords(loc Location) ([]Location, error) {
	term := loc.Format()
	u := fmt.Sprintf("%v/geo/1.0/direct?q=%v&limit=5&appid=%v", g.apiURL, term, g.apiKey)

	resp, err := g.client.Get(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(body))
	}

	var ls []Location
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ls)
	if err != nil {
		return nil, err
	}
	return ls, nil
}
