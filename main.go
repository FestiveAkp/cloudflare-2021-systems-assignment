package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/url"
	"os"
	"sort"
	"time"
)

func main() {
	// Get args from the command line
	fullURL := flag.String("url", "", "URL to make a request to")
	profileRequests := flag.Int("profile", 0, "Run the program in profiling mode and perform the given number of requests")
	flag.Parse()

	// Make sure the URL was supplied
	if *fullURL == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Run the CLI tool in normal or profiling mode
	if *profileRequests == 0 {
		// Make request and print the response
		response := makeRequestTo(*fullURL)
		fmt.Println("Response:")
		fmt.Println(string(response))
	} else {
		// Stuff to track
		times := []float64{}
		errorCodes := []string{}
		byteSizes := []int{}
		success := 0

		for i := 0; i < *profileRequests; i++ {
			// Start timer, make request, end timer
			start := time.Now()
			response := makeRequestTo(*fullURL)
			elapsed := time.Since(start)

			// Keep track of this duration
			times = append(times, float64(elapsed))

			// Get HTTP response code and track it
			code := string(response)[9:12]

			if code != "200" {
				errorCodes = append(errorCodes, code)
			} else {
				success++
			}

			// Keep track of size of response
			byteSizes = append(byteSizes, len(response))
		}

		// Calculate profiling results
		minTime, maxTime := minAndMaxFloats(times)
		minSize, maxSize := minAndMaxInts(byteSizes)
		mean := mean(times)
		median := median(times)

		// Print profiling results
		fmt.Println("Requests sent:", *profileRequests)
		fmt.Println("Fastest response time:", time.Duration(minTime))
		fmt.Println("Slowest response time:", time.Duration(maxTime))
		fmt.Println("Mean response time:", time.Duration(mean))
		fmt.Println("Median response time:", time.Duration(median))
		fmt.Printf("Successful response ratio: %d%%\n", (success / *profileRequests)*100)
		fmt.Println("Error codes encountered:", errorCodes)
		fmt.Println("Smallest response:", minSize, "bytes")
		fmt.Println("Largest response:", maxSize, "bytes")
	}
	os.Exit(0)
}

func makeRequestTo(fullURL string) []byte {
	// Parse the URL
	parsedURL, err := url.Parse(fullURL)
	checkError(err)

	// Get the hostname and lookup the IP
	hostname := parsedURL.Host
	ip, err := net.LookupHost(hostname)
	checkError(err)

	// Form TCP address string
	ipWithPort := ip[0] + ":80"

	// Log before making request
	fmt.Println("GET", ipWithPort, parsedURL.Path)

	// Dial the server
	conn, err := net.Dial("tcp", ipWithPort)
	checkError(err)

	// Ask server for information
	_, err = conn.Write([]byte("GET " + parsedURL.Path + " HTTP/1.0\r\nHost: " + hostname + "\r\n\r\n"))
	checkError(err)

	// Read the response back from the server
	result, err := ioutil.ReadAll(conn)
	checkError(err)

	err = conn.Close()
	checkError(err)
	return result
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
}

func minAndMaxFloats(a []float64) (float64, float64) {
	min := a[0]
	max := a[0]

	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	return min, max
}

func minAndMaxInts(a []int) (int, int) {
	min := a[0]
	max := a[0]

	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	return min, max
}

func mean(a []float64) float64 {
	sum := 0.0

	for _, value := range a {
		sum += value
	}

	len := float64(len(a))

	return math.Round(sum / len)
}

func median(a []float64) float64 {
	sort.Float64s(a)

	medianIndex := len(a) / 2

	if medianIndex%2 == 0 {
		return (a[medianIndex-1] + a[medianIndex]) / 2
	}

	return a[medianIndex]
}
