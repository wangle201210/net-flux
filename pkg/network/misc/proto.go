package misc

import (
	"fmt"

	"github.com/dellinger2023/net-flux/gen"
	"google.golang.org/protobuf/proto"
)

func UnmarshalSystem(subcmd uint8, data []byte) (proto.Message, error) {
	var err error
	switch subcmd {
	case uint8(gen.SCMDSystem_PING):
		ping := &gen.Ping{}
		err = proto.Unmarshal(data, ping)
		return ping, err
	case uint8(gen.SCMDSystem_PONG):
		pong := &gen.Pong{}
		err = proto.Unmarshal(data, pong)
		return pong, err
	default:
		return nil, fmt.Errorf("unknown sys subcmd: %d", subcmd)
	}
}

func UnmarshalDiscovery(subcmd uint8, data []byte) (proto.Message, error) {
	var err error
	switch subcmd {
	case uint8(gen.SCMDDisco_REGISTER):
		register := &gen.Instance{}
		err = proto.Unmarshal(data, register)
		return register, err
	case uint8(gen.SCMDDisco_DEREGISTER):
		deregister := &gen.Deregister{}
		err = proto.Unmarshal(data, deregister)
		return deregister, err
	case uint8(gen.SCMDDisco_LOOKUP):
		lookup := &gen.Lookup{}
		err = proto.Unmarshal(data, lookup)
		return lookup, err
	case uint8(gen.SCMDDisco_LOOKUP_ACK):
		lookupAck := &gen.LookupAck{}
		err = proto.Unmarshal(data, lookupAck)
		return lookupAck, err
	default:
		return nil, fmt.Errorf("unknown disco subcmd: %d", subcmd)
	}
}

func UnmarshalDataReport(subcmd uint8, data []byte) (proto.Message, error) {
	var err error
	switch subcmd {
	case uint8(gen.SCMDDataReport_MACHINE_METRIC):
		machineMetric := &gen.MachineMetric{}
		err = proto.Unmarshal(data, machineMetric)
		return machineMetric, err
	case uint8(gen.SCMDDataReport_NETWORK_METRIC):
		networkMetric := &gen.NetworkMetric{}
		err = proto.Unmarshal(data, networkMetric)
		return networkMetric, err
	case uint8(gen.SCMDDataReport_INSTANCE_METRIC):
		instanceMetric := &gen.InstanceMetric{}
		err = proto.Unmarshal(data, instanceMetric)
		return instanceMetric, err
	case uint8(gen.SCMDDataReport_SESSION_METRIC):
		sessionMetric := &gen.SessionMetric{}
		err = proto.Unmarshal(data, sessionMetric)
		return sessionMetric, err
	case uint8(gen.SCMDDataReport_STREAM_ADD):
		streamAdd := &gen.StreamMetric{}
		err = proto.Unmarshal(data, streamAdd)
		return streamAdd, err
	case uint8(gen.SCMDDataReport_STREAM_DELETE):
		streamDelete := &gen.StreamMetric{}
		err = proto.Unmarshal(data, streamDelete)
		return streamDelete, err
	case uint8(gen.SCMDDataReport_STREAM_STATUS):
		streamStatus := &gen.StreamMetric{}
		err = proto.Unmarshal(data, streamStatus)
		return streamStatus, err
	case uint8(gen.SCMDDataReport_STREAM_FAILED):
		streamFailed := &gen.StreamMetric{}
		err = proto.Unmarshal(data, streamFailed)
		return streamFailed, err
	case uint8(gen.SCMDDataReport_STREAMS_QUERY_REQ):
		streamsQueryReq := &gen.StreamMetric{}
		err = proto.Unmarshal(data, streamsQueryReq)
		return streamsQueryReq, err
	case uint8(gen.SCMDDataReport_STREAMS_QUERY_ACK):
		streamsQueryAck := &gen.StreamMetric{}
		err = proto.Unmarshal(data, streamsQueryAck)
		return streamsQueryAck, err
	default:
		return nil, fmt.Errorf("unknown report subcmd: %d", subcmd)
	}
}

func UnmarshalConfig(subcmd uint8, data []byte) (proto.Message, error) {
	var err error
	switch subcmd {
	case uint8(gen.SCMDConfig_CF_WHITEIP_CHANGED):
		whiteipChanged := &gen.WhiteipChanged{}
		err = proto.Unmarshal(data, whiteipChanged)
		return whiteipChanged, err
	case uint8(gen.SCMDConfig_CF_LIMIT_CHANGED):
		limitChanged := &gen.LimitChanged{}
		err = proto.Unmarshal(data, limitChanged)
		return limitChanged, err
	default:
		return nil, fmt.Errorf("unknown config subcmd: %d", subcmd)
	}
}

func UnmarshalEvent(subcmd uint8, data []byte) (proto.Message, error) {
	var err error
	switch subcmd {
	case uint8(gen.SCMDEvent_SE_ACCESS_DENIED):
		accessDenied := &gen.AccessDeniedEvent{}
		err = proto.Unmarshal(data, accessDenied)
		return accessDenied, err
	case uint8(gen.SCMDEvent_SE_IP_BLOCKED):
		ipBlocked := &gen.IpBlockedEvent{}
		err = proto.Unmarshal(data, ipBlocked)
		return ipBlocked, err
	default:
		return nil, fmt.Errorf("unknown event subcmd: %d", subcmd)
	}
}

func UnmarshalControl(subcmd uint8, data []byte) (proto.Message, error) {
	var err error
	// TODO: implement
	return nil, err
}

func Unbox(cmd, subcmd uint8, data []byte) (proto.Message, error) {

	switch cmd {
	case uint8(gen.CMD_SYSTEM):
		return UnmarshalSystem(subcmd, data)
	case uint8(gen.CMD_DISCOVERY):
		return UnmarshalDiscovery(subcmd, data)
	case uint8(gen.CMD_DATA_REPORT):
		return UnmarshalDataReport(subcmd, data)
	case uint8(gen.CMD_CONFIG):
		return UnmarshalConfig(subcmd, data)
	case uint8(gen.CMD_EVENT):
		return UnmarshalEvent(subcmd, data)
	case uint8(gen.CMD_CONTROL):
		return UnmarshalControl(subcmd, data)
	default:
		return nil, fmt.Errorf("unknown command: %d,%d", cmd, subcmd)
	}
}
