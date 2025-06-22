// cmdclone.go - implementation of the clone command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/bassosimone/clip/pkg/flag"
	"github.com/kballard/go-shellquote"
)

// cmdClone implements the clone command.
type cmdClone struct{}

var _ command = (*cmdClone)(nil)

// Description implements [command].
func (c *cmdClone) Description() string {
	return "Clones a repository into the multirepo."
}

// Run implements [command].
func (c *cmdClone) Run(ctx context.Context, env environ, argv cliArgs) error {
	// Print help if requested to do so by the user.
	if argv.ContainsHelp() {
		return c.help(env.Stdout())
	}

	// Parse command line arguments
	options, err := c.getopt(env, argv)
	if err != nil {
		mustFprintf(env.Stderr(), "multirepo clone: %s\n", err)
		mustFprintf(env.Stderr(), "Try `multirepo clone --help` for help.\n")
		return err
	}

	// Lock the multirepo dir
	dd := defaultDotDir()
	unlock, err := dd.lock(env)
	if err != nil {
		mustFprintf(env.Stderr(), "multirepo clone: %s\n", err)
		return err
	}
	defer unlock()

	// Clone the repository
	if err := c.clone(ctx, env, options, dd, options.Repo); err != nil {
		mustFprintf(env.Stderr(), "multirepo clone: %s\n", err)
		return err
	}

	return nil
}

// help prints the help message for the clone command.
func (c *cmdClone) help(w io.Writer) error {
	mustFprintf(w, "\n")
	mustFprintf(w, "multirepo clone - %s\n", c.Description())
	mustFprintf(w, "\n")
	mustFprintf(w, "This command clones a repository into the multirepo.\n")
	mustFprintf(w, "\n")
	mustFprintf(w, "usage: multirepo clone [-vx] {repo}\n")
	mustFprintf(w, "\n")
	mustFprintf(w, "Flags:\n")
	mustFprintf(w, "  -v, --verbose         Print output of executed commands\n")
	mustFprintf(w, "  -x, --print-commands  Print commands as they are executed\n")
	mustFprintf(w, "\n")
	return nil
}

// cmdCloneOptions contains configuration for the clone command.
type cmdCloneOptions struct {
	// Repo is the repository to clone.
	Repo string

	// VWriterStderr is the writer executed to log the executed commands stderr.
	VWriterStderr io.Writer

	// VWriterStdout is the writer executed to log the executed commands stdout.
	VWriterStdout io.Writer

	// XWriter is the writer used to log executed commands.
	XWriter io.Writer
}

// getopt gets command line options.
func (c *cmdClone) getopt(env environ, argv cliArgs) (*cmdCloneOptions, error) {
	// Initialize the default configuration.
	options := &cmdCloneOptions{
		Repo:          "",
		VWriterStderr: io.Discard,
		VWriterStdout: io.Discard,
		XWriter:       io.Discard,
	}

	// Create empty command line parser.
	clp := flag.NewFlagSet("", flag.ContinueOnError)

	// Add the `-v` flag.
	vflag := clp.Bool("verbose", 'v', false, "")

	// Add the `-x` flag.
	xflag := clp.Bool("print-commands", 'x', false, "")

	// Parse the command line arguments.
	if err := clp.Parse(argv.CommandArgs()); err != nil {
		return nil, err
	}

	args := clp.Args()
	if len(args) != 1 {
		return nil, fmt.Errorf("expected exactly one repository")
	}

	// Set the repository name to clone.
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
func (c *cmdClone) clone(
	ctx context.Context,
	env environ,
	options *cmdCloneOptions,
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
