package hardware

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// Info 硬件信息结构
type Info struct {
	OS           string   `json:"os"`
	Arch         string   `json:"arch"`
	CPUCores     int      `json:"cpu_cores"`
	CPUTotalRAM  float64  `json:"cpu_total_ram_gb"`
	GPUs         []GPU    `json:"gpus"`
	TotalVRAM    float64  `json:"total_vram_gb"`
	HasCUDA      bool     `json:"has_cuda"`
	HasROCm      bool     `json:"has_rocm"`
	HasMetal     bool     `json:"has_metal"`
	HasAVX2      bool     `json:"has_avx2"`
	HasAVX512    bool     `json:"has_avx512"`
	Recommended  string   `json:"recommended_backend"`
}

// GPU 显卡信息
type GPU struct {
	Name       string  `json:"name"`
	Vendor     string  `json:"vendor"`
	VRAM       float64 `json:"vram_gb"`
	IsIntegrated bool  `json:"is_integrated"`
}

// Detector 硬件检测器
type Detector struct{}

// NewDetector 创建新的检测器
func NewDetector() *Detector {
	return &Detector{}
}

// Detect 执行硬件检测
func (d *Detector) Detect() (*Info, error) {
	info := &Info{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
		GPUs: []GPU{},
	}

	// 检测CPU核心数
	info.CPUCores = runtime.NumCPU()

	// 检测内存
	info.CPUTotalRAM = d.detectRAM()

	// 检测GPU
	info.GPUs = d.detectGPUs()
	for _, gpu := range info.GPUs {
		info.TotalVRAM += gpu.VRAM
	}

	// 检测指令集
	info.HasAVX2, info.HasAVX512 = d.detectCPUFeatures()

	// 检测后端支持
	info.HasCUDA = d.hasCUDA()
	info.HasROCm = d.hasROCm()
	info.HasMetal = runtime.GOOS == "darwin" && (runtime.GOARCH == "arm64" || d.hasMetal())

	// 推荐后端
	info.Recommended = d.recommendBackend(info)

	return info, nil
}

// detectRAM 检测系统内存
func (d *Detector) detectRAM() float64 {
	switch runtime.GOOS {
	case "linux":
		return d.detectRAMLinux()
	case "darwin":
		return d.detectRAMMacOS()
	case "windows":
		return d.detectRAMWindows()
	default:
		return 16.0 // 默认值
	}
}

func (d *Detector) detectRAMLinux() float64 {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 16.0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				kb, _ := strconv.ParseFloat(fields[1], 64)
				return kb / 1024 / 1024 // 转换为GB
			}
		}
	}
	return 16.0
}

func (d *Detector) detectRAMMacOS() float64 {
	cmd := exec.Command("sysctl", "-n", "hw.memsize")
	out, err := cmd.Output()
	if err != nil {
		return 16.0
	}
	bytes, _ := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	return bytes / 1024 / 1024 / 1024
}

func (d *Detector) detectRAMWindows() float64 {
	cmd := exec.Command("wmic", "computersystem", "get", "TotalPhysicalMemory", "/value")
	out, err := cmd.Output()
	if err != nil {
		return 16.0
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "TotalPhysicalMemory") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				bytes, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				return bytes / 1024 / 1024 / 1024
			}
		}
	}
	return 16.0
}

// detectGPUs 检测GPU信息
func (d *Detector) detectGPUs() []GPU {
	switch runtime.GOOS {
	case "linux":
		return d.detectGPUsLinux()
	case "darwin":
		return d.detectGPUsMacOS()
	case "windows":
		return d.detectGPUsWindows()
	default:
		return []GPU{}
	}
}

func (d *Detector) detectGPUsLinux() []GPU {
	gpus := []GPU{}

	// 尝试nvidia-smi
	cmd := exec.Command("nvidia-smi", "--query-gpu=name,memory.total", "--format=csv,noheader,nounits")
	out, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range lines {
			parts := strings.Split(line, ", ")
			if len(parts) >= 2 {
				name := strings.TrimSpace(parts[0])
				vram, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				gpus = append(gpus, GPU{
					Name:       name,
					Vendor:     "NVIDIA",
					VRAM:       vram / 1024,
					IsIntegrated: strings.Contains(strings.ToLower(name), "gtx") == false &&
						strings.Contains(strings.ToLower(name), "rtx") == false,
				})
			}
		}
	}

	// 尝试rocm-smi (AMD)
	if len(gpus) == 0 {
		cmd = exec.Command("rocm-smi", "--showproductname", "--showmeminfo", "vram")
		out, err = cmd.Output()
		if err == nil {
			// 简化处理，实际解析更复杂
			gpus = append(gpus, GPU{
				Name:   "AMD GPU (ROCm)",
				Vendor: "AMD",
				VRAM:   16.0,
			})
		}
	}

	// 检测Intel集成显卡
	if len(gpus) == 0 {
		cmd = exec.Command("lspci")
		out, err = cmd.Output()
		if err == nil && strings.Contains(string(out), "VGA") {
			// 有集成显卡
			gpus = append(gpus, GPU{
				Name:         "Integrated Graphics",
				Vendor:       "Intel/AMD",
				VRAM:         0,
				IsIntegrated: true,
			})
		}
	}

	return gpus
}

