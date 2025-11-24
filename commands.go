package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	command map[string]func(*state, command) error
}

// run command, if it exists
func (c *commands) run(s *state, cmd command) error {
	f, exists := c.command[cmd.Name]
	if !exists {
		return fmt.Errorf("command %v does not exist", cmd.Name)
	}

	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.command[name] = f
}
