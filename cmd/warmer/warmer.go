package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/AcroManiac/micropic/internal/adapters/logger"
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
	flag.Usage = Usage
	flag.Parse()

	// Get image URLs from external image source
	urlChan := make(chan string, flagWorkers)
	go func() {
		for i := 0; i < flagNumber; i++ {
			url, err := getRandomDogURL()
			if err != nil {
				logger.Error("failed calling external image source", "error", err)
				continue
			}
			urlChan <- url
		}
	}()

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-done:
			logger.Error("program interrupted by user or OS")
			close(urlChan)
			os.Exit(0)
		case url := <-urlChan:
			go func(url string) {
				urlconv := convertProxyURL(url)

				// Get preview from HTTP proxy
				// This warms up preview cache
				resp, err := http.Get(urlconv) // nolint:gosec
				if err != nil {
					logger.Error("error calling preview proxy", "error", err, "response", resp)
					return
				}

				logger.Debug("preview returned successfully")
				_ = resp.Body.Close()
			}(url)
		}
	}
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
