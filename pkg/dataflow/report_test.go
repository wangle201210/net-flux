package dataflow

import (
	"testing"

	"github.com/dellinger2023/net-flux/pkg/network"
)

func TestInitializeStoresClient(t *testing.T) {
	mutex.Lock()
	cli = nil
	mutex.Unlock()

	testClient := &network.TcpClient{}
	mutex.Lock()
	cli = testClient
	got := cli
	mutex.Unlock()

	if got != testClient {
		t.Fatal("package cli was not assigned")
	}
}
