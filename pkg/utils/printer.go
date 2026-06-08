package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/user/smartmodelpicker-cli/pkg/hardware"
	"github.com/user/smartmodelpicker-cli/pkg/models"
	"github.com/user/smartmodelpicker-cli/pkg/recommender"
)

// PrintBanner 打印欢迎横幅
func PrintBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║           🧠 SmartModelPicker-CLI v1.0.0 🚀                  ║
║     本地LLM硬件适配智能推荐引擎 | Local LLM Hardware Matcher   ║
╚══════════════════════════════════════════════════════════════╝
`
	color.Cyan(banner)
}

// PrintHardwareInfo 打印硬件信息
func PrintHardwareInfo(info *hardware.Info, asJSON bool) {
	if asJSON {
		data, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(data))
		return
	}

	bold := color.New(color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	fmt.Println()
	bold.Println("📊 硬件配置检测")
	fmt.Println(strings.Repeat("─", 50))

	fmt.Printf("   🖥️  操作系统: %s (%s)\n", info.OS, info.Arch)
	fmt.Printf("   🔲 CPU核心: %d 核\n", info.CPUCores)
	fmt.Printf("   💾 系统内存: %.1f GB\n", info.CPUTotalRAM)

	if len(info.GPUs) > 0 {
		fmt.Println("   🎮 显卡信息:")
		for i, gpu := range info.GPUs {
			if gpu.IsIntegrated {
				fmt.Printf("      GPU %d: %s (%s) - 集成显卡\n", i+1, gpu.Name, gpu.Vendor)
			} else {
				fmt.Printf("      GPU %d: %s (%s) - %.1f GB VRAM\n", i+1, gpu.Name, gpu.Vendor, gpu.VRAM)
			}
		}
		fmt.Printf("   📊 总显存: %.1f GB\n", info.TotalVRAM)
	} else {
		yellow.Println("   ⚠️  未检测到独立显卡，将使用CPU推理")
	}

	// 特性支持
	fmt.Println("   ⚡ 加速支持:")
	if info.HasCUDA {
		green.Println("      ✅ NVIDIA CUDA")
	} else {
		fmt.Println("      ❌ NVIDIA CUDA")
	}
	if info.HasROCm {
		green.Println("      ✅ AMD ROCm")
	} else {
		fmt.Println("      ❌ AMD ROCm")
	}
	if info.HasMetal {
		green.Println("      ✅ Apple Metal")
	} else {
		fmt.Println("      ❌ Apple Metal")
	}
	fmt.Printf("      %s AVX2 | %s AVX-512\n",
		map[bool]string{true: "✅", false: "❌"}[info.HasAVX2],
		map[bool]string{true: "✅", false: "❌"}[info.HasAVX512])

	green.Printf("   🎯 推荐后端: %s\n", info.Recommended)
}

// PrintModelList 打印模型列表
func PrintModelList(models []models.Model, asJSON bool) {
	if asJSON {
		data, _ := json.MarshalIndent(models, "", "  ")
		fmt.Println(string(data))
		return
	}

	bold := color.New(color.Bold)
	fmt.Println()
	bold.Printf("📚 模型数据库 (%d 个模型)\n", len(models))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("%-20s %-12s %-8s %-10s %s\n", "模型名称", "组织", "参数量", "上下文", "标签")
	fmt.Println(strings.Repeat("─", 80))

	for _, m := range models {
		tags := strings.Join(m.Tags[:min(3, len(m.Tags))], ", ")
		ctx := formatContextLength(m.ContextLength)
		fmt.Printf("%-20s %-12s %-8s %-10s %s\n",
			truncate(m.Name, 20),
			truncate(m.Organization, 12),
			m.Parameters,
			ctx,
			tags)
	}
}

// PrintRecommendations 打印推荐结果
func PrintRecommendations(recs []recommender.Recommendation, hwInfo *hardware.Info) {
	bold := color.New(color.Bold)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)

	fmt.Println()
	bold.Println("🎯 智能推荐结果")
	fmt.Println(strings.Repeat("═", 70))

	// 分类显示
	var runnable, unrunnable []recommender.Recommendation
	for _, rec := range recs {
		if rec.CanRun {
			runnable = append(runnable, rec)
		} else {
			unrunnable = append(unrunnable, rec)
		}
	}

	// 可运行推荐
	if len(runnable) > 0 {
		fmt.Println()
		green.Println("✅ 推荐运行（按匹配度排序）")
		fmt.Println(strings.Repeat("─", 70))

		for i, rec := range runnable {
			if i >= 10 { // 最多显示10个
				fmt.Printf("\n   ... 还有 %d 个可运行模型\n", len(runnable)-10)
				break
			}

			badge := "🥇"
			if i == 1 {
				badge = "🥈"
			} else if i == 2 {
				badge = "🥉"
			} else {
				badge = "  "
			}

			fmt.Printf("\n   %s %s %s\n", badge, bold.Sprint(rec.Model.Name), rec.PerformanceTier)
			fmt.Printf("      📋 %s\n", rec.Model.Description)
			fmt.Printf("      %s\n", rec.Reason)
			fmt.Printf("      🔧 最佳量化: %s | 需要显存: %.1f GB\n", rec.BestQuant, rec.RequiredVRAM)
			if rec.OllamaCommand != "" {
				fmt.Printf("      💻 一键运行: %s\n", yellow.Sprint(rec.OllamaCommand))
			}
		}
	}

	// 不可运行
	if len(unrunnable) > 0 && len(runnable) < 5 {
		fmt.Println()
		red.Println("❌ 当前硬件无法运行（显存不足）")
		fmt.Println(strings.Repeat("─", 70))

		for i, rec := range unrunnable {
			if i >= 5 {
				fmt.Printf("\n   ... 还有 %d 个不可运行模型\n", len(unrunnable)-5)
				break
			}
			fmt.Printf("   • %s - %s (需要 %.1f GB)\n",
				rec.Model.Name,
				rec.Reason,
				rec.RequiredVRAM)
		}
	}

	// 总结
	fmt.Println()
	fmt.Println(strings.Repeat("═", 70))
	green.Printf("✨ 共找到 %d 个可运行模型，%d 个需要升级硬件\n",
		len(runnable), len(unrunnable))

	if len(runnable) > 0 {
		fmt.Println()
		yellow.Println("💡 提示:")
		fmt.Println("   • 使用 -q 参数指定量化级别（如 -q Q4_0）")
		fmt.Println("   • 使用 -e docker-compose 导出部署配置")
		fmt.Println("   • 使用 --json 获取结构化输出")
	}
}

// PrintRecommendationsJSON 以JSON格式输出推荐
func PrintRecommendationsJSON(recs []recommender.Recommendation, hwInfo *hardware.Info) {
	output := struct {
		Hardware        *hardware.Info                 `json:"hardware"`
		Recommendations []recommender.Recommendation   `json:"recommendations"`
		Summary         map[string]interface{}         `json:"summary"`
	}{
		Hardware:        hwInfo,
		Recommendations: recs,
		Summary: map[string]interface{}{
			"total_models":    len(recs),
			"runnable_models": countRunnable(recs),
			"best_model":      getBestModel(recs),
		},
	}

	data, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(data))
}

// PrintHardwareInfoJSON 以JSON格式输出硬件信息
func PrintHardwareInfoJSON(info *hardware.Info) {
	data, _ := json.MarshalIndent(info, "", "  ")
	fmt.Println(string(data))
}

// 辅助函数
func formatContextLength(length int) string {
	if length >= 100000 {
		return fmt.Sprintf("%.0fK", float64(length)/1000)
	}
	return fmt.Sprintf("%d", length)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func countRunnable(recs []recommender.Recommendation) int {
	count := 0
	for _, rec := range recs {
		if rec.CanRun {
			count++
		}
	}
	return count
}

func getBestModel(recs []recommender.Recommendation) string {
	for _, rec := range recs {
		if rec.CanRun {
			return rec.Model.ID
		}
	}
	return ""
}
