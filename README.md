# Newrelic Go Module

## Requirements
- go get github.com/CemBdc/newrelicmodule
- NEW_RELIC_APP_NAME, NEW_RELIC_KEY env variables

## ExternalSegment Instruments

StartExternalSegment starts the instrumentation of an external call and adds distributed tracing headers to the request. If the Transaction parameter is nil then StartExternalSegment will look for a Transaction in the request's context using FromContext.

When you monitor external http requests and responses, you should define two channels with ExternalSegment type. One of them is for send http requests(ExternalSegment.Request) and the other one is for http responses(ExternalSegment.Response).

## NoticeError Instruments

If you want to log only errors occured in your application, you may use LogError method. Before using this method you should create a channel type 'ErrorLog' and feed it when you get the errors.

## Sample Usage - (main.go)

```go
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/CemBdc/newrelicmodule"
)

var reqChan = make(chan newrelicmodule.ExternalSegment)
var resChan = make(chan newrelicmodule.ExternalSegment)
var errChan = make(chan newrelicmodule.ErrorLog)

func main() {

	go startExternalSegment()
	go startErrorLog()

	doRequest()
}

func startExternalSegment() {

	go newrelicmodule.ProcessExternalSegment(reqChan, resChan)
}

func startErrorLog() {

	go newrelicmodule.LogError(errChan)
}

func doRequest() {

	var jsonStr = []byte(`{"title":"Hello", subject:"World!"}`)
	req, err := http.NewRequest("POST", "http://example.com", bytes.NewBuffer(jsonStr))

	if err != nil {

		errLog := newrelicmodule.ErrorLog{TransactionName: "SampleTrxName", Error: err}
		errChan <- errLog
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	msg := newrelicmodule.ExternalSegment{TransactionName: "SampleTrxName", Request: req}
	reqChan <- msg

	resp := <-resChan

	fmt.Println("response code is ", resp.Response.StatusCode)

	body, _ := ioutil.ReadAll(resp.Response.Body)

	fmt.Println("response Body:", string(body))
}


```