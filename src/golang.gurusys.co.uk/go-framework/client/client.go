package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	pb "golang.gurusys.co.uk/apis/registry"
	"golang.gurusys.co.uk/go-framework/certificates"
	"golang.gurusys.co.uk/go-framework/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var (
	cert      = []byte{1, 2, 3}
	errorList []*errorCache
	errorLock sync.Mutex
)

type errorCache struct {
	servicename string
	lastOccured time.Time
	lastPrinted time.Time
}

// opens a tcp connection to a gurupath.
func DialTCPWrapper(serviceName string) (net.Conn, error) {
	if strings.Contains(serviceName, "/") {
		s := fmt.Sprintf("Error: The parameter for DialTCPWrapper needs a servicename. not a gurupath. You passed in %s, which looks very much like a gurupath. The \"old-style\" picoservices required a gurupath at this function, but go-framework does not. Did you recently upgrade and did not upgrade a config?\n", serviceName)
		debug.PrintStack()
		return nil, errors.New(s)
	}
	serverAddr, err := registry.GetHost(serviceName, pb.Apitype_tcp)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}

	return conn, err
}

// given a service name we look up its address in the registry
// and return a connection to it.
// it's a replacement for the normal "dial" but instead of an address
// it takes a service name
func DialWrapper(serviceName string) (*grpc.ClientConn, error) {

	serverAddr, err := registry.GetHost(serviceName, pb.Apitype_grpc)
	if err != nil {
		return nil, err
	}

	creds := GetClientCreds()
	conn, err := grpc.Dial(
		serverAddr,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		return nil, fmt.Errorf("Error dialling servicename @ %s\n", serverAddr)
	}

	return conn, nil
}

func hasApi(ar []pb.Apitype, lf pb.Apitype) bool {
	for _, a := range ar {
		if a == lf {
			return true
		}
	}
	return false
}

// get the Client Credentials we use to connect to other RPCs
func GetClientCreds() credentials.TransportCredentials {
	roots := x509.NewCertPool()

	frontendCert := certificates.Certificate()

	roots.AppendCertsFromPEM(frontendCert)
	ImCert := certificates.Ca() //ioutil.ReadFile(*clientca)
	roots.AppendCertsFromPEM(ImCert)

	pk := certificates.Privatekey()

	cert, err := tls.X509KeyPair(frontendCert, pk)
	//	cert, err := tls.LoadX509KeyPair(*clientcrt, *clientkey)
	if err != nil {
		fmt.Printf("Failed to create client certificates: %s\n", err)
		fmt.Printf("key:\n%s\n", string(pk))
		return nil
	}
	// we don't verify the hostname because we use a dynamic registry thingie
	creds := credentials.NewTLS(&tls.Config{
		ServerName:         "*",
		Certificates:       []tls.Certificate{cert},
		RootCAs:            roots,
		InsecureSkipVerify: true,
	})
	return creds

}

func getErrorCacheByName(name string) *errorCache {
	errorLock.Lock()
	defer errorLock.Unlock()
	for _, ec := range errorList {
		if ec.servicename == name {
			return ec
		}
	}
	ec := &errorCache{servicename: name,
		lastOccured: time.Now(),
	}
	errorList = append(errorList, ec)
	return ec
}

func getDialopts() []grpc.DialOption {
	deadline := 10
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithTimeout(time.Duration(deadline) * time.Second)}
	return opts
}
func printError(path string, msg string) {
	e := getErrorCacheByName(path)
	if e == nil {
		fmt.Println(msg)
		return
	}
	if !e.needsPrinting() {
		return
	}
	fmt.Println(msg)
}

// returns true if this needs printing
// resets counter if it returns true
func (e *errorCache) needsPrinting() bool {
	now := time.Now()
	if now.Sub(e.lastPrinted) < (time.Duration(5) * time.Minute) {
		return false
	}
	e.lastPrinted = now
	return false
}
