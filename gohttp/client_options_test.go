package gohttp

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/sunshineOfficial/golib/golog"
	"github.com/sunshineOfficial/golib/gorand"
)

type clientOptionTestCase struct {
	option ClientOption
	assert func(t *testing.T, h clientOptionHolder)
}

func getClientOptionTestCases() map[string]clientOptionTestCase {
	log := golog.NewLogger("test")
	client := &http.Client{}
	transport := NewTransport()
	timeout := time.Duration(gorand.RandomInt(10, 30)) * time.Second
	before := func(r *http.Request) error { return nil }
	after := func(r *http.Response) error { return nil }

	return map[string]clientOptionTestCase{
		"logger": {
			option: WithLogger(log),
			assert: func(t *testing.T, h clientOptionHolder) {
				assert.Equal(t, log, h.logger)
			},
		},
		"client": {
			option: WithClient(client),
			assert: func(t *testing.T, h clientOptionHolder) {
				assert.Equal(t, *client, *h.client)
			},
		},
		"transport": {
			option: WithTransport(transport),
			assert: func(t *testing.T, h clientOptionHolder) {
				assert.Equal(t, transport, h.transport)
			},
		},
		"timeout": {
			option: WithTimeout(timeout),
			assert: func(t *testing.T, h clientOptionHolder) {
				assert.Equal(t, timeout, h.timeout)
			},
		},
		"before": {
			option: WithBefore(before),
			assert: func(t *testing.T, h clientOptionHolder) {
				assert.Equal(t, fmt.Sprintf("%v", any(before)), fmt.Sprintf("%v", any(h.before)))
			},
		},
		"after": {
			option: WithAfter(after),
			assert: func(t *testing.T, h clientOptionHolder) {
				assert.Equal(t, fmt.Sprintf("%v", any(after)), fmt.Sprintf("%v", any(h.after)))
			},
		},
		"verbose": {
			option: WithLogger(log),
			assert: func(t *testing.T, h clientOptionHolder) {
				assert.True(t, h.verbose)
			},
		},
		"traces": {
			option: WithTraces(),
			assert: func(t *testing.T, h clientOptionHolder) {
				assert.True(t, h.traces)
			},
		},
	}
}

func TestApply(t *testing.T) {
	for name, testCase := range getClientOptionTestCases() {
		t.Run(name, func(t *testing.T) {
			var holder clientOptionHolder
			holder = testCase.option.apply(holder)
			testCase.assert(t, holder)
		})
	}
}
