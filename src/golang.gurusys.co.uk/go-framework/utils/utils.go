package utils

import (
	"fmt"
	"os"
)

// never returns - if error != nil, print it and exit
func Bail(txt string, err error) {
	if err == nil {
		return
	}
	fmt.Printf("%s: %s", txt, err)
	os.Exit(10)
}
