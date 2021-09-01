# Weather

Weather is a simple command line tool for getting the current weather and forecast 
for different places.

## Goals

This is a little project for me to work on after work. I'm a bit of a weather nerd, so I like to 
have quick access to weather conditions and forecasts. I think it would be handy to have a 
simple command I can use for this rather than going to a website or opening an app.

I'm planning on it having two modes of operation, which I might separate into two commands. One 
would have a full feature TUI to show a lot of data all at once and maybe update it periodically.
The other would take a format string and fill that string with data. That mode would be useful 
for piping weather data into other programs, like tmux's status bar.

## Todo
- [ ] use location timezone for times.
- [ ] 

## Data Source

Weather gets it's data from [OpenWeather](https://openweather.org). It uses the one-call API to 
get weather data and the geocoding API to get the latitude and longitude of locations.

## Configuration

Weather stores its configuration at `~/.weather` in yaml format. This configuration contains the 
information for the weather API and some defaults like location and units.

## Usage

```
    weather Sandy,UT current
    weather Sandy,UT daily=5
    weather Sandy,UT hourly=36
    weather Sandy,UT current daily=3 hourly=12
```

## Output
```
--------------------------------------------------------------------------------
Weather in Sandy, UT
--------------------------------------------------------------------------------
                              Current Conditions

Sunny                                  Rain Last Hour:  .05in

Temperature: 45 F                      Feels Like: 42 F

Wind from the NNW at 24MPH

Humidity 50%                           Dew Point 35 F

```