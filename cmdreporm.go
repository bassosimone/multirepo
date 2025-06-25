// cmdreporm.go - implementation of the 'repo rm' command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/flag"
)

// cmdRepoRm is the static 'repo rm' command
var cmdRepoRm = &clip.LeafCommand[environ]{
	BriefDescriptionText: "Remove a repository from the index.",
	RunFunc:              cmdRepoRmMain,
}

// cmdRepoRmRunner runs the 'repo rm' command.
type cmdRepoRmRunner struct{}

// --- entry & setup ---

// cmdRepRmMain is the entry point for the 'repo rm' command.
func cmdRepoRmMain(ctx context.Context, args *clip.CommandArgs[environ]) error {
	return mustNewCmdRepoRmRunner(args).run(args)
}

// mustNewCmdRepoRmRunner creates a new [*cmdRepoRmRunner].
func mustNewCmdRepoRmRunner(args *clip.CommandArgs[environ]) *cmdRepoRmRunner {
	// initialize the default configuration.
	c := &cmdRepoRmRunner{}

	// Create empty command line parser.
	clp := flag.NewFlagSet(args.CommandName, flag.ExitOnError)
	clp.SetDescription(args.Command.BriefDescription())
	clp.SetArgsDocs("<repo>")

	// Parse the command line arguments.
	clip.Must(args.Env, clp.Parse(args.Args))
	clip.Must(args.Env, clp.PositionalArgsEqualCheck(1))

	return c
}

// --- execution ---

func (c *cmdRepoRmRunner) run(args *clip.CommandArgs[environ]) error {
	// Lock the multirepo dir
	dd := defaultDotDir()
	unlock, err := dd.lock(args.Env)
	if err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo repo rm: %s\n", err)
		return err
	}
	defer unlock()

	// Read the configuration file
	config, err := readConfig(args.Env, dd.configFilePath())
	if err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo repo rm: %s\n", err)
		return err
	}

	// Remove from the configuration
	assert.True(len(args.Args) >= 1, "expected the repository directory path")
	delete(config.Repos, args.Args[0])

	// Write the configuration file back to disk
	if err := config.WriteFile(args.Env, dd.configFilePath()); err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo repo rm: %s\n", err)
		return err
	}

	return nil
}
