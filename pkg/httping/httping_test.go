package httping

import (
	"context"
	"log"
	"net/http"
	"time"
)

func ExampleGetRoundTripTimings() {
	req, _ := http.NewRequest("GET", "https://ipv4.google.com", nil)

	// It's the only way to set global deadline, including all redirects, upgrades, etc
	deadline := 200 * time.Millisecond
	ctx, cancel := context.WithCancel(context.TODO())
	_ = time.AfterFunc(deadline, func() {
		cancel()
	})
	req = req.WithContext(ctx)

	client := &http.Client{Transport: http.DefaultTransport}

	// Granular timeouts example
	// Total request timeout will not be limited by them, use WithCancel context instead!
	/*&http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 150 * time.Millisecond,
			}).DialContext,

			TLSHandshakeTimeout:   300 * time.Millisecond,
			IdleConnTimeout:       300 * time.Millisecond,
			ResponseHeaderTimeout: 300 * time.Millisecond,
			ExpectContinueTimeout: 300 * time.Millisecond,

			ForceAttemptHTTP2: true,
		},
		Timeout: 300 * time.Millisecond,
	}*/

	ht, _, _, err := GetRoundTripTimings(client, req, false)
	if err != nil {
		log.Printf("ExampleGetRoundTripTimings [durations=%+v] http error: %s", ht.Durations(), err)
		return
	}
	log.Printf("ExampleGetRoundTripTimings [durations=%+v]", ht.Durations())

	// Output:
}
