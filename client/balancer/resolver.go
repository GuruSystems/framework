package balancer

import (
//	"fmt"
	//
	"google.golang.org/grpc/naming"
	"github.com/GuruSystems/framework/client/registry"
	pb "github.com/GuruSystems/framework/proto/registrar"
)

type RegistryResolver struct{}

func (r *RegistryResolver) Resolve(serviceName string) (naming.Watcher, error) {

	serverAddresses, err := registry.GetHosts(serviceName, pb.Apitype_grpc)
    if err != nil {
        return nil, err
    }

	var ups []*naming.Update
	for _,a := range serverAddresses {
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
