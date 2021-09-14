package main

import (
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"weather/lib/locations"
	"weather/lib/tui"
)

const usage = `Usage:weather [--location=<location>]`

type Options struct {
	Location string
	Current  bool
}

type Config struct {
	APIURL   string `yaml:"apiURL"`
	APIKey   string `yaml:"apiKey"`
	Location locations.Location
	Units    string
}

func main() {
	/*************************************** Load Config *****************************************/
	var config Config
	file, err := os.OpenFile("./config.yaml", os.O_RDONLY, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	yd := yaml.NewDecoder(file)
	err = yd.Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}

	/*************************************** Process Args ****************************************/
	var opt Options
	args := os.Args[1:]
	usg, err := docopt.ParseArgs(usage, args, "0.1.0")
	if err != nil {
		log.Fatalln(err)
	}

	err = usg.Bind(&opt)
	if err != nil {
		log.Fatalln(err)
	}

	/**************************************** Get Location ***************************************/
	lc := locations.NewGeocoder(config.APIURL, config.APIKey)

	var loc locations.Location
	if opt.Location != "" {
		loc, err = locations.Parse(opt.Location)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		loc = config.Location
	}

	if loc.Lat == 0 || loc.Lng == 0 {
		ls, err := lc.GetCoords(loc)
		if err != nil {
			log.Fatalln(err)
		}

		if len(ls) == 0 {
			log.Fatalln("could not get location")
		}

		if len(ls) > 1 {
			PrintLocations(ls)
			fmt.Print("Choose a city from the options: ")

			var i int
			for i == 0 {
				i, err = GetLocationChoice(len(ls))
				if err != nil {
					fmt.Print("Not a valid choice. Please choose another: ")
				}
			}

			loc = ls[i-1]
		} else {
			loc = ls[0]
		}
	}

	/************************************** Print Results ****************************************/
	t := tui.TUI{}
	err = t.Init(config.APIURL, config.APIKey, config.Units)
	if err != nil {
		log.Fatalln(err)
	}

	t.Start(loc)
}

// PrintLocations prints the locations that matched a search.
func PrintLocations(ls []locations.Location) {
	for i, l := range ls {
		fmt.Printf("%v) %v\n", i+1, l)
	}
}

// GetLocationChoice reads the user's selection from the console.
func GetLocationChoice(max int) (int, error) {
	var i int
	_, err := fmt.Scanf("%d\n", &i)
	if err != nil {
		return 0, err
	}

	if i < 1 || i > max {
		return 0, errors.New("choice out of range")
	}

	return i, nil
}
