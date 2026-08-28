package data

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Downloader is retained for API compatibility. The hardened build never
// downloads or silently replaces data at runtime.
type Downloader struct{}

func NewDownloader() *Downloader { return &Downloader{} }

var expectedFiles = map[string]string{
	"cdn_keywords.txt": "3526689af9ba522084b6bb39f55a525cead6a88555898bda2ae4f2e8d0364626",
	"Country.mmdb":     "577a545e33aa6375d844e28c7becc6f57f40ed435b25cf3687616846ae4f7644",
	"gfwlist.conf":     "de612f34d66f023b7a6c03eb72aa1dcaf1feb147deb8e74f0b24cb8ed7d9f06f",
	"hot_websites.txt": "92e773f5e55e4037d924cabc6c8c45cefd042429cafd8468a764d54c3f9eee4a",
}

func printTimestampedMessage(format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s] %s\n", timestamp, fmt.Sprintf(format, args...))
}

// ResolvePath locates a bundled, read-only data file. An explicit directory
// can be supplied for advanced use, but its contents are still hash checked.
func ResolvePath(name string) (string, error) {
	if _, ok := expectedFiles[name]; !ok {
		return "", fmt.Errorf("未知数据文件: %s", name)
	}

	var candidates []string
	if dir := os.Getenv("REALITYCHECK_DATA_DIR"); dir != "" {
		candidates = append(candidates, filepath.Join(dir, name))
	}
	if executable, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(executable), "data", name))
	}
	candidates = append(candidates, filepath.Join("data", name))

	for _, candidate := range candidates {
		if info, err := os.Lstat(candidate); err == nil && info.Mode().IsRegular() {
			absolute, err := filepath.Abs(candidate)
			if err != nil {
				return "", err
			}
			return absolute, nil
		}
	}
	return "", fmt.Errorf("缺少数据文件 %s；请保留程序旁边的 data 目录", name)
}

// EnsureDataFiles verifies every bundled file before any scan begins.
func (d *Downloader) EnsureDataFiles() error {
	printTimestampedMessage("校验离线数据文件...")
	for name, expected := range expectedFiles {
		path, err := ResolvePath(name)
		if err != nil {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("打开数据文件 %s 失败: %w", name, err)
		}
		hash := sha256.New()
		_, copyErr := io.Copy(hash, file)
		closeErr := file.Close()
		if copyErr != nil {
			return fmt.Errorf("校验数据文件 %s 失败: %w", name, copyErr)
		}
		if closeErr != nil {
			return fmt.Errorf("关闭数据文件 %s 失败: %w", name, closeErr)
		}
		actual := fmt.Sprintf("%x", hash.Sum(nil))
		if actual != expected {
			return fmt.Errorf("数据文件 %s 完整性校验失败（应为 %s，实际为 %s）", name, expected, actual)
		}
	}
	printTimestampedMessage("离线数据完整性校验通过。")
	return nil
}
