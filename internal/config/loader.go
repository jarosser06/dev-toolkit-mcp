package config

import (
	"fmt"
	"os"
)

// LoadManifest loads and validates a task manifest from a file
// It searches for the manifest in the following priority order:
// 1. Custom path (if provided)
// 2. ./mcp-tasks.yaml (project root)
// 3. ./.mcp/tasks.yaml (hidden directory)
//
// Returns:
//   - manifest: The loaded manifest, or an empty manifest if none found
//   - loaded: true if a config file was successfully loaded, false if using default empty config
//   - error: Any error that occurred during parsing or validation (nil if successful or no config found)
func LoadManifest(customPath string) (*Manifest, bool, error) {
	searchPaths := []string{
		customPath,           // CLI flag (if provided)
		"./mcp-tasks.yaml",   // Project root
		"./.mcp/tasks.yaml",  // Hidden directory
	}

	for _, path := range searchPaths {
		if path == "" {
			continue
		}

		// Check if file exists
		if _, err := os.Stat(path); err != nil {
			continue
		}

		// Parse the manifest
		manifest, err := ParseManifest(path)
		if err != nil {
			return nil, false, fmt.Errorf("failed to parse manifest at %s: %w", path, err)
		}

		// Validate the manifest
		if err := Validate(manifest); err != nil {
			return nil, false, fmt.Errorf("invalid manifest at %s: %w", path, err)
		}

		return manifest, true, nil
	}

	// No manifest found - return empty manifest instead of error
	// This allows the server to start and provide the init tool
	emptyManifest := &Manifest{
		Version: "1.0",
		Tasks:   make(map[string]Task),
	}

	return emptyManifest, false, nil
}
