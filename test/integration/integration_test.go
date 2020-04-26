package main // nolint:testpackage

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
)

// Global variables for multiple tests execution
var (
	client = http.DefaultClient
	req    *http.Request
	resp   *http.Response
)

func TestMain(m *testing.M) {
	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:      "progress",
		Paths:       []string{"features"},
		Randomize:   0,
		Concurrency: 0,
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I make a "([^"]*)" request to "([^"]*)"$`, iMakeARequestTo)
	s.Step(`^I get response status code (\d+)$`, iGetResponseStatusCode)
	s.Step(`^the response headers has:$`, theResponseHeadersHas)
	s.Step(`^the response body is "([^"]*)"$`, theResponseBodyIs)
	s.Step(`^preview size is:$`, previewSizeIs)

	s.BeforeScenario(func(*messages.Pickle) {
		// clean the state before every scenario
		if resp != nil {
			_ = resp.Body.Close()
		}
	})
}

func iMakeARequestTo(arg1, arg2 string) error {
	// Make new http request
	var err error
	req, err = http.NewRequest(arg1, arg2, nil)
	if err != nil {
		return errors.Wrap(err, "error allocating http request")
	}
	return nil
}

func iGetResponseStatusCode(arg1 int) error {
	resp, _ = client.Do(req)
	if resp == nil {
		return errors.New("empty response from HTTP proxy")
	}
	if resp.StatusCode != arg1 {
		return errors.New("returned status is not correct")
	}
	return nil
}

func previewSizeIs(arg1 *messages.PickleStepArgument_PickleTable) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error reading response body")
	}

	preview, _, err := image.DecodeConfig(bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "error reading decoded preview")
	}

	width, _ := strconv.Atoi(arg1.Rows[0].Cells[1].Value)
	height, _ := strconv.Atoi(arg1.Rows[1].Cells[1].Value)
	if preview.Width != width || preview.Height != height {
		return errors.New("preview size is wrong")
	}
	return nil
}

func theResponseHeadersHas(arg1 *messages.PickleStepArgument_PickleTable) error {
	key := arg1.Rows[0].Cells[0].Value
	value := arg1.Rows[0].Cells[1].Value
	fmt.Printf("key: %s, value: %s\n", key, value)
	fmt.Printf("Response header: %+v\n", resp.Header)
	if resp.Header.Get(key) != value {
		return errors.New("response header is wrong")
	}
	return nil
}

func theResponseBodyIs(arg1 string) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error reading response body")
	}
	if string(body) != arg1 {
		return errors.New("response body is wrong")
	}
	return nil
}
