package deployer

import (
	"bytes"
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
	"github.com/dotmesh-io/ds-deployer/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestStreamLogs(t *testing.T) {
	logger := logger.GetInstance().Sugar()

	logsGetter := &TestingPodLogsGetter{
		buf: &ClosingBuffer{bytes.NewBufferString("foo")},
	}

	port := 34444
	server, teardown := NewTestingServer(&SrvOpts{Port: port})
	defer teardown()

	server.dispatchLogRequests = []*deployer_v1.LogsRequest{
		{
			TxId:         "55",
			DeploymentId: "100",
			Container:    deployer_v1.LogsRequest_MODEL,
		},
	}

	client := New(&Opts{
		PodLogsGetter: logsGetter,
		RequireTLS:    false,
		Addr:          "localhost:" + strconv.Itoa(port),
		Logger:        logger,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.dial(ctx)
	assert.Nil(t, err)

	go func() {
		err = client.getLogRequests(ctx, &deployer_v1.LogsFilter{})
		assert.Nil(t, err)
	}()

	for {
		select {
		case <-ctx.Done():
			t.Error("context deadline exceeded, logs not retrieved")
			return
		default:
			if len(server.receivedLogs) != 1 {
				t.Logf("received %d", len(server.receivedLogs))
				time.Sleep(100 * time.Millisecond)
				continue
			}
			// success
			assert.Equal(t, true, server.receivedLogs[0].Eof)
			assert.Equal(t, "55", server.receivedLogs[0].TxId)
			return
		}
	}
}
