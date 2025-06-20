// command.go - Command definition.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"fmt"
)

// command implements a given command.
type command interface {
	// Description returns a short description of the command.
	Description() string

	// Run runs this command using the specified [environ].
	Run(ctx context.Context, env environ, argv cliArgs) error
}

// commandTable maps a name to a [command].
type commandTable map[string]command

var _ command = (commandTable)(nil)

// Description implements [command].
func (c commandTable) Description() string {
	assertTrue(false, "The commandTable should never invoked as an `--help target`.")
	return ""
}

// Run implements [command].
func (c commandTable) Run(ctx context.Context, env environ, argv cliArgs) error {
	// Obtain the command name
	cmdName := argv.CommandName()

	// Find the command inside the command table.
	cmd := c[cmdName]

	// Handle the case where the command is not found.
	if cmd == nil {
		err := fmt.Errorf("command %q not found", cmdName)
		mustFprintf(env.Stderr(), "multirepo: %s\n", err.Error())
		mustFprintf(env.Stderr(), "Try `multirepo --help` for help.\n")
		return err
	}

	// Run the command.
	return cmd.Run(ctx, env, argv)
}
