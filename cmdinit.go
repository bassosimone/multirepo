// cmdinit.go - implementation of the init command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"io"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nflag"
	"github.com/kballard/go-shellquote"
)

// cmdInit is the static init command
var cmdInit = &clip.LeafCommand[environ]{
	BriefDescriptionText: "Initialize a multirepo.",
	RunFunc:              cmdInitMain,
}

// cmdInitRunner runs the init command.
type cmdInitRunner struct {
	// Style is the nil-safe libgloss style to use.
	Style *nilSafeLipglossStyle

	// XWriter is the writer used to log executed commands.
	XWriter io.Writer
}

// --- entry & setup ---

// cmdInitMain is the entry point for the init command.
func cmdInitMain(ctx context.Context, args *clip.CommandArgs[environ]) error {
	return mustNewCmdInitRunner(args).run(args)
}

// mustNewCmdInitRunner creates a new [*cmdInitRunner].
func mustNewCmdInitRunner(args *clip.CommandArgs[environ]) *cmdInitRunner {
	// Initialize the default configuration.
	c := &cmdInitRunner{
		Style:   nil,
		XWriter: io.Discard,
	}

	// Create empty command line parser.
	fset := nflag.NewFlagSet(args.CommandName, nflag.ExitOnError)
	fset.Description = args.Command.BriefDescription()
	fset.PositionalArgumentsUsage = ""
	fset.MinPositionalArgs = 0
	fset.MaxPositionalArgs = 0

	// Add the `-h, --help` flag.
	fset.AutoHelp("help", 'h', "Show this help message and exit.")

	// Add the `-x` flag.
	xflag := fset.Bool("print-commands", 'x', "Log the commands we execute.")

	// Parse the command line arguments.
	assert.NotError(fset.Parse(args.Args))

	// Honour the `-x` flag.
	if *xflag {
		c.XWriter = args.Env.Stderr()
		c.Style = newNilSafeLipglossStyle()
	}

	return c
}

// --- execution ---

func (c *cmdInitRunner) run(args *clip.CommandArgs[environ]) error {
	// Create the `.multirepo` directory
	dd := defaultDotDir()
	mustFprintf(c.XWriter, "%s\n", c.Style.Renderf("+ mkdir -p %s", shellquote.Join(dd.String())))
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
		mustFprintf(c.XWriter, "%s\n", c.Style.Renderf("+ echo '{}' > %s", shellquote.Join(dd.configFilePath())))
		if err := args.Env.WriteFile(dd.configFilePath(), []byte("{}\n"), 0600); err != nil {
			mustFprintf(args.Env.Stderr(), "multirepo init: %s\n", err.Error())
			return err
		}
	}

	return nil
}
