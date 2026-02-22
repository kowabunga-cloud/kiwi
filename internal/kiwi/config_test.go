/*
 * Copyright (c) The Kowabunga Project
 * Apache License, Version 2.0 (see LICENSE or https://www.apache.org/licenses/LICENSE-2.0.txt)
 * SPDX-License-Identifier: Apache-2.0
 */

package kiwi

import (
	"os"
	"testing"

	"github.com/kowabunga-cloud/common/agents"
)

func TestKiwiConfigParser(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		wantErr     bool
		validate    func(*testing.T, *KiwiAgentConfig)
	}{
		{
			name: "valid config with DNS settings",
			yamlContent: `global:
  id: test-agent
  endpoint: http://localhost:8080
  apiKey: test-key
  logLevel: INFO
dns:
  port: 5353
  recursors:
    - 8.8.8.8
    - 8.8.4.4
`,
			wantErr: false,
			validate: func(t *testing.T, cfg *KiwiAgentConfig) {
				if cfg.Global.ID != "test-agent" {
					t.Errorf("expected ID 'test-agent', got '%s'", cfg.Global.ID)
				}
				if cfg.DNS.Port != 5353 {
					t.Errorf("expected port 5353, got %d", cfg.DNS.Port)
				}
				if len(cfg.DNS.Recursors) != 2 {
					t.Errorf("expected 2 recursors, got %d", len(cfg.DNS.Recursors))
				}
			},
		},
		{
			name: "minimal valid config",
			yamlContent: `global:
  id: minimal-agent
  endpoint: http://localhost:8080
  apiKey: test-key
`,
			wantErr: false,
			validate: func(t *testing.T, cfg *KiwiAgentConfig) {
				if cfg.Global.ID != "minimal-agent" {
					t.Errorf("expected ID 'minimal-agent', got '%s'", cfg.Global.ID)
				}
				if cfg.DNS.Port != 0 {
					t.Errorf("expected default port 0, got %d", cfg.DNS.Port)
				}
			},
		},
		{
			name:        "invalid yaml",
			yamlContent: `invalid: yaml: content: [`,
			wantErr:     true,
			validate:    nil,
		},
		{
			name: "empty dns config",
			yamlContent: `global:
  id: test-agent
  endpoint: http://localhost:8080
  apiKey: test-key
dns:
`,
			wantErr: false,
			validate: func(t *testing.T, cfg *KiwiAgentConfig) {
				if cfg.DNS.Port != 0 {
					t.Errorf("expected port 0, got %d", cfg.DNS.Port)
				}
				if len(cfg.DNS.Recursors) != 0 {
					t.Errorf("expected 0 recursors, got %d", len(cfg.DNS.Recursors))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpFile, err := os.CreateTemp("", "kiwi-config-*.yaml")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer func() {
				_ = os.Remove(tmpFile.Name())
			}()

			// Write test content
			if _, err := tmpFile.Write([]byte(tt.yamlContent)); err != nil {
				t.Fatalf("failed to write temp file: %v", err)
			}

			// Seek back to beginning for reading
			if _, err := tmpFile.Seek(0, 0); err != nil {
				t.Fatalf("failed to seek temp file: %v", err)
			}

			// Parse config
			cfg, err := KiwiConfigParser(tmpFile)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("KiwiConfigParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Run validation if provided
			if tt.validate != nil && err == nil {
				tt.validate(t, cfg)
			}
		})
	}
}

func TestKiwiAgentConfig_Structure(t *testing.T) {
	cfg := KiwiAgentConfig{
		Global: agents.KowabungaAgentGlobalConfig{
			ID:       "test",
			Endpoint: "http://test",
			APIKey:   "key",
			LogLevel: "DEBUG",
		},
		DNS: KiwiAgentDnsConfig{
			Port:      53,
			Recursors: []string{"8.8.8.8"},
		},
	}

	if cfg.Global.ID != "test" {
		t.Errorf("expected ID 'test', got '%s'", cfg.Global.ID)
	}
	if cfg.DNS.Port != 53 {
		t.Errorf("expected port 53, got %d", cfg.DNS.Port)
	}
}
