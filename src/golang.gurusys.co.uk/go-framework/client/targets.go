package client

import (
	"flag"
	"fmt"
	pb "golang.gurusys.co.uk/apis/registry"
	"golang.gurusys.co.uk/go-framework/cmdline"
	"golang.gurusys.co.uk/go-framework/registry"
	"google.golang.org/grpc/naming"
	"math/rand"
	"sync"
	"time"
)

// we keep a local "mirror" copy of the registry
// so that, in case the registry goes away for some time
// we can still do our load balancing.
// we also do that because we need to regularly requery the registry
// to find updates to our target locations.

// to enable periodic updates, call notifyTargetChanges(servicename/channel)
// and this thing will inform the channel of changes

// notifyTargetChanges() will block until at least one target becomes
// available.
// This should prevent services from coming up and immediately failing
// because they're missing a backend.
// this might cause trouble during testing, if so a variable is provided to
// change that behavior

var (
	randomer          = rand.New(rand.NewSource(time.Now().Unix()))
	dialer_wait_first = flag.Bool("dial_wait_for_target", true, "initially block until at least one target becomes available")
	targetNotifiers   []*targetNotifier
	notifierLock      sync.Mutex
)

type targetNotifier struct {
	serviceName string
	lock        sync.Mutex // lock updates to watchers
	watchers    []*targetWatcher
}

func GetOrCreateNotifier(service string) *targetNotifier {
	notifierLock.Lock()
	defer notifierLock.Unlock()
	for _, tn := range targetNotifiers {
		if tn.serviceName == service {
			return tn
		}
	}
	tn := &targetNotifier{serviceName: service}
	targetNotifiers = append(targetNotifiers, tn)
	go tn.requery()
	return tn
}

func notifyTargetChanges(ch *targetWatcher) error {
	if *dialer_debug {
		fmt.Printf("New notifier requested for %s\n", ch.serviceName)
	}
	tn := GetOrCreateNotifier(ch.serviceName)
	tn.AddWatcher(ch)
	return nil
}

// get a list of instances
func queryForActiveInstances(serviceName string) []string {

	if *dialer_debug {
		fmt.Printf("Resolving service address \"%s\" via registry %s...\n", serviceName, cmdline.GetRegistryAddress())
	}
	var err error
	var serverAddresses []string
	serverAddresses, err = registry.GetHosts(serviceName, pb.Apitype_grpc)
	// error. so what?
	if err != nil {
		if *dialer_debug {
			fmt.Printf("error retrieving hosts for %s: %s\n", serviceName, err)
		}
	} else {
		if *dialer_debug {
			fmt.Printf("Instances of %s: %s\n", serviceName, serverAddresses)
		}
	}

	return serverAddresses

}

// this thing runs in the background, one per servicename
func (tn *targetNotifier) requery() {
	for {
		serverAddresses := queryForActiveInstances(tn.serviceName)
		var ups []*naming.Update
		for _, a := range serverAddresses {
			ups = append(
				ups,
				&naming.Update{naming.Add, a, ""},
			)
		}
		tn.lock.Lock()
		curWatchers := make([]*targetWatcher, len(tn.watchers))
		i := 0
		for _, w := range tn.watchers {
			curWatchers[i] = w
			i++
		}
		tn.lock.Unlock()
		for _, watcher := range curWatchers {
			sendDiff(serverAddresses, watcher)
			//			watcher.updates <- ups
		}
		time.Sleep(time.Duration(randomer.Intn(30)) * time.Second)
	}
}

// work out changes (compared to this particular watcher)
func sendDiff(serviceAddresses []string, tw *targetWatcher) {
	// take care of NEW services
	var ups []*naming.Update
	for _, sa := range serviceAddresses {
		if !isInArray(tw.currentTargets, sa) {
			ups = append(ups, &naming.Update{naming.Add, sa, ""})
			if *dialer_debug {
				fmt.Printf("notifying watcher of new target %s\n", sa)
				tw.currentTargets = append(tw.currentTargets, sa)
			}

		}
	}

	// take care of REMOVED services
	for i, sa := range tw.currentTargets {
		if !isInArray(serviceAddresses, sa) {
			ups = append(ups, &naming.Update{naming.Delete, sa, ""})
			if *dialer_debug {
				fmt.Printf("notifying watcher of REMOVED target %s\n", sa)
				tw.currentTargets = append(tw.currentTargets[:i], tw.currentTargets[i+1:]...)
			}

		}
	}

	// this should be truly async I guess
	tw.updates <- ups

}

func isInArray(s []string, find string) bool {
	for _, x := range s {
		if x == find {
			return true
		}
	}
	return false
}

func (tn *targetNotifier) AddWatcher(watcher *targetWatcher) {
	tn.lock.Lock()
	tn.watchers = append(tn.watchers, watcher)
	tn.lock.Unlock()
}
