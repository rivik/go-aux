package httping

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
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

func (ht HTTPRoundTripTimings) Durations() HTTPRoundTripDurations {
	return HTTPRoundTripDurations{
		Resolve:           ht.DNSDone.Sub(ht.DNSStart),
		Connect:           ht.ConnectDone.Sub(ht.ConnectStart),
		TLS:               ht.TLSDone.Sub(ht.TLSStart),
		FirstResponseByte: ht.GotFirstResponseByte.Sub(ht.RequestStart),
		Headers:           ht.ReadHeaders.Sub(ht.GotFirstResponseByte),
		Body:              ht.ReadBody.Sub(ht.ReadHeaders),
		Total:             ht.RequestDone.Sub(ht.RequestStart),
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
		return timings, resp, nil, err
	}
	defer resp.Body.Close()
	timings.ReadHeaders = time.Now()

	var body []byte
	if readBody {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return timings, resp, body, err
		}
		timings.ReadBody = time.Now()
	}

	timings.RequestDone = time.Now()
	return timings, resp, body, err
}
