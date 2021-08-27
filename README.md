# Weather

Weather is a simple command line tool for getting the current weather and forecast 
for different places.

## Goals

This is a little project for me to work on after work. I like weather and I like doing stuff on 
the command line, so this is a pretty fun project that I'll probably end up actually using a lot.
I'd also like to be able to show the output in a tmux status bar.


## Data Source

Weather gets it's data from [OpenWeather](https://openweather.org). It uses the one-call API to 
get weather data and the geocoding API to get the latitude and longitude of locations.

## Configuration

Weather stores its configuration at `~/.weather` in yaml format.

## Usage

```
    weather Sandy,UT current
    weather Sandy,UT daily=5
    weather Sandy,UT hourly=36
    weather Sandy,UT current daily=3 hourly=12
```

## Output

```
Weather in Sandy, UT
```