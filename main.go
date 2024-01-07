package main

import (
	"bufio"
	"errors"
	"fmt"
	"internal/maps"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type config struct {
	pokeapiClient    maps.Client
	nextLocationsURL *string
	prevLocationsURL *string
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	pokeClient := maps.NewClient(5 * time.Second)
	cfg := &config{
		pokeapiClient: pokeClient,
	}

	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		cmd := words[0]
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}

		result, ok := getCommand()[cmd]

		if ok {
			err := result.callback(cfg, args...)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			fmt.Print("Command does not exist\n")
		}
	}
}

func cleanInput(text string) []string {
	temp := strings.ToLower(text)
	res := strings.Fields(temp)
	return res
}

func getCommand() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Prints (the next) 20 locations",
			callback:    commandMapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Prints (the previous) 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Prints pokemons in a location",
			callback:    commandExplore,
		},
	}
}

func commandHelp(cfg *config, args ...string) error {
	v := getCommand()
	for _, value := range v {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	return nil
}

func commandExit(cfg *config, args ...string) error {
	os.Exit(0)
	return nil
}

func commandMapf(cfg *config, args ...string) error {
	locationsResp, err := cfg.pokeapiClient.ListLocations(cfg.nextLocationsURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locationsResp.Next
	cfg.prevLocationsURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.prevLocationsURL == nil {
		return errors.New("you're on the first page")
	}

	locationResp, err := cfg.pokeapiClient.ListLocations(cfg.prevLocationsURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locationResp.Next
	cfg.prevLocationsURL = locationResp.Previous

	for _, loc := range locationResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a location name")
	}
	name := args[0]
	fmt.Printf("Exploring %s\n", name)

	locationResp, err := cfg.pokeapiClient.ExploreLocation(name)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, poks := range locationResp.PokemonEncounters {
		fmt.Printf("-%s\n", poks.Pokemon.Name)
	}

	return nil
}
