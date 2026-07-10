package gather

type trafficGather struct {
}

func (g *trafficGather) Gather() (interface{}, error) {
	return nil, nil
}

func (g *trafficGather) ID() string {
	return TrafficGatherID
}
