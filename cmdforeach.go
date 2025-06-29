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
	"os/exec"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nflag"
	"github.com/kballard/go-shellquote"
)

// cmdForeach is the static foreach command
var cmdForeach = &clip.LeafCommand[environ]{
	BriefDescriptionText: "Execute a command in each repository.",
	RunFunc:              cmdForeachMain,
}

// cmdForeachRunner runs the foreach command.
type cmdForeachRunner struct {
	// Argv contains the command and its arguments.
	Argv []string

	// KeepGoing indicates whether to continue executing commands even if one fails.
	KeepGoing bool

	// Style is the nil-safe lipgloss style to use.
	Style *nilSafeLipglossStyle

	// XWriter is the writer used to log executed commands.
	XWriter io.Writer
}

// --- entry & setup ---

// cmdForeachMain is the entry point for the foreach command.
func cmdForeachMain(ctx context.Context, args *clip.CommandArgs[environ]) error {
	return mustNewCmdForeachRunner(args).run(ctx, args)
}

// mustNewCmdForeachRunner creates a new [*cmdForeachRunner].
func mustNewCmdForeachRunner(args *clip.CommandArgs[environ]) *cmdForeachRunner {
	// Initialize the default configuration.
	c := &cmdForeachRunner{
		Argv:      []string{},
		KeepGoing: false,
		Style:     nil,
		XWriter:   io.Discard,
	}

	// Create empty command line parser.
	fset := nflag.NewFlagSet(args.CommandName, nflag.ExitOnError)
	fset.Description = args.Command.BriefDescription()
	fset.PositionalArgumentsUsage = "<command> [args...]"
	fset.DisablePermute = true // Disable option permutaion to allow passing options to subcommands
	fset.MinPositionalArgs = 1
	fset.MaxPositionalArgs = math.MaxInt

	// Add the `-h, --help` flag.
	fset.AutoHelp("help", 'h', "Show this help message and exit.")

	// Add the `-k` flag.
	kflag := fset.Bool("keep-going", 'k', "Continue iterating even if the subcommand fails.")

	// Add the `-x` flag.
	xflag := fset.Bool("print-commands", 'x', "Log the commands we execute.")

	// Parse the command line arguments.
	assert.NotError(fset.Parse(args.Args))

	// Add the command to execute.
	c.Argv = fset.Args()

	// Honour the `-k` flag.
	if *kflag {
		c.KeepGoing = true
	}

	// Honour the `-x` flag.
	if *xflag {
		c.XWriter = args.Env.Stderr()
		c.Style = newNilSafeLipglossStyle()
	}

	return c
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
	config, err := readConfig(args.Env, dd.configFilePath())
	if err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo foreach: %s\n", err)
		return err
	}

	// Execute command in each repository
	errlist := []error{}
	for repo := range config.Repos {
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
	environ := env.Environ()

	// Conditionally add the `MULTIREPO_ROOT` environment variable.
	if _, found := env.LookupEnv("MULTIREPO_ROOT"); !found {
		wdir, err := env.Getwd()
		if err != nil {
			return err
		}
		variable := fmt.Sprintf("MULTIREPO_ROOT=%s", wdir)
		environ = append(environ, variable)
	}

	// Conditionally add the `MULTIREPO_EXECUTABLE` environment variable.
	if _, found := env.LookupEnv("MULTIREPO_EXECUTABLE"); !found {
		exe, err := env.Executable()
		if err != nil {
			return err
		}
		exe, err = env.AbsFilepath(exe)
		if err != nil {
			return err
		}
		variable := fmt.Sprintf("MULTIREPO_EXECUTABLE=%s", exe)
		environ = append(environ, variable)
	}

	// As an optimization, if we're running `git`, prepend the
	// `--no-pager` option otherwise... it's painful!
	assert.True(len(c.Argv) >= 1, "expected at least the command name")
	reargv := []string{c.Argv[0]}
	if reargv[0] == "git" {
		reargv = append(reargv, "--no-pager")
	}
	reargv = append(reargv, c.Argv[1:]...)

	// Create the subcommand to execute.
	cmd := exec.CommandContext(ctx, reargv[0], reargv[1:]...)
	cmd.Stdin = io.NopCloser(bytes.NewReader(nil))
	cmd.Stdout = env.Stdout()
	cmd.Stderr = env.Stderr()
	cmd.Dir = repo
	cmd.Env = environ

	// Log that we're executing the command.
	//
	// Add a newline before each entry so that it stands out when
	// skimming the terminal. Note that we cannot make `-x` the
	// default, since it would be quite annoying when reading diffs
	mustFprintf(c.XWriter, "%s\n", c.Style.Renderf("+ (cd %s && %s)", shellquote.Join(repo), shellquote.Join(cmd.Args...)))

	// Execute the command
	return env.RunCommand(cmd)
}
