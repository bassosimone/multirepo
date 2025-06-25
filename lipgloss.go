// lipgloss.go - Code to interact with lipgloss.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// nilSaybeLipglossStyle is a nil-safe container for a libgloss style.
type nilSafeLipglossStyle struct {
	s lipgloss.Style
}

// NilSafeLipglossStyle creates a [*nilSafeLipglossStyle].
func newNilSafeLipglossStyle() *nilSafeLipglossStyle {
	const (
		white = "#FFFFFF"
		blue  = "#0000FF"
	)
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(blue)).Bold(true)
	return &nilSafeLipglossStyle{style}
}

// Render is a nil-safe colored-text renderer.
func (style *nilSafeLipglossStyle) Render(message string) string {
	if style != nil {
		message = style.s.Render(message)
	}
	return message
}

// Renderf is a convenience function to render and color.
func (style *nilSafeLipglossStyle) Renderf(format string, v ...any) string {
	message := fmt.Sprintf(format, v...)
	message = style.Render(message) // this is nil safe
	return message
}
