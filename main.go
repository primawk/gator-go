package main

import (
	"fmt"
	"github.com/primawk/gator-go/internal/config"
	"os"
)

type state struct {
	state *config.Config
}

func main() {
	cfg := config.Read()

	stateConfig := state{
		state: &cfg,
	}

	commandsInit := commands{
		commandList: map[string]func(*state, command) error{},
	}

	commandsInit.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Please provide a command")
		os.Exit(1)
	}

	var commandInput command

	if os.Args[1] == "login" && len(os.Args) > 2 {
		commandInput.name = os.Args[1]
		commandInput.args = []string{os.Args[2]}
	} else {
		fmt.Println("username is required.")
		os.Exit(1)
	}

	commandsInit.run(&stateConfig, commandInput)
}
