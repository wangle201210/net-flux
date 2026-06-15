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
