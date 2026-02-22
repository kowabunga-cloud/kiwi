/*
 * Copyright (c) The Kowabunga Project
 * Apache License, Version 2.0 (see LICENSE or https://www.apache.org/licenses/LICENSE-2.0.txt)
 * SPDX-License-Identifier: Apache-2.0
 */

package kiwi

import (
	"testing"

	"github.com/kowabunga-cloud/common/agents"
	"github.com/kowabunga-cloud/common/proto"
)

func TestCapabilities(t *testing.T) {
	// Create a minimal agent config for testing
	cfg := &KiwiAgentConfig{
		Global: agents.KowabungaAgentGlobalConfig{
			ID:       "test-agent",
			Endpoint: "http://localhost:8080",
			APIKey:   "test-key",
			LogLevel: "ERROR", // Use ERROR to minimize log output during tests
		},
		DNS: KiwiAgentDnsConfig{
			Port:      15353, // Use non-privileged port for testing
			Recursors: []string{"8.8.8.8"},
		},
	}

	agent, err := NewKiwiAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to create KiwiAgent: %v", err)
	}
	defer agent.Shutdown()

	kiwi := newKiwi(agent)

	args := &agents.CapabilitiesArgs{}
	reply := &agents.CapabilitiesReply{}

	err = kiwi.Capabilities(args, reply)
	if err != nil {
		t.Errorf("Capabilities() error = %v", err)
	}

	if reply.Version != KiwiVersion {
		t.Errorf("Expected version %s, got %s", KiwiVersion, reply.Version)
	}

	if reply.Methods == nil {
		t.Error("Expected methods to be populated")
	}
}

func TestReload(t *testing.T) {
	cfg := &KiwiAgentConfig{
		Global: agents.KowabungaAgentGlobalConfig{
			ID:       "test-agent",
			Endpoint: "http://localhost:8080",
			APIKey:   "test-key",
			LogLevel: "ERROR",
		},
		DNS: KiwiAgentDnsConfig{
			Port:      15354,
			Recursors: []string{"8.8.8.8"},
		},
	}

	agent, err := NewKiwiAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to create KiwiAgent: %v", err)
	}
	defer agent.Shutdown()

	kiwi := newKiwi(agent)

	tests := []struct {
		name string
		args *proto.KiwiReloadArgs
	}{
		{
			name: "reload with single domain and A record",
			args: &proto.KiwiReloadArgs{
				Domains: []proto.KiwiReloadArgsDomain{
					{
						Name: "example.com",
						Records: []proto.KiwiReloadArgsRecord{
							{
								Name:      "www",
								Type:      "A",
								Addresses: []string{"192.168.1.1"},
							},
						},
					},
				},
			},
		},
		{
			name: "reload with multiple domains",
			args: &proto.KiwiReloadArgs{
				Domains: []proto.KiwiReloadArgsDomain{
					{
						Name: "example.com",
						Records: []proto.KiwiReloadArgsRecord{
							{
								Name:      "www",
								Type:      "A",
								Addresses: []string{"192.168.1.1"},
							},
							{
								Name:      "mail",
								Type:      "A",
								Addresses: []string{"192.168.1.2"},
							},
						},
					},
					{
						Name: "test.com",
						Records: []proto.KiwiReloadArgsRecord{
							{
								Name:      "api",
								Type:      "A",
								Addresses: []string{"192.168.2.1"},
							},
						},
					},
				},
			},
		},
		{
			name: "reload with non-A records (should be ignored)",
			args: &proto.KiwiReloadArgs{
				Domains: []proto.KiwiReloadArgsDomain{
					{
						Name: "example.com",
						Records: []proto.KiwiReloadArgsRecord{
							{
								Name:      "www",
								Type:      "CNAME",
								Addresses: []string{"example.com"},
							},
						},
					},
				},
			},
		},
		{
			name: "reload with empty domains",
			args: &proto.KiwiReloadArgs{
				Domains: []proto.KiwiReloadArgsDomain{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reply := &proto.KiwiReloadReply{}
			err := kiwi.Reload(tt.args, reply)
			if err != nil {
				t.Errorf("Reload() error = %v", err)
			}
		})
	}
}

func TestCreateDnsZone(t *testing.T) {
	cfg := &KiwiAgentConfig{
		Global: agents.KowabungaAgentGlobalConfig{
			ID:       "test-agent",
			Endpoint: "http://localhost:8080",
			APIKey:   "test-key",
			LogLevel: "ERROR",
		},
		DNS: KiwiAgentDnsConfig{
			Port:      15355,
			Recursors: []string{"8.8.8.8"},
		},
	}

	agent, err := NewKiwiAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to create KiwiAgent: %v", err)
	}
	defer agent.Shutdown()

	kiwi := newKiwi(agent)

	args := &proto.KiwiCreateDnsZoneArgs{}
	reply := &proto.KiwiCreateDnsZoneReply{}

	err = kiwi.CreateDnsZone(args, reply)
	if err != nil {
		t.Errorf("CreateDnsZone() error = %v", err)
	}
}

func TestDeleteDnsZone(t *testing.T) {
	cfg := &KiwiAgentConfig{
		Global: agents.KowabungaAgentGlobalConfig{
			ID:       "test-agent",
			Endpoint: "http://localhost:8080",
			APIKey:   "test-key",
			LogLevel: "ERROR",
		},
		DNS: KiwiAgentDnsConfig{
			Port:      15356,
			Recursors: []string{"8.8.8.8"},
		},
	}

	agent, err := NewKiwiAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to create KiwiAgent: %v", err)
	}
	defer agent.Shutdown()

	kiwi := newKiwi(agent)

	args := &proto.KiwiDeleteDnsZoneArgs{}
	reply := &proto.KiwiDeleteDnsZoneReply{}

	err = kiwi.DeleteDnsZone(args, reply)
	if err != nil {
		t.Errorf("DeleteDnsZone() error = %v", err)
	}
}

