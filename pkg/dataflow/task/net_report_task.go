package task

import (
	"time"

	"github.com/dellinger2023/net-flux/pkg/dataflow/core"
)

const netReportTaskName = "net_report"

type netReportTask struct {
	channel  core.DataflowChannel
	next     time.Time
	interval time.Duration
	nodeInfo *core.NodeInfo
}

func NewNetReportTask(channel core.DataflowChannel, interval time.Duration, nodeInfo *core.NodeInfo) Task {
	return &netReportTask{channel: channel, interval: interval, nodeInfo: nodeInfo}
}

func (t *netReportTask) ID() string {
	return netReportTaskName
}

func (t *netReportTask) Run() error {
	return nil
}

func (t *netReportTask) Next() time.Time {
	return t.next
}
