package newrelicmodule

import (
	"fmt"
	"net/http"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var app *newrelic.Application

func getEnv(key string) string {
	var env string
	var ok bool
	if env, ok = os.LookupEnv(key); !ok {
		return ""
	}
	return env

}

func createApplication() {

	if app != nil {
		return
	}

	var err error

	app, err = newrelic.NewApplication(
		newrelic.ConfigAppName(getEnv("NEW_RELIC_APP_NAME")),
		newrelic.ConfigLicense(getEnv("NEW_RELIC_KEY")),
		newrelic.ConfigDebugLogger(os.Stdout),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if nil != err {
		fmt.Printf("Unable to create New Relic application..... Error %+v\n", err)
		os.Exit(1)
	}

}

func ProcessExternalSegment(externalSegmentRequestCh chan ExternalSegment, externalSegmentResponseCh chan ExternalSegment) {

	createApplication()

	for {
		msg := <-externalSegmentRequestCh
		txn := app.StartTransaction(msg.TransactionName)

		client := &http.Client{}

		seg := newrelic.StartExternalSegment(txn, msg.Request)

		resp, err := client.Do(msg.Request)

		if err != nil {
			txn.NoticeError(err)
		}

		externalSegmentResponseCh <- ExternalSegment{TransactionName: msg.TransactionName, Request: msg.Request, Response: resp}

		seg.End()
		txn.End()
	}
}

func LogError(errChan chan ErrorLog) {

	createApplication()

	for {
		errLog := <-errChan

		txn := app.StartTransaction(errLog.TransactionName)

		txn.NoticeError(errLog.Error)

		txn.End()
	}
}
