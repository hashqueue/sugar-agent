package internal

import (
	"fmt"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"

	"sugar-agent/pkg/utils"
)

type HostInfo struct {
	Hostname        string `json:"hostname"`
	Uptime          string `json:"uptime"`
	OS              string `json:"os"`              // ex: freebsd, linux
	Platform        string `json:"platform"`        // ex: ubuntu, linux mint
	PlatformFamily  string `json:"platformFamily"`  // ex: debian, rhel
	PlatformVersion string `json:"platformVersion"` // version of the complete OS
	KernelVersion   string `json:"kernelVersion"`   // version of the OS kernel (if available)
	KernelArch      string `json:"kernelArch"`      // native cpu architecture queried at runtime, as returned by `uname -m` or empty string in case of error
	HostID          string `json:"hostId"`          // ex: uuid
}

type CpuInfo struct {
	PhysicalCoresCount uint16 `json:"physicalCoresCount"` // physical cores count
	LogicalCoresCount  uint16 `json:"logicalCoresCount"`  // logical cores count
	ModelName          string `json:"modelName"`          // cpu model name
}

type DiskInfo struct {
	Total       float64 `json:"total"`       // total disk size in GB
	Free        float64 `json:"free"`        // free disk size in GB
	Used        float64 `json:"used"`        // used disk size in GB
	UsedPercent float64 `json:"usedPercent"` // used disk size in percent
}

type MemoryInfo struct {
	Total       float64 `json:"total"`       // total memory size in GB
	Available   float64 `json:"available"`   // available memory size in GB
	Used        float64 `json:"used"`        // used memory size in GB
	UsedPercent float64 `json:"usedPercent"` // used memory size in percent
	Free        float64 `json:"free"`        // free memory size in GB
	Cached      float64 `json:"cached"`      // cached memory size in GB
}

type LoadInfo struct {
	Load1  float64 `json:"load1"`  // 1 minute load average
	Load5  float64 `json:"load5"`  // 5 minute load average
	Load15 float64 `json:"load15"` // 15 minute load average
}

type DynamicDataSummary struct {
	TimeStamp  string     `json:"timeStamp"`
	CpuPercent float64    `json:"cpuPercent"`
	MemInfo    MemoryInfo `json:"memInfo"`
	DiskInfo   DiskInfo   `json:"diskInfo"`
	LoadInfo   LoadInfo   `json:"loadInfo"`
}

type PropertiesSummary struct {
	HostInfo HostInfo `json:"hostInfo"`
	CpuInfo  CpuInfo  `json:"cpuInfo"`
}

type PerfData struct {
	Properties PropertiesSummary    `json:"properties"`
	Data       []DynamicDataSummary `json:"perfData"`
}

func getCpuProperties() *CpuInfo {
	cpuPhysicalCoresCount, err := cpu.Counts(false)
	utils.LogOnError(err, "get cpu physical cores count failed")
	cpuLogicalCoresCount, err := cpu.Counts(true)
	utils.LogOnError(err, "get cpu logical cores count failed")
	info, err := cpu.Info()
	utils.LogOnError(err, "get cpu info failed")
	ModelName := info[0].ModelName
	cpuInfo := CpuInfo{
		PhysicalCoresCount: uint16(cpuPhysicalCoresCount),
		LogicalCoresCount:  uint16(cpuLogicalCoresCount),
		ModelName:          ModelName,
	}
	return &cpuInfo
}

func getCpuPercent() float64 {
	cpuPercent, err := cpu.Percent(time.Second, false)
	utils.LogOnError(err, "get cpu percent failed")
	val := humanizePercent(cpuPercent[0])
	return val
}

func getDiskInfo() *DiskInfo {
	diskInfoData, err := disk.Usage("/")
	utils.LogOnError(err, "get disk info failed")
	diskInfo := DiskInfo{
		Total:       humanizeGB(float64(diskInfoData.Total)),
		Free:        humanizeGB(float64(diskInfoData.Free)),
		Used:        humanizeGB(float64(diskInfoData.Used)),
		UsedPercent: humanizePercent(diskInfoData.UsedPercent),
	}
	return &diskInfo
}

