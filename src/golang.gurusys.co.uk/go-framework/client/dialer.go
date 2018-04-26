package client

import (
	"flag"
	"fmt"
	"sync"

	"google.golang.org/grpc"
)

var (
	dialer_debug = flag.Bool("dialer_debug", false, "set to true to debug the grpc dialer")
)

type Dialer struct {
	conn *grpc.ClientConn
	sync.RWMutex
}

// opens a tcp connection to a gurupath.
func (dialer *Dialer) DialService(serviceName string) error {

	dialer.Lock()
	defer dialer.Unlock()
	if *dialer_debug {
		fmt.Println("protoclient.DialService: Dialling " + serviceName + " and blocking until successful connection...")
	}
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
		if *dialer_debug {
			fmt.Printf("Connected to %s\n", serviceName)
		}

	}

	return nil
}

func (dialer *Dialer) GetConn() *grpc.ClientConn {
	return dialer.conn
}
