package task

import (
	"time"

	"github.com/dellinger2023/net-flux/pkg/dataflow/core"
	"github.com/dellinger2023/net-flux/pkg/logger"
)

const netCollectTaskName = "net_collect"

type netCollectTask struct {
	channel  core.DataflowChannel
	next     time.Time
	interval time.Duration
	nodeInfo *core.NodeInfo
	addrList []string
}

func NewNetCollectTask(channel core.DataflowChannel, interval time.Duration,
	nodeInfo *core.NodeInfo, addrList []string) Task {

	return &netCollectTask{channel: channel, interval: interval,
		nodeInfo: nodeInfo, addrList: addrList}
}

func (t *netCollectTask) ID() string {
	return netCollectTaskName
}

func (t *netCollectTask) Run() error {
	defer func() {
		t.next = time.Now().Add(t.interval)
	}()

	for _, addr := range t.addrList {
		logger.Debugf("net collect task: collect network info from %s", addr)
		prober := newProber(addr, t.nodeInfo)
		if err := prober.Probe(); err != nil {
			logger.Errorf("net collect task: collect network info from %s failed: %v", addr, err)
			continue
		}
		result := prober.Result()
		if result.GetExtra()["reachable"] == "false" {
			logger.Warningf("net collect task: %s unreachable, rtt=%d jitter=%d packet_loss=%.0f%%",
				addr, result.GetRtt(), result.GetJitter(), result.GetPacketLoss()*100)
		}
		if err := t.channel.SendNetworkInfo(result); err != nil {
			logger.Errorf("net collect task: send network info to channel failed: %v", err)
			continue
		}
	}
	return nil
}

func (t *netCollectTask) Next() time.Time {
	return t.next
}
