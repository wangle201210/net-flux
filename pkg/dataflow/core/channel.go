package core

import (
	"errors"
	"sync"

	"github.com/dellinger2023/net-flux/gen"
	"github.com/dellinger2023/net-flux/pkg/logger"
	"github.com/dellinger2023/net-flux/pkg/network"
)

type channel struct {
	sync.RWMutex
	cli *network.TcpClient
}

// SendMachineInfo implements DataflowChannel.
func (c *channel) SendMachineInfo(metrics ...*gen.MachineMetric) error {

	if c.cli == nil || c.cli.IsClosed() {
		return errors.New("tcp client is not initialized")
	}

	for _, metric := range metrics {
		if metric == nil {
			continue
		}
		if err := c.cli.Write(uint8(gen.CMD_DATA_REPORT),
			uint8(gen.SCMDDataReport_MACHINE_METRIC), metric); err != nil {
			logger.Errorf("failed to send machine info: %v, %v", metric, err)
		}
	}

	return nil
}

// SendNetworkInfo implements DataflowChannel.
func (c *channel) SendNetworkInfo(metrics ...*gen.NetworkMetric) error {
	if c.cli == nil || c.cli.IsClosed() {
		return errors.New("tcp client is not initialized")
	}

	for _, metric := range metrics {
		if metric == nil {
			continue
		}
		if err := c.cli.Write(uint8(gen.CMD_DATA_REPORT),
			uint8(gen.SCMDDataReport_NETWORK_METRIC), metric); err != nil {
			logger.Errorf("failed to send network info: %v, %v", metric, err)
		}
	}
	return nil
}

func NewChannel(cli *network.TcpClient) (DataflowChannel, error) {

	if cli == nil || cli.IsClosed() {
		return nil, errors.New("tcp client is not initialized")
	}

	return &channel{
		cli: cli,
	}, nil
}
