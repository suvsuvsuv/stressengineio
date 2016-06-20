package boomer

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"stressengineio/engineio2"
	"strings"
	"sync"
	"time"
)

type result struct {
	err           error
	statusCode    int
	duration      time.Duration
	contentLength int64
}

type Boomer struct {
	// Request is the request to be made.
	Request *http.Request

	RequestBody string

	// N is the total number of requests to make.
	N int

	// C is the concurrency level, the number of concurrent workers to run.
	C int

	// Timeout in seconds.
	Timeout int

	// Qps is the rate limit.
	Qps int

	// DisableCompression is an option to disable compression in response
	DisableCompression bool

	// DisableKeepAlives is an option to prevents re-use of TCP connections between different HTTP requests
	DisableKeepAlives bool

	EnableEngineIo bool // sunny

	// Output represents the output type. If "csv" is provided, the
	// output will be dumped as a csv stream.
	Output string

	// ProxyAddr is the address of HTTP proxy server in the format on "host:port".
	// Optional.
	ProxyAddr *url.URL
	RawURL    string
	results   chan *result

	clients    []*engineioclient2.Client
	startTimes []time.Time
	durations  []time.Duration
}

var totalRequests int
var requestCountChan chan int
var quit chan bool

// Run makes all the requests, prints the summary. It blocks until
// all work is done.
func (b *Boomer) Run() {
	b.results = make(chan *result, b.N)

	totalRequests = 0
	requestCountChan = make(chan int)
	quit = make(chan bool)

	start := time.Now()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		fmt.Print("quit\n")
		// TODO(jbd): Progress bar should not be finalized.
		newReport(b.N, b.results, b.Output, time.Now().Sub(start)).finalize()

		close(requestCountChan)
		os.Exit(1)
	}()
	go showProgress(b.C) // sunny

	if b.EnableEngineIo {
		go b.readConsole()
		b.clients = make([]*engineioclient2.Client, b.C)
		b.startTimes = make([]time.Time, b.C)
		b.durations = make([]time.Duration, b.C)
	}
	b.runWorkers()

	fmt.Printf("Finished %d requests\n\n", totalRequests)
	close(requestCountChan)

	newReport(b.N, b.results, b.Output, time.Now().Sub(start)).finalize()
	close(b.results)
	//<-quit
}

func showProgress(c int) {
	for s := range requestCountChan {
		totalRequests += s
		if totalRequests%c == 0 {
			fmt.Printf("Completed request %d\n", totalRequests)
		}
	}
}

func (b *Boomer) makeRequest(c *http.Client) {
	s := time.Now()
	var size int64
	var code int

	resp, err := c.Do(cloneRequest(b.Request, b.RequestBody))
	if err == nil {
		size = resp.ContentLength
		code = resp.StatusCode
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	b.results <- &result{
		statusCode:    code,
		duration:      time.Now().Sub(s),
		err:           err,
		contentLength: size,
	}
}

func (b *Boomer) runWorker(n int) {
	var throttle <-chan time.Time
	if b.Qps > 0 {
		throttle = time.Tick(time.Duration(1e6/(b.Qps)) * time.Microsecond)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: b.DisableCompression,
		DisableKeepAlives:  b.DisableKeepAlives,
		// TODO(jbd): Add dial timeout.
		TLSHandshakeTimeout: time.Duration(b.Timeout) * time.Millisecond,
		Proxy:               http.ProxyURL(b.ProxyAddr),
	}
	client := &http.Client{Transport: tr}
	for i := 0; i < n; i++ {
		if b.Qps > 0 {
			<-throttle
		}
		requestCountChan <- 1 // sunny
		//fmt.Printf(".")
		b.makeRequest(client)
	}
}

func (b *Boomer) runWorkers() {
	fmt.Printf("Benchmarking %s (be patient)\n", b.Request.Host)
	var wg sync.WaitGroup
	wg.Add(b.C)

	// Ignore the case where b.N % b.C != 0.
	for i := 0; i < b.C; i++ {
		go func(id int) {
			nPerC := b.N / b.C
			if !b.EnableEngineIo {
				b.runWorker(nPerC)
			} else {
				//log.Printf("%d %d\n", b.C, id)
				b.runWorkerEngineIo(id) // defined in boomer_engineio.go
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request, body string) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	r2.Body = ioutil.NopCloser(strings.NewReader(body))
	return r2
}
