package stresser

import (
	"fmt"
	"load-tester/pkg/types"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"
)

var (
	successfulRequests int
	failedRequests     int
	totalRequestTime   time.Duration
	requestDurations   []time.Duration
)

// StartStresser initiates the stress testing process.
func StartStresser(url types.Uri, method types.HttpMethod, virtualUsersCount, requestsCount int, cookies []http.Cookie, proxy string, liveFeed bool) {
	responses := make(chan *http.Response, requestsCount)
	startTime := time.Now()

	// Launch goroutines for sending requests
	for i := 0; i < virtualUsersCount; i++ {
		go sendRequests(url, requestsCount/virtualUsersCount, responses, method, cookies, proxy)
	}

	// Collect responses
	for i := 0; i < requestsCount; i++ {
		response := <-responses
		handleResponse(response, i, liveFeed, url, virtualUsersCount, requestsCount)
	}

	// Sorting durations for percentile calculations
	sort.Slice(requestDurations, func(i, j int) bool {
		return requestDurations[i] < requestDurations[j]
	})

	duration := time.Since(startTime)
	ClearScreen()
	PrintStatistics(url, virtualUsersCount, requestsCount, duration)
	PrintPercentiles()
}

// sendRequests sends a specified number of requests and collects their responses.
func sendRequests(url types.Uri, numRequests int, responses chan *http.Response, method types.HttpMethod, cookies []http.Cookie, proxy string) {
	for i := 0; i < numRequests; i++ {
		sendSingleRequest(url, method, cookies, proxy, responses)
	}
}

// sendSingleRequest handles the sending of an individual request.
func sendSingleRequest(url types.Uri, method types.HttpMethod, cookies []http.Cookie, proxy string, responses chan *http.Response) {
	requestStart := time.Now()

	request, err := http.NewRequest(string(method), string(url), nil)
	if err != nil {
		fmt.Println("Error creating request: ", err)
		responses <- nil
		return
	}

	// Custom signature
	request.Header.Set("User-Agent", "BrickHauler/0.0.1")

	for _, cookie := range cookies {
		request.AddCookie(&cookie)
	}

	// Handling the proxy
	var transport *http.Transport
	var errTransport error
	transport, errTransport = useProxyHttp(proxy)
	if errTransport != nil {
		fmt.Println("Error when using the proxy")
		os.Exit(1)
	}

	client := &http.Client{
		Transport: transport,
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		failedRequests++
		responses <- nil
		return
	}
	defer response.Body.Close()

	recordRequestMetrics(response, requestStart)
	responses <- response
}

func useProxyHttp(urlProxy string) (*http.Transport, error) {
	if urlProxy == "" {
		return &http.Transport{}, nil
	}

	proxyURL, err := url.Parse(urlProxy)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("invalid proxy")
	}

	return &http.Transport{Proxy: http.ProxyURL(proxyURL)}, nil
}

// recordRequestMetrics records the metrics of each request.
func recordRequestMetrics(response *http.Response, requestStart time.Time) {
	requestDuration := time.Since(requestStart)
	totalRequestTime += requestDuration
	requestDurations = append(requestDurations, requestDuration)

	if response.StatusCode < 400 {
		successfulRequests++
	} else {
		failedRequests++
	}
}

// handleResponse processes each received response.
func handleResponse(response *http.Response, requestIndex int, liveFeed bool, url types.Uri, virtualUsersCount, requestsCount int) {
	if response != nil {
		if liveFeed {
			ClearScreen()
			PrintPrematureStats(url, virtualUsersCount, requestsCount)
		} else {
			ClearScreen()
			fmt.Println("To monitor real-time data from the load test as it runs, add the --feed parameter when executing the program.")
		}
	}
}
