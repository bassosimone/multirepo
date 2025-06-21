// value.go - Value implementation.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

// Value represents a flag value.
type Value interface {
	// OptionType returns the type of the option.
	OptionType() OptionType

	// String returns the string representation of the value.
	String() string

	// Set sets the value of the flag.
	Set(values []string) error
}
