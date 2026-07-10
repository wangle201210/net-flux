package util

import (
	"os"
	"strconv"

	"github.com/dellinger2023/net-flux/pkg/logger"
)

const (
	ENV_MB_MODE          = "MB_MODE"
	ENV_MB_EXTERNAL_IP   = "MB_EXTERNAL_IP"
	ENV_MB_EXTERNAL_PORT = "MB_EXTERNAL_PORT"
	ENV_MB_INTERNAL_IP   = "MB_INTERNAL_IP"
	ENV_MB_INTERNAL_PORT = "MB_INTERNAL_PORT"
	ENV_MB_NODE          = "MB_NODE"
	ENV_MB_REDIS_HOST    = "MB_REDIS_HOST"
	ENV_MB_REDIS_PORT    = "MB_REDIS_PORT"
	ENV_HOSTNAME         = "HOSTNAME"
)

func EnvHostname() string {
	return os.Getenv(ENV_HOSTNAME)
}

func EnvNode() int {
	val := os.Getenv(ENV_MB_NODE)
	if IsEmptyStr(val) {
		// 未配置则返回默认节点 1
		return 1
	}
	node, err := strconv.Atoi(val)
	if err != nil {
		logger.Errorf("fail to convert node to int: %v", err)
		return 1
	}
	return node
}

func EnvExternalIP() string {
	return os.Getenv(ENV_MB_EXTERNAL_IP)
}

func EnvInternalIP() string {
	return os.Getenv(ENV_MB_INTERNAL_IP)
}

func EnvExternalPort() int {
	val := os.Getenv(ENV_MB_EXTERNAL_PORT)
	if IsEmptyStr(val) {
		return 0
	}
	port, err := strconv.Atoi(val)
	if err != nil {
		logger.Errorf("fail to convert external port to int: %v", err)
		return 0
	}
	return port
}

func EnvInternalPort() int {
	val := os.Getenv(ENV_MB_INTERNAL_PORT)
	if IsEmptyStr(val) {
		return 0
	}
	port, err := strconv.Atoi(val)
	if err != nil {
		logger.Errorf("fail to convert internal port to int: %v", err)
		return 0
	}
	return port
}

func EnvMode() string {
	return os.Getenv(ENV_MB_MODE)
}
