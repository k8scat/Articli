package segmentfault

import (
	"os"
	"testing"
)

func TestListNotes(t *testing.T) {
	client, err := NewClient(os.Getenv(TestEnvSegmentfaultToken))
	if err != nil {
		t.Fatal(err)
	}

	opts := &ListOptions{}
	resp, err := client.ListNotes(opts)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range resp.Rows {
		t.Logf("%+v", row)
	}
}