func (d *Detector) detectGPUsMacOS() []GPU {
	gpus := []GPU{}

	if runtime.GOARCH == "arm64" {
		// Apple Silicon
		cmd := exec.Command("system_profiler", "SPDisplaysDataType")
		out, err := cmd.Output()
		if err == nil {
			output := string(out)
			var vram float64 = 8.0 // 默认共享内存

			if strings.Contains(output, "M3 Max") || strings.Contains(output, "M2 Max") {
				vram = 32.0
			} else if strings.Contains(output, "M3 Pro") || strings.Contains(output, "M2 Pro") {
				vram = 18.0
			} else if strings.Contains(output, "M3") || strings.Contains(output, "M2") {
				vram = 8.0
			} else if strings.Contains(output, "M1 Max") {
				vram = 32.0
			} else if strings.Contains(output, "M1 Pro") {
				vram = 16.0
			} else if strings.Contains(output, "M1") {
				vram = 8.0
			}

			gpus = append(gpus, GPU{
				Name:         "Apple Silicon",
				Vendor:       "Apple",
				VRAM:         vram,
				IsIntegrated: true,
			})
		}
	} else {
		// Intel Mac
		gpus = append(gpus, GPU{
			Name:         "Intel Integrated",
			Vendor:       "Intel",
			VRAM:         0,
			IsIntegrated: true,
		})
	}

	return gpus
}

func (d *Detector) detectGPUsWindows() []GPU {
	gpus := []GPU{}

	cmd := exec.Command("wmic", "path", "win32_VideoController", "get", "Name,AdapterRAM", "/value")
	out, err := cmd.Output()
	if err == nil {
		// 简化处理
		output := string(out)
		if strings.Contains(output, "NVIDIA") {
			gpus = append(gpus, GPU{Name: "NVIDIA GPU", Vendor: "NVIDIA", VRAM: 8.0})
		} else if strings.Contains(output, "AMD") {
			gpus = append(gpus, GPU{Name: "AMD GPU", Vendor: "AMD", VRAM: 8.0})
		} else if strings.Contains(output, "Intel") {
			gpus = append(gpus, GPU{Name: "Intel GPU", Vendor: "Intel", VRAM: 0, IsIntegrated: true})
		}
	}

	return gpus
}

// detectCPUFeatures 检测CPU指令集
func (d *Detector) detectCPUFeatures() (hasAVX2, hasAVX512 bool) {
	switch runtime.GOOS {
	case "linux":
		return d.detectCPUFeaturesLinux()
	case "darwin":
		return d.detectCPUFeaturesMacOS()
	case "windows":
		return d.detectCPUFeaturesWindows()
	default:
		return false, false
	}
}

func (d *Detector) detectCPUFeaturesLinux() (bool, bool) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return false, false
	}
	content := string(data)
	return strings.Contains(content, "avx2"), strings.Contains(content, "avx512")
}

func (d *Detector) detectCPUFeaturesMacOS() (bool, bool) {
	if runtime.GOARCH == "arm64" {
		return true, false // Apple Silicon支持NEON
	}
	cmd := exec.Command("sysctl", "-a")
	out, err := cmd.Output()
	if err != nil {
		return false, false
	}
	content := string(out)
	return strings.Contains(content, "AVX2"), strings.Contains(content, "AVX512")
}

func (d *Detector) detectCPUFeaturesWindows() (bool, bool) {
	// Windows检测较复杂，简化处理
	return true, false
}

// hasCUDA 检测CUDA是否可用
func (d *Detector) hasCUDA() bool {
	_, err := exec.LookPath("nvidia-smi")
	return err == nil
}

// hasROCm 检测ROCm是否可用
func (d *Detector) hasROCm() bool {
	_, err := exec.LookPath("rocm-smi")
	return err == nil
}

// hasMetal 检测Metal是否可用
func (d *Detector) hasMetal() bool {
	if runtime.GOOS != "darwin" {
		return false
	}
	_, err := exec.Command("system_profiler", "SPDisplaysDataType").Output()
	return err == nil
}

// recommendBackend 推荐最佳后端
func (d *Detector) recommendBackend(info *Info) string {
	if info.HasCUDA {
		return "llama.cpp (CUDA)"
	}
	if info.HasMetal && runtime.GOOS == "darwin" {
		return "llama.cpp (Metal)"
	}
	if info.HasROCm {
		return "llama.cpp (ROCm)"
	}
	if info.CPUTotalRAM >= 16 {
		return "llama.cpp (CPU, AVX2)"
	}
	return "ollama (推荐新手)"
}

// GetGPUName 获取主GPU名称
func (info *Info) GetGPUName() string {
	if len(info.GPUs) > 0 {
		return info.GPUs[0].Name
	}
	return "无独立显卡"
}

// GetTotalVRAM 获取总显存
func (info *Info) GetTotalVRAM() float64 {
	if info.TotalVRAM > 0 {
		return info.TotalVRAM
	}
	// 对于集成显卡，使用系统内存的一部分
	return info.CPUTotalRAM * 0.5
}
