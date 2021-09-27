package newrelicmodule

import (
	"net/http"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

var app *newrelic.Application
var Sugar *zap.SugaredLogger

func getEnv(key string) string {
	var env string
	var ok bool
	if env, ok = os.LookupEnv(key); !ok {
		return ""
	}
	return env

}

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	Sugar = logger.Sugar()
}

func createApplication() error {

	var err error

	app, err = newrelic.NewApplication(
		newrelic.ConfigAppName(getEnv("NEW_RELIC_APP_NAME")),
		newrelic.ConfigLicense(getEnv("NEW_RELIC_KEY")),
		newrelic.ConfigDebugLogger(os.Stdout),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if nil != err {
		Sugar.Infof("Unable to create New Relic application..... Error %+v\n", err)
	}

	return err
}

func ProcessExternalSegment(externalSegmentRequestCh chan ExternalSegment, externalSegmentResponseCh chan ExternalSegment) {
	err := createApplication()

	if err != nil {
		return
	}

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
