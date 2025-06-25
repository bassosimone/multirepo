// cmdrepoadd.go - implementation of the 'repo add' command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"context"
	"io"
	"math"
	"os/exec"
	"strings"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/flag"
	"github.com/kballard/go-shellquote"
)

// cmdRepoAdd is the static 'repo add' command
var cmdRepoAdd = &clip.LeafCommand[environ]{
	BriefDescriptionText: "Add an existing repository to the multirepo.",
	RunFunc:              cmdRepoAddMain,
}

// cmdRepoAddRunner runs the 'repo add' command.
type cmdRepoAddRunner struct {
	// Repo is the repository directory name to add.
	Repos []string

	// Style is the nil-safe lipgloss style to use.
	Style *nilSafeLipglossStyle

	// XWriter is the writer used to log executed commands.
	XWriter io.Writer
}

// --- entry & setup ---

// cmdRepoAddMain is the entry point for the 'repo add' command.
func cmdRepoAddMain(ctx context.Context, args *clip.CommandArgs[environ]) error {
	return mustNewCmdRepoAddRunner(args).run(ctx, args)
}

// mustNewCmdRepoAddRunner creates a new [*cmdRepoAddRunner].
func mustNewCmdRepoAddRunner(args *clip.CommandArgs[environ]) *cmdRepoAddRunner {
	// Initialize the default configuration.
	c := &cmdRepoAddRunner{
		Repos:   []string{},
		Style:   nil,
		XWriter: io.Discard,
	}

	// Create empty command line parser.
	fset := flag.NewFlagSet(args.CommandName, flag.ExitOnError)
	fset.SetDescription(args.Command.BriefDescription())
	fset.SetArgsDocs("<repo>")

	// Add the `-x` flag.
	xflag := fset.Bool("print-commands", 'x', "Log the commands we execute.")

	// Parse the command line arguments.
	_ = fset.Parse(args.Args)
	_ = fset.PositionalArgsRangeCheck(1, math.MaxInt)

	// Honour the `-x` flag.
	if *xflag {
		c.XWriter = args.Env.Stderr()
		c.Style = newNilSafeLipglossStyle()
	}

	// Add the repo to add to the multirepo index.
	c.Repos = fset.Args()

	return c
}

// --- execution ---

func (c *cmdRepoAddRunner) run(ctx context.Context, args *clip.CommandArgs[environ]) error {
	// Lock the multirepo dir
	dd := defaultDotDir()
	unlock, err := dd.lock(args.Env)
	if err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo repo add: %s\n", err)
		return err
	}
	defer unlock()

	// Read the configuration file
	config, err := readConfig(args.Env, dd.configFilePath())
	if err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo repo add: %s\n", err)
		return err
	}

	// Iterate over the repositories
	for _, repo := range c.Repos {
		// Obtain the URL
		URL, err := c.getrepourl(ctx, args.Env, repo)
		if err != nil {
			mustFprintf(args.Env.Stderr(), "multirepo repo add: %s\n", err)
			return err
		}

		// Update the config
		config.AddRepo(repo, URL)
	}

	// Write the configuration file back to disk
	if err := config.WriteFile(args.Env, dd.configFilePath()); err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo repo add: %s\n", err)
		return err
	}

	return nil
}

// getrepourl obtains the repository URL.
func (c *cmdRepoAddRunner) getrepourl(ctx context.Context, env environ, repo string) (string, error) {
	// Create the subcommand to execute.
	var captured strings.Builder
	cmd := exec.CommandContext(ctx, "git", "config", "--get", "remote.origin.url")
	cmd.Stdin = io.NopCloser(bytes.NewReader(nil))
	cmd.Stdout = &captured
	cmd.Stderr = env.Stderr()
	cmd.Dir = repo

	// Log that we're executing the command.
	//
	// Add a newline before each entry so that it stands out when
	// skimming the terminal. Note that we cannot make `-x` the
	// default, since it would be quite annoying when reading diffs
	mustFprintf(c.XWriter, "%s\n", c.Style.Renderf("+ (cd %s && %s)", shellquote.Join(repo), shellquote.Join(cmd.Args...)))

	// Execute the command
	if err := env.RunCommand(cmd); err != nil {
		return "", err
	}

	// Obtain the trimmed stdout
	URL := strings.TrimSpace(captured.String())

	// Return the URL to the caller
	//
	// Note: the URL would be actually empty if there is no
	// configured remote or its name is not origin. Yet, I am
	// not sold on emitting an error when this happens.
	return URL, nil
}
