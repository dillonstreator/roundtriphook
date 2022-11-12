package roundtriphook

import (
	"net/http"
)

type BeforeFn func(req *http.Request) *http.Request
type AfterFn func(req *http.Request, res *http.Response, err error)

type transport struct {
	base    http.RoundTripper
	befores []BeforeFn
	afters  []AfterFn
}

var _ http.RoundTripper = (*transport)(nil)

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, fn := range t.befores {
		req = fn(req)
	}

	res, err := t.base.RoundTrip(req)

	for _, fn := range t.afters {
		fn(req, res, err)
	}

	return res, err
}

func NewTransport(opts ...option) *transport {
	t := &transport{}
	for _, opt := range opts {
		opt(t)
	}

	if t.base == nil {
		t.base = http.DefaultTransport
	}

	return t
}
