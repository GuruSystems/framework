package client

import (
	"fmt"
	"golang.gurusys.co.uk/go-framework/client/registry"
	"golang.gurusys.co.uk/go-framework/cmdline"
	pb "golang.gurusys.co.uk/go-framework/proto/registrar"
	"google.golang.org/grpc/naming"
)

type RegistryResolver struct{}

func (r *RegistryResolver) Resolve(serviceName string) (naming.Watcher, error) {

	fmt.Printf("Resolving service address \"%s\" via registry %s...\n", serviceName, cmdline.GetRegistryAddress())

	serverAddresses, err := registry.GetHosts(serviceName, pb.Apitype_grpc)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Instances of %s: %s\n", serviceName, serverAddresses)

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
