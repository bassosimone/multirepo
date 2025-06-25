// cmdrepols.go - implementation of the 'repo ls' command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/flag"
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
	fset := flag.NewFlagSet(args.CommandName, flag.ExitOnError)
	fset.SetDescription(args.Command.BriefDescription())
	fset.SetArgsDocs("<repo>")

	// Parse the command line arguments.
	clip.Must(args.Env, fset.Parse(args.Args))
	clip.Must(args.Env, fset.PositionalArgsEqualCheck(0))

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
