package gather

const (
	CPUGatherID     = "cpu"
	MemGatherID     = "mem"
	DiskGatherID    = "disk"
	TrafficGatherID = "traffic"
)

type CPUInfo struct {
	Count uint16
	Usage float64
}

type MemInfo struct {
	Used  uint64
	Total uint64
}

type DiskInfo struct {
	ReadCount  uint64
	WriteCount uint64
	ReadBytes  uint64
	WriteBytes uint64
	ReadTime   uint64
	WriteTime  uint64
	Used       uint64
	Total      uint64
}
