package security

import (
	"net"
	"testing"
)

func TestIsPublicIP(t *testing.T) {
	tests := []struct {
		address string
		public  bool
	}{
		{"8.8.8.8", true},
		{"2606:4700:4700::1111", true},
		{"127.0.0.1", false},
		{"10.0.0.1", false},
		{"172.16.0.1", false},
		{"192.168.1.1", false},
		{"169.254.169.254", false},
		{"100.64.0.1", false},
		{"198.18.0.1", false},
		{"::1", false},
		{"fc00::1", false},
		{"fe80::1", false},
	}
	for _, test := range tests {
		t.Run(test.address, func(t *testing.T) {
			if actual := IsPublicIP(net.ParseIP(test.address)); actual != test.public {
				t.Fatalf("IsPublicIP(%s) = %v, want %v", test.address, actual, test.public)
			}
		})
	}
}

func TestLookupPublicIPsRejectsPrivateLiteral(t *testing.T) {
	if _, err := LookupPublicIPs(t.Context(), "127.0.0.1"); err == nil {
		t.Fatal("expected loopback address to be rejected")
	}
}
