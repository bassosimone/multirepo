// main.go - Main file.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"os"
)

// errorToExitCode maps an error to an exit code.
func errorToExitCode(err error) int {
	emap := map[bool]int{
		true:  1,
		false: 0,
	}
	return emap[err != nil]
}

// osExit allows mocking the [os.Exit] function.
var osExit = os.Exit

func main() {
	// Initialize the environment
	env := &stdlibEnviron{}

	// Initialize the context
	ctx := context.Background()

	// Initialize the main command
	cmd := &cmdMain{}

	// Initialize the CLI args
	argv := cliArgs(os.Args)

	// Run the main command
	osExit(errorToExitCode(cmd.Run(ctx, env, argv)))
}