func getMemoryInfo() *MemoryInfo {
	memInfoData, err := mem.VirtualMemory()
	utils.LogOnError(err, "get memory info failed")
	memInfo := MemoryInfo{
		Total:       humanizeGB(float64(memInfoData.Total)),
		Available:   humanizeGB(float64(memInfoData.Available)),
		Used:        humanizeGB(float64(memInfoData.Used)),
		UsedPercent: humanizePercent(memInfoData.UsedPercent),
		Free:        humanizeGB(float64(memInfoData.Free)),
		Cached:      humanizeGB(float64(memInfoData.Cached)),
	}
	return &memInfo
}

func getLoadInfo() *LoadInfo {
	loadInfoData, err := load.Avg()
	utils.LogOnError(err, "get load info failed")
	loadInfo := LoadInfo{
		Load1:  humanizePercent(loadInfoData.Load1),
		Load5:  humanizePercent(loadInfoData.Load5),
		Load15: humanizePercent(loadInfoData.Load15),
	}
	return &loadInfo
}

func getHostInfo() *HostInfo {
	hostInfoData, err := host.Info()
	utils.LogOnError(err, "get host info failed")
	hostInfo := HostInfo{
		Hostname:        hostInfoData.Hostname,
		Uptime:          humanizeDuration(time.Duration(hostInfoData.Uptime) * time.Second),
		OS:              hostInfoData.OS,
		Platform:        hostInfoData.Platform,
		PlatformFamily:  hostInfoData.PlatformFamily,
		PlatformVersion: hostInfoData.PlatformVersion,
		KernelVersion:   hostInfoData.KernelVersion,
		KernelArch:      hostInfoData.KernelArch,
		HostID:          hostInfoData.HostID,
	}
	return &hostInfo
}

// humanizeGB converts bytes to GB
// 1GB = 1024MB = 1024KB = 1024B
// return GB
func humanizeGB(bytes float64) float64 {
	val := strconv.FormatFloat(bytes/1024/1024/1024, 'f', 2, 64)
	valF, err := strconv.ParseFloat(val, 64)
	utils.LogOnError(err, "parse float failed")
	return valF
}

// humanizePercent converts float64 to 2 decimal places
// return percent
func humanizePercent(percent float64) float64 {
	if percent >= 100 {
		return 100
	}
	val := strconv.FormatFloat(percent, 'f', 2, 64)
	valF, err := strconv.ParseFloat(val, 64)
	utils.LogOnError(err, "parse float failed")
	return valF
}

// humanizeDuration converts time.Duration to human-readable string
// return human-readable string
func humanizeDuration(duration time.Duration) string {
	var res string
	seconds := int(duration.Seconds())
	days := seconds / 86400
	seconds -= days * 86400
	hours := seconds / 3600
	seconds -= hours * 3600
	minutes := seconds / 60
	seconds -= minutes * 60
	if days > 0 {
		res += fmt.Sprintf("%dd ", days)
	}
	if hours > 0 {
		res += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 {
		res += fmt.Sprintf("%dm ", minutes)
	}
	if seconds > 0 {
		res += fmt.Sprintf("%ds ", seconds)
	}
	if len(res) == 0 {
		return "0s"
	}
	return res[:len(res)-1] // remove trailing space
}

func StartGetPerfDataTask(duration uint64, count uint64) *PerfData {
	var dynamicData []DynamicDataSummary
	cpuInfo := getCpuProperties()
	hostInfo := getHostInfo()
	for i := 0; i < int(count); i++ {
		diskInfo := getDiskInfo()
		memInfo := getMemoryInfo()
		loadAvg := getLoadInfo()
		dynamicData = append(dynamicData, DynamicDataSummary{
			TimeStamp:  time.Now().Format("2006-01-02 15:04:05"),
			CpuPercent: getCpuPercent(),
			MemInfo:    *memInfo,
			DiskInfo:   *diskInfo,
			LoadInfo:   *loadAvg,
		})
		time.Sleep(time.Second * time.Duration(duration))
	}
	perfData := PerfData{
		PropertiesSummary{
			HostInfo: *hostInfo,
			CpuInfo:  *cpuInfo,
		},
		dynamicData[:],
	}
	return &perfData
}
