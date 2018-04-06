package logger

import (
	"fmt"
	"flag"
	"sync"
	"time"
	"errors"
	//
	"golang.gurusys.co.uk/go-framework/client/tokens"
	"golang.gurusys.co.uk/guru/proto/client"
	"golang.gurusys.co.uk/guru/proto/logservice"
)

var (
	log_debug = flag.Bool("logger_debug", false, "set to true to debug logging")
)

type QueueEntry struct {
	sent    bool
	created int64
	line    string
}

type AsyncLogQueue struct {
	grpcClient     logservice.LogServiceClient
	appDef         *logservice.LogAppDef
	entries        []*QueueEntry
	lastErrPrinted time.Time
	MaxSize        int
	sync.Mutex
}

func NewAsyncLogQueue(appname, repo, group, namespace, deplid string) (*AsyncLogQueue, error) {

	grpcClient, err := protoclient.LogServiceClient()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Logqueue flush error: %s", err))
	}

	alq := &AsyncLogQueue{
		grpcClient: grpcClient,
		appDef: &logservice.LogAppDef{
			Appname: appname,
			Repository: repo,
			Groupname: group,
			Namespace: namespace,
			DeploymentID: deplid,
		},
		MaxSize: 5000,
	}

	t := time.NewTicker(2 * time.Second)

	go func(a *AsyncLogQueue) {
		for _ = range t.C {
			err := a.Flush()
			if (*log_debug) && (err != nil) {
				fmt.Printf("Error flushing logqueue:%s\n", err)
			}
		}
	}(alq)

	return alq, nil
}

func (alq *AsyncLogQueue) LogCommandStdout(line string, status string) error {

	alq.Lock()
	defer alq.Unlock()

	qe := QueueEntry{
		sent: false,
		created: time.Now().Unix(),
		line:    line,
	}

	if len(alq.entries) > alq.MaxSize {
		if *log_debug {
			fmt.Printf("queue size larger than %d (it is %d) - discarding log entries\n", alq.MaxSize, len(alq.entries))
		}
		alq.entries = alq.entries[0:]
	}

	alq.entries = append(alq.entries, &qe)

	return nil
}

func (alq *AsyncLogQueue) Flush() error {

	if *log_debug {
		fmt.Printf("Logqueue flush\n")
	}

	alq.Lock()
	defer alq.Unlock()

	if len(alq.entries) == 0 {
		// save ourselves from dialing and stuff
		return nil
	}

	logRequest := &logservice.LogRequest{
		AppDef: alq.appDef,
	}

	for _, qe := range alq.entries {
		logRequest.Lines = append(
			logRequest.Lines,
			&logservice.LogLine{
				Time: qe.created,
				Line: qe.line,
			},
		)
	}

	_, err := alq.grpcClient.LogCommandStdout(tokens.SetAuthToken(), logRequest)
	if err != nil {
		if time.Since(alq.lastErrPrinted) > (10 * time.Second) {
			fmt.Printf("Failed to send log: %s\n", err)
			alq.lastErrPrinted = time.Now()
		}
	}

	// all done, so clear the array so we free up the memory
	alq.entries = alq.entries[:0]

	return nil
}
