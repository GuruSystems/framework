package registry

import (
	"fmt"
	pb "golang.gurusys.co.uk/apis/registry"
	"golang.gurusys.co.uk/go-framework/cmdline"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

const (
	CONST_CALL_TIMEOUT = 2
)

// This function gets a host and returns it's host:port address.
func GetHost(serviceName string, apiType pb.Apitype) (string, error) {

	list, err := HostList(serviceName, apiType)
	if err != nil {
		return "", err
	}

	if len(list) == 0 {
		return "", fmt.Errorf("No grpc target found for name %s", serviceName)
	}

	if len(list[0].Location.Address) == 0 {
		return "", fmt.Errorf("No grpc location found for name %s - is it running?", serviceName)
	}

	return fmt.Sprintf(
		"%s:%d",
		list[0].Location.Address[0].Host,
		list[0].Location.Address[0].Port,
	), nil
}

// This function gets the hosts and returns their host:port address.
func GetHosts(serviceName string, apiType pb.Apitype) ([]string, error) {

	list, err := HostList(serviceName, apiType)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("No grpc target found for name %s", serviceName)
	}

	if len(list[0].Location.Address) == 0 {
		return nil, fmt.Errorf("No grpc location found for name %s - is it running?", serviceName)
	}

	addresses := []string{}

	for _, item := range list {
		addresses = append(
			addresses,
			fmt.Sprintf(
				"%s:%d",
				item.Location.Address[0].Host,
				item.Location.Address[0].Port,
			),
		)
	}

	return addresses, nil
}

// This function returns a list of hosts.
func HostList(serviceName string, apiType pb.Apitype) ([]*pb.GetResponse, error) {

	registryAddress := cmdline.GetRegistryAddress()

	conn, err := grpc.Dial(
		registryAddress,
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Duration(CONST_CALL_TIMEOUT)*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("HostList - Error dialling servicename %s @ %s\n", serviceName, registryAddress)
	}
	defer conn.Close()

	list, err := pb.NewRegistryClient(conn).GetTarget(
		context.Background(),
		&pb.GetTargetRequest{
			Name:    serviceName,
			ApiType: apiType,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("HostList - Error getting grpc service address %s: %s\n", serviceName, err)
	}

	return list.Service, nil
}
