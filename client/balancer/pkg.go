package balancer

import (
    "fmt"
    "sync"
    //
    "golang.org/x/net/context"
    "google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
    //
    "github.com/GuruSystems/framework/client"
    "github.com/GuruSystems/framework/client/registry"
    "gitlab.gurusys.co.uk/guru/proto/slackgateway"
	pb "github.com/GuruSystems/framework/proto/registrar"
)

type ServiceManager struct {
    waitGroup *sync.WaitGroup
    grpcClient interface{}
    connections map[*Host]struct{}
    sync.RWMutex
}

type Host struct {
    manager *ServiceManager
    waitGroup *sync.WaitGroup
    address string
    clientConn *grpc.ClientConn
    killChannel chan struct{}
    sync.RWMutex
}

// opens a tcp connection to a gurupath.
func NewServiceManager(serviceName string) (*ServiceManager, error) {

    manager := &ServiceManager{
        waitGroup: &sync.WaitGroup{},
        connections: map[*Host]struct{}{},
    }

	serverAddresses, err := registry.GetHosts(serviceName, pb.Apitype_tcp)
    if err != nil {
        return nil, err
    }

    for _, address := range serverAddresses {

        host := &Host{
            manager: manager,
            waitGroup: &sync.WaitGroup{},
            address: address,
            killChannel: make(chan struct{}),
        }

        manager.Lock()
        manager.connections[host] = struct{}{}
        manager.Unlock()

        go host.Connection()

    }

    return manager, nil
}

func (manager *ServiceManager) SlackGatewayClient() *slackgateway.SlackGatewayClient {

    manager.waitGroup.Wait()
    return manager.grpcClient.(*slackgateway.SlackGatewayClient)
}

func (manager *ServiceManager) RemoveHost(host *Host) {
    manager.Lock()
    delete(manager.connections, host)
    manager.Unlock()
}

func (host *Host) Remove() {
    host.manager.RemoveHost(host)
}

func (host *Host) Connection() {

    var s connectivity.State
    creds := client.GetClientCreds()

    for {
        select {
            case <- host.killChannel:

                fmt.Println("STOPPING GOROUTINE FOR HOST: "+host.address)
                return

            default:
        }

        host.waitGroup.Add(1)

        clientConn, err := grpc.Dial(
            host.address,
            grpc.WithTransportCredentials(creds),
        )
        if err != nil {
            fmt.Println(
                fmt.Errorf("Error dialling servicename @ %s\n", host.address),
            )
            continue
        }

        host.waitGroup.Done()

    	for clientConn != nil {

            host.Lock()
            host.clientConn = clientConn
            host.Unlock()

    		clientConn.WaitForStateChange(context.Background(), s)

			s = clientConn.GetState()
            fmt.Println("state has changed", s)

			if s == connectivity.TransientFailure || s == connectivity.Shutdown {
				break
			}

        }

        host.waitGroup.Add(1)
    }

}
