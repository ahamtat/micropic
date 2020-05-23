package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ahamtat/micropic/internal/adapters/logger"
)

var (
	flagNumber  int
	flagWorkers int
	flagProxy   string
)

func init() {
	rand.Seed(time.Now().UnixNano())
	flag.IntVar(&flagNumber, "number", 0, "number of images for preview making")
	flag.IntVar(&flagWorkers, "workers", 1, "number of worker goroutines")
	flag.StringVar(&flagProxy, "proxy", "http://localhost:8080/fill", "proxy URL")
	// Initialize log parameters
	logger.Init("debug", "")
}

var Usage = func() {
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	start := time.Now()

	flag.Usage = Usage
	flag.Parse()

	// Get image URLs from external image source
	var wg sync.WaitGroup
	urlChan := make(chan string, flagWorkers)
	for w := 0; w < flagWorkers; w++ {
		go func() {
			for i := 0; i < flagNumber; i++ {
				wg.Add(1)
				url, err := getRandomDogURL()
				if err != nil {
					logger.Error("failed calling external image source", "error", err)
					continue
				}
				urlChan <- url
			}
		}()
	}

	// Get preview from HTTP proxy
	// This warms up preview cache
	go func() {
		for url := range urlChan {
			go func(url string) {
				defer wg.Done()
				urlconv := convertProxyURL(url)
				resp, err := http.Get(urlconv) // nolint:gosec
				if err != nil || resp == nil {
					logger.Error("error calling preview proxy", "error", err, "response", resp)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode == 200 {
					logger.Debug("preview returned successfully")
				} else {
					body, _ := ioutil.ReadAll(resp.Body)
					logger.Error("returned error response",
						"code", resp.StatusCode,
						"body", body)
				}
			}(url)
		}
	}()

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	wg.Wait()
	close(urlChan)

	elapsed := time.Since(start)
	logger.Info(fmt.Sprintf("Application working time is %s", elapsed))
}

var sizes = []string{
	"50", "100", "200", "300", "500", "1024", "2000",
}

func convertProxyURL(url string) string {
	pos := strings.Index(url, ":") + 3
	result := flagProxy + "/" +
		sizes[rand.Intn(len(sizes))] + "/" +
		sizes[rand.Intn(len(sizes))] + "/" +
		url[pos:]
	return result
}
