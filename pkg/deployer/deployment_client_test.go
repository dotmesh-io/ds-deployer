package deployer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
	"google.golang.org/grpc"
)

var _ deployer_v1.DeployerServer = &TestingsServer{}

type SrvOpts struct {
	// PodLogsGetter PodLogsGetter

	// meant for testing, for example reject new connection due to auth failure
	ConnectionError error

	Port int
}

type TestingsServer struct {
	opts *SrvOpts

	// this list of requests will be dispatched to a client that connects
	dispatchLogRequests []*deployer_v1.LogsRequest

	receivedLogs []*deployer_v1.Logs
}

func (s *TestingsServer) StreamDeployments(*deployer_v1.DeploymentFilter, deployer_v1.Deployer_StreamDeploymentsServer) error {
	return nil
}

func (s *TestingsServer) ListDeployments(context.Context, *deployer_v1.DeploymentFilter) (*deployer_v1.GetDeploymentsResponse, error) {
	return nil, nil
}

func (s *TestingsServer) UpdateDeployment(context.Context, *deployer_v1.UpdateDeploymentRequest) (*deployer_v1.UpdateDeploymentResponse, error) {
	return nil, nil
}

func (s *TestingsServer) UpdateDeployer(context.Context, *deployer_v1.UpdateDeployerRequest) (*deployer_v1.UpdateDeployerResponse, error) {
	return nil, nil
}

func (s *TestingsServer) StreamLogRequests(f *deployer_v1.LogsFilter, stream deployer_v1.Deployer_StreamLogRequestsServer) error {

	for _, req := range s.dispatchLogRequests {
		err := stream.Send(req)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *TestingsServer) SendLogs(stream deployer_v1.Deployer_SendLogsServer) error {

	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}

		s.receivedLogs = append(s.receivedLogs, msg)
	}
}

type TestingPodLogsGetter struct {
	buf *ClosingBuffer
}

func (g *TestingPodLogsGetter) Logs(request *deployer_v1.LogsRequest) (io.ReadCloser, error) {
	return g.buf, nil
}

// ClosingBuffer is a helper wrapper to make bytes.Buffer be compatible with a ReadCloser that we get
// from kube client-go
type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() (err error) {
	return
}

func NewTestingServer(testingOpts *SrvOpts) (srv *TestingsServer, teardown func()) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", testingOpts.Port))
	if err != nil {
		log.Fatalf("failed to start TCP listener on port %d: %s", testingOpts.Port, err)
	}
	var opts []grpc.ServerOption
	s := &TestingsServer{
		opts: testingOpts,
	}

	grpcSrv := grpc.NewServer(opts...)
	deployer_v1.RegisterDeployerServer(grpcSrv, s)

	go func() {
		err := grpcSrv.Serve(listener)
		if err != nil {
			log.Fatalf("failed to start grpc server: %s", err)
		}
	}()

	time.Sleep(60 * time.Millisecond)

	return s, grpcSrv.Stop
}
