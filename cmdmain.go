// cmdmain.go - implementation of the main command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"io"
)

// cmdMain implements the main command.
type cmdMain struct{}

var _ command = (*cmdMain)(nil)

// Description implements [command].
func (c *cmdMain) Description() string {
	assertTrue(false, "The cmdMain should never invoked as an `--help target`.")
	return ""
}

// Run implements [command].
func (c *cmdMain) Run(ctx context.Context, env environ, argv cliArgs) error {
	// Construct the command table.
	ct := c.buildCommandsTable()

	// Print help if requested to do so by the user.
	if argv.IsEmpty() || argv.FirstArgIsHelp() {
		return c.help(env.Stdout(), ct)
	}

	// Obtain the arguments passed to the command.
	argv = argv.CommandArgs()

	// Execute the selected subcommand.
	return ct.Run(ctx, env, argv)
}

// help prints the help message.
func (c *cmdMain) help(w io.Writer, ct commandTable) error {
	mustFprintf(w, "\n")
	mustFprintf(w, "multirepo - manage multiple git repositories as a monorepo.\n")
	mustFprintf(w, "\n")
	mustFprintf(w, "usage: multirepo {command} [args...]\n")
	mustFprintf(w, "\n")
	mustFprintf(w, "commands:\n")
	mustFprintf(w, "\n")
	for cname, cmd := range ct {
		mustFprintf(w, "\t%-10s\t%s\n", cname, cmd.Description())
	}
	mustFprintf(w, "\n")
	mustFprintf(w, "Use `multirepo {command} --help` for help on `{command}`.\n")
	mustFprintf(w, "\n")
	return nil
}

// buildCommandsTable builds the [command] table.
func (c *cmdMain) buildCommandsTable() commandTable {
	return map[string]command{
		"clone":   &cmdClone{},
		"foreach": &cmdForeach{},
		"init":    &cmdInit{},
	}
}
