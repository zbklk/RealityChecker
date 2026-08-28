package security

import (
	"context"
	"fmt"
	"net"
	"time"
)

var reservedNetworks = mustParseCIDRs(
	"100.64.0.0/10",
	"192.0.0.0/24",
	"192.0.2.0/24",
	"198.18.0.0/15",
	"198.51.100.0/24",
	"203.0.113.0/24",
	"2001:db8::/32",
)

func mustParseCIDRs(values ...string) []*net.IPNet {
	result := make([]*net.IPNet, 0, len(values))
	for _, value := range values {
		_, network, err := net.ParseCIDR(value)
		if err != nil {
			panic(err)
		}
		result = append(result, network)
	}
	return result
}

func IsPublicIP(ip net.IP) bool {
	if ip == nil || !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
		return false
	}
	for _, network := range reservedNetworks {
		if network.Contains(ip) {
			return false
		}
	}
	return true
}

func LookupPublicIPs(ctx context.Context, host string) ([]net.IP, error) {
	if parsed := net.ParseIP(host); parsed != nil {
		if !IsPublicIP(parsed) {
			return nil, fmt.Errorf("拒绝访问非公网地址 %s", host)
		}
		return []net.IP{parsed}, nil
	}
	addresses, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
	if err != nil {
		return nil, err
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("域名 %s 没有可用地址", host)
	}
	for _, address := range addresses {
		if !IsPublicIP(address) {
			return nil, fmt.Errorf("域名 %s 解析到非公网地址 %s，已拒绝", host, address)
		}
	}
	return addresses, nil
}

// DialContextPublic resolves once, validates every answer, and dials the
// validated IP directly. This closes the DNS-rebinding gap between checking
// and connecting.
func DialContextPublic(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	if port != "80" && port != "443" {
		return nil, fmt.Errorf("拒绝访问端口 %s；仅允许 80/443", port)
	}
	addresses, err := LookupPublicIPs(ctx, host)
	if err != nil {
		return nil, err
	}
	dialer := &net.Dialer{}
	var lastErr error
	for _, ip := range addresses {
		conn, dialErr := dialer.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
		if dialErr == nil {
			return conn, nil
		}
		lastErr = dialErr
	}
	return nil, lastErr
}

func DialTimeoutPublic(network, address string, timeout time.Duration) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return DialContextPublic(ctx, network, address)
}
