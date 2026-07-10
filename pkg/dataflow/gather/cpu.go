package gather

import (
	"time"

	"github.com/dellinger2023/net-flux/pkg/logger"
	"github.com/shirou/gopsutil/v4/cpu"
)

func CollectCPUInfo() (*CPUInfo, error) {
	info := &CPUInfo{}
	cpuCount, err := cpu.Counts(true)
	if err != nil {
		logger.Errorf("get cpu count failed: %v", err)
		return nil, err
	}
	info.Count = uint16(cpuCount)

	// 第一次采样
	times1, _ := cpu.Times(true)
	timer := time.NewTimer(time.Second)
	select {
	case <-timer.C:
	}
	// 第二次采样
	times2, _ := cpu.Times(true)

	count := len(times1)
	usageTotal := 0.0
	for i := range times1 {
		// 计算每个核的时间差
		total1 := times1[i].User + times1[i].System + times1[i].Idle
		total2 := times2[i].User + times2[i].System + times2[i].Idle

		deltaTotal := total2 - total1
		if deltaTotal == 0 {
			continue
		}

		deltaIdle := times2[i].Idle - times1[i].Idle
		usage := (1.0 - float64(deltaIdle)/float64(deltaTotal)) * 100

		usageTotal += usage
	}
	info.Usage = usageTotal / float64(count)

	return info, nil
}
