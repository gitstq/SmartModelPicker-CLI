package models

import (
	"encoding/json"
	"fmt"
)

// Model 模型信息
type Model struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Organization  string   `json:"organization"`
	Description   string   `json:"description"`
	Parameters    string   `json:"parameters"`
	ParamCount    float64  `json:"param_count_billion"`
	ContextLength int      `json:"context_length"`
	VRAMByQuant   map[string]float64 `json:"vram_by_quantization_gb"`
	Tags          []string `json:"tags"`
	License       string   `json:"license"`
	URL           string   `json:"url"`
	Backend       []string `json:"supported_backends"`
}

// Registry 模型注册表
type Registry struct {
	models []Model
}

// NewRegistry 创建模型注册表
func NewRegistry() *Registry {
	r := &Registry{}
	r.loadBuiltinModels()
	return r
}

// GetAllModels 获取所有模型
func (r *Registry) GetAllModels() []Model {
	return r.models
}

// GetModelByID 通过ID获取模型
func (r *Registry) GetModelByID(id string) (*Model, error) {
	for i := range r.models {
		if r.models[i].ID == id {
			return &r.models[i], nil
		}
	}
	return nil, fmt.Errorf("模型 %s 未找到", id)
}

// loadBuiltinModels 加载内置模型数据库
func (r *Registry) loadBuiltinModels() {
	r.models = []Model{
		// Meta Llama 系列
		{
			ID:           "llama3.1-8b",
			Name:         "Llama 3.1 8B",
			Organization: "Meta",
			Description:  "Meta最新开源大模型，8B参数，性能优异，支持128K上下文",
			Parameters:   "8B",
			ParamCount:   8.0,
			ContextLength: 128000,
			VRAMByQuant: map[string]float64{
				"Q4_0":  4.7,
				"Q4_K_M": 5.1,
				"Q5_K_M": 5.7,
				"Q6_K":  6.6,
				"Q8_0":  8.5,
				"FP16":  16.0,
			},
			Tags:     []string{"chat", "general", "multilingual"},
			License:  "Llama 3.1 License",
			URL:      "https://huggingface.co/meta-llama/Meta-Llama-3.1-8B",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
		{
			ID:           "llama3.1-70b",
			Name:         "Llama 3.1 70B",
			Organization: "Meta",
			Description:  "Meta大型开源模型，70B参数，接近GPT-4级别性能",
			Parameters:   "70B",
			ParamCount:   70.0,
			ContextLength: 128000,
			VRAMByQuant: map[string]float64{
				"Q4_0":  40.0,
				"Q4_K_M": 42.0,
				"Q5_K_M": 47.0,
				"Q6_K":  54.0,
				"Q8_0":  70.0,
				"FP16":  140.0,
			},
			Tags:     []string{"chat", "general", "multilingual", "advanced"},
			License:  "Llama 3.1 License",
			URL:      "https://huggingface.co/meta-llama/Meta-Llama-3.1-70B",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
		{
			ID:           "llama3.2-3b",
			Name:         "Llama 3.2 3B",
			Organization: "Meta",
			Description:  "Meta轻量级模型，3B参数，适合边缘设备和移动端",
			Parameters:   "3B",
			ParamCount:   3.0,
			ContextLength: 128000,
			VRAMByQuant: map[string]float64{
				"Q4_0":  1.9,
				"Q4_K_M": 2.0,
				"Q5_K_M": 2.3,
				"Q6_K":  2.6,
				"Q8_0":  3.4,
				"FP16":  6.0,
			},
			Tags:     []string{"chat", "edge", "mobile", "fast"},
			License:  "Llama 3.2 License",
			URL:      "https://huggingface.co/meta-llama/Llama-3.2-3B",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
		{
			ID:           "llama3.2-1b",
			Name:         "Llama 3.2 1B",
			Organization: "Meta",
			Description:  "Meta超轻量模型，1B参数，极低资源占用",
			Parameters:   "1B",
			ParamCount:   1.0,
			ContextLength: 128000,
			VRAMByQuant: map[string]float64{
				"Q4_0":  0.7,
				"Q4_K_M": 0.8,
				"Q5_K_M": 0.9,
				"Q6_K":  1.0,
				"Q8_0":  1.3,
				"FP16":  2.2,
			},
			Tags:     []string{"chat", "edge", "ultra-light", "fast"},
			License:  "Llama 3.2 License",
			URL:      "https://huggingface.co/meta-llama/Llama-3.2-1B",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},

		// Mistral 系列
		{
			ID:           "mistral-7b",
			Name:         "Mistral 7B",
			Organization: "Mistral AI",
			Description:  "Mistral AI高质量7B模型，性能超越Llama 2 13B",
			Parameters:   "7B",
			ParamCount:   7.0,
			ContextLength: 32768,
			VRAMByQuant: map[string]float64{
				"Q4_0":  4.1,
				"Q4_K_M": 4.4,
				"Q5_K_M": 5.0,
				"Q6_K":  5.7,
				"Q8_0":  7.3,
				"FP16":  14.0,
			},
			Tags:     []string{"chat", "general", "efficient"},
			License:  "Apache 2.0",
			URL:      "https://huggingface.co/mistralai/Mistral-7B-v0.3",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
		{
			ID:           "mixtral-8x7b",
			Name:         "Mixtral 8x7B",
			Organization: "Mistral AI",
			Description:  "Mistral MoE模型，47B有效参数，推理时仅使用13B活跃参数",
			Parameters:   "47B (13B active)",
			ParamCount:   47.0,
			ContextLength: 32768,
			VRAMByQuant: map[string]float64{
				"Q4_0":  26.0,
				"Q4_K_M": 28.0,
				"Q5_K_M": 31.0,
				"Q6_K":  36.0,
				"Q8_0":  47.0,
				"FP16":  94.0,
			},
			Tags:     []string{"chat", "moe", "advanced", "efficient"},
			License:  "Apache 2.0",
			URL:      "https://huggingface.co/mistralai/Mixtral-8x7B-v0.1",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
		{
			ID:           "mistral-nemo",
			Name:         "Mistral Nemo 12B",
			Organization: "Mistral AI / NVIDIA",
			Description:  "Mistral与NVIDIA合作模型，12B参数，128K上下文",
			Parameters:   "12B",
			ParamCount:   12.0,
			ContextLength: 128000,
			VRAMByQuant: map[string]float64{
				"Q4_0":  7.0,
				"Q4_K_M": 7.5,
				"Q5_K_M": 8.4,
				"Q6_K":  9.7,
				"Q8_0":  12.4,
				"FP16":  24.0,
			},
			Tags:     []string{"chat", "multilingual", "long-context"},
			License:  "Apache 2.0",
			URL:      "https://huggingface.co/mistralai/Mistral-Nemo-Instruct-2407",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},

		// Qwen 系列
		{
			ID:           "qwen2.5-7b",
			Name:         "Qwen 2.5 7B",
			Organization: "Alibaba",
			Description:  "通义千问2.5，7B参数，中文能力突出，支持128K上下文",
			Parameters:   "7B",
			ParamCount:   7.0,
			ContextLength: 128000,
			VRAMByQuant: map[string]float64{
				"Q4_0":  4.1,
				"Q4_K_M": 4.4,
				"Q5_K_M": 5.0,
				"Q6_K":  5.7,
				"Q8_0":  7.3,
				"FP16":  14.0,
			},
			Tags:     []string{"chat", "chinese", "multilingual", "coding"},
			License:  "Apache 2.0",
			URL:      "https://huggingface.co/Qwen/Qwen2.5-7B-Instruct",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
		{
			ID:           "qwen2.5-14b",
			Name:         "Qwen 2.5 14B",
			Organization: "Alibaba",
			Description:  "通义千问2.5，14B参数，更强的推理和编码能力",
			Parameters:   "14B",
			ParamCount:   14.0,
			ContextLength: 128000,
			VRAMByQuant: map[string]float64{
				"Q4_0":  8.2,
				"Q4_K_M": 8.8,
				"Q5_K_M": 9.9,
				"Q6_K":  11.4,
				"Q8_0":  14.6,
				"FP16":  28.0,
			},
			Tags:     []string{"chat", "chinese", "coding", "advanced"},
			License:  "Apache 2.0",
			URL:      "https://huggingface.co/Qwen/Qwen2.5-14B-Instruct",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
		{
			ID:           "qwen2.5-32b",
			Name:         "Qwen 2.5 32B",
			Organization: "Alibaba",
			Description:  "通义千问2.5，32B参数，接近70B模型的性能",
			Parameters:   "32B",
			ParamCount:   32.0,
			ContextLength: 128000,
			VRAMByQuant: map[string]float64{
				"Q4_0":  18.5,
				"Q4_K_M": 19.8,
				"Q5_K_M": 22.3,
				"Q6_K":  25.7,
				"Q8_0":  32.9,
				"FP16":  64.0,
			},
			Tags:     []string{"chat", "chinese", "coding", "advanced"},
			License:  "Apache 2.0",
			URL:      "https://huggingface.co/Qwen/Qwen2.5-32B-Instruct",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
		{
			ID:           "qwen2.5-coder-7b",
			Name:         "Qwen 2.5 Coder 7B",
			Organization: "Alibaba",
			Description:  "通义千问代码专用模型，7B参数，代码生成能力优异",
			Parameters:   "7B",
			ParamCount:   7.0,
			ContextLength: 128000,
			VRAMByQuant: map[string]float64{
				"Q4_0":  4.1,
				"Q4_K_M": 4.4,
				"Q5_K_M": 5.0,
				"Q6_K":  5.7,
				"Q8_0":  7.3,
				"FP16":  14.0,
			},
			Tags:     []string{"coding", "chinese", "specialized"},
			License:  "Apache 2.0",
			URL:      "https://huggingface.co/Qwen/Qwen2.5-Coder-7B-Instruct",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},

		// DeepSeek 系列
		{
			ID:           "deepseek-coder-v2",
			Name:         "DeepSeek Coder V2",
			Organization: "DeepSeek",
			Description:  "DeepSeek代码模型，16B有效参数，支持338种编程语言",
			Parameters:   "16B (236B total)",
			ParamCount:   16.0,
			ContextLength: 128000,
			VRAMByQuant: map[string]float64{
				"Q4_0":  9.0,
				"Q4_K_M": 9.6,
				"Q5_K_M": 10.8,
				"Q6_K":  12.4,
				"Q8_0":  15.9,
				"FP16":  30.0,
			},
			Tags:     []string{"coding", "multilingual", "specialized"},
			License:  "DeepSeek License",
			URL:      "https://huggingface.co/deepseek-ai/DeepSeek-Coder-V2-Lite-Instruct",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
		{
			ID:           "deepseek-v2.5",
			Name:         "DeepSeek V2.5",
			Organization: "DeepSeek",
			Description:  "DeepSeek通用对话模型，21B有效参数，MoE架构",
			Parameters:   "21B (236B total)",
			ParamCount:   21.0,
			ContextLength: 128000,
			VRAMByQuant: map[string]float64{
				"Q4_0":  12.0,
				"Q4_K_M": 12.8,
				"Q5_K_M": 14.4,
				"Q6_K":  16.5,
				"Q8_0":  21.1,
				"FP16":  40.0,
			},
			Tags:     []string{"chat", "coding", "advanced"},
			License:  "DeepSeek License",
			URL:      "https://huggingface.co/deepseek-ai/DeepSeek-V2.5",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},

		// Google Gemma 系列
		{
			ID:           "gemma2-9b",
			Name:         "Gemma 2 9B",
			Organization: "Google",
			Description:  "Google Gemma 2，9B参数，知识截止2024年",
			Parameters:   "9B",
			ParamCount:   9.0,
			ContextLength: 8192,
			VRAMByQuant: map[string]float64{
				"Q4_0":  5.3,
				"Q4_K_M": 5.7,
				"Q5_K_M": 6.4,
				"Q6_K":  7.3,
				"Q8_0":  9.4,
				"FP16":  18.0,
			},
			Tags:     []string{"chat", "general", "google"},
			License:  "Gemma License",
			URL:      "https://huggingface.co/google/gemma-2-9b-it",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
		{
			ID:           "gemma2-2b",
			Name:         "Gemma 2 2B",
			Organization: "Google",
			Description:  "Google Gemma 2轻量版，2B参数，适合边缘设备",
			Parameters:   "2B",
			ParamCount:   2.0,
			ContextLength: 8192,
			VRAMByQuant: map[string]float64{
				"Q4_0":  1.3,
				"Q4_K_M": 1.4,
				"Q5_K_M": 1.5,
				"Q6_K":  1.7,
				"Q8_0":  2.2,
				"FP16":  4.2,
			},
			Tags:     []string{"chat", "edge", "ultra-light"},
			License:  "Gemma License",
			URL:      "https://huggingface.co/google/gemma-2-2b-it",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},

		// Microsoft Phi 系列
		{
			ID:           "phi3-mini",
			Name:         "Phi-3 Mini 3.8B",
			Organization: "Microsoft",
			Description:  "Microsoft小型模型，3.8B参数，性能媲美大模型",
			Parameters:   "3.8B",
			ParamCount:   3.8,
			ContextLength: 128000,
			VRAMByQuant: map[string]float64{
				"Q4_0":  2.3,
				"Q4_K_M": 2.5,
				"Q5_K_M": 2.8,
				"Q6_K":  3.2,
				"Q8_0":  4.1,
				"FP16":  7.6,
			},
			Tags:     []string{"chat", "efficient", "microsoft"},
			License:  "MIT",
			URL:      "https://huggingface.co/microsoft/Phi-3-mini-4k-instruct",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
		{
			ID:           "phi3-medium",
			Name:         "Phi-3 Medium 14B",
			Organization: "Microsoft",
			Description:  "Microsoft中型模型，14B参数，平衡性能与资源",
			Parameters:   "14B",
			ParamCount:   14.0,
			ContextLength: 128000,
			VRAMByQuant: map[string]float64{
				"Q4_0":  8.2,
				"Q4_K_M": 8.8,
				"Q5_K_M": 9.9,
				"Q6_K":  11.4,
				"Q8_0":  14.6,
				"FP16":  28.0,
			},
			Tags:     []string{"chat", "general", "microsoft"},
			License:  "MIT",
			URL:      "https://huggingface.co/microsoft/Phi-3-medium-4k-instruct",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},

		// Yi 系列
		{
			ID:           "yi-1.5-9b",
			Name:         "Yi 1.5 9B",
			Organization: "01.AI",
			Description:  "零一万物Yi 1.5，9B参数，中文和英文能力均衡",
			Parameters:   "9B",
			ParamCount:   9.0,
			ContextLength: 32768,
			VRAMByQuant: map[string]float64{
				"Q4_0":  5.3,
				"Q4_K_M": 5.7,
				"Q5_K_M": 6.4,
				"Q6_K":  7.3,
				"Q8_0":  9.4,
				"FP16":  18.0,
			},
			Tags:     []string{"chat", "chinese", "general"},
			License:  "Apache 2.0",
			URL:      "https://huggingface.co/01-ai/Yi-1.5-9B-Chat",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
		{
			ID:           "yi-1.5-34b",
			Name:         "Yi 1.5 34B",
			Organization: "01.AI",
			Description:  "零一万物Yi 1.5，34B参数，更强的推理能力",
			Parameters:   "34B",
			ParamCount:   34.0,
			ContextLength: 32768,
			VRAMByQuant: map[string]float64{
				"Q4_0":  19.5,
				"Q4_K_M": 20.9,
				"Q5_K_M": 23.5,
				"Q6_K":  27.0,
				"Q8_0":  34.6,
				"FP16":  68.0,
			},
			Tags:     []string{"chat", "chinese", "advanced"},
			License:  "Apache 2.0",
			URL:      "https://huggingface.co/01-ai/Yi-1.5-34B-Chat",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},

		// Command R 系列
		{
			ID:           "command-r",
			Name:         "Cohere Command R",
			Organization: "Cohere",
			Description:  "Cohere Command R，35B参数，128K上下文，RAG优化",
			Parameters:   "35B",
			ParamCount:   35.0,
			ContextLength: 128000,
			VRAMByQuant: map[string]float64{
				"Q4_0":  20.0,
				"Q4_K_M": 21.4,
				"Q5_K_M": 24.1,
				"Q6_K":  27.7,
				"Q8_0":  35.5,
				"FP16":  70.0,
			},
			Tags:     []string{"chat", "rag", "long-context"},
			License:  "CC BY-NC 4.0",
			URL:      "https://huggingface.co/CohereForAI/c4ai-command-r-v01",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},

		// 专用模型
		{
			ID:           "codellama-7b",
			Name:         "CodeLlama 7B",
			Organization: "Meta",
			Description:  "Meta代码专用模型，7B参数，支持多种编程语言",
			Parameters:   "7B",
			ParamCount:   7.0,
			ContextLength: 16384,
			VRAMByQuant: map[string]float64{
				"Q4_0":  4.1,
				"Q4_K_M": 4.4,
				"Q5_K_M": 5.0,
				"Q6_K":  5.7,
				"Q8_0":  7.3,
				"FP16":  14.0,
			},
			Tags:     []string{"coding", "specialized"},
			License:  "Llama 2 License",
			URL:      "https://huggingface.co/codellama/CodeLlama-7b-Instruct-hf",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
		{
			ID:           "starcoder2-7b",
			Name:         "StarCoder2 7B",
			Organization: "BigCode",
			Description:  "代码生成模型，7B参数，支持600+编程语言",
			Parameters:   "7B",
			ParamCount:   7.0,
			ContextLength: 16384,
			VRAMByQuant: map[string]float64{
				"Q4_0":  4.1,
				"Q4_K_M": 4.4,
				"Q5_K_M": 5.0,
				"Q6_K":  5.7,
				"Q8_0":  7.3,
				"FP16":  14.0,
			},
			Tags:     []string{"coding", "fill-in-middle", "specialized"},
			License:  "BigCode OpenRAIL-M v1",
			URL:      "https://huggingface.co/bigcode/starcoder2-7b",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
		{
			ID:           "nomic-embed-text",
			Name:         "Nomic Embed Text v1.5",
			Organization: "Nomic AI",
			Description:  "文本嵌入模型，137M参数，适合RAG和语义搜索",
			Parameters:   "137M",
			ParamCount:   0.137,
			ContextLength: 8192,
			VRAMByQuant: map[string]float64{
				"Q4_0":  0.1,
				"Q4_K_M": 0.1,
				"Q5_K_M": 0.1,
				"Q6_K":  0.1,
				"Q8_0":  0.2,
				"FP16":  0.3,
			},
			Tags:     []string{"embedding", "rag", "ultra-light"},
			License:  "Apache 2.0",
			URL:      "https://huggingface.co/nomic-ai/nomic-embed-text-v1.5",
			Backend:  []string{"ollama", "lmstudio", "llamacpp"},
		},
	}
}

// ToJSON 导出为JSON
func (r *Registry) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r.models, "", "  ")
}
