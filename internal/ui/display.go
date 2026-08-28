package ui

import (
	"fmt"
	"time"

	"RealityChecker/internal/version"
)

// GitHubRelease GitHub发布信息结构
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

// PrintUsage 打印使用说明
func PrintUsage() {
	fmt.Printf("Reality协议目标网站检测器 %s\n\n", version.GetVersion())
	fmt.Println("用法:")
	fmt.Println("  reality-checker check <domain>          检测单个域名")
	fmt.Println("  reality-checker batch <domain1> <domain2> <domain3> ...  批量检测域名")
	fmt.Println("  reality-checker csv <csv_file>          从CSV文件批量检测域名")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  reality-checker check apple.com")
	fmt.Println("  reality-checker batch apple.com tesla.com microsoft.com")
	fmt.Println("  reality-checker csv file.csv")
}

// PrintTimestampedMessage 打印带时间戳的消息
func PrintTimestampedMessage(format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	message := fmt.Sprintf(format, args...)
	fmt.Printf("[%s] %s\n", timestamp, message)
}

// PrintError 打印错误信息（带空行间距）
func PrintError(message string) {
	fmt.Println()
	fmt.Println(message)
	fmt.Println()
}

// PrintErrorWithDetails 打印错误信息和详细信息
func PrintErrorWithDetails(message string, details ...string) {
	fmt.Println()
	fmt.Println(message)
	for _, detail := range details {
		fmt.Println(detail)
	}
	fmt.Println()
}

// getLatestVersion 获取GitHub最新版本号
func getLatestVersion() string {
	return ""
}

// getVersionInfo 获取版本信息字符串
func getVersionInfo() string {
	currentVersion := version.GetVersion() // 使用动态版本号

	latestVersion := getLatestVersion()

	if latestVersion == "" {
		// 无法获取最新版本，只显示当前版本
		return currentVersion
	}

	if latestVersion == currentVersion {
		// 版本相同，只显示当前版本
		return currentVersion
	}

	// 版本不同，显示当前版本和最新版本
	return fmt.Sprintf("%s (最新: %s)", currentVersion, latestVersion)
}

// getDisplayWidth 计算字符串的显示宽度（中文字符占2个位置）
func getDisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		if r < 128 {
			// ASCII字符占1个位置
			width++
		} else {
			// 中文字符占2个位置
			width += 2
		}
	}
	return width
}

// PrintAdvertisement 打印广告信息
func PrintAdvertisement() {
	// Hardened build intentionally contains no advertising or tracking links.
}
