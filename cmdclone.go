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

// cmdClone is the static clone command.
var cmdClone = &clip.LeafCommand[environ]{
	BriefDescriptionText: "Clone a repository into the multirepo.",
	RunFunc:              cmdCloneMain,
}

// cmdCloneRunner runs the clone command.
type cmdCloneRunner struct {
	// Repo is the repository to clone.
	Repo string

	// Style is the nil-safe lipgloss style to use.
	Style *nilSafeLipglossStyle

	// VWriterStderr is the writer executed to log the executed commands stderr.
	VWriterStderr io.Writer

	// VWriterStdout is the writer executed to log the executed commands stdout.
	VWriterStdout io.Writer

	// XWriter is the writer used to log executed commands.
	XWriter io.Writer
}

// --- entry & setup ---

// cmdCloneMain is the entry point for the clone command.
func cmdCloneMain(ctx context.Context, args *clip.CommandArgs[environ]) error {
	return mustNewCmdCloneRunner(args).run(ctx, args)
}

// mustNewCmdCloneRunner creates a new [*cmdCloneRunner].
func mustNewCmdCloneRunner(args *clip.CommandArgs[environ]) *cmdCloneRunner {
	// Initialize the default configuration.
	c := &cmdCloneRunner{
		Repo:          "",
		Style:         nil,
		VWriterStderr: io.Discard,
		VWriterStdout: io.Discard,
		XWriter:       io.Discard,
	}

	// Create empty command line parser.
	fset := flag.NewFlagSet(args.CommandName, flag.ExitOnError)
	fset.SetDescription(args.Command.BriefDescription())
	fset.SetArgsDocs("git@github.com:user/repo")

	// Add the `-v` flag.
	vflag := fset.Bool("verbose", 'v', "Show the output of git clone.")

	// Add the `-x` flag.
	xflag := fset.Bool("print-commands", 'x', "Log the commands we execute.")

	// Parse the command line arguments.
	_ = fset.Parse(args.Args)
	_ = fset.PositionalArgsEqualCheck(1)

	// Set the repository name to clone.
	c.Repo = fset.Args()[0]

	// Honour the `-v` flag.
	if *vflag {
		c.VWriterStderr = args.Env.Stderr()
		c.VWriterStdout = args.Env.Stdout()
	}

	// Honour the `-x` flag.
	if *xflag {
		c.XWriter = args.Env.Stderr()
		c.Style = newNilSafeLipglossStyle()
	}

	return c
}

// --- execution ---

func (c *cmdCloneRunner) run(ctx context.Context, args *clip.CommandArgs[environ]) error {
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

// clone clones a repository.
func (c *cmdCloneRunner) clone(ctx context.Context, env environ, dd dotDir) error {
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
	mustFprintf(c.XWriter, "%s\n", c.Style.Renderf("+ %s", shellquote.Join(cmd.Args...)))

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
