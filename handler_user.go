package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username is required.")
	}

	s.state.SetUser(cmd.args[0])
	fmt.Println("new user has logged in.")
	return nil
}
