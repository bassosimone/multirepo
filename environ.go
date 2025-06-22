// environ.go - Environment definition and implementation.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rogpeppe/go-internal/lockedfile"
)

// lockReleaser is a function that releases a lock on the dot directory.
type lockReleaser func()

// environ is the environment in which a [command] runs.
type environ interface {
	// AbsFilepath returns the absolute path of the given path.
	AbsFilepath(path string) (string, error)

	// CreateLockFile creates a lockfile at the given path.
	CreateLockFile(path string) (lockReleaser, error)

	// Executable returns the path to the executable file.
	Executable() (string, error)

	// FileExists checks if a file exists and is a regular file.
	FileExists(path string) (bool, error)

	// Getwd returns the current working directory.
	Getwd() (string, error)

	// MkdirAll creates a directory and all its parents if they do not exist.
	MkdirAll(path string, perm os.FileMode) error

	// ReadFile reads the given file and returns its contents.
	ReadFile(filename string) ([]byte, error)

	// RunCommand runs the given [*exec.Cmd].
	RunCommand(cmd *exec.Cmd) error

	// Stdin returns the standard input.
	Stdin() io.Reader

	// Stdout returns the standard output.
	Stdout() io.Writer

	// Stderr returns the standard error.
	Stderr() io.Writer

	// WriteFile writes the given data to the given file using the given permissions.
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

// stdlibEnviron implements the [environ] interface using the standard library.
type stdlibEnviron struct{}

var _ environ = (*stdlibEnviron)(nil)

// AbsFilepath returns the absolute path of the given path.
func (stdlibEnviron) AbsFilepath(path string) (string, error) {
	return filepath.Abs(path)
}

// CreateLockFile implements the [environ] interface.
func (stdlibEnviron) CreateLockFile(path string) (lockReleaser, error) {
	return lockedfile.MutexAt(path).Lock()
}

// Executable implements the [environ] interface.
func (stdlibEnviron) Executable() (string, error) {
	return os.Executable()
}

// FileExists implements the [environ] interface.
func (stdlibEnviron) FileExists(path string) (bool, error) {
	return fsxFileExists(path)
}

// Getwd implements the [environ] interface.
func (stdlibEnviron) Getwd() (string, error) {
	return os.Getwd()
}

// MkdirAll implements the [environ] interface.
func (stdlibEnviron) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// ReadFile implements the [environ] interface.
func (stdlibEnviron) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// RunCommand implements the [environ] interface.
func (stdlibEnviron) RunCommand(cmd *exec.Cmd) error {
	return cmd.Run()
}

// Stdin implements the [environ] interface.
func (stdlibEnviron) Stdin() io.Reader {
	return os.Stdin
}

// Stdout implements the [environ] interface.
func (stdlibEnviron) Stdout() io.Writer {
	return os.Stdout
}

// Stderr implements the [environ] interface.
func (stdlibEnviron) Stderr() io.Writer {
	return os.Stderr
}

// WriteFile implements the [environ] interface.
func (stdlibEnviron) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}
