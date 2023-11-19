package main

import (
	"fmt"
	"load-tester/pkg/stresser"
	"load-tester/pkg/types"
	"math"
	"net/http"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	banner "moul.io/banner"
)

// TODO: Create help command. (see pflag docs)
// TODO: Add version flag. (see pflag docs)

func main() {
	var (
		method            string
		uri               string
		virtualUsersCount int
		requestsCount     int
		cookies           []string
		liveFeed          bool
		proxy             string
	)

	// Terminal flags
	versionFlag := flag.BoolP("version", "v", false, "Display the version")
	flag.StringVar(&method, "verb", "GET", "Specifies the HTTP verb to be used (GET, POST, PUT, PATCH, DELETE, etc.)")
	flag.StringVar(&uri, "uri", "", "The URL where the tests will be performed (e.g., https://example.com)")
	flag.IntVar(&virtualUsersCount, "concurrent", 0, "The number of virtual users to launch requests concurrently")
	flag.IntVar(&requestsCount, "request", 0, "The total number of requests to be sent by all users")
	flag.StringSliceVar(&cookies, "cookie", []string{}, "Cookie to be included in the requests (format: cookieName=cookieValue)")
	flag.StringVar(&proxy, "proxy", "", "to the proxy that is going to take all the request")
	flag.BoolVar(&liveFeed, "feed", false, "Display real-time logs of the test")

	// Customize help message
	flag.ErrHelp = nil
	flag.Usage = func() {
		fmt.Println(banner.Inline("brick hauler") + "    0.1.0\n") //TODO: Move all this banners to a global place
		fmt.Println("Brick Hauler [flags]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.Parse()
	// Show only version
	if *versionFlag {
		fmt.Println("Brick Hauler 0.1.0") //TODO: Refactor this version and move to a global place
		os.Exit(0)
	}

	// Input validation
	if err := validateInputs(method, uri, virtualUsersCount, requestsCount); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Converting to cookies
	cookiesValidated, err := validateCookies(cookies)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start the stresser
	stresser.StartStresser(types.Uri(uri), types.HttpMethod(method), virtualUsersCount, requestsCount, cookiesValidated, proxy, liveFeed)
}

func validateInputs(method, uri string, virtualUsersCount, requestsCount int) error {
	if !types.IsHttpMethod(types.HttpMethod(method)) {
		return fmt.Errorf("invalid HTTP method")
	}

	if !types.IsUri(types.Uri(uri)) {
		return fmt.Errorf("invalid URI")
	}

	if virtualUsersCount <= 0 {
		return fmt.Errorf("virtualUsersCount is required and must be greater than 0")
	}

	if requestsCount <= 0 {
		return fmt.Errorf("requestsCount is required and must be greater than 0")
	}

	if requestPerVirtualUsers := float64(requestsCount) / float64(virtualUsersCount); requestPerVirtualUsers != math.Trunc(requestPerVirtualUsers) {
		return fmt.Errorf("the division of requests and virtual users is not an integer: %v", requestPerVirtualUsers)
	}

	return nil
}

func validateCookies(cookies []string) ([]http.Cookie, error) {
	// TODO : Move the creation of the cookies to the stresser package.
	var cookiesValidated []http.Cookie

	for _, cookie := range cookies {
		if cookie != "" {
			if !strings.Contains(cookie, "=") {
				return nil, fmt.Errorf("invalid cookie format")
			}

			cookieName, cookieBody, _ := strings.Cut(cookie, "=")
			cookiesValidated = append(cookiesValidated, http.Cookie{Name: cookieName, Value: cookieBody})
		}
	}

	return cookiesValidated, nil
}
