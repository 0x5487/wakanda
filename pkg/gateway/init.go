package gateway

import (
	"context"

	"github.com/jasonsoft/wakanda/internal/config"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/wakanda/pkg/dispatcher/proto"
	"google.golang.org/grpc"
)

var (
	_manager          *Manager
	_dispatcherClient proto.DispatcherClient
)

// customCredential 自定義認證
type customCredential struct{}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"user_id": "jason",
		"roles":   "admin",
	}, nil
}

func (c customCredential) RequireTransportSecurity() bool {
	return false
}

func Initialize(config *config.Configuration) {
	_manager = NewManager()

	// Set up a connection to the server.
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	// 使用自定義認證
	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))

	conn, err := grpc.Dial(config.Dispatcher.AdvertiseAddr, opts...)
	if err != nil {
		log.Fatalf("gateway: can't connect to messenger grpc service: %v", err)
	}
	//defer conn.Close()

	_dispatcherClient = proto.NewDispatcherClient(conn)
}
