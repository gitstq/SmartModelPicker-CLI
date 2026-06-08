package recommender

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/user/smartmodelpicker-cli/pkg/hardware"
	"github.com/user/smartmodelpicker-cli/pkg/models"
)

// Options 推荐选项
type Options struct {
	Backend  string
	Quantize string
	MinVRAM  float64
	MaxVRAM  float64
	ShowAll  bool
}

// Recommendation 推荐结果
type Recommendation struct {
	Model           models.Model `json:"model"`
	BestQuant       string       `json:"best_quantization"`
	RequiredVRAM    float64      `json:"required_vram_gb"`
	CanRun          bool         `json:"can_run"`
	Score           float64      `json:"score"`
	Reason          string       `json:"reason"`
	PerformanceTier string       `json:"performance_tier"`
	OllamaCommand   string       `json:"ollama_command,omitempty"`
}

// Engine 推荐引擎
type Engine struct {
	hwInfo  *hardware.Info
	options *Options
}

// NewEngine 创建推荐引擎
func NewEngine(hwInfo *hardware.Info, options *Options) *Engine {
	return &Engine{
		hwInfo:  hwInfo,
		options: options,
	}
}

// Recommend 生成推荐
func (e *Engine) Recommend(allModels []models.Model) ([]Recommendation, error) {
	var recommendations []Recommendation

	// 计算可用VRAM
	availableVRAM := e.hwInfo.GetTotalVRAM()
	if e.options.MaxVRAM > 0 && e.options.MaxVRAM < availableVRAM {
		availableVRAM = e.options.MaxVRAM
	}

	for _, model := range allModels {
		rec := e.evaluateModel(model, availableVRAM)
		if e.options.ShowAll || rec.CanRun {
			recommendations = append(recommendations, rec)
		}
	}

	// 按评分排序
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	return recommendations, nil
}

// evaluateModel 评估单个模型
func (e *Engine) evaluateModel(model models.Model, availableVRAM float64) Recommendation {
	rec := Recommendation{
		Model: model,
		CanRun: false,
		Score:  0,
	}

	// 确定最佳量化级别
	bestQuant, requiredVRAM := e.findBestQuantization(model, availableVRAM)
	rec.BestQuant = bestQuant
	rec.RequiredVRAM = requiredVRAM

	// 检查是否可以运行
	if requiredVRAM > 0 && requiredVRAM <= availableVRAM {
		rec.CanRun = true
	}

	// 如果指定了最小VRAM，检查是否满足
	if e.options.MinVRAM > 0 && requiredVRAM < e.options.MinVRAM {
		rec.CanRun = false
	}

	// 计算评分
	rec.Score = e.calculateScore(model, rec.CanRun, requiredVRAM, availableVRAM)
	rec.PerformanceTier = e.getPerformanceTier(model)
	rec.Reason = e.generateReason(model, rec.CanRun, bestQuant, requiredVRAM, availableVRAM)

	// 生成Ollama命令
	if rec.CanRun {
		rec.OllamaCommand = fmt.Sprintf("ollama run %s:%s", model.ID, strings.ToLower(bestQuant))
	}

	return rec
}

// findBestQuantization 找到最佳量化级别
func (e *Engine) findBestQuantization(model models.Model, availableVRAM float64) (string, float64) {
	// 量化级别优先级（质量从高到低）
	quantPriority := []string{"FP16", "Q8_0", "Q6_K", "Q5_K_M", "Q4_K_M", "Q4_0"}

	// 如果指定了量化级别
	if e.options.Quantize != "auto" {
		if vram, ok := model.VRAMByQuant[e.options.Quantize]; ok {
			return e.options.Quantize, vram
		}
	}

	// 自动选择最佳量化
	var bestQuant string
	var bestVRAM float64

	for _, quant := range quantPriority {
		if vram, ok := model.VRAMByQuant[quant]; ok {
			// 选择质量最高且能运行的
			if vram <= availableVRAM {
				return quant, vram
			}
			// 记录最低要求
			if bestVRAM == 0 || vram < bestVRAM {
				bestQuant = quant
				bestVRAM = vram
			}
		}
	}

	return bestQuant, bestVRAM
}

