package client

import (
	"google.golang.org/grpc/naming"
)

type RegistryResolver struct{}

func (r *RegistryResolver) Resolve(serviceName string) (naming.Watcher, error) {
	//panic("wtf")

	var ch chan []*naming.Update = make(chan []*naming.Update, 1)
	sw := &targetWatcher{updates: ch, serviceName: serviceName}
	notifyTargetChanges(sw)
	return sw, nil
}

type targetWatcher struct {
	updates        chan []*naming.Update
	currentTargets []string
	serviceName    string
}

func (w *targetWatcher) Next() ([]*naming.Update, error) {
	return <-w.updates, nil
}

func (w *targetWatcher) Close() {
	close(w.updates)
}
