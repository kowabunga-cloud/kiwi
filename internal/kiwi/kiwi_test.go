/*
 * Copyright (c) The Kowabunga Project
 * Apache License, Version 2.0 (see LICENSE or https://www.apache.org/licenses/LICENSE-2.0.txt)
 * SPDX-License-Identifier: Apache-2.0
 */

package kiwi

import (
	"testing"
	"time"

	"github.com/kowabunga-cloud/common/agents"
)

func TestKiwiConstants(t *testing.T) {
	if KiwiVersion != "1.1" {
		t.Errorf("KiwiVersion = %s, want 1.1", KiwiVersion)
	}
	if KiwiAppNmame != "kowabunga-kiwi" {
		t.Errorf("KiwiAppNmame = %s, want kowabunga-kiwi", KiwiAppNmame)
	}
}

func TestNewKiwiAgent(t *testing.T) {
	tests := []struct {
		name    string
		config  *KiwiAgentConfig
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: &KiwiAgentConfig{
				Global: agents.KowabungaAgentGlobalConfig{
					ID:       "test-agent",
					Endpoint: "http://localhost:8080",
					APIKey:   "test-key",
					LogLevel: "ERROR",
				},
				DNS: KiwiAgentDnsConfig{
					Port:      15361,
					Recursors: []string{"8.8.8.8"},
				},
			},
			wantErr: false,
		},
		{
			name: "minimal configuration",
			config: &KiwiAgentConfig{
				Global: agents.KowabungaAgentGlobalConfig{
					ID:       "minimal-agent",
					Endpoint: "http://localhost:8080",
					APIKey:   "test-key",
					LogLevel: "ERROR",
				},
				DNS: KiwiAgentDnsConfig{
					Port: 15362,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := NewKiwiAgent(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKiwiAgent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if agent == nil {
					t.Error("NewKiwiAgent() returned nil agent")
					return
				}

				if agent.dns == nil {
					t.Error("DNS server not initialized")
				}

				if agent.KowabungaAgent == nil {
					t.Error("KowabungaAgent not initialized")
				}

				if agent.PostFlight == nil {
					t.Error("PostFlight not set")
				}

				// Clean up
				agent.Shutdown()
			}
		})
	}
}

func TestKiwiAgent_Shutdown(t *testing.T) {
	cfg := &KiwiAgentConfig{
		Global: agents.KowabungaAgentGlobalConfig{
			ID:       "test-agent",
			Endpoint: "http://localhost:8080",
			APIKey:   "test-key",
			LogLevel: "ERROR",
		},
		DNS: KiwiAgentDnsConfig{
			Port:      15363,
			Recursors: []string{"8.8.8.8"},
		},
	}

	agent, err := NewKiwiAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to create KiwiAgent: %v", err)
	}

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Shutdown should not panic
	agent.Shutdown()

	// Multiple shutdowns should not panic
	agent.Shutdown()
}

func TestKiwiAgent_Integration(t *testing.T) {
	cfg := &KiwiAgentConfig{
		Global: agents.KowabungaAgentGlobalConfig{
			ID:       "integration-test",
			Endpoint: "http://localhost:8080",
			APIKey:   "test-key",
			LogLevel: "ERROR",
		},
		DNS: KiwiAgentDnsConfig{
			Port:      15364,
			Recursors: []string{"8.8.8.8", "8.8.4.4"},
		},
	}

	agent, err := NewKiwiAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to create KiwiAgent: %v", err)
	}
	defer agent.Shutdown()

	// Verify DNS server is running
	if agent.dns == nil {
		t.Error("DNS server is nil")
	}

	// Verify DNS server has correct configuration
	if agent.dns.Port != 15364 {
		t.Errorf("DNS port = %d, want 15364", agent.dns.Port)
	}

	if len(agent.dns.Recursors) != 2 {
		t.Errorf("Recursors count = %d, want 2", len(agent.dns.Recursors))
	}

	// Test adding a DNS record
	err = agent.dns.AddRecord("test.example.com.", "192.168.1.1")
	if err != nil {
		t.Errorf("Failed to add DNS record: %v", err)
	}

	// Verify record was added
	agent.dns.m.Lock()
	value, exists := agent.dns.records["test.example.com."]
	agent.dns.m.Unlock()

	if !exists {
		t.Error("DNS record was not added")
	}

	if value != "192.168.1.1" {
		t.Errorf("DNS record value = %s, want 192.168.1.1", value)
	}
}

func TestKiwiAgent_DNSServerLifecycle(t *testing.T) {
	cfg := &KiwiAgentConfig{
		Global: agents.KowabungaAgentGlobalConfig{
			ID:       "lifecycle-test",
			Endpoint: "http://localhost:8080",
			APIKey:   "test-key",
			LogLevel: "ERROR",
		},
		DNS: KiwiAgentDnsConfig{
			Port:      15365,
			Recursors: []string{"8.8.8.8"},
		},
	}

	// Create agent (starts DNS server)
	agent, err := NewKiwiAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to create KiwiAgent: %v", err)
	}

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Verify DNS server is accessible
	if agent.dns.srv == nil {
		t.Error("DNS server not started")
	}

	// Shutdown (stops DNS server)
	agent.Shutdown()

	// Verify cleanup
	// Note: The actual server may still be shutting down asynchronously
	// We just verify the Shutdown call doesn't panic
}