// calculateScore 计算模型评分
func (e *Engine) calculateScore(model models.Model, canRun bool, requiredVRAM, availableVRAM float64) float64 {
	if !canRun {
		// 不能运行的模型给低分，但按接近程度排序
		gap := requiredVRAM - availableVRAM
		return math.Max(0, 10.0-gap*5)
	}

	score := 50.0 // 基础分

	// 参数规模加分（适中偏好）
	switch {
	case model.ParamCount >= 30:
		score += 25 // 大型模型
	case model.ParamCount >= 14:
		score += 20 // 中型模型
	case model.ParamCount >= 7:
		score += 15 // 标准模型
	case model.ParamCount >= 3:
		score += 10 // 小型模型
	default:
		score += 5  // 微型模型
	}

	// 上下文长度加分
	if model.ContextLength >= 128000 {
		score += 10
	} else if model.ContextLength >= 32768 {
		score += 5
	}

	// VRAM利用率加分（偏好高效利用）
	utilization := requiredVRAM / availableVRAM
	if utilization >= 0.5 && utilization <= 0.9 {
		score += 10 // 理想利用率
	} else if utilization > 0.9 {
		score += 5  // 接近上限
	} else {
		score += utilization * 10
	}

	// 量化质量加分
	switch model.VRAMByQuant {
	case nil:
		// 不做处理
	default:
		if _, ok := model.VRAMByQuant["FP16"]; ok && requiredVRAM == model.VRAMByQuant["FP16"] {
			score += 5 // FP16质量最高
		} else if _, ok := model.VRAMByQuant["Q8_0"]; ok && requiredVRAM == model.VRAMByQuant["Q8_0"] {
			score += 3 // Q8_0质量很好
		}
	}

	// 开源协议加分
	if model.License == "Apache 2.0" || model.License == "MIT" {
		score += 3
	}

	return score
}

// getPerformanceTier 获取性能等级
func (e *Engine) getPerformanceTier(model models.Model) string {
	switch {
	case model.ParamCount >= 30:
		return "🏆 旗舰级"
	case model.ParamCount >= 14:
		return "🥇 高端"
	case model.ParamCount >= 7:
		return "🥈 主流"
	case model.ParamCount >= 3:
		return "🥉 入门"
	default:
		return "💎 轻量"
	}
}

// generateReason 生成推荐理由
func (e *Engine) generateReason(model models.Model, canRun bool, quant string, requiredVRAM, availableVRAM float64) string {
	if !canRun {
		gap := requiredVRAM - availableVRAM
		return fmt.Sprintf("❌ 显存不足，需要 %.1fGB，缺口 %.1fGB", requiredVRAM, gap)
	}

	reasons := []string{}

	// 量化说明
	if quant == "FP16" {
		reasons = append(reasons, "🎯 无损精度")
	} else if quant == "Q8_0" {
		reasons = append(reasons, "✨ 接近无损")
	} else if quant == "Q4_K_M" || quant == "Q5_K_M" {
		reasons = append(reasons, "⚖️ 平衡质量与速度")
	} else if quant == "Q4_0" {
		reasons = append(reasons, "🚀 极速推理")
	}

	// 显存利用
	utilization := requiredVRAM / availableVRAM * 100
	if utilization >= 80 {
		reasons = append(reasons, fmt.Sprintf("💾 高显存利用 (%.0f%%)", utilization))
	} else {
		reasons = append(reasons, fmt.Sprintf("💾 显存占用 %.1fGB (%.0f%%)", requiredVRAM, utilization))
	}

	// 特性标签
	if contains(model.Tags, "chinese") {
		reasons = append(reasons, "🇨🇳 中文优化")
	}
	if contains(model.Tags, "coding") {
		reasons = append(reasons, "💻 代码专长")
	}
	if model.ContextLength >= 128000 {
		reasons = append(reasons, "📜 超长上下文")
	}

	return strings.Join(reasons, " | ")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
