// assert.go - Runtime assertions.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import "errors"

// assertTrue panics with the given message if the condition is false.
func assertTrue(condition bool, message string) {
	if !condition {
		panic(errors.New(message))
	}
}
