// cmdforeach.go - implementation of the foreach command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/bassosimone/clip/pkg/flag"
	"github.com/bassosimone/clip/pkg/parser"
	"github.com/kballard/go-shellquote"
)

// cmdForeach implements the foreach command.
type cmdForeach struct{}

var _ command = (*cmdForeach)(nil)

// Description implements [command].
func (c *cmdForeach) Description() string {
	return "Executes a command in each repository."
}

// Run implements [command].
func (c *cmdForeach) Run(ctx context.Context, env environ, argv cliArgs) error {
	// Print help if requested to do so by the user.
	if argv.ContainsHelp() {
		return c.help(env.Stdout())
	}

	// Parse command line arguments
	options, err := c.getopt(env, argv)
	if err != nil {
		mustFprintf(env.Stderr(), "multirepo foreach: %s\n", err)
		mustFprintf(env.Stderr(), "Try `multirepo foreach --help` for help.\n")
		return err
	}

	// Lock the multirepo dir
	dd := defaultDotDir()
	unlock, err := dd.lock(env)
	if err != nil {
		mustFprintf(env.Stderr(), "multirepo foreach: %s\n", err)
		return err
	}
	defer unlock()

	// Read the configuration file
	cinfo, err := readConfig(env, dd.configFilePath())
	if err != nil {
		mustFprintf(env.Stderr(), "multirepo foreach: %s\n", err)
		return err
	}

	// Execute command in each repository
	errlist := []error{}
	for repo := range cinfo.Repos {
		if err := c.execute(ctx, env, options, repo); err != nil {
			mustFprintf(env.Stderr(), "multirepo foreach: %s\n", err)
			errlist = append(errlist, err)
			if !options.KeepGoing {
				break
			}
		}
	}

	return errors.Join(errlist...)
}

// help prints the help message for the foreach command.
func (c *cmdForeach) help(w io.Writer) error {
	mustFprintf(w, "\n")
	mustFprintf(w, "multirepo foreach - %s\n", c.Description())
	mustFprintf(w, "\n")
	mustFprintf(w, "This command executes a command in each repository.\n")
	mustFprintf(w, "\n")
	mustFprintf(w, "usage: multirepo foreach [-x] {repo}\n")
	mustFprintf(w, "\n")
	mustFprintf(w, "Flags:\n")
	mustFprintf(w, "  -k, --keep-going      Continue executing commands even if one fails\n")
	mustFprintf(w, "  -x, --print-commands  Print commands as they are executed\n")
	mustFprintf(w, "\n")
	return nil
}

// cmdForeachOptions contains configuration for the foreach command.
type cmdForeachOptions struct {
	// Argv contains the command and its arguments.
	Argv []string

	// KeepGoing indicates whether to continue executing commands even if one fails.
	KeepGoing bool

	// XWriter is the writer used to log executed commands.
	XWriter io.Writer
}

// getopt gets command line options.
func (c *cmdForeach) getopt(env environ, argv cliArgs) (*cmdForeachOptions, error) {
	// Initialize the default configuration.
	options := &cmdForeachOptions{
		Argv:      []string{},
		KeepGoing: false,
		XWriter:   io.Discard,
	}

	// Create empty command line parser.
	clp := flag.NewFlagSet("", flag.ContinueOnError)

	// Add the `-k` flag.
	kflag := clp.Bool("keep-going", 'k', false, "")

	// Add the `-x` flag.
	xflag := clp.Bool("print-commands", 'x', false, "")

	// Disable option permuation to allow passing options to subcommands
	clp.Parser().Flags |= parser.FlagNoPermute

	// Parse the command line arguments.
	if err := clp.Parse(argv.CommandArgs()); err != nil {
		return nil, err
	}

	args := clp.Args()
	if len(args) < 1 {
		return nil, fmt.Errorf("expected at least the command name")
	}

	// Add the command to execute.
	options.Argv = args

	// Honour the `-k` flag.
	if *kflag {
		options.KeepGoing = true
	}

	// Honour the `-x` flag.
	if *xflag {
		options.XWriter = env.Stderr()
	}

	// Return the configuration.
	return options, nil
}

// execute executes the command in a given repository.
func (c *cmdForeach) execute(
	ctx context.Context,
	env environ,
	options *cmdForeachOptions,
	repo string,
) error {
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
	assertTrue(len(options.Argv) >= 1, "expected at least the command name")
	cmd := exec.CommandContext(ctx, options.Argv[0], options.Argv[1:]...)
	cmd.Stdin = io.NopCloser(bytes.NewReader(nil))
	cmd.Stdout = env.Stdout()
	cmd.Stderr = env.Stderr()
	cmd.Dir = repo
	cmd.Env = environ

	// Log that we're executing the command.
	mustFprintf(options.XWriter, "+ (cd %s && %s)\n", shellquote.Join(repo),
		shellquote.Join(cmd.Args...))

	// Execute the command
	return env.RunCommand(cmd)
}
