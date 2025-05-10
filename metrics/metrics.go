package metrics

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

type ProcessDetail struct {
	PID        int32
	PPID       int32
	Name       string
	Username   string
	Status     string
	CreateTime int64
	CPUPercent float64
	NumThreads int32
	MemoryRSS  uint64
	MemoryVMS  uint64
}

type SystemMonitor interface {
	GetProcessDetails() ([]ProcessDetail, error)
	GetCpuUsage() (float64, error)
	GetMemInfo() (*mem.VirtualMemoryStat, error)
}

type GopsutilMonitor struct {
	refreshInterval time.Duration
}

func NewGopsutilMonitor(refreshInterval time.Duration) *GopsutilMonitor {
	return &GopsutilMonitor{
		refreshInterval: refreshInterval,
	}
}

func (g *GopsutilMonitor) GetCpuUsage() (float64, error) {
	cpuPercentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Printf("Error getting CPU usage: %v\n", err)
	}

	if len(cpuPercentages) > 0 {
		return cpuPercentages[0], nil
	}
	return 0, fmt.Errorf("no CPU usage data available")

}

func (g *GopsutilMonitor) GetProcessDetails() ([]ProcessDetail, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var details []ProcessDetail
	for _, p := range procs {
		name, _ := p.Name()
		cpuPercent, _ := p.CPUPercent() // This needs to be called after an initial call or with a sleep
		memInfo, _ := p.MemoryInfo()
		status, _ := p.Status()
		username, _ := p.Username()
		createTime, _ := p.CreateTime() // ms since epoch
		ppid, _ := p.Ppid()
		numThreads, _ := p.NumThreads()

		detail := ProcessDetail{
			PID:        p.Pid,
			Name:       name,
			CPUPercent: cpuPercent,
			MemoryRSS:  0, // Default if not available
			MemoryVMS:  0, // Default if not available
			Status:     status,
			Username:   username,
			CreateTime: createTime,
			PPID:       ppid,
			NumThreads: numThreads,
		}
		if memInfo != nil {
			detail.MemoryRSS = memInfo.RSS
			detail.MemoryVMS = memInfo.VMS
		}
		details = append(details, detail)
	}

	return details, nil
}

// XXX: Basic wrapper for now.
func (g *GopsutilMonitor) GetMemInfo() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}
