// cmdclone.go - implementation of the clone command.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/flag"
	"github.com/kballard/go-shellquote"
)

// cloneCmd is the static clone command.
var cloneCmd = &clip.LeafCommand[environ]{
	BriefDescriptionText: "Clone a repository into the multirepo.",
	RunFunc:              (&cmdClone{}).Run,
}

// cmdClone implements the clone command.
type cmdClone struct {
	// Repo is the repository to clone.
	Repo string

	// VWriterStderr is the writer executed to log the executed commands stderr.
	VWriterStderr io.Writer

	// VWriterStdout is the writer executed to log the executed commands stdout.
	VWriterStdout io.Writer

	// XWriter is the writer used to log executed commands.
	XWriter io.Writer
}

// Run is the entry point for the clone command.
func (c *cmdClone) Run(ctx context.Context, args *clip.CommandArgs[environ]) error {
	// Parse command line arguments
	c.mustGetopt(args)

	// Lock the multirepo dir
	dd := defaultDotDir()
	unlock, err := dd.lock(args.Env)
	if err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo clone: %s\n", err)
		return err
	}
	defer unlock()

	// Clone the repository
	if err := c.clone(ctx, args.Env, dd); err != nil {
		mustFprintf(args.Env.Stderr(), "multirepo clone: %s\n", err)
		return err
	}

	return nil
}

// mustGetopt gets command line options.
func (c *cmdClone) mustGetopt(args *clip.CommandArgs[environ]) {
	// Initialize the default configuration.
	c.Repo = ""
	c.VWriterStderr = io.Discard
	c.VWriterStdout = io.Discard
	c.XWriter = io.Discard

	// Create empty command line parser.
	clp := flag.NewFlagSet(args.CommandName, flag.ExitOnError)
	clp.SetDescription(args.Command.BriefDescription())
	clp.SetArgsDocs("git@github.com:user/repo")

	// Add the `-v` flag.
	vflag := clp.Bool("verbose", 'v', false, "Show the output of git clone.")

	// Add the `-x` flag.
	xflag := clp.Bool("print-commands", 'x', false, "Log the commands we execute.")

	// Parse the command line arguments.
	clip.Must(args.Env, clp.Parse(args.Args))
	clip.Must(args.Env, clp.PositionalArgsEqualCheck(1))

	// Set the repository name to clone.
	c.Repo = clp.Args()[0]

	// Honour the `-v` flag.
	if *vflag {
		c.VWriterStderr = args.Env.Stderr()
		c.VWriterStdout = args.Env.Stdout()
	}

	// Honour the `-x` flag.
	if *xflag {
		c.XWriter = args.Env.Stderr()
	}
}

// clone clones a repository.
func (c *cmdClone) clone(ctx context.Context, env environ, dd dotDir) error {
	// Read the configuration file.
	cinfo, err := readConfig(env, dd.configFilePath())
	if err != nil {
		return err
	}

	// Parses the scp-like URL.
	scpInfo, good := scpLikeParse(c.Repo)
	if !good {
		return fmt.Errorf("invalid repository URL: %s", c.Repo)
	}

	// Create the subcommand to execute.
	cmd := exec.CommandContext(ctx, "git", "clone", scpInfo.String(), scpInfo.Name())
	cmd.Stdin = io.NopCloser(bytes.NewReader(nil))
	cmd.Stdout = c.VWriterStdout
	cmd.Stderr = c.VWriterStderr

	// Log that we're executing the command.
	mustFprintf(c.XWriter, "+ %s\n", shellquote.Join(cmd.Args...))

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
