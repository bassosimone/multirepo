// cmdversion.go - implementation of the version command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/flag"
)

// versionCmd is the static version command
var versionCmd = &clip.LeafCommand[environ]{
	BriefDescriptionText: "Display the version of the tool.",
	RunFunc:              (&cmdVersion{}).Run,
}

// cmdVersion implements the version command.
type cmdVersion struct{}

// Run is the entry point for the version command.
func (c *cmdVersion) Run(ctx context.Context, args *clip.CommandArgs[environ]) error {
	// Parse command line arguments
	if err := c.getopt(args); err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo version: %s\n", err)
		return err
	}

	// Print version information
	mustFprintf(args.Env.Stdout(), "%s\n", Version)
	return nil
}

// getopt gets command line options.
func (c *cmdVersion) getopt(args *clip.CommandArgs[environ]) error {
	// Create empty command line parser.
	clp := flag.NewFlagSet(args.CommandName, flag.ContinueOnError)
	clp.SetDescription(args.Command.BriefDescription())
	clp.SetArgsDocs("")

	// Parse the command line arguments.
	if err := clp.Parse(args.Args); err != nil {
		return err
	}
	if err := clp.PositionalArgsEqualCheck(0); err != nil {
		return err
	}

	return nil
}
