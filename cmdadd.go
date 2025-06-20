// cmdadd.go - implementation of the add command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/kballard/go-shellquote"
	"github.com/spf13/pflag"
)

// cmdAdd implements the add command.
type cmdAdd struct{}

var _ command = (*cmdAdd)(nil)

// Description implements [command].
func (c *cmdAdd) Description() string {
	return "Adds one or more repositories to the multirepo."
}

// Run implements [command].
func (c *cmdAdd) Run(ctx context.Context, env environ, argv cliArgs) error {
	// Print help if requested to do so by the user.
	if argv.ContainsHelp() {
		return c.help(env.Stdout())
	}

	// Parse command line arguments
	options, err := c.getopt(env, argv)
	if err != nil {
		mustFprintf(env.Stderr(), "multirepo add: %s\n", err)
		mustFprintf(env.Stderr(), "Try `multirepo add --help` for help.\n")
		return err
	}

	// Lock the multirepo dir
	dd := defaultDotDir()
	unlock, err := dd.lock(env)
	if err != nil {
		mustFprintf(env.Stderr(), "multirepo add: %s\n", err)
		return err
	}
	defer unlock()

	// Clone the repository
	if err := c.clone(ctx, env, options, dd, options.Repo); err != nil {
		mustFprintf(env.Stderr(), "multirepo add: %s\n", err)
		return err
	}

	return nil
}

// help prints the help message for the add command.
func (c *cmdAdd) help(w io.Writer) error {
	mustFprintf(w, "\n")
	mustFprintf(w, "add - %s\n", c.Description())
	mustFprintf(w, "\n")
	mustFprintf(w, "usage: multirepo add [-vx] {repo}\n")
	mustFprintf(w, "\n")
	mustFprintf(w, "Flags:\n")
	mustFprintf(w, "  -v, --verbose         Print output of executed commands\n")
	mustFprintf(w, "  -x, --print-commands  Print commands as they are executed\n")
	mustFprintf(w, "\n")
	mustFprintf(w, "This command adds a repository to the multirepo.\n")
	mustFprintf(w, "\n")
	return nil
}

// cmdAddOptions contains configuration for the add command.
type cmdAddOptions struct {
	// Repo is the repository to add.
	Repo string

	// VWriterStderr is the writer executed to log the executed commands stderr.
	VWriterStderr io.Writer

	// VWriterStdout is the writer executed to log the executed commands stdout.
	VWriterStdout io.Writer

	// XWriter is the writer used to log executed commands.
	XWriter io.Writer
}

// getopt gets command line options.
func (c *cmdAdd) getopt(env environ, argv cliArgs) (*cmdAddOptions, error) {
	// Initialize the default configuration.
	options := &cmdAddOptions{
		Repo:          "",
		VWriterStderr: io.Discard,
		VWriterStdout: io.Discard,
		XWriter:       io.Discard,
	}

	// Create empty command line parser.
	clip := pflag.NewFlagSet("", pflag.ContinueOnError)

	// Add the `-v` flag.
	vflag := clip.BoolP("verbose", "v", false, "")

	// Add the `-x` flag.
	xflag := clip.BoolP("print-commands", "x", false, "")

	// Parse the command line arguments.
	if err := clip.Parse(argv.CommandArgs()); err != nil {
		return nil, err
	}

	args := clip.Args()
	if len(args) != 1 {
		return nil, fmt.Errorf("expected exactly one repository")
	}

	// Add the repository names to the list of repositories to add.
	options.Repo = args[0]

	// Honour the `-v` flag.
	if *vflag {
		options.VWriterStderr = env.Stderr()
		options.VWriterStdout = env.Stdout()
	}

	// Honour the `-x` flag.
	if *xflag {
		options.XWriter = env.Stderr()
	}

	// Return the configuration.
	return options, nil
}

// clone clones a repository.
func (c *cmdAdd) clone(
	ctx context.Context,
	env environ,
	options *cmdAddOptions,
	dd dotDir,
	repo string,
) error {
	// Read the configuration file.
	cinfo, err := readConfig(env, dd.configFilePath())
	if err != nil {
		return err
	}

	// Parses the scp-like repository URL.
	scpInfo, good := scpLikeParse(repo)
	if !good {
		return fmt.Errorf("invalid repository URL: %s", repo)
	}

	// Create the subcommand to execute.
	cmd := exec.CommandContext(ctx, "git", "clone", scpInfo.String(), scpInfo.Name())
	cmd.Stdin = io.NopCloser(bytes.NewReader(nil))
	cmd.Stdout = options.VWriterStdout
	cmd.Stderr = options.VWriterStderr

	// Log that we're executing the command.
	mustFprintf(options.XWriter, "+ %s\n", shellquote.Join(cmd.Args...))

	// Execute the command
	if err := env.RunCommand(cmd); err != nil {
		return err
	}

	// Update the configuration file.
	cinfo.AddRepo(scpInfo.Name(), scpInfo.String())
	if err := cinfo.WriteFile(env, dd.configFilePath()); err != nil {
		return err
	}

	return nil
}
