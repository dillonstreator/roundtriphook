package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dillonstreator/roundtriphook"
)

type rthCtxKey string

var (
	timeStartKey = rthCtxKey("startTime")
	idKey        = rthCtxKey("id")
)

var loggingTransport = roundtriphook.NewTransport(
	// This call to roundtriphook.WithBaseRoundTripper is unnecessary
	// since the default behavior is to set the base round tripper to http.DefaultTransport if none is provided
	roundtriphook.WithBaseRoundTripper(http.DefaultTransport),
	roundtriphook.WithBefore(func(req *http.Request) *http.Request {
		startTime := time.Now()
		id := startTime.UnixNano()

		fmt.Printf("[%d] -> %s %s\n", id, req.Method, req.URL)

		ctx := req.Context()
		ctx = context.WithValue(ctx, timeStartKey, startTime)
		ctx = context.WithValue(ctx, idKey, id)

		return req.WithContext(ctx)
	}),
	roundtriphook.WithAfter(func(req *http.Request, res *http.Response, err error) {
		startTime := req.Context().Value(timeStartKey).(time.Time)
		id := req.Context().Value(idKey).(int64)

		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("[%d] <- %s %s %s", id, req.Method, req.URL, time.Since(startTime)))

		if res != nil {
			sb.WriteString(" " + res.Status)
		}

		if err != nil {
			sb.WriteString(fmt.Sprintf(" %s", err.Error()))
		}

		fmt.Printf("%s\n", sb.String())
	}),
)

func main() {
	httpClient := &http.Client{
		Transport: loggingTransport,
	}

	httpClient.Get("https://www.google.com")
}
