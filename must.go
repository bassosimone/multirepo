// must.go - Functions that panic on error
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"encoding/json"
	"fmt"
	"io"
)

// mustFprintf is like [fmt.Fprintf] but panics in case of error.
func mustFprintf(w io.Writer, format string, args ...any) {
	if _, err := fmt.Fprintf(w, format, args...); err != nil {
		panic(err)
	}
}

// mustMarshalIndentJSON is like [json.MarshalIndent] but panics in case of error.
func mustMarshalIndentJSON(v any, prefix string, indent string) []byte {
	data, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		panic(err)
	}
	return data
}
