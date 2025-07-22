# goproxy 简洁的socks5和http代理服务


- [*] 支持socks5无认证方式代理
- [*] 支持http代理，兼容http和https协议

运行示例：
``` sh
goproxy.exe -L 127.0.0.1:1080
```
-L 参数为监听地址，默认127.0.0.1:11088

提示：
1. 建议运行在可靠的内部环境
2. socks5暂不支持ipv6代理

go模块依赖代理：
``` sh
# 启用 Go Modules 功能
go env -w GO111MODULE=on
# 配置 GOPROXY 环境变量 阿里云
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

多平台构建
``` sh
# Linux 64位
GOOS=linux GOARCH=amd64 go build -o goproxy-linux-amd64
# Windows 64位 
GOOS=windows GOARCH=amd64 go build -o goproxy-windows-amd64.exe
# MacOS 64位
GOOS=darwin GOARCH=amd64 go build -o goproxy-darwin-amd64
# ARM架构 (MacOS）
GOOS=darwin GOARCH=arm64 go build -o goproxy-darwin-arm64
# ARM架构 (Linux)
GOOS=linux GOARCH=arm64 go build -o goproxy-linux-arm64
```

参考链接：

- socks5协议：
https://zhuanlan.zhihu.com/p/11284611035
https://www.rfc-editor.org/rfc/rfc1928
- flag模块：
https://blog.csdn.net/qq_38378384/article/details/121720787
- go模块依赖代理：
https://learnku.com/go/wikis/38122
- http协议：
https://www.rfc-editor.org/rfc/rfc1928



