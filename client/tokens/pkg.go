package tokens

import (
	"os"
	"fmt"
	"time"
	"errors"
	"strings"
	"os/user"
	"io/ioutil"
	//
	"google.golang.org/grpc/metadata"
	//
	"golang.org/x/net/context"
)

var (
	displayedTokenInfo = false
)

func SaveToken(tk string) error {

	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Unable to get current user: %s\n", err)
		return err
	}
	cfgdir := fmt.Sprintf("%s/.picoservices/tokens", usr.HomeDir)
	fname := fmt.Sprintf("%s/%s", cfgdir, tk)
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

func GetToken(token string) string {
	var tok string
	var btok []byte
	var fname string
	fname = "n/a"
	usr, err := user.Current()
	if err == nil {
		fname = fmt.Sprintf("%s/.picoservices/tokens/%s", usr.HomeDir, token)
		btok, _ = ioutil.ReadFile(fname)
	}
	if (err != nil) || (len(btok) == 0) {
		tok = token
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

func SetAuthToken(token string) context.Context {
	tok := GetToken(token)
	md := metadata.Pairs("token", tok,
		"clid", "itsme",
	)
	millis := 5000
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(millis)*time.Millisecond)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}
