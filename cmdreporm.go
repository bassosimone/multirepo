// cmdreporm.go - implementation of the 'repo rm' command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nflag"
)

// cmdRepoRm is the static 'repo rm' command
var cmdRepoRm = &clip.LeafCommand[environ]{
	BriefDescriptionText: "Remove a repository from the index.",
	RunFunc:              cmdRepoRmMain,
}

// cmdRepoRmRunner runs the 'repo rm' command.
type cmdRepoRmRunner struct {
	// Repo is the name of the repository directory to remove.
	Repo string
}

// --- entry & setup ---

// cmdRepRmMain is the entry point for the 'repo rm' command.
func cmdRepoRmMain(ctx context.Context, args *clip.CommandArgs[environ]) error {
	return mustNewCmdRepoRmRunner(args).run(args)
}

// mustNewCmdRepoRmRunner creates a new [*cmdRepoRmRunner].
func mustNewCmdRepoRmRunner(args *clip.CommandArgs[environ]) *cmdRepoRmRunner {
	// initialize the default configuration.
	c := &cmdRepoRmRunner{
		Repo: "",
	}

	// Create empty command line parser.
	fset := nflag.NewFlagSet(args.CommandName, nflag.ExitOnError)
	fset.Description = args.Command.BriefDescription()
	fset.PositionalArgumentsUsage = "<repo>"
	fset.MinPositionalArgs = 1
	fset.MaxPositionalArgs = 1

	// Add the `-h, --help` flag.
	fset.AutoHelp("help", 'h', "Show this help message and exit.")

	// Parse the command line arguments.
	assert.NotError(fset.Parse(args.Args))

	// Add the repo to remove
	c.Repo = fset.Args()[0]
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
	delete(config.Repos, c.Repo)

	// Write the configuration file back to disk
	if err := config.WriteFile(args.Env, dd.configFilePath()); err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo repo rm: %s\n", err)
		return err
	}

	return nil
}
