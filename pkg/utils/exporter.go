package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/user/smartmodelpicker-cli/pkg/hardware"
	"github.com/user/smartmodelpicker-cli/pkg/recommender"
)

// Exporter 配置导出器
type Exporter struct {
	hwInfo          *hardware.Info
	recommendations []recommender.Recommendation
}

// NewExporter 创建导出器
func NewExporter(hwInfo *hardware.Info, recommendations []recommender.Recommendation) *Exporter {
	return &Exporter{
		hwInfo:          hwInfo,
		recommendations: recommendations,
	}
}

// Export 导出配置
func (e *Exporter) Export(format string) (string, error) {
	switch format {
	case "docker-compose":
		return e.exportDockerCompose(), nil
	case "shell":
		return e.exportShell(), nil
	case "json":
		return e.exportJSON(), nil
	default:
		return "", fmt.Errorf("不支持的导出格式: %s", format)
	}
}

// exportDockerCompose 导出docker-compose配置
func (e *Exporter) exportDockerCompose() string {
	var models []recommender.Recommendation
	for _, rec := range e.recommendations {
		if rec.CanRun {
			models = append(models, rec)
		}
	}

	var sb strings.Builder
	sb.WriteString("# SmartModelPicker 自动生成的 Docker Compose 配置\n")
	sb.WriteString("# 基于硬件: ")
	sb.WriteString(e.hwInfo.GetGPUName())
	sb.WriteString(fmt.Sprintf(" (%.1f GB VRAM)\n", e.hwInfo.GetTotalVRAM()))
	sb.WriteString("# 生成时间: 2024\n\n")
	sb.WriteString("version: '3.8'\n\n")
	sb.WriteString("services:\n")

	// Ollama 服务
	sb.WriteString("  ollama:\n")
	sb.WriteString("    image: ollama/ollama:latest\n")
	sb.WriteString("    container_name: ollama\n")
	sb.WriteString("    ports:\n")
	sb.WriteString("      - \"11434:11434\"\n")
	sb.WriteString("    volumes:\n")
	sb.WriteString("      - ollama-data:/root/.ollama\n")

	// GPU配置
	if e.hwInfo.HasCUDA {
		sb.WriteString("    deploy:\n")
		sb.WriteString("      resources:\n")
		sb.WriteString("        reservations:\n")
		sb.WriteString("          devices:\n")
		sb.WriteString("            - driver: nvidia\n")
		sb.WriteString("              count: all\n")
		sb.WriteString("              capabilities: [gpu]\n")
	}

	sb.WriteString("    environment:\n")
	sb.WriteString("      - OLLAMA_ORIGINS=*\n")
	sb.WriteString("    restart: unless-stopped\n\n")

	// Open WebUI
	sb.WriteString("  open-webui:\n")
	sb.WriteString("    image: ghcr.io/open-webui/open-webui:main\n")
	sb.WriteString("    container_name: open-webui\n")
	sb.WriteString("    ports:\n")
	sb.WriteString("      - \"3000:8080\"\n")
	sb.WriteString("    volumes:\n")
	sb.WriteString("      - open-webui-data:/app/backend/data\n")
	sb.WriteString("    environment:\n")
	sb.WriteString("      - OLLAMA_BASE_URL=http://ollama:11434\n")
	sb.WriteString("    depends_on:\n")
	sb.WriteString("      - ollama\n")
	sb.WriteString("    restart: unless-stopped\n\n")

	sb.WriteString("volumes:\n")
	sb.WriteString("  ollama-data:\n")
	sb.WriteString("  open-webui-data:\n\n")

	// 推荐模型注释
	sb.WriteString("# 🎯 推荐模型（按顺序运行）:\n")
	for i, rec := range models {
		if i >= 5 {
			break
		}
		sb.WriteString(fmt.Sprintf("# %d. %s (%s)\n", i+1, rec.Model.Name, rec.OllamaCommand))
	}

	return sb.String()
}

// exportShell 导出Shell脚本
func (e *Exporter) exportShell() string {
	var models []recommender.Recommendation
	for _, rec := range e.recommendations {
		if rec.CanRun {
			models = append(models, rec)
		}
	}

	var sb strings.Builder
	sb.WriteString("#!/bin/bash\n\n")
	sb.WriteString("# SmartModelPicker 自动生成的安装脚本\n")
	sb.WriteString(fmt.Sprintf("# 基于硬件: %s (%.1f GB VRAM)\n\n", e.hwInfo.GetGPUName(), e.hwInfo.GetTotalVRAM()))

	sb.WriteString("set -e\n\n")

	// 检查依赖
	sb.WriteString("echo \"🔍 检查依赖...\"\n")
	sb.WriteString("if ! command -v ollama &> /dev/null; then\n")
	sb.WriteString("    echo \"📦 安装 Ollama...\"\n")
	sb.WriteString("    curl -fsSL https://ollama.com/install.sh | sh\n")
	sb.WriteString("fi\n\n")

	// 拉取推荐模型
	sb.WriteString("echo \"🚀 拉取推荐模型...\"\n")
	for i, rec := range models {
		if i >= 5 {
			break
		}
		sb.WriteString(fmt.Sprintf("echo \"  [%d/%d] %s...\"\n", i+1, min(5, len(models)), rec.Model.Name))
		sb.WriteString(fmt.Sprintf("ollama pull %s:%s\n", rec.Model.ID, strings.ToLower(rec.BestQuant)))
	}

	sb.WriteString("\necho \"✅ 安装完成！\"\n")
	sb.WriteString("echo \"\"\n")
	sb.WriteString("echo \"💡 使用示例:\"\n")
	if len(models) > 0 {
		sb.WriteString(fmt.Sprintf("echo \"   ollama run %s:%s\"\n", models[0].Model.ID, strings.ToLower(models[0].BestQuant)))
	}

	return sb.String()
}

// exportJSON 导出JSON配置
func (e *Exporter) exportJSON() string {
	var models []recommender.Recommendation
	for _, rec := range e.recommendations {
		if rec.CanRun {
			models = append(models, rec)
		}
	}

	output := map[string]interface{}{
		"hardware": map[string]interface{}{
			"os":        e.hwInfo.OS,
			"arch":      e.hwInfo.Arch,
			"cpu_cores": e.hwInfo.CPUCores,
			"total_ram": e.hwInfo.CPUTotalRAM,
			"gpu":       e.hwInfo.GetGPUName(),
			"vram_gb":   e.hwInfo.GetTotalVRAM(),
		},
		"recommended_models": models,
		"commands": map[string]string{},
	}

	commands := output["commands"].(map[string]string)
	for _, rec := range models {
		commands[rec.Model.ID] = rec.OllamaCommand
	}

	data, _ := json.MarshalIndent(output, "", "  ")
	return string(data)
}
