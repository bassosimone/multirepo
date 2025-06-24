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

// cmdInit is the static init command
var cmdInit = &clip.LeafCommand[environ]{
	BriefDescriptionText: "Initialize a multirepo.",
	RunFunc:              (&cmdInitRunner{}).Run,
}

// cmdInitRunner runs the init command.
type cmdInitRunner struct {
	// XWriter is the writer used to log executed commands.
	XWriter io.Writer
}

// --- entry & setup ---

// Run is the entry point for the init command.
func (c *cmdInitRunner) Run(ctx context.Context, args *clip.CommandArgs[environ]) error {
	c.mustGetopt(args)
	return c.run(args)
}

// mustGetopt gets command line options.
func (c *cmdInitRunner) mustGetopt(args *clip.CommandArgs[environ]) {
	// Initialize the default configuration.
	c.XWriter = io.Discard

	// Create empty command line parser.
	clp := flag.NewFlagSet(args.CommandName, flag.ExitOnError)
	clp.SetDescription(args.Command.BriefDescription())
	clp.SetArgsDocs("")

	// Add the `-x` flag.
	xflag := clp.Bool("print-commands", 'x', false, "Log the commands we execute.")

	// Parse the command line arguments.
	clip.Must(args.Env, clp.Parse(args.Args))
	clip.Must(args.Env, clp.PositionalArgsEqualCheck(0))

	// Honour the `-x` flag.
	if *xflag {
		c.XWriter = args.Env.Stderr()
	}
}

// --- execution ---

func (c *cmdInitRunner) run(args *clip.CommandArgs[environ]) error {
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
