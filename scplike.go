// scplike.go - Parse SCP-like URLs.
// Adapted-From: https://github.com/go-git/go-git/blob/v4.7.0/plumbing/transport/common.go#L232
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// scpLikeEndpoint is an SCP-like git endpoint.
type scpLikeEndpoint struct {
	// Protocol is the protocol to use.
	Protocol string

	// User is the user.
	User string

	// Password is the password.
	Password string

	// Host is the host.
	Host string

	// Port is the port to connect, if 0 the default port for
	// the given protocol wil be used.
	Port int

	// Path is the repository path.
	Path string
}

// scpLikeDefaultPorts is a map of default ports for each protocol.
var scpLikeDefaultPorts = map[string]int{
	"http":  80,
	"https": 443,
	"git":   9418,
	"ssh":   22,
}

// String returns the endpoint as a string.
func (epnt *scpLikeEndpoint) String() string {
	var buf bytes.Buffer
	if epnt.Protocol != "" {
		buf.WriteString(epnt.Protocol)
		buf.WriteByte(':')
	}

	if epnt.Protocol != "" || epnt.Host != "" || epnt.User != "" || epnt.Password != "" {
		buf.WriteString("//")

		if epnt.User != "" || epnt.Password != "" {
			buf.WriteString(url.PathEscape(epnt.User))
			if epnt.Password != "" {
				buf.WriteByte(':')
				buf.WriteString(url.PathEscape(epnt.Password))
			}

			buf.WriteByte('@')
		}

		if epnt.Host != "" {
			buf.WriteString(epnt.Host)

			if epnt.Port != 0 {
				port, ok := scpLikeDefaultPorts[strings.ToLower(epnt.Protocol)]
				if !ok || ok && port != epnt.Port {
					fmt.Fprintf(&buf, ":%d", epnt.Port)
				}
			}
		}
	}

	if epnt.Path != "" && epnt.Path[0] != '/' && epnt.Host != "" {
		buf.WriteByte('/')
	}

	buf.WriteString(epnt.Path)
	return buf.String()
}

// Name returns the repository name derived from the path.
func (epnt *scpLikeEndpoint) Name() string {
	// "If s does not contain sep and sep is not empty, Split returns
	// a slice of length 1 whose only element is s."
	values := strings.Split(epnt.Path, "/")
	return values[len(values)-1]
}

var (
	// isSchemeRegExp is a regular expression that matches a scheme.
	isSchemeRegExp = regexp.MustCompile(`^[^:]+://`)

	// scpLikeUrlRegExp is a regular expression that matches an SCP-like URL.
	scpLikeUrlRegExp = regexp.MustCompile(`^(?:(?P<user>[^@]+)@)?(?P<host>[^:\s]+):(?:(?P<port>[0-9]{1,5})/)?(?P<path>[^\\].*)$`)
)

// scpLikeParse parses an SCP-like git endpoint.
func scpLikeParse(endpoint string) (*scpLikeEndpoint, bool) {
	if isSchemeRegExp.MatchString(endpoint) || !scpLikeUrlRegExp.MatchString(endpoint) {
		return nil, false
	}

	m := scpLikeUrlRegExp.FindStringSubmatch(endpoint)

	port, err := strconv.Atoi(m[3])
	if err != nil {
		port = 22
	}

	epnt := &scpLikeEndpoint{
		Protocol: "ssh",
		User:     m[1],
		Host:     m[2],
		Port:     port,
		Path:     m[4],
	}

	return epnt, true
}
