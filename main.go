package main

import (
	"database/sql"
	"fmt"
	"github.com/primawk/gator-go/internal/config"
	"github.com/primawk/gator-go/internal/database"
	"os"
)

import _ "github.com/lib/pq"

type state struct {
	state *config.Config
	db    *database.Queries
}

func main() {
	cfg := config.Read()
	db, err := sql.Open("postgres", cfg.Db_url)

	if err != nil {
		fmt.Printf("error reading database: %s\n", err)
		os.Exit(1)
	}

	newDb := database.New(db)

	stateConfig := state{
		state: &cfg,
		db:    newDb,
	}

	commandsInit := commands{
		commandList: map[string]func(*state, command) error{},
	}

	commandsInit.register("login", handlerLogin)
	commandsInit.register("register", handlerRegister)
	commandsInit.register("reset", handlerReset)
	commandsInit.register("users", handlerFetchUsers)
	commandsInit.register("agg", handlerFetchFeed)

	if len(os.Args) < 2 {
		fmt.Println("Please provide a command")
		os.Exit(1)
	}

	var commandInput command

	if (os.Args[1] == "login" || os.Args[1] == "register") && len(os.Args) > 2 {
		commandInput.name = os.Args[1]
		commandInput.args = []string{os.Args[2]}
	} else if os.Args[1] == "reset" || os.Args[1] == "users" || os.Args[1] == "agg" {
		commandInput.name = os.Args[1]
	} else {
		fmt.Println("an argument is required.")
		os.Exit(1)
	}

	commandsInit.run(&stateConfig, commandInput)
}
