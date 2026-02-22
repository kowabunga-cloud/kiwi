/*
 * Copyright (c) The Kowabunga Project
 * Apache License, Version 2.0 (see LICENSE or https://www.apache.org/licenses/LICENSE-2.0.txt)
 * SPDX-License-Identifier: Apache-2.0
 */

package kiwi

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/miekg/dns"
)

func TestNewDnsServer(t *testing.T) {
	tests := []struct {
		name   string
		config KiwiAgentDnsConfig
		want   *DnsServer
	}{
		{
			name: "default port and recursors",
			config: KiwiAgentDnsConfig{
				Port:      0,
				Recursors: []string{},
			},
			want: &DnsServer{
				Port: DnsDefaultPort,
				Recursors: []string{
					DnsDefaultRecursorPrimary,
					DnsDefaultRecursorSecondary,
				},
			},
		},
		{
			name: "custom port and recursors",
			config: KiwiAgentDnsConfig{
				Port:      5353,
				Recursors: []string{"8.8.8.8", "8.8.4.4"},
			},
			want: &DnsServer{
				Port:      5353,
				Recursors: []string{"8.8.8.8", "8.8.4.4"},
			},
		},
		{
			name: "custom port only",
			config: KiwiAgentDnsConfig{
				Port:      8053,
				Recursors: []string{},
			},
			want: &DnsServer{
				Port: 8053,
				Recursors: []string{
					DnsDefaultRecursorPrimary,
					DnsDefaultRecursorSecondary,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDnsServer(tt.config)
			if err != nil {
				t.Errorf("NewDnsServer() error = %v", err)
				return
			}
			if got.Port != tt.want.Port {
				t.Errorf("Port = %v, want %v", got.Port, tt.want.Port)
			}
			if len(got.Recursors) != len(tt.want.Recursors) {
				t.Errorf("Recursors length = %v, want %v", len(got.Recursors), len(tt.want.Recursors))
			}
			for i := range got.Recursors {
				if got.Recursors[i] != tt.want.Recursors[i] {
					t.Errorf("Recursors[%d] = %v, want %v", i, got.Recursors[i], tt.want.Recursors[i])
				}
			}
			if got.records == nil {
				t.Error("records map should be initialized")
			}
		})
	}
}

func TestDnsServer_AddRecord(t *testing.T) {
	srv, _ := NewDnsServer(KiwiAgentDnsConfig{})

	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
		setup   func()
	}{
		{
			name:    "add new record",
			key:     "example.com.",
			value:   "192.168.1.1",
			wantErr: false,
			setup:   func() {},
		},
		{
			name:    "add duplicate record",
			key:     "duplicate.com.",
			value:   "192.168.1.1",
			wantErr: true,
			setup: func() {
				srv.records["duplicate.com."] = "192.168.1.1"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := srv.AddRecord(tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if srv.records[tt.key] != tt.value {
					t.Errorf("Record not added correctly: got %v, want %v", srv.records[tt.key], tt.value)
				}
			}
		})
	}
}

func TestDnsServer_UpdateRecord(t *testing.T) {
	srv, _ := NewDnsServer(KiwiAgentDnsConfig{})

	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
		setup   func()
	}{
		{
			name:    "update existing record",
			key:     "existing.com.",
			value:   "192.168.1.2",
			wantErr: false,
			setup: func() {
				srv.records["existing.com."] = "192.168.1.1"
			},
		},
		{
			name:    "update non-existing record",
			key:     "nonexistent.com.",
			value:   "192.168.1.1",
			wantErr: true,
			setup:   func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := srv.UpdateRecord(tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if srv.records[tt.key] != tt.value {
					t.Errorf("Record not updated correctly: got %v, want %v", srv.records[tt.key], tt.value)
				}
			}
		})
	}
}

func TestDnsServer_DeleteRecord(t *testing.T) {
	srv, _ := NewDnsServer(KiwiAgentDnsConfig{})
	srv.records["test.com."] = "192.168.1.1"

	srv.DeleteRecord("test.com.")

	if _, exists := srv.records["test.com."]; exists {
		t.Error("Record should have been deleted")
	}

	// Deleting non-existent record should not panic
	srv.DeleteRecord("nonexistent.com.")
}

func TestDnsServer_UpdateAllRecords(t *testing.T) {
	srv, _ := NewDnsServer(KiwiAgentDnsConfig{})
	srv.records["old.com."] = "192.168.1.1"

	newRecords := map[string]string{
		"new1.com.": "192.168.1.2",
		"new2.com.": "192.168.1.3",
	}

	srv.UpdateAllRecords(newRecords)

	if _, exists := srv.records["old.com."]; exists {
		t.Error("Old record should have been removed")
	}

	if srv.records["new1.com."] != "192.168.1.2" {
		t.Error("New record 1 not added correctly")
	}
	if srv.records["new2.com."] != "192.168.1.3" {
		t.Error("New record 2 not added correctly")
	}
}

