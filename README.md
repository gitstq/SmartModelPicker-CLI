<div align="center">

# 🧠 SmartModelPicker-CLI

**本地LLM硬件适配智能推荐引擎 | Local LLM Hardware Matcher**

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)](https://github.com/gitstq/SmartModelPicker-CLI/releases)
[![Release](https://img.shields.io/badge/Release-v1.0.0-orange)](https://github.com/gitstq/SmartModelPicker-CLI/releases)

[简体中文](#简体中文) | [繁體中文](#繁體中文) | [English](#english)

</div>

---

## 简体中文

### 🎉 项目介绍

**SmartModelPicker-CLI** 是一款专为本地大语言模型（LLM）部署设计的智能硬件适配推荐工具。它能自动检测你的电脑硬件配置（CPU、GPU、内存、指令集支持），并从内置的 **23+ 主流开源模型数据库** 中，智能推荐最适合你硬件运行的模型及最佳量化方案。

**灵感来源**：受到 [whichllm](https://github.com/Andyyyy64/whichllm) 项目的启发，但我们做了全面自研重构：
- 用 **Go** 重写（vs Python），零依赖、单二进制文件分发
- 支持 **Ollama / LM Studio / llama.cpp** 多后端
- 集成 **GPU显存智能计算** + **量化级别自动推荐**
- 支持 **Apple Silicon / AMD / NVIDIA** 全平台硬件检测
- 一键生成 **docker-compose** 部署配置

### ✨ 核心特性

| 特性 | 说明 |
|------|------|
| 🔍 **智能硬件检测** | 自动识别 CPU核心数、系统内存、GPU型号与显存、AVX2/AVX-512指令集支持 |
| 🧠 **23+ 模型数据库** | 覆盖 Llama、Mistral、Qwen、DeepSeek、Gemma、Phi、Yi 等主流开源模型 |
| 🎯 **智能量化推荐** | 根据可用显存自动选择 FP16/Q8_0/Q6_K/Q5_K_M/Q4_K_M/Q4_0 最佳量化方案 |
| 🖥️ **多平台支持** | Linux / macOS (Intel & Apple Silicon) / Windows 全平台兼容 |
| 🚀 **一键运行命令** | 自动生成 `ollama run` 命令，复制即可运行 |
| 📦 **配置导出** | 支持导出 docker-compose / Shell脚本 / JSON 配置 |
| 🌈 **彩色终端** | 美观的彩色输出，同时支持 `--json` 结构化输出 |
| ⚡ **零依赖** | 单二进制文件，无需 Python/Node.js 运行时 |

### 🚀 快速开始

#### 环境要求
- **操作系统**: Linux / macOS 10.15+ / Windows 10+
- **架构**: amd64 或 arm64
- **可选**: 已安装 [Ollama](https://ollama.com) 或 [LM Studio](https://lmstudio.ai)

#### 安装

**方式一：直接下载二进制（推荐）**

```bash
# Linux/macOS
curl -fsSL https://github.com/gitstq/SmartModelPicker-CLI/releases/latest/download/smartmodelpicker-linux-amd64 -o smartmodelpicker
chmod +x smartmodelpicker
sudo mv smartmodelpicker /usr/local/bin/

# macOS Apple Silicon
curl -fsSL https://github.com/gitstq/SmartModelPicker-CLI/releases/latest/download/smartmodelpicker-darwin-arm64 -o smartmodelpicker
chmod +x smartmodelpicker
sudo mv smartmodelpicker /usr/local/bin/
```

**方式二：从源码编译**

```bash
git clone https://github.com/gitstq/SmartModelPicker-CLI.git
cd SmartModelPicker-CLI
make build
sudo make install
```

#### 基本使用

```bash
# 运行完整检测与推荐
smartmodelpicker

# 仅检测硬件信息
smartmodelpicker check

# 列出所有支持的模型
smartmodelpicker list

# 指定量化级别
smartmodelpicker -q Q4_K_M

# 导出 docker-compose 配置
smartmodelpicker -e docker-compose

# JSON 输出（适合脚本集成）
smartmodelpicker --json
```

### 📖 详细使用指南

#### 命令行参数

```
Flags:
  -b, --backend string   指定LLM后端 (auto|ollama|lmstudio|llamacpp)
  -q, --quantize string  量化级别 (auto|Q4_0|Q4_K_M|Q5_K_M|Q6_K|Q8_0|FP16)
      --min-vram float   最小显存要求(GB)
      --max-vram float   最大显存限制(GB)
  -a, --all              显示所有模型（不限于可运行）
  -e, --export string    导出配置 (docker-compose|shell|json)
      --json             以JSON格式输出
  -v, --verbose          详细输出模式
  -h, --help             帮助信息
```

#### 典型使用场景

**场景1：新手首次使用**
```bash
$ smartmodelpicker

╔══════════════════════════════════════════════════════════════╗
║           🧠 SmartModelPicker-CLI v1.0.0 🚀                  ║
╚══════════════════════════════════════════════════════════════╝

📊 硬件配置检测
──────────────────────────────────────────────────
   🖥️  操作系统: linux (amd64)
   🔲 CPU核心: 8 核
   💾 系统内存: 32.0 GB
   🎮 显卡信息:
      GPU 1: NVIDIA GeForce RTX 3060 (NVIDIA) - 12.0 GB VRAM
   📊 总显存: 12.0 GB
   ⚡ 加速支持:
      ✅ NVIDIA CUDA
      ❌ AMD ROCm
      ❌ Apple Metal
      ✅ AVX2 | ❌ AVX-512
   🎯 推荐后端: llama.cpp (CUDA)

🎯 智能推荐结果
══════════════════════════════════════════════════════

✅ 推荐运行（按匹配度排序）
──────────────────────────────────────────────────────

   🥇 Llama 3.1 8B 🥈 主流
      📋 Meta最新开源大模型，8B参数，性能优异，支持128K上下文
      🎯 无损精度 | 💾 显存占用 16.0GB (133%)
      🔧 最佳量化: FP16 | 需要显存: 16.0 GB
      💻 一键运行: ollama run llama3.1-8b:fp16

   🥈 Qwen 2.5 7B 🥈 主流
      📋 通义千问2.5，7B参数，中文能力突出，支持128K上下文
      ⚖️ 平衡质量与速度 | 💾 显存占用 5.1GB (43%) | 🇨🇳 中文优化
      🔧 最佳量化: Q4_K_M | 需要显存: 5.1 GB
      💻 一键运行: ollama run qwen2.5-7b:q4_k_m
```

**场景2：导出部署配置**
```bash
$ smartmodelpicker -e docker-compose
# 生成 smartmodelpicker-docker-compose.yaml
# 包含 Ollama + Open WebUI 的完整配置

$ smartmodelpicker -e shell
# 生成 smartmodelpicker-shell.sh
# 自动安装 Ollama 并拉取推荐模型
```

**场景3：JSON输出集成**
```bash
$ smartmodelpicker --json | jq '.recommendations[0]'
{
  "model": { "id": "llama3.1-8b", "name": "Llama 3.1 8B", ... },
  "best_quantization": "Q4_K_M",
  "required_vram_gb": 5.1,
  "can_run": true,
  "score": 85,
  "ollama_command": "ollama run llama3.1-8b:q4_k_m"
}
```

### 💡 设计思路与迭代规划

#### 技术选型原因
- **Go语言**: 编译为单二进制文件，零依赖，跨平台分发简单
- **Cobra框架**: 业界标准的CLI框架，支持子命令和丰富参数
- **Color库**: 终端彩色输出，提升用户体验
- **内置模型数据库**: 无需联网即可使用，保证隐私和速度

#### 后续迭代计划
- [ ] 支持从 HuggingFace API 实时拉取最新模型
- [ ] 添加模型下载速度测试与推荐
- [ ] 支持多GPU环境下的模型并行推荐
- [ ] 集成性能基准测试（tokens/s 预测）
- [ ] Web UI 版本开发
- [ ] 支持更多后端（vLLM、TGI等）

### 📦 打包与部署指南

#### 本地构建
```bash
make build        # 当前平台
make build-all    # 全平台 (Linux/macOS/Windows)
make test         # 运行测试
make install      # 安装到 $GOPATH/bin
```

#### 跨平台编译
```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o smartmodelpicker-linux-amd64 .

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o smartmodelpicker-darwin-arm64 .

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o smartmodelpicker-windows-amd64.exe .
```

### 🤝 贡献指南

欢迎提交 Issue 和 PR！

- **Bug反馈**: 请使用 `bug` 标签，附上硬件信息和复现步骤
- **功能建议**: 请使用 `enhancement` 标签
- **模型补充**: 欢迎提交 PR 补充新模型数据
- **代码规范**: 遵循 Go 官方代码规范，提交前运行 `make fmt`

### 📄 开源协议

本项目采用 [MIT License](LICENSE) 开源协议。

---

## 繁體中文

### 🎉 專案介紹

**SmartModelPicker-CLI** 是一款專為本地大型語言模型（LLM）部署設計的智慧硬體適配推薦工具。它能自動檢測你的電腦硬體配置（CPU、GPU、記憶體、指令集支援），並從內建的 **23+ 主流開源模型資料庫** 中，智慧推薦最適合你硬體運行的模型及最佳量化方案。

**靈感來源**：受到 [whichllm](https://github.com/Andyyyy64/whichllm) 專案的啟發，但我們做了全面自研重構：
- 用 **Go** 重寫（vs Python），零依賴、單二進制檔案分發
- 支援 **Ollama / LM Studio / llama.cpp** 多後端
- 整合 **GPU顯存智慧計算** + **量化級別自動推薦**
- 支援 **Apple Silicon / AMD / NVIDIA** 全平台硬體檢測
- 一鍵生成 **docker-compose** 部署配置

### ✨ 核心特性

| 特性 | 說明 |
|------|------|
| 🔍 **智慧硬體檢測** | 自動識別 CPU核心數、系統記憶體、GPU型號與顯存、AVX2/AVX-512指令集支援 |
| 🧠 **23+ 模型資料庫** | 覆蓋 Llama、Mistral、Qwen、DeepSeek、Gemma、Phi、Yi 等主流開源模型 |
| 🎯 **智慧量化推薦** | 根據可用顯存自動選擇 FP16/Q8_0/Q6_K/Q5_K_M/Q4_K_M/Q4_0 最佳量化方案 |
| 🖥️ **多平台支援** | Linux / macOS (Intel & Apple Silicon) / Windows 全平台相容 |
| 🚀 **一鍵運行命令** | 自動生成 `ollama run` 命令，複製即可運行 |
| 📦 **配置匯出** | 支援匯出 docker-compose / Shell腳本 / JSON 配置 |
| 🌈 **彩色終端** | 美觀的彩色輸出，同時支援 `--json` 結構化輸出 |
| ⚡ **零依賴** | 單二進制檔案，無需 Python/Node.js 運行時 |

### 🚀 快速開始

#### 環境要求
- **作業系統**: Linux / macOS 10.15+ / Windows 10+
- **架構**: amd64 或 arm64
- **可選**: 已安裝 [Ollama](https://ollama.com) 或 [LM Studio](https://lmstudio.ai)

#### 安裝

**方式一：直接下載二進制（推薦）**

```bash
# Linux/macOS
curl -fsSL https://github.com/gitstq/SmartModelPicker-CLI/releases/latest/download/smartmodelpicker-linux-amd64 -o smartmodelpicker
chmod +x smartmodelpicker
sudo mv smartmodelpicker /usr/local/bin/

# macOS Apple Silicon
curl -fsSL https://github.com/gitstq/SmartModelPicker-CLI/releases/latest/download/smartmodelpicker-darwin-arm64 -o smartmodelpicker
chmod +x smartmodelpicker
sudo mv smartmodelpicker /usr/local/bin/
```

**方式二：從原始碼編譯**

```bash
git clone https://github.com/gitstq/SmartModelPicker-CLI.git
cd SmartModelPicker-CLI
make build
sudo make install
```

#### 基本使用

```bash
# 運行完整檢測與推薦
smartmodelpicker

# 僅檢測硬體資訊
smartmodelpicker check

# 列出所有支援的模型
smartmodelpicker list

# 指定量化級別
smartmodelpicker -q Q4_K_M

# 匯出 docker-compose 配置
smartmodelpicker -e docker-compose

# JSON 輸出（適合腳本整合）
smartmodelpicker --json
```

### 📖 詳細使用指南

#### 命令列參數

```
Flags:
  -b, --backend string   指定LLM後端 (auto|ollama|lmstudio|llamacpp)
  -q, --quantize string  量化級別 (auto|Q4_0|Q4_K_M|Q5_K_M|Q6_K|Q8_0|FP16)
      --min-vram float   最小顯存要求(GB)
      --max-vram float   最大顯存限制(GB)
  -a, --all              顯示所有模型（不限於可運行）
  -e, --export string    匯出配置 (docker-compose|shell|json)
      --json             以JSON格式輸出
  -v, --verbose          詳細輸出模式
  -h, --help             幫助資訊
```

### 💡 設計思路與迭代規劃

#### 技術選型原因
- **Go語言**: 編譯為單二進制檔案，零依賴，跨平台分發簡單
- **Cobra框架**: 業界標準的CLI框架，支援子命令和豐富參數
- **Color庫**: 終端彩色輸出，提升使用者體驗
- **內建模型資料庫**: 無需連網即可使用，保證隱私和速度

#### 後續迭代計劃
- [ ] 支援從 HuggingFace API 即時拉取最新模型
- [ ] 添加模型下載速度測試與推薦
- [ ] 支援多GPU環境下的模型並行推薦
- [ ] 整合性能基準測試（tokens/s 預測）
- [ ] Web UI 版本開發
- [ ] 支援更多後端（vLLM、TGI等）

### 📦 打包與部署指南

#### 本地構建
```bash
make build        # 當前平台
make build-all    # 全平台 (Linux/macOS/Windows)
make test         # 運行測試
make install      # 安裝到 $GOPATH/bin
```

### 🤝 貢獻指南

歡迎提交 Issue 和 PR！

- **Bug反饋**: 請使用 `bug` 標籤，附上硬體資訊和復現步驟
- **功能建議**: 請使用 `enhancement` 標籤
- **模型補充**: 歡迎提交 PR 補充新模型資料
- **代碼規範**: 遵循 Go 官方代碼規範，提交前運行 `make fmt`

### 📄 開源協議

本專案採用 [MIT License](LICENSE) 開源協議。

---

## English

### 🎉 Introduction

**SmartModelPicker-CLI** is an intelligent hardware matching tool designed for local Large Language Model (LLM) deployment. It automatically detects your computer's hardware configuration (CPU, GPU, RAM, instruction set support) and intelligently recommends the best-suited open-source models and optimal quantization schemes from its built-in **23+ mainstream model database**.

**Inspired by**: [whichllm](https://github.com/Andyyyy64/whichllm), but completely reimagined and rebuilt from scratch:
- Written in **Go** (vs Python), zero dependencies, single binary distribution
- Supports **Ollama / LM Studio / llama.cpp** multiple backends
- Integrated **GPU VRAM smart calculation** + **automatic quantization recommendation**
- Supports **Apple Silicon / AMD / NVIDIA** full-platform hardware detection
- One-click **docker-compose** deployment config generation

### ✨ Key Features

| Feature | Description |
|---------|-------------|
| 🔍 **Smart Hardware Detection** | Auto-detect CPU cores, system RAM, GPU model & VRAM, AVX2/AVX-512 support |
| 🧠 **23+ Model Database** | Covers Llama, Mistral, Qwen, DeepSeek, Gemma, Phi, Yi and more |
| 🎯 **Intelligent Quantization** | Auto-select best FP16/Q8_0/Q6_K/Q5_K_M/Q4_K_M/Q4_0 based on available VRAM |
| 🖥️ **Multi-Platform** | Linux / macOS (Intel & Apple Silicon) / Windows support |
| 🚀 **One-Click Commands** | Auto-generate `ollama run` commands, copy and run |
| 📦 **Config Export** | Export docker-compose / Shell script / JSON configurations |
| 🌈 **Colorful Terminal** | Beautiful colored output, plus `--json` structured output |
| ⚡ **Zero Dependencies** | Single binary, no Python/Node.js runtime required |

### 🚀 Quick Start

#### Requirements
- **OS**: Linux / macOS 10.15+ / Windows 10+
- **Architecture**: amd64 or arm64
- **Optional**: [Ollama](https://ollama.com) or [LM Studio](https://lmstudio.ai) installed

#### Installation

**Option 1: Download Binary (Recommended)**

```bash
# Linux/macOS
curl -fsSL https://github.com/gitstq/SmartModelPicker-CLI/releases/latest/download/smartmodelpicker-linux-amd64 -o smartmodelpicker
chmod +x smartmodelpicker
sudo mv smartmodelpicker /usr/local/bin/

# macOS Apple Silicon
curl -fsSL https://github.com/gitstq/SmartModelPicker-CLI/releases/latest/download/smartmodelpicker-darwin-arm64 -o smartmodelpicker
chmod +x smartmodelpicker
sudo mv smartmodelpicker /usr/local/bin/
```

**Option 2: Build from Source**

```bash
git clone https://github.com/gitstq/SmartModelPicker-CLI.git
cd SmartModelPicker-CLI
make build
sudo make install
```

#### Basic Usage

```bash
# Run full detection and recommendation
smartmodelpicker

# Check hardware only
smartmodelpicker check

# List all supported models
smartmodelpicker list

# Specify quantization level
smartmodelpicker -q Q4_K_M

# Export docker-compose config
smartmodelpicker -e docker-compose

# JSON output (for scripting)
smartmodelpicker --json
```

### 📖 Detailed Usage

#### CLI Flags

```
Flags:
  -b, --backend string   LLM backend (auto|ollama|lmstudio|llamacpp)
  -q, --quantize string  Quantization level (auto|Q4_0|Q4_K_M|Q5_K_M|Q6_K|Q8_0|FP16)
      --min-vram float   Minimum VRAM requirement (GB)
      --max-vram float   Maximum VRAM limit (GB)
  -a, --all              Show all models (not just runnable)
  -e, --export string    Export config (docker-compose|shell|json)
      --json             Output in JSON format
  -v, --verbose          Verbose output mode
  -h, --help             Help information
```

#### Example Scenarios

**Scenario 1: First-time User**
```bash
$ smartmodelpicker

╔══════════════════════════════════════════════════════════════╗
║           🧠 SmartModelPicker-CLI v1.0.0 🚀                  ║
╚══════════════════════════════════════════════════════════════╝

📊 Hardware Detection
──────────────────────────────────────────────────
   🖥️  OS: linux (amd64)
   🔲 CPU Cores: 8
   💾 System RAM: 32.0 GB
   🎮 GPU: NVIDIA GeForce RTX 3060 - 12.0 GB VRAM
   ⚡ Acceleration:
      ✅ NVIDIA CUDA
      ✅ AVX2
   🎯 Recommended Backend: llama.cpp (CUDA)

🎯 Recommendations (sorted by match score)

   🥇 Llama 3.1 8B
      📋 Meta's latest open model, 8B params, 128K context
      🎯 Lossless precision | 💾 VRAM: 16.0GB
      💻 Run: ollama run llama3.1-8b:fp16
```

**Scenario 2: Export Deployment Config**
```bash
$ smartmodelpicker -e docker-compose
# Generates smartmodelpicker-docker-compose.yaml
# Includes Ollama + Open WebUI complete setup
```

**Scenario 3: JSON Integration**
```bash
$ smartmodelpicker --json | jq '.recommendations[0].ollama_command'
"ollama run llama3.1-8b:q4_k_m"
```

### 💡 Design & Roadmap

#### Tech Stack Rationale
- **Go**: Compiles to single binary, zero dependencies, easy cross-platform distribution
- **Cobra**: Industry-standard CLI framework with subcommand support
- **Color**: Terminal colored output for better UX
- **Built-in Model DB**: Works offline, ensures privacy and speed

#### Roadmap
- [ ] Real-time model fetching from HuggingFace API
- [ ] Model download speed testing
- [ ] Multi-GPU environment recommendations
- [ ] Performance benchmark predictions (tokens/s)
- [ ] Web UI version
- [ ] Additional backend support (vLLM, TGI, etc.)

### 📦 Build & Deploy

#### Local Build
```bash
make build        # Current platform
make build-all    # All platforms (Linux/macOS/Windows)
make test         # Run tests
make install      # Install to $GOPATH/bin
```

#### Cross-Platform Compilation
```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o smartmodelpicker-linux-amd64 .

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o smartmodelpicker-darwin-arm64 .

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o smartmodelpicker-windows-amd64.exe .
```

### 🤝 Contributing

Issues and PRs are welcome!

- **Bug Reports**: Use `bug` label, include hardware info and reproduction steps
- **Feature Requests**: Use `enhancement` label
- **Model Additions**: Submit PRs to add new models
- **Code Style**: Follow Go official style guide, run `make fmt` before committing

### 📄 License

This project is licensed under the [MIT License](LICENSE).

---

<div align="center">

**Made with ❤️ by SmartModelPicker Team**

[⭐ Star us on GitHub](https://github.com/gitstq/SmartModelPicker-CLI) | [🐛 Report Issue](https://github.com/gitstq/SmartModelPicker-CLI/issues) | [💡 Feature Request](https://github.com/gitstq/SmartModelPicker-CLI/issues)

</div>
