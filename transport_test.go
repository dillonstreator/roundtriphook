package roundtriphook

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTransport struct {
	mock.Mock
}

var _ http.RoundTripper = (*mockTransport)(nil)

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestNewTransport(t *testing.T) {
	assert := assert.New(t)

	tpt := NewTransport()

	assert.Equal(http.DefaultTransport, tpt.base)
}

func TestNewTransport_Options(t *testing.T) {
	assert := assert.New(t)

	baseTpt := &http.Transport{}

	var before1 BeforeFn = func(req *http.Request) *http.Request { return req }
	var before2 BeforeFn = func(req *http.Request) *http.Request { return req }
	var after1 AfterFn = func(req *http.Request, res *http.Response, err error) {}
	var after2 AfterFn = func(req *http.Request, res *http.Response, err error) {}

	tpt := NewTransport(
		WithBaseRoundTripper(baseTpt),
		WithBefore(before1),
		WithBefore(before2),
		WithAfter(after1),
		WithAfter(after2),
	)

	assert.Equal(baseTpt, tpt.base)
	assert.Len(tpt.befores, 2)
	assert.Len(tpt.afters, 2)
}

func TestRoundTrip(t *testing.T) {
	assert := assert.New(t)

	mockTpt := &mockTransport{}
	req1 := &http.Request{
		Header: http.Header{
			"header1": []string{"value1"},
		},
	}
	req2 := &http.Request{
		Header: http.Header{
			"header1": []string{"value1"},
			"header2": []string{"value2"},
		},
	}
	expectedRes := &http.Response{}

	mockTpt.On("RoundTrip", req2).Return(expectedRes, nil).Once()

	before1CalledAt := time.Time{}
	before2CalledAt := time.Time{}
	after1CalledAt := time.Time{}
	after2CalledAt := time.Time{}
	before1 := func(req *http.Request) *http.Request {
		time.Sleep(time.Millisecond)
		before1CalledAt = time.Now()
		assert.Equal(req1, req)
		return req2
	}
	before2 := func(req *http.Request) *http.Request {
		time.Sleep(time.Millisecond)
		before2CalledAt = time.Now()
		assert.Equal(req2, req)
		return req2
	}
	after1 := func(req *http.Request, res *http.Response, err error) {
		time.Sleep(time.Millisecond)
		after1CalledAt = time.Now()
		assert.Equal(req2, req)
		assert.Equal(expectedRes, res)
	}
	after2 := func(req *http.Request, res *http.Response, err error) {
		time.Sleep(time.Millisecond)
		after2CalledAt = time.Now()
		assert.Equal(req2, req)
		assert.Equal(expectedRes, res)
	}
	tpt := NewTransport(
		WithBaseRoundTripper(mockTpt),
		WithBefore(before1),
		WithBefore(before2),
		WithAfter(after1),
		WithAfter(after2),
	)

	res, err := tpt.RoundTrip(req1)
	assert.NoError(err)

	assert.Equal(expectedRes, res)
	assert.NotZero(before1CalledAt)
	assert.NotZero(before2CalledAt)
	assert.NotZero(after1CalledAt)
	assert.NotZero(after2CalledAt)
	assert.True(before1CalledAt.Before(before2CalledAt))
	assert.True(before2CalledAt.Before(after1CalledAt))
	assert.True(after1CalledAt.Before(after2CalledAt))
}
