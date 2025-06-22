// cmdinit.go - implementation of the init command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"io"

	"github.com/bassosimone/clip/pkg/flag"
	"github.com/kballard/go-shellquote"
)

// cmdInit implements the init command.
type cmdInit struct{}

var _ command = (*cmdInit)(nil)

// Description implements [command].
func (c *cmdInit) Description() string {
	return "Initialize a new multirepo."
}

// Run implements [command].
func (c *cmdInit) Run(ctx context.Context, env environ, argv cliArgs) error {
	// Print help if requested to do so by the user.
	if argv.ContainsHelp() {
		return c.help(env.Stdout())
	}

	// Parse command line arguments
	options, err := c.getopt(env, argv)
	if err != nil {
		mustFprintf(env.Stderr(), "multirepo init: %s\n", err)
		mustFprintf(env.Stderr(), "Try `multirepo init --help` for help.\n")
		return err
	}

	// Create the `.multirepo` directory
	dd := defaultDotDir()
	mustFprintf(options.XWriter, "+ mkdir -p %s\n", shellquote.Join(dd.String()))
	if err := env.MkdirAll(dd.String(), 0700); err != nil {
		mustFprintf(env.Stderr(), "multirepo init: %s\n", err.Error())
		return err
	}

	// Lock the multirepo dir
	unlock, err := dd.lock(env)
	if err != nil {
		mustFprintf(env.Stderr(), "multirepo init: %s\n", err)
		return err
	}
	defer unlock()

	// Check whether the configuration file already exists
	exists, err := env.FileExists(dd.configFilePath())
	if err != nil {
		mustFprintf(env.Stderr(), "multirepo init: %s\n", err.Error())
		return err
	}

	// Write the initial configuration file
	if !exists {
		mustFprintf(options.XWriter, "+ echo '{}' > %s\n", shellquote.Join(dd.configFilePath()))
		if err := env.WriteFile(dd.configFilePath(), []byte("{}\n"), 0600); err != nil {
			mustFprintf(env.Stderr(), "multirepo init: %s\n", err.Error())
			return err
		}
	}

	return nil
}

// help prints the help message for the init command.
func (c *cmdInit) help(w io.Writer) error {
	mustFprintf(w, "\n")
	mustFprintf(w, "multirepo init - %s\n", c.Description())
	mustFprintf(w, "\n")
	mustFprintf(w, "This command initializes a multirepo in the current directory.\n")
	mustFprintf(w, "\n")
	mustFprintf(w, "usage: multirepo init [-x]\n")
	mustFprintf(w, "\n")
	mustFprintf(w, "Flags:\n")
	mustFprintf(w, "  -x, --print-commands  Print commands as they are executed\n")
	mustFprintf(w, "\n")
	return nil
}

// cmdInitOptions contains configuration for the init command.
type cmdInitOptions struct {
	// XWriter is the writer used to log executed commands.
	XWriter io.Writer
}

// getopt gets command line options.
func (c *cmdInit) getopt(env environ, argv cliArgs) (*cmdInitOptions, error) {
	// Initialize the default configuration.
	options := &cmdInitOptions{
		XWriter: io.Discard,
	}

	// Create empty command line parser.
	clp := flag.NewFlagSet("", flag.ContinueOnError)

	// Add the `-x` flag.
	xflag := clp.Bool("print-commands", 'x', false, "")

	// Parse the command line arguments.
	if err := clp.Parse(argv.CommandArgs()); err != nil {
		return nil, err
	}

	// Honour the `-x` flag.
	if *xflag {
		options.XWriter = env.Stderr()
	}

	// Return the configuration.
	return options, nil
}
