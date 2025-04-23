# gbuild — Go Cross-Platform Builder

`gbuild` 是一个轻量级的命令行工具，用于在不同平台（Linux、Windows、macOS）和架构（amd64、arm64）下快速编译 Go 应用。

---

## 🚀 安装

### 源码安装（开发或本地测试）
```bash
# 进入项目目录（包含 gbuild.go）
cd gbuild
# 初始化模块
go mod init github.com/ljx619/gbuild
# 安装
GO111MODULE=on
go install
```
安装后，默认二进制会放在:
- Unix/macOS: `$HOME/go/bin/gbuild`
- Windows: `%USERPROFILE%\go\bin\gbuild`

> 确保把 `$GOPATH/bin`（或 `$GOBIN`） 加到 `PATH` 环境变量。

### 一键安装（正式发布后）
```bash
go install github.com/ljx619/gbuild@latest
```

---

## ⚡ 快速开始

```bash
# 构建当前系统默认平台/架构
gbuild

# 构建指定平台
gbuild -os linux -arch amd64

# 构建并自定义输出路径及开启 CGO
gbuild -os windows -arch arm64 -o dist/myapp.exe -cgo
```

---

## 🔧 命令行参数

| 参数           | 类型    | 默认值               | 描述                                                         |
|----------------|---------|----------------------|--------------------------------------------------------------|
| `-os`          | string  | `runtime.GOOS`       | 目标操作系统：`linux`, `windows`, `darwin`                  |
| `-arch`        | string  | `runtime.GOARCH`     | 目标架构：`amd64`, `arm64`                                  |
| `-o`           | string  | `./bin/build-<os>-<arch>` | 输出文件路径，可指定目录及文件名                              |
| `-cgo`         | bool    | `false`              | 启用 CGO                                                    |
| `-tags`        | string  | `""`               | Go 构建标签，多个标签用空格隔开                               |
| `-ldflags`     | string  | `""`               | Go 链接器参数，如 `-s -w`                                   |
| `-hash`        | bool    | `false`              | 构建后打印二进制文件的 SHA256 哈希值                          |
| `-version`     | bool    | `false`              | 打印 `gbuild` 版本号、构建时间及当前 Go 版本                 |

---

## 🖥️ 支持平台

| 操作系统 | 架构         |
|----------|--------------|
| linux    | amd64, arm64 |
| windows  | amd64, arm64 |
| darwin   | amd64, arm64 |

> ⚠️ 不再支持 `386` 架构

---

## 💡 常见示例

- **构建 Linux/amd64 并打印 SHA256**
  ```bash
  gbuild -os linux -arch amd64 -hash
  ```

- **自定义构建标签与链接器参数**
  ```bash
  gbuild -tags "prod netgo" -ldflags "-s -w"
  ```

- **查看版本信息**
  ```bash
  gbuild -version
  ```
  
---

## 📜 License

本项目使用 **MIT License**，详见 [LICENSE](LICENSE)。

---

欢迎提交 **issues** 和 **pull requests**！如果你有新的功能建议或 bug 修复，感谢参与贡献 🙏

****