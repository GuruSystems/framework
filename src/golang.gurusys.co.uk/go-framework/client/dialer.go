package client

import (
	"fmt"
	"google.golang.org/grpc"
	"sync"
)

type Dialer struct {
	conn *grpc.ClientConn
	sync.RWMutex
}

// opens a tcp connection to a gurupath.
func (dialer *Dialer) DialService(serviceName string) error {

	dialer.Lock()
	defer dialer.Unlock()

	fmt.Println("protoclient.DialService: Dialling " + serviceName + " and blocking until successful connection...")

	if dialer.conn == nil {

		conn, err := grpc.Dial(
			serviceName,
			grpc.WithBlock(),
			grpc.WithBalancer(
				grpc.RoundRobin(&RegistryResolver{}),
			),
			grpc.WithTransportCredentials(
				GetClientCreds(),
			),
		)
		if err != nil {
			return err
		}

		dialer.conn = conn

		fmt.Println("protoclient.DialService: Connected to address(es)...")

	}

	return nil
}
