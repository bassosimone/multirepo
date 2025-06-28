// main.go - Main file.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/nflag"
)

func main() {
	// Initialize the environment
	env := newStdlibEnviron()

	// Create root command and the related subcommands
	root := &clip.RootCommand[environ]{
		Command: &clip.DispatcherCommand[environ]{
			BriefDescriptionText: "Manage multiple git repositories as a monorepo.",
			Commands: map[string]clip.Command[environ]{
				"clone":   cmdClone,
				"foreach": cmdForeach,
				"init":    cmdInit,
				"repo": &clip.DispatcherCommand[environ]{
					BriefDescriptionText: "Add/remove repositories from the multirepo index.",
					Commands: map[string]clip.Command[environ]{
						"add": cmdRepoAdd,
						"ls":  cmdRepoLs,
						"rm":  cmdRepoRm,
					},
					ErrorHandling:             nflag.ExitOnError,
					Version:                   Version,
					OptionPrefixes:            []string{"--", "-"},
					OptionsArgumentsSeparator: "--",
				},
			},
			ErrorHandling:             nflag.ExitOnError,
			Version:                   Version,
			OptionPrefixes:            []string{"--", "-"},
			OptionsArgumentsSeparator: "--",
		},
		AutoCancel: true,
	}

	// Run the root command
	root.Main(env)
}
