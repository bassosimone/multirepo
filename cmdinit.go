// cmdinit.go - implementation of the init command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"io"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/flag"
	"github.com/kballard/go-shellquote"
)

// initCmd is the static init command
var initCmd = &clip.LeafCommand[environ]{
	BriefDescriptionText: "Initialize a multirepo.",
	RunFunc:              (&cmdInit{}).Run,
}

// cmdInit implements the init command.
type cmdInit struct {
	// XWriter is the writer used to log executed commands.
	XWriter io.Writer
}

// Run is the entry point for the init command.
func (c *cmdInit) Run(ctx context.Context, args *clip.CommandArgs[environ]) error {
	// Parse command line arguments
	if err := c.getopt(args); err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo init: %s\n", err)
		return err
	}

	// Create the `.multirepo` directory
	dd := defaultDotDir()
	mustFprintf(c.XWriter, "+ mkdir -p %s\n", shellquote.Join(dd.String()))
	if err := args.Env.MkdirAll(dd.String(), 0700); err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo init: %s\n", err.Error())
		return err
	}

	// Lock the multirepo dir
	unlock, err := dd.lock(args.Env)
	if err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo init: %s\n", err)
		return err
	}
	defer unlock()

	// Check whether the configuration file already exists
	exists, err := args.Env.FileExists(dd.configFilePath())
	if err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo init: %s\n", err.Error())
		return err
	}

	// Write the initial configuration file
	if !exists {
		mustFprintf(c.XWriter, "+ echo '{}' > %s\n", shellquote.Join(dd.configFilePath()))
		if err := args.Env.WriteFile(dd.configFilePath(), []byte("{}\n"), 0600); err != nil {
			mustFprintf(args.Env.Stderr(), "multirepo init: %s\n", err.Error())
			return err
		}
	}

	return nil
}

// getopt gets command line options.
func (c *cmdInit) getopt(args *clip.CommandArgs[environ]) error {
	// Initialize the default configuration.
	c.XWriter = io.Discard

	// Create empty command line parser.
	clp := flag.NewFlagSet(args.CommandName, flag.ContinueOnError)
	clp.SetDescription(args.Command.BriefDescription())
	clp.SetArgsDocs("")

	// Add the `-x` flag.
	xflag := clp.Bool("print-commands", 'x', false, "Log the commands we execute.")

	// Parse the command line arguments.
	if err := clp.Parse(args.Args); err != nil {
		return err
	}
	if err := clp.PositionalArgsEqualCheck(0); err != nil {
		return err
	}

	// Honour the `-x` flag.
	if *xflag {
		c.XWriter = args.Env.Stderr()
	}

	return nil
}
