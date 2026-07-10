package gather

import (
	"github.com/dellinger2023/net-flux/pkg/logger"
	"github.com/shirou/gopsutil/v4/disk"
)

func CollectDiskInfo() (*DiskInfo, error) {
	var totalReadTime uint64 = 0
	var totalWriteTime uint64 = 0
	inf := &DiskInfo{}
	diskInfo, err := disk.IOCounters()
	if err != nil {
		logger.Errorf("get disk usage failed: %v", err)
		return nil, err
	}

	counter := 0

	for _, info := range diskInfo {
		totalReadTime += info.ReadTime
		totalWriteTime += info.WriteTime
		inf.ReadCount += info.MergedReadCount
		inf.WriteCount += info.MergedWriteCount
		inf.ReadBytes += info.ReadBytes
		inf.WriteBytes += info.WriteBytes
		counter++
	}

	if counter > 0 {
		inf.ReadTime = totalReadTime / uint64(counter)
		inf.WriteTime = totalWriteTime / uint64(counter)
	}
	return inf, nil
}
