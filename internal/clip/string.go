// string.go - stringValue implementation.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

// stringValue implements [Value] for string.
type stringValue struct {
	value *string
}

var _ Value = stringValue{}

// newStringValue creates a new StringValue with the given default value.
func newStringValue(value *string) stringValue {
	return stringValue{value: value}
}

// OptionType implements [Value].
func (v stringValue) OptionType() OptionType {
	return WithArgOption
}

// String implements [Value].
func (v stringValue) String() string {
	return *v.value
}

// Set implements [Value].
func (v stringValue) Set(values []string) error {
	for _, value := range values {
		*v.value = value
	}
	return nil
}
