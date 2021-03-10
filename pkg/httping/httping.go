package httping

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"time"
)

// Absolute timings since start of request
type HTTPRoundTripTimings struct {
	RequestStart time.Time

	DNSStart time.Time
	DNSDone  time.Time

	ConnectStart time.Time
	ConnectDone  time.Time

	TLSStart time.Time
	TLSDone  time.Time

	GotFirstResponseByte time.Time
	ReadHeaders          time.Time
	ReadBody             time.Time

	RequestDone time.Time
}

type HTTPRoundTripDurations struct {
	// duration of round-trip parts
	Resolve time.Duration
	Connect time.Duration
	TLS     time.Duration

	// since FirstResponseByte
	Headers time.Duration
	// since ReadHeaders
	Body time.Duration

	// since start of request, for ease
	FirstResponseByte time.Duration
	Total             time.Duration
}

func subIfNotZero(a, b time.Time) time.Duration {
	if a.IsZero() {
		return time.Duration(0)
	} else {
		return a.Sub(b)
	}
}

func (ht HTTPRoundTripTimings) Durations() HTTPRoundTripDurations {
	return HTTPRoundTripDurations{
		Resolve:           subIfNotZero(ht.DNSDone, ht.DNSStart),
		Connect:           subIfNotZero(ht.ConnectDone, ht.ConnectStart),
		TLS:               subIfNotZero(ht.TLSDone, ht.TLSStart),
		FirstResponseByte: subIfNotZero(ht.GotFirstResponseByte, ht.RequestStart),
		Headers:           subIfNotZero(ht.ReadHeaders, ht.GotFirstResponseByte),
		Body:              subIfNotZero(ht.ReadBody, ht.ReadHeaders),
		Total:             subIfNotZero(ht.RequestDone, ht.RequestStart),
	}
}

func GetRoundTripTimings(c *http.Client, req *http.Request, readBody bool) (HTTPRoundTripTimings, *http.Response, []byte, error) {
	timings := HTTPRoundTripTimings{}
	if c.Transport == nil {
		return timings, nil, nil, errors.New("http.Client.Transport must not be 'nil'")
	}

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { timings.DNSStart = time.Now() },
		DNSDone:  func(ddi httptrace.DNSDoneInfo) { timings.DNSDone = time.Now() },

		TLSHandshakeStart: func() { timings.TLSStart = time.Now() },
		TLSHandshakeDone:  func(cs tls.ConnectionState, err error) { timings.TLSDone = time.Now() },

		ConnectStart: func(network, addr string) { timings.ConnectStart = time.Now() },
		ConnectDone:  func(network, addr string, err error) { timings.ConnectDone = time.Now() },

		GotFirstResponseByte: func() { timings.GotFirstResponseByte = time.Now() },
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	timings.RequestStart = time.Now()

	resp, err := c.Transport.RoundTrip(req)
	if err != nil {
		return timings, resp, nil, fmt.Errorf("round trip failed: %w", err)
	}
	defer resp.Body.Close()
	timings.ReadHeaders = time.Now()

	var body []byte
	if readBody {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return timings, resp, body, fmt.Errorf("read body failed: %w", err)
		}
		timings.ReadBody = time.Now()
	}

	timings.RequestDone = time.Now()
	return timings, resp, body, err
}
