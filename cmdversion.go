// cmdversion.go - implementation of the version command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/flag"
)

// cmdVersion is the static version command
var cmdVersion = &clip.LeafCommand[environ]{
	BriefDescriptionText: "Display the version of the tool.",
	RunFunc:              (&cmdVersionRunner{}).Run,
}

// cmdVersionRunner runs the version command.
type cmdVersionRunner struct{}

// Run is the entry point for the version command.
func (c *cmdVersionRunner) Run(ctx context.Context, args *clip.CommandArgs[environ]) error {
	c.mustGetopt(args)
	mustFprintf(args.Env.Stdout(), "%s\n", Version)
	return nil
}

// mustGetopt gets command line options.
func (c *cmdVersionRunner) mustGetopt(args *clip.CommandArgs[environ]) {
	// Create empty command line parser.
	clp := flag.NewFlagSet(args.CommandName, flag.ExitOnError)
	clp.SetDescription(args.Command.BriefDescription())
	clp.SetArgsDocs("")

	// Parse the command line arguments.
	clip.Must(args.Env, clp.Parse(args.Args))
	clip.Must(args.Env, clp.PositionalArgsEqualCheck(0))
}
