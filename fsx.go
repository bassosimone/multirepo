// fsx.go - File-system extensions.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"errors"
	"os"
)

// errUnexpectedFileType is returned when a file is not a regular file.
var errUnexpectedFileType = errors.New("unexpected file type")

// fsxFileExists checks if a file exists and is a regular file. You SHOULD only
// invoke this function when the `.multirepo` directory has been locked.
func fsxFileExists(path string) (bool, error) {
	// Get information about the file
	sbuf, err := os.Stat(path)

	// Handle the error case distinguishing between the file
	// not existing and other kinds of errors
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return true, err
	}

	// Handle the case where the file is not a regular file
	if !sbuf.Mode().IsRegular() {
		return false, errUnexpectedFileType
	}

	// Handle the successful case
	return true, nil
}
