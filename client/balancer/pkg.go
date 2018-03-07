package balancer

import (
    "fmt"
    //
    "google.golang.org/grpc"
    //
    "github.com/GuruSystems/framework/client"
)

// opens a tcp connection to a gurupath.
func DialService(serviceName string) (*grpc.ClientConn, error) {

    creds := client.GetClientCreds()

    conn, err := grpc.Dial(
        serviceName,
        grpc.WithBlock(),
        grpc.WithBalancer(
            grpc.RoundRobin(&RegistryResolver{}),
        ),
        grpc.WithTransportCredentials(creds),
    )
    if err != nil {
        return nil, err
    }

    fmt.Println("STARTED SERVICE CLIENT: "+serviceName)
    return conn, nil
}
