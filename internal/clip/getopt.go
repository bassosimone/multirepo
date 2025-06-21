// getopt.go - Internal getopt implementation.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"errors"
	"fmt"
	"strings"
)

// getoptResult contains the result of parsing command-line options.
type getoptResult struct {
	// Command is the program name, provided as argv[0].
	Command string

	// Options maps of option names to their values.
	Options map[string][]string

	// Positional contains the post-options positional arguments.
	Positional []string
}

// OptionType represents the type of an option.
type OptionType int

const (
	// noSuchOption represents an option that does not exist.
	noSuchOption = OptionType(iota)

	// BoolOption represents an option that does not require an argument.
	BoolOption = OptionType(iota)

	// WithArgOption represents an option that requires an argument.
	WithArgOption
)

// optSpec maps an option name to [withArgOption] if the option requires
// an argument, or [boolOption] if it does not.
type optSpec map[string]OptionType

// errMissingProgramName is returned when the program name is missing.
var errMissingProgramName = errors.New("missing program name")

// getoptLong parses the command line arguments according to the given specification.
func getoptLong(spec optSpec, argv []string) (*getoptResult, error) {
	// Create the initial empty result
	result := &getoptResult{
		Command:    "",
		Options:    map[string][]string{},
		Positional: []string{},
	}

	// Handle the case of empty input
	if len(argv) <= 0 {
		return nil, errMissingProgramName
	}

	// Remember the program name
	result.Command = argv[0]
	argv = argv[1:]

	// Process the options
	argv, err := parseOptions(spec, argv, result)
	if err != nil {
		return nil, err
	}

	// Finish copying positional arguments
	result.Positional = append(result.Positional, argv...)
	return result, nil
}

// parseOptions parses both long and short options and updates the result and the argv.
func parseOptions(spec optSpec, argv []string, result *getoptResult) ([]string, error) {
	for len(argv) > 0 {
		// Get the current entry
		cur := argv[0]

		// Skip the `--` separator and stop processing
		if cur == "--" {
			argv = argv[1:]
			break
		}

		// Stop at first non-option argument
		if !strings.HasPrefix(cur, "-") {
			break
		}

		// We can now advance argv
		argv = argv[1:]

		// Skip the leading `-`
		cur = cur[1:]

		// Select the proper parser to use
		var pf optionStyleParser
		if strings.HasPrefix(cur, "-") {
			cur = cur[1:]
			pf = parseStyleLong
		} else {
			pf = parseStyleShort
		}

		// Parse the option or options
		var err error
		argv, err = pf(spec, argv, result, cur)

		// Handle errors
		if err != nil {
			return nil, err
		}
	}

	// Return the updated arguments vector
	return argv, nil
}

// ErrUnknownOption is returned when an option is not found in the specification.
var ErrUnknownOption = errors.New("unknown option")

// ErrOptionRequiresValue is returned when an option requires a value but none is provided.
var ErrOptionRequiresValue = errors.New("option requires a value")

// ErrInvalidOptionValue is returned when an option value is invalid.
var ErrInvalidOptionValue = errors.New("invalid option value")

// optionStyleParser is the type of function that parses a given option style.
type optionStyleParser func(spec optSpec, argv []string, result *getoptResult, cur string) ([]string, error)

// parseStyleLong parse the long option at cur and updates the result and the argv.
func parseStyleLong(spec optSpec, argv []string, result *getoptResult, cur string) ([]string, error) {
	// The option may contain a value, account for this
	var optname, optvalue string
	index := strings.Index(cur, "=")
	if index > 0 {
		optname = cur[:index]
		optvalue = cur[index+1:]
	} else {
		optname = cur
	}

	// Determine what to do based on the option kind
	optkind := spec[optname]
	switch optkind {

	// Handle the case of boolean option
	case BoolOption:
		switch optvalue {
		case "true", "":
			result.Options[optname] = append(result.Options[optname], "true")
			return argv, nil
		case "false":
			result.Options[optname] = append(result.Options[optname], "false")
			return argv, nil
		default:
			return nil, fmt.Errorf("%w for option %s: %s", ErrInvalidOptionValue, optname, optvalue)
		}

	// Handle the case of an option with value
	case WithArgOption:
		if optvalue != "" {
			result.Options[optname] = append(result.Options[optname], optvalue)
			return argv, nil
		}

		// Otherwise try to use the next entry in the argv
		if len(argv) < 1 {
			return nil, fmt.Errorf("%w: %s", ErrOptionRequiresValue, optname)
		}
		result.Options[optname] = append(result.Options[optname], argv[0])
		argv = argv[1:]
		return argv, nil

	// Otherwise, it does not exist
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownOption, optname)
	}
}

// parseStyleShort parses the short options at cur and updates the result and the argv.
func parseStyleShort(spec optSpec, argv []string, result *getoptResult, cur string) ([]string, error) {
	// Process each character in the option string
	for len(cur) > 0 {
		// Get the option character and advance
		optname := string(cur[0])
		cur = cur[1:]

		// Determine what to do based on the option kind
		optkind := spec[optname]
		switch optkind {

		// If the option does not need an argument, advance
		case BoolOption:
			result.Options[optname] = append(result.Options[optname], "true")
			continue

		// If the option needs an argument, fetch it
		case WithArgOption:
			if len(cur) > 0 {
				result.Options[optname] = append(result.Options[optname], cur)
				return argv, nil
			}

			// Otherwise try to use the next entry in the argv
			if len(argv) < 1 {
				return nil, fmt.Errorf("%w: %s", ErrOptionRequiresValue, optname)
			}
			result.Options[optname] = append(result.Options[optname], argv[0])
			argv = argv[1:]
			return argv, nil

		// Otherwise, it does not exist
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnknownOption, optname)
		}
	}

	// Return the updated arguments vector
	return argv, nil
}
