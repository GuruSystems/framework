package balancer

import (
    "fmt"
    "flag"
    "strconv"
    "testing"
    //
	"github.com/GuruSystems/framework/client"
    slackgateway "gitlab.gurusys.co.uk/guru/proto/slackgateway"
)

func TestNewServiceManager(t *testing.T) {

    flag.Parse()

    service, err := SlackGatewayClient()
    if err != nil {
        t.Error(err)
        return
    }

    for x := 0; x < 10; x ++ {

        req := &slackgateway.PublishMessageRequest{
            Channel: "bot",
            Text: "framework/client/balancer "+strconv.Itoa(x),
        }

        _, err = service.PublishMessage(client.SetAuthToken(), req)
        if err != nil {
            fmt.Println("FAILED TO PUBLISH", req, err)
            t.Error(err)
            return
        }

    }
}
