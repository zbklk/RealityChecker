# RealityChecker 安全加固说明

## 版本与来源

- 加固版本：`hardened-20260828`
- 基础提交：`8a152d4014b3781fa1e004f51e37dabcb1ffcbcf`
- fork 与 `V2RaySSR/RealityChecker` 上游 `main` 在审查时为同一提交，差异为 `0 ahead / 0 behind`。
- 构建工具链：Go 1.27.0，官方 Windows AMD64 归档 SHA-256 为 `f0c0a0d33ba94f4d2c5dbc887334ce678b21813504ddb3aafcb06e60a5a667c4`。

## 加固内容

1. 删除运行时自动下载以及每三天静默更新规则数据的逻辑。
2. 固定四个离线数据文件的 SHA-256；程序启动时先校验，任何文件不匹配都会拒绝运行。
3. 优先从可执行文件旁边的 `data` 目录加载文件，避免依赖启动目录。
4. 取消启动时的 GitHub Release 查询和推广广告，不再产生这些额外联网请求。
5. 所有主动网络连接只允许公网地址的 TCP 80/443；拒绝回环、RFC1918 私网、链路本地、共享地址空间、云元数据和常见保留网段。
6. 网络层先解析并检查所有 DNS 结果，然后直接连接已检查的 IP，降低 DNS rebinding 风险。
7. 更新 `golang.org/x/sys` 至 v0.47.0、`golang.org/x/text` 至 v0.41.0。
8. 增加离线数据完整性和非公网地址拦截单元测试。

## 验证结果

- `go test ./...`：通过。
- `govulncheck ./...`：`No vulnerabilities found`。
- `govulncheck -mode=binary reality-checker.exe`：`No vulnerabilities found`。
- Windows AMD64、Linux AMD64、Linux ARM64 均连续构建两次，二进制逐字节一致。
- 最终 Windows ZIP 重新解压后，`version` 命令正常，数据完整性校验通过。
- 对 `127.0.0.1` 的测试没有建立连接，由安全网络层拒绝。

## 构建参数

安全构建使用：

```text
CGO_ENABLED=0
go build -trimpath -buildvcs=false -ldflags "-s -w -buildid= ..."
```

Release 必须同时包含可执行文件、完整 `data/` 目录、使用说明和 SHA-256 清单。

## 使用限制

- 本工具会主动连接和探测目标网站。只能用于你有权检测的目标，并应遵守当地法律、服务条款和网络使用政策。
- Windows 构建没有商业 Authenticode 证书签名，SmartScreen 可能显示未知发布者；请使用 Release 附带的 SHA-256 清单核对文件。
- 安全版数据是固定快照，不会自动更新。若需要更新数据，应在源码中同步更新固定哈希、重新审查并重新构建。
- 上游仓库未提供明确的软件许可证。本分支及构建产物不应在未确认许可的情况下再次分发或商用。

## 报告问题

安全问题请不要在公开 Issue 中附带可直接利用的细节。建议先通过仓库所有者公开提供的私密联系方式进行报告；确认修复窗口后再公开披露。
