package client

import (
	"os"
	"fmt"
	"net"
	"sync"
	"time"
	"flag"
	"errors"
	"os/user"
	"strings"
	"io/ioutil"
	"crypto/tls"
	"crypto/x509"
	//
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	//
	"github.com/GuruSystems/framework/cmdline"
	pb "github.com/GuruSystems/framework/proto/registrar"
)

var (
	cert               = []byte{1, 2, 3}
	displayedTokenInfo = false

	token     = flag.String("token", "user_token", "The authentication token (cookie) to authenticate with. May be name of a file in ~/.picoservices/tokens/, if so file contents shall be used as cookie")
	errorList []*errorCache
	errorLock sync.Mutex
)

type errorCache struct {
	servicename string
	lastOccured time.Time
	lastPrinted time.Time
}

func SaveToken(tk string) error {

	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Unable to get current user: %s\n", err)
		return err
	}
	cfgdir := fmt.Sprintf("%s/.picoservices/tokens", usr.HomeDir)
	fname := fmt.Sprintf("%s/%s", cfgdir, *token)
	if _, err := os.Stat(fname); !os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("File %s exists already", fname))
	}
	os.MkdirAll(cfgdir, 0700)
	fmt.Printf("Saving new token to %s\n", fname)
	err = ioutil.WriteFile(fname, []byte(tk), 0600)
	if err != nil {
		fmt.Printf("Failed to save token to %s: %s\n", fname, err)
	}
	return err
}

// opens a tcp connection to a gurupath.
func DialTCPWrapper(serviceName string) (net.Conn, error) {

	serverAddr, err := ServiceHost(serviceName, pb.Apitype_tcp)
    if err != nil {
        return nil, err
    }

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
        return nil, err
    }

	return conn, err
}

func ServiceHost(serviceName string, apiType pb.Apitype) (string, error) {

    list, err := ServiceList(serviceName, apiType)
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

func ServiceList(serviceName string, apiType pb.Apitype) ([]*pb.GetResponse, error) {

	registryAddress := cmdline.GetRegistryAddress()

    conn, err := grpc.Dial(
        registryAddress,
        getDialopts()...,
    )
    if err != nil {
        return nil, fmt.Errorf("Error dialling servicename %s @ %s\n", serviceName, registryAddress)
    }
    defer conn.Close()

    grpcRegistryClient := pb.NewRegistryClient(conn)
    list, err := grpcRegistryClient.GetTarget(
        context.Background(),
        &pb.GetTargetRequest{
            Name: serviceName,
            ApiType: apiType,
        },
    )
    if err != nil {
        return nil, fmt.Errorf("Error getting grpc service address %s: %s\n", serviceName, err)
    }

    return list.Service, nil
}

// given a service name we look up its address in the registry
// and return a connection to it.
// it's a replacement for the normal "dial" but instead of an address
// it takes a service name
func DialWrapper(serviceName string) (*grpc.ClientConn, error) {

    serverAddr, err := ServiceHost(serviceName, pb.Apitype_grpc)
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
	FrontendCert := Certificate //ioutil.ReadFile(*clientcrt)
	roots.AppendCertsFromPEM(FrontendCert)
	ImCert := Ca //ioutil.ReadFile(*clientca)
	roots.AppendCertsFromPEM(ImCert)
	cert, err := tls.X509KeyPair(Certificate, Privatekey)
	//	cert, err := tls.LoadX509KeyPair(*clientcrt, *clientkey)
	if err != nil {
		fmt.Printf("Failed to create client certificates: %s\n", err)
		fmt.Printf("key:\n%s\n", string(Privatekey))
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
func GetToken() string {
	var tok string
	var btok []byte
	var fname string
	fname = "n/a"
	usr, err := user.Current()
	if err == nil {
		fname = fmt.Sprintf("%s/.picoservices/tokens/%s", usr.HomeDir, *token)
		btok, _ = ioutil.ReadFile(fname)
	}
	if (err != nil) || (len(btok) == 0) {
		tok = *token
	} else {
		tok = string(btok)
		if displayedTokenInfo {
			fmt.Printf("Using token from %s\n", fname)
			displayedTokenInfo = true
		}
	}
	tok = strings.TrimSpace(tok)

	return tok
}

func SetAuthToken() context.Context {
	tok := GetToken()
	md := metadata.Pairs("token", tok,
		"clid", "itsme",
	)
	millis := 5000
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(millis)*time.Millisecond)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
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
	deadline := 2
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
