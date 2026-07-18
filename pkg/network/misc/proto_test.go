package misc

import (
	"testing"

	"github.com/dellinger2023/net-flux/gen"
	"google.golang.org/protobuf/proto"
)

func TestUnmarshalDataReportSessionLifecycle(t *testing.T) {
	payload, err := proto.Marshal(&gen.SessionMetric{SessionId: "session-a"})
	if err != nil {
		t.Fatalf("marshal session metric: %v", err)
	}

	for _, subcmd := range []gen.SCMDDataReport{
		gen.SCMDDataReport_SESSION_ADD,
		gen.SCMDDataReport_SESSION_DELETE,
		gen.SCMDDataReport_SESSION_STATUS,
		gen.SCMDDataReport_SESSION_FAILED,
		gen.SCMDDataReport_SESSION_QUERY_REQ,
		gen.SCMDDataReport_SESSION_QUERY_ACK,
	} {
		message, err := UnmarshalDataReport(uint8(subcmd), payload)
		if err != nil {
			t.Fatalf("unmarshal %s: %v", subcmd, err)
		}
		session, ok := message.(*gen.SessionMetric)
		if !ok {
			t.Fatalf("expected session metric for %s, got %T", subcmd, message)
		}
		if session.GetSessionId() != "session-a" {
			t.Fatalf("expected session id session-a for %s, got %s", subcmd, session.GetSessionId())
		}
	}
}

func TestReportMetricsRoundTripDeviceID(t *testing.T) {
	streamPayload, err := proto.Marshal(&gen.StreamMetric{
		StreamId: "stream-a",
		DeviceId: "device-a",
	})
	if err != nil {
		t.Fatalf("marshal stream metric: %v", err)
	}
	streamMessage, err := UnmarshalDataReport(uint8(gen.SCMDDataReport_STREAM_ADD), streamPayload)
	if err != nil {
		t.Fatalf("unmarshal stream metric: %v", err)
	}
	stream, ok := streamMessage.(*gen.StreamMetric)
	if !ok {
		t.Fatalf("expected stream metric, got %T", streamMessage)
	}
	if stream.GetDeviceId() != "device-a" {
		t.Fatalf("expected stream device id device-a, got %s", stream.GetDeviceId())
	}

	sessionPayload, err := proto.Marshal(&gen.SessionMetric{
		SessionId: "session-a",
		DeviceId:  "device-a",
	})
	if err != nil {
		t.Fatalf("marshal session metric: %v", err)
	}
	sessionMessage, err := UnmarshalDataReport(uint8(gen.SCMDDataReport_SESSION_ADD), sessionPayload)
	if err != nil {
		t.Fatalf("unmarshal session metric: %v", err)
	}
	session, ok := sessionMessage.(*gen.SessionMetric)
	if !ok {
		t.Fatalf("expected session metric, got %T", sessionMessage)
	}
	if session.GetDeviceId() != "device-a" {
		t.Fatalf("expected session device id device-a, got %s", session.GetDeviceId())
	}
}
