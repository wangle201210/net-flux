package task

import (
	"time"

	"github.com/dellinger2023/net-flux/gen"
	"github.com/dellinger2023/net-flux/pkg/dataflow/core"
	"github.com/dellinger2023/net-flux/pkg/dataflow/gather"
	"github.com/dellinger2023/net-flux/pkg/logger"
)

const machineTaskName = "machine"

type machineTask struct {
	nodeInfo *core.NodeInfo
	interval time.Duration
	next     time.Time
	channel  core.DataflowChannel
}

// Next implements Task.
func (t *machineTask) Next() time.Time {
	return t.next
}

func NewMachineTask(nodeInfo *core.NodeInfo, interval time.Duration, channel core.DataflowChannel) Task {
	return &machineTask{nodeInfo: nodeInfo, interval: interval, channel: channel}
}

func (t *machineTask) ID() string {
	return machineTaskName
}

func (t *machineTask) Run() error {

	defer func() {
		t.next = time.Now().Add(t.interval)
	}()

	cpuInfo, err := gather.CollectCPUInfo()
	if err != nil {
		return err
	}

	memInfo, err := gather.CollectMemInfo()
	if err != nil {
		return err
	}

	diskInfo, err := gather.CollectDiskInfo()
	if err != nil {
		return err
	}

	metric := &gen.MachineMetric{
		CpuUsage:       cpuInfo.Usage,
		CpuCount:       int32(cpuInfo.Count),
		MemUsed:        int64(memInfo.Used),
		MemTotal:       int64(memInfo.Total),
		DiskReadBytes:  int64(diskInfo.ReadBytes),
		DiskWriteBytes: int64(diskInfo.WriteBytes),
		DiskReadTime:   int64(diskInfo.ReadTime),
		DiskWriteTime:  int64(diskInfo.WriteTime),
		DiskReadCount:  int64(diskInfo.ReadCount),
		DiskWriteCount: int64(diskInfo.WriteCount),
		Timestamp:      time.Now().Unix(),
	}

	if t.nodeInfo != nil {
		metric.MachineId = t.nodeInfo.ID
		metric.InstanceId = t.nodeInfo.ID
		metric.InstanceName = t.nodeInfo.Name
		metric.NodeId = t.nodeInfo.ID
		metric.NodeName = t.nodeInfo.Name
		metric.NodeStatus = t.nodeInfo.Status
		metric.Status = t.nodeInfo.Status
		metric.Version = t.nodeInfo.Version
		metric.StartTime = t.nodeInfo.StartTime
	}

	if err := t.channel.SendMachineInfo(metric); err != nil {
		logger.Errorf("send machine info failed: %v", err)
		return err
	}

	logger.Debugf("send machine info: %+v", metric)
	return nil
}
