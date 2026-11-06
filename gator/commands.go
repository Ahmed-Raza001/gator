package main

import "fmt"

// command represents a CLI command.
type command struct {
	name string
	args []string
}

// commands holds all registered command handlers.
type commands struct {
	handlers map[string]func(*state, command) error
}

// run executes a command if it exists.
func (c *commands) run(s *state, cmd command) error {

	handler, ok := c.handlers[cmd.name]

	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}

	return handler(s, cmd)
}

// register adds a new command handler.
func (c *commands) register(
	name string,
	f func(*state, command) error,
) {
	c.handlers[name] = f
}
