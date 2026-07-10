package task

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/dellinger2023/net-flux/pkg/dataflow/core"
)

func TestNetProber_Unreachable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	prober := newProber(server.URL, &core.NodeInfo{ID: "test", IP: "127.0.0.1", Node: 1})
	prober.probeCount = 5
	prober.timeout = 5 * time.Second
	prober.client = prober.buildClient()

	start := time.Now()
	if err := prober.Probe(); err != nil {
		t.Fatalf("probe failed: %v", err)
	}
	elapsed := time.Since(start)
	if elapsed > 6*time.Second {
		t.Fatalf("probe took too long: %v", elapsed)
	}

	result := prober.Result()
	if result == nil {
		t.Fatal("probe result is nil")
	}
	if result.Rtt <= 1000 {
		t.Fatalf("expected rtt > 1000, got %d", result.Rtt)
	}
	if result.Jitter <= 5000 {
		t.Fatalf("expected jitter > 5000, got %d", result.Jitter)
	}
	if result.PacketLoss != unreachablePacketLoss {
		t.Fatalf("expected packet loss 100%%, got %v", result.PacketLoss)
	}
	if result.Extra["reachable"] != "false" {
		t.Fatalf("expected unreachable result, got extra=%v", result.Extra)
	}
}

func TestNetProber_Probe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	prober := newProber(server.URL, &core.NodeInfo{ID: "test", IP: "127.0.0.1", Node: 1})
	prober.probeCount = 3
	prober.timeout = 2 * time.Second
	prober.client = prober.buildClient()

	if err := prober.Probe(); err != nil {
		t.Fatalf("probe failed: %v", err)
	}

	result := prober.Result()
	if result == nil {
		t.Fatal("probe result is nil")
	}
	if result.Rtt <= 0 {
		t.Fatalf("expected positive rtt, got %d", result.Rtt)
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", result.StatusCode)
	}
	if result.PacketLoss != 0 {
		t.Fatalf("expected zero packet loss, got %v", result.PacketLoss)
	}
}

func TestSchedulerStopNoPanic(t *testing.T) {
	s := &scheduler{}
	if err := s.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	s.Stop()
	s.Stop()
}

func TestSchedulerWaitGroup(t *testing.T) {
	s := &scheduler{}
	s.AddTask(&stubScheduleTask{id: "t1", next: time.Now()})

	if err := s.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}

	done := make(chan struct{})
	go func() {
		s.Stop()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("stop did not return in time")
	}
}

type stubScheduleTask struct {
	id   string
	next time.Time
	mu   sync.Mutex
}

func (t *stubScheduleTask) ID() string   { return t.id }
func (t *stubScheduleTask) Next() time.Time { return t.next }
func (t *stubScheduleTask) Run() error {
	t.mu.Lock()
	t.next = time.Now().Add(time.Minute)
	t.mu.Unlock()
	return nil
}
