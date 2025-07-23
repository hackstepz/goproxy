# GoProxy - 轻量级SOCKS5/HTTP代理服务器

## 功能特性

- ✅ 支持SOCKS5代理协议（无认证方式）
- ✅ 支持HTTP/HTTPS代理协议
- 🚀 轻量高效，单文件可执行

## 快速开始

### 基本使用
```bash
goproxy.exe -L 127.0.0.1:1080
```

**参数说明**:
- `-L` 指定监听地址，默认为 `127.0.0.1:11088`

### 注意事项
1. 建议仅在可信的内部网络环境中使用
2. 当前版本socks5协议暂不支持IPv6代理

## 开发环境配置

### Go Modules 设置
```bash
# 启用Go Modules
go env -w GO111MODULE=on

# 配置国内镜像源（阿里云）
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

## 多平台构建

```bash
# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o build/goproxy-linux-amd64

# Windows (amd64)
GOOS=windows GOARCH=amd64 go build -o build/goproxy-windows-amd64.exe

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o build/goproxy-darwin-amd64

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o build/goproxy-darwin-arm64

# Linux (ARM64)
GOOS=linux GOARCH=arm64 go build -o build/goproxy-linux-arm64
```

## 参考文档

### 协议规范
- [SOCKS5协议 (RFC 1928)](https://www.rfc-editor.org/rfc/rfc1928)
- [HTTP协议](https://www.rfc-editor.org/rfc/rfc2616)

### 开发相关
- [Go Flag模块使用指南](https://pkg.go.dev/flag)
- [Go Modules官方文档](https://go.dev/ref/mod)
