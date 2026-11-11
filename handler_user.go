package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/primawk/gator-go/internal/database"
	"time"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username is required.")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])

	if err != nil {
		return fmt.Errorf("error login as: %s", cmd.args[0])
	}

	s.state.SetUser(user.Name)
	fmt.Println("new user has logged in.")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username is required.")
	}

	payload := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	user, err := s.db.CreateUser(context.Background(), payload)

	if err != nil {
		return fmt.Errorf("error creating user: %s\n", err)
	}

	s.state.SetUser(cmd.args[0])
	fmt.Println(user)

	return nil
}
