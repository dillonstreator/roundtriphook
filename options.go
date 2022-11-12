package roundtriphook

import (
	"net/http"
)

type option func(t *transport)

func WithBaseRoundTripper(base http.RoundTripper) option {
	return func(t *transport) {
		t.base = base
	}
}

func WithBefore(fn ...BeforeFn) option {
	return func(t *transport) {
		t.befores = append(t.befores, fn...)
	}
}

func WithAfter(fn ...AfterFn) option {
	return func(t *transport) {
		t.afters = append(t.afters, fn...)
	}
}
