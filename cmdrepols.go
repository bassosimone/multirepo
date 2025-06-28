// cmdrepols.go - implementation of the 'repo ls' command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nflag"
)

// cmdRepoLs is the static 'repo ls' command
var cmdRepoLs = &clip.LeafCommand[environ]{
	BriefDescriptionText: "Remove a repository from the index.",
	RunFunc:              cmdRepoLsMain,
}

// cmdRepoLsRunner runs the 'repo ls' command.
type cmdRepoLsRunner struct{}

// --- entry & setup ---

// cmdRepRmMain is the entry point for the 'repo ls' command.
func cmdRepoLsMain(ctx context.Context, args *clip.CommandArgs[environ]) error {
	return mustNewCmdRepoLsRunner(args).run(args)
}

// mustNewCmdRepoLsRunner creates a new [*cmdRepoLsRunner].
func mustNewCmdRepoLsRunner(args *clip.CommandArgs[environ]) *cmdRepoLsRunner {
	// initialize the default configuration.
	c := &cmdRepoLsRunner{}

	// Create empty command line parser.
	fset := nflag.NewFlagSet(args.CommandName, nflag.ExitOnError)
	fset.Description = args.Command.BriefDescription()
	fset.PositionalArgumentsUsage = "<repo>"
	fset.MinPositionalArgs = 0
	fset.MaxPositionalArgs = 0

	// Add the `-h, --help` flag.
	fset.AutoHelp("help", 'h', "Show this help message and exit.")

	// Parse the command line arguments.
	assert.NotError(fset.Parse(args.Args))

	return c
}

// --- execution ---

func (c *cmdRepoLsRunner) run(args *clip.CommandArgs[environ]) error {
	// Lock the multirepo dir
	dd := defaultDotDir()
	unlock, err := dd.lock(args.Env)
	if err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo repo ls: %s\n", err)
		return err
	}
	defer unlock()

	// Read the configuration file
	config, err := readConfig(args.Env, dd.configFilePath())
	if err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo repo ls: %s\n", err)
		return err
	}

	// Print each entry in the configuration file.
	for _, name := range slices.Sorted(maps.Keys(config.Repos)) {
		fmt.Fprintf(args.Env.Stdout(), "%-24s %s\n", name, config.Repos[name].URL)
	}
	return nil
}
