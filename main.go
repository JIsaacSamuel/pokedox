package main

import (
	"bufio"
	"fmt"
	"internal/maps"
	"os"
	"strings"
	// "internal/maps"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type config struct {
	next string
	prev string
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		cmd := words[0]

		result, ok := getCommand()[cmd]

		if ok {
			err := result.callback()
			if err != nil {
				continue
			}
		} else {
			fmt.Print("Command does not exist")
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
			callback:    maps.Map,
		},
		"mapb": {
			name:        "mapb",
			description: "Prints (the previous) 20 locations",
			callback:    maps.Mapb,
		},
	}
}

func commandHelp() error {
	v := getCommand()
	for _, value := range v {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}
