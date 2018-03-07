package balancer

import (
    "gitlab.gurusys.co.uk/guru/proto/slackgateway"
)

// Initialises a client
func SlackGatewayClient() (slackgateway.SlackGatewayClient, error) {

    conn, err := DialService("slackgateway.SlackGateway")
    if err != nil {
        return nil, err
    }

    return slackgateway.NewSlackGatewayClient(conn), nil
}
