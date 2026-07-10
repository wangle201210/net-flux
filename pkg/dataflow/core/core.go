package core

import (
	"time"

	"github.com/dellinger2023/net-flux/pkg/logger"
	"github.com/dellinger2023/net-flux/pkg/util"
)

func ReadBaseNodeInfo(node int) (*NodeInfo, error) {

	hostname, err := util.Hostname()
	if err != nil {
		logger.Errorf("get hostname failed: %v", err)
		return nil, err
	}

	ip := util.EnvInternalIP()
	if util.IsEmptyStr(ip) {
		ips, err := util.HostIP()
		if err != nil {
			logger.Errorf("get host ip failed: %v", err)
			ip = "127.0.0.1"
		} else {
			ip = ips[0]
		}
	}

	nodeInfo := &NodeInfo{
		ID:        hostname,
		IP:        ip,
		Port:      8080,
		Node:      node,
		Version:   "1.0.0",
		Status:    "running",
		StartTime: time.Now().Unix(),
	}

	return nodeInfo, nil
}
