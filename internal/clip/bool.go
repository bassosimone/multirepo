// bool.go - boolValue implementation.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import "strconv"

// boolValue implements [Value] for bool.
type boolValue struct {
	value *bool
}

var _ Value = boolValue{}

// newBoolValue creates a new BoolValue with the given default value.
func newBoolValue(value *bool) boolValue {
	return boolValue{value: value}
}

// OptionType implements [Value].
func (v boolValue) OptionType() OptionType {
	return BoolOption
}

// String implements [Value].
func (v boolValue) String() string {
	return strconv.FormatBool(*v.value)
}

// Set implements [Value].
func (v boolValue) Set(values []string) error {
	for _, value := range values {
		*v.value = value == "true"
	}
	return nil
}
