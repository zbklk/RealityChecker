package detectors

import (
	"context"
	"fmt"
	"net"
	"time"

	"RealityChecker/internal/security"
	"RealityChecker/internal/types"
)

// IPResolverStage IP解析阶段
type IPResolverStage struct{}

// NewIPResolverStage 创建IP解析阶段
func NewIPResolverStage() *IPResolverStage {
	return &IPResolverStage{}
}

// Execute 执行IP解析
func (irs *IPResolverStage) Execute(ctx *types.PipelineContext) error {

	// 解析IP地址
	ip, err := irs.resolveIP(ctx.Domain)
	if err != nil {
		return fmt.Errorf("IP解析失败: %v", err)
	}

	// 快速连通性测试
	if !irs.quickConnectivityTest(ip) {
		return fmt.Errorf("网络不可达")
	}

	// 设置IP地址到Location结果中
	if ctx.Result.Location == nil {
		ctx.Result.Location = &types.LocationResult{}
	}
	ctx.Result.Location.IPAddress = ip

	return nil
}

// quickConnectivityTest 快速连通性测试
func (irs *IPResolverStage) quickConnectivityTest(ip string) bool {
	// 测试HTTPS端口443的连通性
	conn, err := security.DialTimeoutPublic("tcp", net.JoinHostPort(ip, "443"), 2*time.Second)
	if err != nil {
		// 如果HTTPS不可达，尝试HTTP端口80
		conn, err = security.DialTimeoutPublic("tcp", net.JoinHostPort(ip, "80"), 2*time.Second)
		if err != nil {
			return false
		}
	}
	conn.Close()
	return true
}

// resolveIP 解析IP地址
func (irs *IPResolverStage) resolveIP(domain string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ips, err := security.LookupPublicIPs(ctx, domain)
	if err != nil {
		return "", err
	}

	if len(ips) == 0 {
		return "", fmt.Errorf("未找到IP地址")
	}

	// 优先选择IPv4地址
	for _, ipAddr := range ips {
		if ipAddr.To4() != nil {
			return ipAddr.String(), nil
		}
	}

	// 如果没有IPv4，使用IPv6
	return ips[0].String(), nil
}

// CanEarlyExit 是否可以早期退出
func (irs *IPResolverStage) CanEarlyExit() bool {
	return true
}

// Priority 优先级
func (irs *IPResolverStage) Priority() int {
	return 3 // IP解析第三优先级
}

// Name 阶段名称
func (irs *IPResolverStage) Name() string {
	return "ip_resolver"
}
