package gather

import (
	"github.com/dellinger2023/net-flux/pkg/logger"
	"github.com/shirou/gopsutil/v4/mem"
)

func CollectMemInfo() (*MemInfo, error) {
	info := &MemInfo{}
	m, err := mem.VirtualMemory()
	if err != nil {
		logger.Errorf("get memory usage failed: %v", err)
		return nil, err
	}
	info.Used = m.Used
	info.Total = m.Total

	return info, nil
}
