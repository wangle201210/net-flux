package dataflow

import (
	"errors"
	"sync"

	"github.com/dellinger2023/net-flux/pkg/dataflow/core"
	"github.com/dellinger2023/net-flux/pkg/dataflow/task"
	"github.com/dellinger2023/net-flux/pkg/network"
)

var (
	mutex         sync.RWMutex
	cli           *network.TcpClient
	configuration *core.DataflowConfig
)

func Initialize(tcpClient *network.TcpClient, config *core.DataflowConfig) error {
	if config == nil {
		return errors.New("config is nil")
	}

	if tcpClient == nil || tcpClient.IsClosed() {
		return errors.New("tcp client is not initialized")
	}

	mutex.Lock()
	cli = tcpClient
	configuration = config
	mutex.Unlock()

	channel, err := core.NewChannel(tcpClient)
	if err != nil {
		return err
	}

	if config.ReportMachineInterval > 0 {
		t := task.NewMachineTask(config.NodeInfo, config.ReportMachineInterval, channel)
		task.Manager().AddTask(t)
	}

	if config.ReportNetworkInterval > 0 && len(config.ProbeAddresses) > 0 {
		t := task.NewNetCollectTask(channel, config.ReportNetworkInterval,
			config.NodeInfo, config.ProbeAddresses)
		task.Manager().AddTask(t)
	}

	return task.Manager().Start()
}

func Shutdown() error {
	task.Manager().Stop()

	mutex.Lock()
	cli = nil
	configuration = nil
	mutex.Unlock()
	return nil
}
