package gohttp

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"time"
)

type transportOptionHolder struct {
	proxy      func(r *http.Request) (*url.URL, error)
	tlsConfig  *tls.Config
	retryCount int
}

type TransportOption interface {
	apply(transport transportOptionHolder) transportOptionHolder
}

type transportOptionFunc func(transport transportOptionHolder) transportOptionHolder

func (f transportOptionFunc) apply(transport transportOptionHolder) transportOptionHolder {
	return f(transport)
}

func WithProxy(proxyURL string) TransportOption {
	parsedURL, _ := url.Parse(proxyURL)
	proxy := http.ProxyURL(parsedURL)

	return transportOptionFunc(func(transport transportOptionHolder) transportOptionHolder {
		transport.proxy = proxy
		return transport
	})
}

func WithTLS(config *tls.Config) TransportOption {
	return transportOptionFunc(func(transport transportOptionHolder) transportOptionHolder {
		transport.tlsConfig = config
		return transport
	})
}

func WithRetries(retryCount int) TransportOption {
	return transportOptionFunc(func(transport transportOptionHolder) transportOptionHolder {
		transport.retryCount = retryCount
		return transport
	})
}

type Transport struct {
	roundTripper http.RoundTripper
	retryCount   int
}

func NewTransport(options ...TransportOption) *Transport {
	var holder transportOptionHolder

	for _, opt := range options {
		holder = opt.apply(holder)
	}

	baseTransport := http.DefaultTransport.(*http.Transport).Clone()
	if holder.proxy != nil {
		baseTransport.Proxy = holder.proxy
	}
	if holder.tlsConfig != nil {
		baseTransport.TLSClientConfig = holder.tlsConfig
	}

	return &Transport{
		roundTripper: baseTransport,
		retryCount:   holder.retryCount,
	}
}

func (t *Transport) RoundTrip(request *http.Request) (*http.Response, error) {
	if t.retryCount == 0 {
		return t.roundTripper.RoundTrip(request)
	}

	var (
		err       error
		bodyBytes []byte
	)

	if request.Body != nil {
		bodyBytes, err = io.ReadAll(request.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading request body: %v", err)
		}

		request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	response, err := t.roundTripper.RoundTrip(request)

	retries := 0
	for shouldRetry(err, response) && retries < t.retryCount {
		time.Sleep(backoff(retries))

		err = drainBody(response)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		if request.Body != nil {
			request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		response, err = t.roundTripper.RoundTrip(request)
		retries++
	}

	return response, err
}

func backoff(retries int) time.Duration {
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}

func shouldRetry(err error, resp *http.Response) bool {
	if err != nil {
		return true
	}

	if resp != nil &&
		(resp.StatusCode == http.StatusBadGateway ||
			resp.StatusCode == http.StatusServiceUnavailable ||
			resp.StatusCode == http.StatusGatewayTimeout) {
		return true
	}

	return false
}

func drainBody(resp *http.Response) (err error) {
	if resp != nil && resp.Body != nil {
		_, err = io.Copy(io.Discard, resp.Body)
		defer func() {
			if tempErr := resp.Body.Close(); tempErr != nil {
				err = errors.Join(err, fmt.Errorf("failed to close response body: %w", tempErr))
			}
		}()
		if err != nil {
			err = fmt.Errorf("failed to drain response body: %w", err)
			return
		}
	}

	return
}
