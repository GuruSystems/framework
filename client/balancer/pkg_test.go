package balancer

import (
    "testing"
)

func TestNewServiceManager(t *testing.T) {

    service, err := NewServiceManager("slackgateway.SlackGateway")
    if err != nil {
        t.Error(err)
    }

    req := &pb.PublishMessageRequest{
        Channel: "bot",
        Text: *flag_message + " @AlexB " + strconv.Itoa(x),
    }

    resp, err := service.SlackGatewayClient().PublishMessage(client.SetAuthToken(), req)
    if err != nil {
        fmt.Println("FAILED TO PUBLISH", req, err)
        continue
    }

}
