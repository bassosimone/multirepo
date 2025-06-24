// cmdforeach.go - implementation of the foreach command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/flag"
	"github.com/bassosimone/clip/pkg/parser"
	"github.com/kballard/go-shellquote"
)

// cmdForeach is the static foreach command
var cmdForeach = &clip.LeafCommand[environ]{
	BriefDescriptionText: "Execute a command in each repository.",
	RunFunc:              (&cmdForeachRunner{}).Run,
}

// cmdForeachRunner runs the foreach command.
type cmdForeachRunner struct {
	// Argv contains the command and its arguments.
	Argv []string

	// KeepGoing indicates whether to continue executing commands even if one fails.
	KeepGoing bool

	// XWriter is the writer used to log executed commands.
	XWriter io.Writer
}

// --- entry & setup ---

// Run is the entry point for the foreach command.
func (c *cmdForeachRunner) Run(ctx context.Context, args *clip.CommandArgs[environ]) error {
	c.mustGetopt(args)
	return c.run(ctx, args)
}

// mustGetopt gets command line options.
func (c *cmdForeachRunner) mustGetopt(args *clip.CommandArgs[environ]) {
	// Initialize the default configuration.
	c.Argv = []string{}
	c.KeepGoing = false
	c.XWriter = io.Discard

	// Create empty command line parser.
	clp := flag.NewFlagSet(args.CommandName, flag.ExitOnError)
	clp.SetDescription(args.Command.BriefDescription())
	clp.SetArgsDocs("command [args...]")

	// Add the `-k` flag.
	kflag := clp.Bool("keep-going", 'k', false, "Continue iterating even if the subcommand fails.")

	// Add the `-x` flag.
	xflag := clp.Bool("print-commands", 'x', false, "Log the commands we execute.")

	// Disable option permuation to allow passing options to subcommands
	clp.Parser().Flags |= parser.FlagNoPermute

	// Parse the command line arguments.
	clip.Must(args.Env, clp.Parse(args.Args))
	clip.Must(args.Env, clp.PositionalArgsRangeCheck(1, math.MaxInt))

	// Add the command to execute.
	c.Argv = clp.Args()

	// Honour the `-k` flag.
	if *kflag {
		c.KeepGoing = true
	}

	// Honour the `-x` flag.
	if *xflag {
		c.XWriter = args.Env.Stderr()
	}
}

// --- execution ---

func (c *cmdForeachRunner) run(ctx context.Context, args *clip.CommandArgs[environ]) error {
	// Lock the multirepo dir
	dd := defaultDotDir()
	unlock, err := dd.lock(args.Env)
	if err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo foreach: %s\n", err)
		return err
	}
	defer unlock()

	// Read the configuration file
	cinfo, err := readConfig(args.Env, dd.configFilePath())
	if err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo foreach: %s\n", err)
		return err
	}

	// Execute command in each repository
	errlist := []error{}
	for repo := range cinfo.Repos {
		if err := c.execute(ctx, args.Env, repo); err != nil {
			mustFprintf(args.Env.Stderr(), "multirepo foreach: %s\n", err)
			errlist = append(errlist, err)
			if !c.KeepGoing {
				break
			}
		}
	}

	return errors.Join(errlist...)
}

// execute executes the command in a given repository.
func (c *cmdForeachRunner) execute(ctx context.Context, env environ, repo string) error {
	// Preparing for adding to the environment variables.
	environ := os.Environ()

	// Conditionally add the `MULTIREPO_ROOT` environment variable.
	if os.Getenv("MULTIREPO_ROOT") == "" {
		wdir, err := env.Getwd()
		if err != nil {
			return err
		}
		environ = append(environ, fmt.Sprintf("MULTIREPO_ROOT=%s", wdir))
	}

	// Conditionally add the `MULTIREPO_EXECUTABLE` environment variable.
	if os.Getenv("MULTIREPO_EXECUTABLE") == "" {
		exe, err := env.Executable()
		if err != nil {
			return err
		}
		exe, err = env.AbsFilepath(exe)
		if err != nil {
			return err
		}
		environ = append(environ, fmt.Sprintf("MULTIREPO_EXECUTABLE=%s", exe))
	}

	// Create the subcommand to execute.
	assert.True(len(c.Argv) >= 1, "expected at least the command name")
	cmd := exec.CommandContext(ctx, c.Argv[0], c.Argv[1:]...)
	cmd.Stdin = io.NopCloser(bytes.NewReader(nil))
	cmd.Stdout = env.Stdout()
	cmd.Stderr = env.Stderr()
	cmd.Dir = repo
	cmd.Env = environ

	// Log that we're executing the command.
	mustFprintf(c.XWriter, "+ (cd %s && %s)\n", shellquote.Join(repo), shellquote.Join(cmd.Args...))

	// Execute the command
	return env.RunCommand(cmd)
}
