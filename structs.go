package newrelicmodule

import "net/http"

type ExternalSegment struct {
	TransactionName string
	Request         *http.Request
	Response        *http.Response
}

type ErrorLog struct {
	TransactionName string
	Error           error
}
