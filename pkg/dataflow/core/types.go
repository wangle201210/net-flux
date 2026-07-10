package core

import (
	"time"

	"github.com/dellinger2023/net-flux/gen"
)

type NodeInfo struct {
	ID        string
	Name      string
	IP        string
	Port      int
	Node      int
	Version   string
	Status    string
	StartTime int64
}

type DataflowConfig struct {
	NodeInfo              *NodeInfo
	ProbeAddresses        []string
	ReportMachineInterval time.Duration
	ReportNetworkInterval time.Duration
}

type DataflowChannel interface {
	SendMachineInfo(metrics ...*gen.MachineMetric) error
	SendNetworkInfo(metrics ...*gen.NetworkMetric) error
	// SendStreamAdd(metric *gen.StreamMetric) error
	// SendStreamDelete(metric *gen.StreamMetric) error
	// SendStreamStatus(metric *gen.StreamMetric) error
	// SendStreamFailed(metric *gen.StreamMetric) error
	// SendStreamsQueryReq(metric *gen.StreamMetric) error
}
