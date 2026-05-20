# pswitch

[English](README.md) | [中文](README_ZH.md)

`pswitch` 是一个本地多 Provider 代理，适合把多个上游模型服务统一收口到一个稳定的本地入口，并提供自动切换、熔断恢复和可视化后台。

## 截图

![pswitch dashboard](docs/assets/dashboard.png)

## 核心能力

- 多个上游 Provider 自动切换
- 熔断与后台健康恢复探测
- 三种路由模式：
  - `round_robin`
  - `sequential`
  - `least_failures`
- 默认支持 OpenAI 风格路由
- 可选 Anthropic 风格协议适配
- 内置 `/dashboard/` 管理后台
- 持久化统计请求数、Token、失败次数、按模型使用量
- 运行配置支持页面编辑，尽量热更新
- 运行态文件持久化到当前目录：
  - `settings.json`
  - `metrics.json`

## 快速开始

### 下载 Release

打开 [Releases](https://github.com/wlynxg/pswitch/releases/latest)，下载与你平台对应的压缩包：

- Linux x86_64：`pswitch_vX.Y.Z_linux_amd64.tar.gz`
- Linux ARM64：`pswitch_vX.Y.Z_linux_arm64.tar.gz`
- macOS Intel：`pswitch_vX.Y.Z_darwin_amd64.tar.gz`
- macOS Apple Silicon：`pswitch_vX.Y.Z_darwin_arm64.tar.gz`
- Windows x86_64：`pswitch_vX.Y.Z_windows_amd64.zip`

解压后，在解压目录中直接运行二进制即可。

### 启动服务

首次启动不需要配置文件。

macOS / Linux：

```bash
./pswitch
```

Windows PowerShell：

```powershell
.\pswitch.exe
```

默认行为：

- 监听 `0.0.0.0:8080`
- 路由模式为 `round_robin`
- 默认路由只有 `/codex`
- 不预置 Provider

### 打开网页后台

本机访问：

```text
http://127.0.0.1:8080/dashboard/
```

服务器或局域网访问：

```text
http://<服务器IP>:8080/dashboard/
```

你可以直接在 `Config` 页面里添加 Provider 并保存运行配置。程序会在当前运行目录写入 `settings.json` 和 `metrics.json`。

### 配置客户端

把客户端指向本地代理：

```text
http://127.0.0.1:8080/codex
```

Codex 风格配置示例：

```toml
[model_providers.OpenAI]
base_url = "http://127.0.0.1:8080/codex"
wire_api = "responses"
requires_openai_auth = true
```

### 从源码构建

```bash
make build
```

然后运行：

```bash
./bin/pswitch
```

## 文档

- [配置说明](docs/config.md)
- [使用方法](docs/usage.md)
- [日志](docs/logging.md)
- [故障排查](docs/troubleshooting.md)
- [开发](docs/development.md)

## 自动发布

推送版本标签后，GitHub Actions 会自动构建各平台包并发布到 Releases：

```bash
git tag v0.1.0
git push origin v0.1.0
```
