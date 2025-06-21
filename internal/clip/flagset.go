// flagset.go - FlagSet implementation.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

// ErrorHandling defines the behavior of the flag set when an error occurs.
type ErrorHandling int

const (
	// ContinueOnError causes the flag set to continue parsing after an error occurs.
	ContinueOnError = ErrorHandling(iota)
)

// FlagSet represents a set of command-line flags.
//
// The zero value is invalid. Construct using [NewFlagSet].
type FlagSet struct {
	// args contains the positional arguments.
	args []string

	// progname contains the name of the program.
	progname string

	// values maps flags to their values.
	values map[string]Value
}

// NewFlagSet creates a new flag set with the given program name and error handling.
func NewFlagSet(progname string, handling ErrorHandling) *FlagSet {
	return &FlagSet{
		args:     []string{},
		progname: progname,
		values:   make(map[string]Value),
	}
}

// Bool creates a new boolean flag with the given name, default value, and usage.
//
// This method MUST be invoked before parsing the arguments.
func (fs *FlagSet) Bool(longName string, shortName byte, value bool, usage string) *bool {
	if longName != "" {
		fs.values[longName] = newBoolValue(&value)
	}
	if shortName != 0 {
		fs.values[string(shortName)] = newBoolValue(&value)
	}
	return &value
}

// String creates a new string flag with the given name, default value, and usage.
//
// This method MUST be invoked before parsing the arguments.
func (fs *FlagSet) String(longName string, shortName byte, value string, usage string) *string {
	if longName != "" {
		fs.values[longName] = newStringValue(&value)
	}
	if shortName != 0 {
		fs.values[string(shortName)] = newStringValue(&value)
	}
	return &value
}

// Args returns the positional arguments.
//
// This method MUST be invoked after parsing the arguments.
func (fs *FlagSet) Args() []string {
	return fs.args
}

// Parse parses the command line arguments.
//
// This method MUST be invoked after setting up the flags.
func (fs *FlagSet) Parse(arguments []string) error {
	// Generate the argv from the program name and the arguments.
	argv := append([]string{fs.progname}, arguments...)

	// Generate the option specification
	spec, err := fs.newOptSpec()
	if err != nil {
		return err
	}

	// Parse the command line
	res, err := getoptLong(spec, argv)
	if err != nil {
		return err
	}

	// Save the positional arguments
	fs.args = res.Positional

	// Copy the resulting values back
	return fs.saveResults(res)
}

// newOptSpec creates a new OptSpec from the FlagSet.
func (fs *FlagSet) newOptSpec() (optSpec, error) {
	spec := make(optSpec)
	for name, value := range fs.values {
		spec[name] = value.OptionType()
	}
	return spec, nil
}

// saveResults saves the parsed results back to the FlagSet.
func (fs *FlagSet) saveResults(res *getoptResult) error {
	for name, value := range fs.values {
		// Attempt to access the parsed option
		results := res.Options[name]
		if len(results) <= 0 {
			continue
		}

		// Attempt to assign the parsed value
		if err := value.Set(results); err != nil {
			return err
		}
	}
	return nil
}
