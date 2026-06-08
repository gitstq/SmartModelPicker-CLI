package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/user/smartmodelpicker-cli/pkg/hardware"
	"github.com/user/smartmodelpicker-cli/pkg/models"
	"github.com/user/smartmodelpicker-cli/pkg/recommender"
	"github.com/user/smartmodelpicker-cli/pkg/utils"
)

var (
	version   = "1.0.0"
	buildTime = "unknown"
	gitCommit = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "smartmodelpicker",
	Short: "🧠 SmartModelPicker - 本地LLM硬件适配智能推荐引擎",
	Long: `SmartModelPicker-CLI 🚀
智能检测你的硬件配置，推荐最适合本地运行的开源大语言模型。

支持 Ollama / LM Studio / llama.cpp 多后端，
覆盖 NVIDIA / AMD / Apple Silicon / Intel 全平台硬件。`,
	RunE: runMain,
}

var (
	flagBackend    string
	flagQuantize   string
	flagMinVRAM    float64
	flagMaxVRAM    float64
	flagShowAll    bool
	flagExport     string
	flagJSON       bool
	flagVerbose    bool
)

func init() {
	rootCmd.Flags().StringVarP(&flagBackend, "backend", "b", "auto", "指定LLM后端 (auto|ollama|lmstudio|llamacpp)")
	rootCmd.Flags().StringVarP(&flagQuantize, "quantize", "q", "auto", "量化级别 (auto|Q4_0|Q4_K_M|Q5_K_M|Q6_K|Q8_0|FP16)")
	rootCmd.Flags().Float64Var(&flagMinVRAM, "min-vram", 0, "最小显存要求(GB)")
	rootCmd.Flags().Float64Var(&flagMaxVRAM, "max-vram", 0, "最大显存限制(GB)")
	rootCmd.Flags().BoolVarP(&flagShowAll, "all", "a", false, "显示所有模型（不限于可运行）")
	rootCmd.Flags().StringVarP(&flagExport, "export", "e", "", "导出配置 (docker-compose|shell|json)")
	rootCmd.Flags().BoolVar(&flagJSON, "json", false, "以JSON格式输出")
	rootCmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "详细输出模式")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listModelsCmd)
	rootCmd.AddCommand(checkCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("SmartModelPicker v%s\n", version)
		fmt.Printf("  Git Commit: %s\n", gitCommit)
		fmt.Printf("  Build Time: %s\n", buildTime)
		fmt.Printf("  Go Version: %s\n", runtime.Version())
		fmt.Printf("  OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

var listModelsCmd = &cobra.Command{
	Use:   "list",
	Short: "列出支持的模型数据库",
	RunE: func(cmd *cobra.Command, args []string) error {
		registry := models.NewRegistry()
		utils.PrintModelList(registry.GetAllModels(), flagJSON)
		return nil
	},
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "仅检测硬件信息",
	RunE: func(cmd *cobra.Command, args []string) error {
		detector := hardware.NewDetector()
		info, err := detector.Detect()
		if err != nil {
			return fmt.Errorf("硬件检测失败: %w", err)
		}
		utils.PrintHardwareInfo(info, flagJSON)
		return nil
	},
}

func runMain(cmd *cobra.Command, args []string) error {
	utils.PrintBanner()

	// 1. 检测硬件
	fmt.Println("🔍 正在检测硬件配置...")
	detector := hardware.NewDetector()
	hwInfo, err := detector.Detect()
	if err != nil {
		return fmt.Errorf("硬件检测失败: %w", err)
	}
	utils.PrintHardwareInfo(hwInfo, false)

	// 2. 加载模型数据库
	fmt.Println("\n📚 加载模型数据库...")
	registry := models.NewRegistry()
	allModels := registry.GetAllModels()
	fmt.Printf("   已加载 %d 个模型\n", len(allModels))

	// 3. 智能推荐
	fmt.Println("\n🎯 正在生成推荐...")
	engine := recommender.NewEngine(hwInfo, &recommender.Options{
		Backend:  flagBackend,
		Quantize: flagQuantize,
		MinVRAM:  flagMinVRAM,
		MaxVRAM:  flagMaxVRAM,
		ShowAll:  flagShowAll,
	})

	recommendations, err := engine.Recommend(allModels)
	if err != nil {
		return fmt.Errorf("推荐生成失败: %w", err)
	}

	// 4. 输出结果
	if flagJSON {
		utils.PrintRecommendationsJSON(recommendations, hwInfo)
	} else {
		utils.PrintRecommendations(recommendations, hwInfo)
	}

	// 5. 导出配置
	if flagExport != "" {
		fmt.Printf("\n📦 正在导出 %s 配置...\n", flagExport)
		exporter := utils.NewExporter(hwInfo, recommendations)
		content, err := exporter.Export(flagExport)
		if err != nil {
			return fmt.Errorf("导出失败: %w", err)
		}
		filename := fmt.Sprintf("smartmodelpicker-%s.%s", flagExport, getExtension(flagExport))
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			return fmt.Errorf("写入文件失败: %w", err)
		}
		fmt.Printf("   ✅ 已保存至 %s\n", filename)
	}

	fmt.Println("\n✨ 完成！使用 --help 查看更多选项")
	return nil
}

func getExtension(format string) string {
	switch format {
	case "docker-compose":
		return "yaml"
	case "shell":
		return "sh"
	case "json":
		return "json"
	default:
		return "txt"
	}
}
