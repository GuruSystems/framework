package client

import (
	"fmt"
	pb "golang.gurusys.co.uk/apis/registry"
	"golang.gurusys.co.uk/go-framework/cmdline"
	"golang.gurusys.co.uk/go-framework/registry"
	"google.golang.org/grpc/naming"
	"math/rand"
	"time"
)

var (
	randomer = rand.New(rand.NewSource(time.Now().Unix()))
)

type RegistryResolver struct{}

func (r *RegistryResolver) Resolve(serviceName string) (naming.Watcher, error) {
	if *dialer_debug {
		fmt.Printf("Resolving service address \"%s\" via registry %s...\n", serviceName, cmdline.GetRegistryAddress())
	}
	var err error
	var serverAddresses []string
	for {
		// if we find 0 addresses (or have an error), e.g. registry is unavailable
		// we got to sleep and try again every 5 seconds(ish)
		serverAddresses, err = registry.GetHosts(serviceName, pb.Apitype_grpc)
		if err == nil {
			break
		}
		time.Sleep(time.Second * time.Duration(randomer.Intn(10)))

		fmt.Printf("error retrieving hosts for %s: %s\n", serviceName, err)

	}
	if *dialer_debug {
		fmt.Printf("Instances of %s: %s\n", serviceName, serverAddresses)
	}
	// we randomize the server addresses we got initially.
	// otherwise we always end up hitting the same thing

	for i := len(serverAddresses) - 1; i > 0; i-- {
		j := randomer.Intn(i + 1)
		serverAddresses[i], serverAddresses[j] = serverAddresses[j], serverAddresses[i]
	}
	var ups []*naming.Update
	for _, a := range serverAddresses {
		ups = append(
			ups,
			&naming.Update{naming.Add, a, ""},
		)
	}

	var ch chan []*naming.Update = make(chan []*naming.Update, 1)
	ch <- ups

	return &staticWatcher{ch}, nil
}

type staticWatcher struct {
	updates chan []*naming.Update
}

func (w *staticWatcher) Next() ([]*naming.Update, error) {
	return <-w.updates, nil
}

func (w *staticWatcher) Close() {
	close(w.updates)
}
