package client

import (
	"fmt"
	"golang.gurusys.co.uk/apis/auth"
	"golang.gurusys.co.uk/apis/autodeployer"
	"golang.gurusys.co.uk/apis/buildrepo"
	"golang.gurusys.co.uk/apis/certificatemanager"
	"golang.gurusys.co.uk/apis/deploymonkey"
	"golang.gurusys.co.uk/apis/echoservice"
	"golang.gurusys.co.uk/apis/httpkpi"
	"golang.gurusys.co.uk/apis/hubextractor"
	"golang.gurusys.co.uk/apis/kpitracker"
	"golang.gurusys.co.uk/apis/lbproxy"
	"golang.gurusys.co.uk/apis/logservice"
	"golang.gurusys.co.uk/apis/paypointendpoint"
	"golang.gurusys.co.uk/apis/rfccreator"
	"golang.gurusys.co.uk/apis/sensorapi"
	"golang.gurusys.co.uk/apis/slackgateway"
	"golang.gurusys.co.uk/apis/testservice"
	"golang.gurusys.co.uk/apis/vpnmanager"
	"google.golang.org/grpc"
)

func NewDialer() *Dialer {
	return &Dialer{}
}

// Initialises a client
func (dialer *Dialer) SlackGatewayClient() (slackgateway.SlackGatewayClient, error) {

	err := dialer.DialService("slackgateway.SlackGateway")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return slackgateway.NewSlackGatewayClient(dialer.conn), nil
}

// Initialises a client
func (dialer *Dialer) VPNManagerClient() (vpnmanager.VpnManagerClient, error) {

	err := dialer.DialService("vpnmanager.VPNManager")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return vpnmanager.NewVpnManagerClient(dialer.conn), nil
}

// Initialises a client
func (dialer *Dialer) AutoDeployerClient() (autodeployer.AutoDeployerClient, error) {

	err := dialer.DialService("autodeployer.AutoDeployer")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return autodeployer.NewAutoDeployerClient(dialer.conn), nil
}

// Initialises a client
func (dialer *Dialer) DeployMonkeyClient() (deploymonkey.DeployMonkeyClient, error) {

	err := dialer.DialService("deploymonkey.DeployMonkey")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return deploymonkey.NewDeployMonkeyClient(dialer.conn), nil
}

// Initialises a client
func (dialer *Dialer) LBProxyServiceClient() (lbproxy.LBProxyServiceClient, error) {

	err := dialer.DialService("lbproxy.LBProxyService")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return lbproxy.NewLBProxyServiceClient(dialer.conn), nil
}

// Initialises a client
func (dialer *Dialer) PaypointListenerClient() (paypointendpoint.PaypointListenerClient, error) {

	err := dialer.DialService("paypointendpoint.PaypointListener")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return paypointendpoint.NewPaypointListenerClient(dialer.conn), nil
}

// Initialises a client
func (dialer *Dialer) RFCManagerClient() (rfccreator.RFCManagerClient, error) {

	err := dialer.DialService("rfccreator.RFCManager")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return rfccreator.NewRFCManagerClient(dialer.conn), nil
}

// Initialises a client
func (dialer *Dialer) BuildRepoManagerClient(serverAddr string) (buildrepo.BuildRepoManagerClient, error) {

	conn, err := grpc.Dial(
		serverAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	//dialer.conn = conn

	return buildrepo.NewBuildRepoManagerClient(conn), nil
}

// Initialises a client
func (dialer *Dialer) KPITrackerClient() (kpitracker.KPITrackerClient, error) {

	err := dialer.DialService("kpitracker.KPITracker")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return kpitracker.NewKPITrackerClient(dialer.conn), nil
}

// Initialises a client
func (dialer *Dialer) HTTPKPITrackerClient() (httpkpi.HTTPKPITrackerClient, error) {

	err := dialer.DialService("httpkpi.HTTPKPITracker")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return httpkpi.NewHTTPKPITrackerClient(dialer.conn), nil
}

// Initialises a client
func (dialer *Dialer) LogServiceClient() (logservice.LogServiceClient, error) {

	err := dialer.DialService("logservice.LogService")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return logservice.NewLogServiceClient(dialer.conn), nil
}

// Initialises a client
func (dialer *Dialer) AuthenticationServiceClient() (auth.AuthenticationServiceClient, error) {

	err := dialer.DialService("auth.AuthenticationService")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return auth.NewAuthenticationServiceClient(dialer.conn), nil
}

// Initialises a client
func (dialer *Dialer) CertificateManagerClient() (certificatemanager.CertificateManagerClient, error) {

	err := dialer.DialService("certificatemanager.CertificateManager")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return certificatemanager.NewCertificateManagerClient(dialer.conn), nil
}

func (dialer *Dialer) EchoServiceClient() (echoservice.EchoServiceClient, error) {

	err := dialer.DialService("echoservice.EchoService")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return echoservice.NewEchoServiceClient(dialer.conn), nil
}

func (dialer *Dialer) SensorAPIClient() (sensorapi.SensorStoreServiceClient, error) {

	err := dialer.DialService("sensorapi.SensorStoreService")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return sensorapi.NewSensorStoreServiceClient(dialer.conn), nil
}

// Initialises a client
func (dialer *Dialer) HubExtractorClient() (hubextractor.HubExtractorClient, error) {

	err := dialer.DialService("hubextractor.HubExtractor")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return hubextractor.NewHubExtractorClient(dialer.conn), nil
}

// Initialises a client
func (dialer *Dialer) TestServiceClient() (testservice.TestServiceClient, error) {

	err := dialer.DialService("testservice.TestServiceClient")
	if err != nil {
		return nil, err
	}

	dialer.RLock()
	defer dialer.RUnlock()

	if dialer.conn == nil {
		return nil, fmt.Errorf("DIALER CONNECTION IS NIL")
	}

	return testservice.NewTestServiceClient(dialer.conn), nil
}
