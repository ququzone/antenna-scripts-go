package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	mainnetRPC     = "api.iotex.one:443"
	testnetRPC     = "api.testnet.iotex.one:443"
	mainnetChainID = 1
	testnetChainID = 2
)

func NewDefaultGRPCConn(endpoint string) (*grpc.ClientConn, error) {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Second)),
		grpc_retry.WithMax(3),
	}
	return grpc.Dial(endpoint,
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
}

func main() {
	conn, err := NewDefaultGRPCConn(mainnetRPC)
	if err != nil {
		log.Fatalf("Create connect error: %v", err)
	}

	client := iotexapi.NewAPIServiceClient(conn)

	action, err := client.GetActions(context.Background(), &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByHash{
			ByHash: &iotexapi.GetActionByHashRequest{
				ActionHash: "1075bc78ea75f7e1e8a14bec239d2c1708f161856d4f5b1244a5d13389be80cf",
			},
		},
	})
	if err != nil {
		log.Fatalf("Query actions error: %v", err)
	}

	fmt.Println(action.ActionInfo[0].Action.Encoding)
}
