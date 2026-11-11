package main

import (
	"fmt"
	"os"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commandList map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) {
	if err := c.commandList[cmd.name](s, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	if _, ok := c.commandList[name]; ok {
		fmt.Printf("command %s is already exist.\n", name)
	} else {
		c.commandList[name] = f
	}
}
