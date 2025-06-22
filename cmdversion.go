// cmdversion.go - implementation of the version command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"fmt"
	"io"

	"github.com/bassosimone/clip/pkg/flag"
)

// cmdVersion implements the version command.
type cmdVersion struct{}

var _ command = (*cmdVersion)(nil)

// Description implements [command].
func (c *cmdVersion) Description() string {
	return "Prints the tool version."
}

// Run implements [command].
func (c *cmdVersion) Run(ctx context.Context, env environ, argv cliArgs) error {
	// Print help if requested to do so by the user.
	if argv.ContainsHelp() {
		return c.help(env.Stdout())
	}

	// Parse command line arguments
	if err := c.getopt(env, argv); err != nil {
		mustFprintf(env.Stderr(), "multirepo version: %s\n", err)
		mustFprintf(env.Stderr(), "Try `multirepo version --help` for help.\n")
		return err
	}

	// Print version information
	mustFprintf(env.Stdout(), "%s\n", Version)
	return nil
}

// help prints the help message for the version command.
func (c *cmdVersion) help(w io.Writer) error {
	mustFprintf(w, "\n")
	mustFprintf(w, "multirepo version - %s\n", c.Description())
	mustFprintf(w, "\n")
	mustFprintf(w, "This command prints the tool version.\n")
	mustFprintf(w, "\n")
	mustFprintf(w, "usage: multirepo version\n")
	mustFprintf(w, "\n")
	return nil
}

// getopt gets command line options.
func (c *cmdVersion) getopt(env environ, argv cliArgs) error {
	// Create empty command line parser.
	clp := flag.NewFlagSet("", flag.ContinueOnError)

	// Parse the command line arguments.
	if err := clp.Parse(argv.CommandArgs()); err != nil {
		return err
	}

	// Ensure there are no positional arguments
	if len(clp.Args()) > 0 {
		return fmt.Errorf("unexpected positional arguments: %v", clp.Args())
	}

	// Return the configuration.
	return nil
}
