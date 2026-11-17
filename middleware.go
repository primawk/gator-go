package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/primawk/gator-go/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		if len(s.state.Current_user_name) == 0 {
			return errors.New("login is required.")
		}

		user, err := s.db.GetUser(context.Background(), s.state.Current_user_name)

		if err != nil {
			return fmt.Errorf("error getting the user id: %s", err)
		}

		return handler(s, cmd, user)
	}
}
