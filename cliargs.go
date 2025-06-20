// cliargs.go - Code to manage command-line arguments.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

// cliArgs contains the command-line arguments.
type cliArgs []string

// FirstArgIsHelp returns true if the first argument is either "-h" or "--help".
//
// This method PANICS if the cliArgs is empty.
func (args cliArgs) FirstArgIsHelp() bool {
	assertTrue(len(args) >= 1, "empty argv")
	return len(args) >= 2 && (args[1] == "-h" || args[1] == "--help")
}

// ContainsHelp returns true if the cliArgs contains "-h" or "--help".
//
// This method PANICS if the cliArgs is empty.
func (args cliArgs) ContainsHelp() bool {
	assertTrue(len(args) >= 1, "empty argv")
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

// IsEmpty returns true if the cliArgs is empty.
//
// This method PANICS if the cliArgs is empty.
func (args cliArgs) IsEmpty() bool {
	assertTrue(len(args) >= 1, "empty argv")
	return len(args) == 1
}

// CommandArgs returns the arguments passed to a command.
//
// This method PANICS if the cliArgs is empty.
func (args cliArgs) CommandArgs() cliArgs {
	assertTrue(len(args) >= 1, "no subcommand")
	return args[1:]
}

// CommandName returns the name of the subcommand to execute.
//
// This method PANICS if the cliArgs is empty.
func (args cliArgs) CommandName() string {
	assertTrue(len(args) >= 1, "empty argv")
	return args[0]
}