func TestCreateDnsRecord(t *testing.T) {
	cfg := &KiwiAgentConfig{
		Global: agents.KowabungaAgentGlobalConfig{
			ID:       "test-agent",
			Endpoint: "http://localhost:8080",
			APIKey:   "test-key",
			LogLevel: "ERROR",
		},
		DNS: KiwiAgentDnsConfig{
			Port:      15357,
			Recursors: []string{"8.8.8.8"},
		},
	}

	agent, err := NewKiwiAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to create KiwiAgent: %v", err)
	}
	defer agent.Shutdown()

	kiwi := newKiwi(agent)

	tests := []struct {
		name    string
		args    *proto.KiwiCreateDnsRecordArgs
		wantErr bool
	}{
		{
			name: "create new record",
			args: &proto.KiwiCreateDnsRecordArgs{
				Domain:    "example.com",
				Entry:     "www",
				Addresses: []string{"192.168.1.1"},
			},
			wantErr: false,
		},
		{
			name: "create record with multiple IPs",
			args: &proto.KiwiCreateDnsRecordArgs{
				Domain:    "example.com",
				Entry:     "api",
				Addresses: []string{"192.168.1.1", "192.168.1.2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reply := &proto.KiwiCreateDnsRecordReply{}
			err := kiwi.CreateDnsRecord(tt.args, reply)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDnsRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateDnsRecord(t *testing.T) {
	cfg := &KiwiAgentConfig{
		Global: agents.KowabungaAgentGlobalConfig{
			ID:       "test-agent",
			Endpoint: "http://localhost:8080",
			APIKey:   "test-key",
			LogLevel: "ERROR",
		},
		DNS: KiwiAgentDnsConfig{
			Port:      15358,
			Recursors: []string{"8.8.8.8"},
		},
	}

	agent, err := NewKiwiAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to create KiwiAgent: %v", err)
	}
	defer agent.Shutdown()

	kiwi := newKiwi(agent)

	// First create a record
	createArgs := &proto.KiwiCreateDnsRecordArgs{
		Domain:    "example.com",
		Entry:     "test",
		Addresses: []string{"192.168.1.1"},
	}
	createReply := &proto.KiwiCreateDnsRecordReply{}
	_ = kiwi.CreateDnsRecord(createArgs, createReply)

	tests := []struct {
		name    string
		args    *proto.KiwiUpdateDnsRecordArgs
		wantErr bool
	}{
		{
			name: "update existing record",
			args: &proto.KiwiUpdateDnsRecordArgs{
				Domain:    "example.com",
				Entry:     "test",
				Addresses: []string{"192.168.1.2"},
			},
			wantErr: false,
		},
		{
			name: "update non-existent record",
			args: &proto.KiwiUpdateDnsRecordArgs{
				Domain:    "example.com",
				Entry:     "nonexistent",
				Addresses: []string{"192.168.1.1"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reply := &proto.KiwiUpdateDnsRecordReply{}
			err := kiwi.UpdateDnsRecord(tt.args, reply)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateDnsRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteDnsRecord(t *testing.T) {
	cfg := &KiwiAgentConfig{
		Global: agents.KowabungaAgentGlobalConfig{
			ID:       "test-agent",
			Endpoint: "http://localhost:8080",
			APIKey:   "test-key",
			LogLevel: "ERROR",
		},
		DNS: KiwiAgentDnsConfig{
			Port:      15359,
			Recursors: []string{"8.8.8.8"},
		},
	}

	agent, err := NewKiwiAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to create KiwiAgent: %v", err)
	}
	defer agent.Shutdown()

	kiwi := newKiwi(agent)

	// First create a record
	createArgs := &proto.KiwiCreateDnsRecordArgs{
		Domain:    "example.com",
		Entry:     "delete-test",
		Addresses: []string{"192.168.1.1"},
	}
	createReply := &proto.KiwiCreateDnsRecordReply{}
	_ = kiwi.CreateDnsRecord(createArgs, createReply)

	// Now delete it
	args := &proto.KiwiDeleteDnsRecordArgs{
		Domain: "example.com",
		Entry:  "delete-test",
	}
	reply := &proto.KiwiDeleteDnsRecordReply{}

	err = kiwi.DeleteDnsRecord(args, reply)
	if err != nil {
		t.Errorf("DeleteDnsRecord() error = %v", err)
	}

	// Deleting non-existent record should not error
	args2 := &proto.KiwiDeleteDnsRecordArgs{
		Domain: "example.com",
		Entry:  "nonexistent",
	}
	reply2 := &proto.KiwiDeleteDnsRecordReply{}

	err = kiwi.DeleteDnsRecord(args2, reply2)
	if err != nil {
		t.Errorf("DeleteDnsRecord() for non-existent record should not error, got: %v", err)
	}
}

func TestNewKiwi(t *testing.T) {
	cfg := &KiwiAgentConfig{
		Global: agents.KowabungaAgentGlobalConfig{
			ID:       "test-agent",
			Endpoint: "http://localhost:8080",
			APIKey:   "test-key",
			LogLevel: "ERROR",
		},
		DNS: KiwiAgentDnsConfig{
			Port:      15360,
			Recursors: []string{"8.8.8.8"},
		},
	}

	agent, err := NewKiwiAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to create KiwiAgent: %v", err)
	}
	defer agent.Shutdown()

	kiwi := newKiwi(agent)
	if kiwi == nil {
		t.Fatal("newKiwi() returned nil")
	}

	if kiwi.agent != agent {
		t.Error("kiwi.agent does not match provided agent")
	}
}
