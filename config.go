// config.go - JSON configuration file management.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import "encoding/json"

// config contains the configuration.
type config struct {
	// Repos maps repository names to their information.
	Repos map[string]repoInfo `json:"repos"`
}

// repoInfo contains information about a repository.
type repoInfo struct {
	// URL contains the scp-like URL of the repository.
	URL string `json:"url"`
}

// readConfig reads the configuration from a file.
func readConfig(env environ, filename string) (*config, error) {
	// read the file from the disk
	data, err := env.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// attempt to parse to JSON
	var cfg config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// ensure the map is not empty
	if cfg.Repos == nil {
		cfg.Repos = make(map[string]repoInfo)
	}

	return &cfg, nil
}

// WriteFile writes the configuration to a file.
func (cfg *config) WriteFile(env environ, filename string) error {
	data := append(mustMarshalIndentJSON(cfg, "", "  "), '\n')
	return env.WriteFile(filename, data, 0644)
}

// AddRepo is a convenience method to add a repository to the configuration.
func (cfg *config) AddRepo(name, url string) error {
	cfg.Repos[name] = repoInfo{URL: url}
	return nil
}
