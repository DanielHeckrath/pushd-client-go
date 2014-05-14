package client

import (
	"github.com/streadway/handy/breaker"
	"net/http"
)

const (
	UNEXPECTED_ERROR = "Unexpected response from pushd service: status=%d, body=%s"
)

type Request struct {
	circuitBreaker breaker.Breaker
	client         http.Client
	endpoint       string
}

// Impact client version one.
type V1 struct {
	request *Request
}

// Client object holding versioned pushd clients.
type Client struct {
	V1 V1
}

func NewHttpClient(endpoint string) *Client {
	// The circuitbreaker defaults to the following (currently not-configurable values)
	// (See https://github.com/streadway/handy/issues/10)
	//  Window = 5 Seconds
	//  CoolDown = 1 Second
	//  MinObservations = 10
	failureRate := 0.5

	// NOTE: A circuitbreaker is _open_, when too many errors occured and _closed_ when operations can be performed.
	// => The breaker will not _open_, if there are less then 10 {MinObservations} values to judge on (even if 100% failed)
	// => If 50% {failureRatio} of all requests in the last 5 seconds {window} fails, the breaker will _open_
	// => After 1 second {coolDown}, the breaker will allow one request again, to check availability

	circuitBreaker := breaker.NewBreaker(failureRate)
	endpoint = parseEndpoint(endpoint)
	next := http.DefaultTransport
	transport := breaker.Transport(circuitBreaker, breaker.DefaultResponseValidator, next)
	client := http.Client{Transport: transport}

	request := &Request{
		circuitBreaker: circuitBreaker,
		endpoint:       endpoint,
		client:         client,
	}

	return &Client{
		V1: V1{request},
	}
}

func parseEndpoint(endpoint string) string {
	return "http://" + endpoint
}
