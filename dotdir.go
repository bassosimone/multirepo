// dotdir.go - Code to manage the `.multirepo` dot directory.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import "path/filepath"

// dotDir is the multirepo dot directory.
type dotDir string

// defaultDotDir is the default multirepo dot directory.
func defaultDotDir() dotDir {
	return ".multirepo"
}

// String returns the string representation of the dot directory.
func (dd dotDir) String() string {
	return string(dd)
}

// configFilePath returns the path to the configuration file.
func (dd dotDir) configFilePath() string {
	return filepath.Join(dd.String(), "config.json")
}

// lock locks the dot directory until it is released.
func (dd dotDir) lock(env environ) (lockReleaser, error) {
	lpath := filepath.Join(dd.String(), "lock")
	return env.CreateLockFile(lpath)
}
