// doc.go - Public documentation for the clip package.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package clip implements command-line options parsing.
//
// The API mimics the standard library but uses GNU getopt-long style
// parsing with short flags, that can be combined, and long flags.
//
// Long flags work both in the `--foo=bar` and `--foo bar` way.
//
// Boolean long flags optionally accept `--foo=true` and `--foo=false`
// as well as `--foo` and `--foo true` and `--foo false`.
//
// We do not reorder flags, so the following CLI would work as intended:
//
//	command -xyz subcommand -abc
//
// The `--` separator stops processing early.
package clip