func TestDnsServer_ServeDNS_LocalRecord(t *testing.T) {
	srv, _ := NewDnsServer(KiwiAgentDnsConfig{})
	srv.records["example.com."] = "192.168.1.1"

	// Create a DNS query
	m := new(dns.Msg)
	m.SetQuestion("example.com.", dns.TypeA)

	// Create a mock response writer
	writer := &mockResponseWriter{}

	// Handle the request
	srv.ServeDNS(writer, m)

	// Verify response
	if writer.msg == nil {
		t.Fatal("No response written")
	}

	if len(writer.msg.Answer) != 1 {
		t.Fatalf("Expected 1 answer, got %d", len(writer.msg.Answer))
	}

	if aRecord, ok := writer.msg.Answer[0].(*dns.A); ok {
		if aRecord.A.String() != "192.168.1.1" {
			t.Errorf("Expected IP 192.168.1.1, got %s", aRecord.A.String())
		}
	} else {
		t.Error("Answer is not an A record")
	}
}

func TestDnsServer_ServeDNS_MultipleIPs(t *testing.T) {
	srv, _ := NewDnsServer(KiwiAgentDnsConfig{})
	srv.records["multi.com."] = "192.168.1.1,192.168.1.2"

	m := new(dns.Msg)
	m.SetQuestion("multi.com.", dns.TypeA)

	writer := &mockResponseWriter{}
	srv.ServeDNS(writer, m)

	if writer.msg == nil {
		t.Fatal("No response written")
	}

	if len(writer.msg.Answer) != 2 {
		t.Fatalf("Expected 2 answers, got %d", len(writer.msg.Answer))
	}
}

func TestDnsServer_ServeDNS_EmptyQuestion(t *testing.T) {
	srv, _ := NewDnsServer(KiwiAgentDnsConfig{})

	m := new(dns.Msg)
	m.Question = []dns.Question{} // Empty questions

	writer := &mockResponseWriter{}
	srv.ServeDNS(writer, m)

	if writer.msg == nil {
		t.Fatal("No response written")
	}

	if len(writer.msg.Answer) != 0 {
		t.Error("Expected no answers for empty question")
	}
}

func TestDnsServer_StartStop(t *testing.T) {
	// Use a non-privileged port for testing
	srv, _ := NewDnsServer(KiwiAgentDnsConfig{Port: 15353})

	err := srv.Start()
	if err != nil {
		t.Fatalf("Failed to start DNS server: %v", err)
	}

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	err = srv.Stop()
	if err != nil {
		t.Errorf("Failed to stop DNS server: %v", err)
	}
}

// mockResponseWriter implements dns.ResponseWriter for testing
type mockResponseWriter struct {
	msg *dns.Msg
}

func (m *mockResponseWriter) LocalAddr() net.Addr {
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:53")
	return addr
}
func (m *mockResponseWriter) RemoteAddr() net.Addr {
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:12345")
	return addr
}
func (m *mockResponseWriter) WriteMsg(msg *dns.Msg) error {
	m.msg = msg
	return nil
}
func (m *mockResponseWriter) Write([]byte) (int, error) { return 0, nil }
func (m *mockResponseWriter) Close() error              { return nil }
func (m *mockResponseWriter) TsigStatus() error         { return nil }
func (m *mockResponseWriter) TsigTimersOnly(bool)       {}
func (m *mockResponseWriter) Hijack()                   {}

func TestDnsConstants(t *testing.T) {
	if DnsDefaultPort != 53 {
		t.Errorf("DnsDefaultPort should be 53, got %d", DnsDefaultPort)
	}
	if DnsDefaultRecursorPrimary != "9.9.9.9" {
		t.Errorf("DnsDefaultRecursorPrimary incorrect: %s", DnsDefaultRecursorPrimary)
	}
	if DnsDefaultRecursorSecondary != "149.112.112.112" {
		t.Errorf("DnsDefaultRecursorSecondary incorrect: %s", DnsDefaultRecursorSecondary)
	}
}

func TestDnsServer_ConcurrentAccess(t *testing.T) {
	srv, _ := NewDnsServer(KiwiAgentDnsConfig{})

	// Test concurrent writes
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			key := fmt.Sprintf("test%d.com.", idx)
			_ = srv.AddRecord(key, "192.168.1.1")
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Test concurrent reads
	for i := 0; i < 10; i++ {
		go func(idx int) {
			m := new(dns.Msg)
			m.SetQuestion(fmt.Sprintf("test%d.com.", idx), dns.TypeA)
			writer := &mockResponseWriter{}
			srv.ServeDNS(writer, m)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
