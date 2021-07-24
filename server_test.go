package main

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ewenquim/horkruxes/service"
	"github.com/stretchr/testify/assert"
)

func mockService() service.Service {
	// Database setup
	db := initDatabase(dbOptions{test: true})

	// Service init
	s := service.Service{
		GeneralConfig: service.GeneralConfig{Name: "Test HK", URL: "localhost"},
		ServerConfig:  service.ServerConfig{Enabled: true, Port: 8888},
		ClientConfig:  service.ClientConfig{Enabled: false},
	}
	s.DB = db
	s.Regexes = service.InitializeDetectors()
	return s
}

func TestRoutes(t *testing.T) {
	// Define a structure for specifying input and output
	// data of a single test case. This structure is then used
	// to create a so called test map, which contains all test
	// cases, that should be run for testing this function
	tests := []struct {
		description string

		// Test input
		route string

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "ping",
			route:         "/ping",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "pong",
		},
		{
			description:   "non existing route",
			route:         "/i-dont-exist",
			expectedError: false,
			expectedCode:  404,
			expectedBody:  "404 error: wrong URL",
		},
	}

	mockService := mockService()

	// Setup the app as it is done in the main function
	app, _ := setupServer(mockService)

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route
		// from the test case
		req, _ := http.NewRequest(
			"GET",
			test.route,
			nil,
		)

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, 100)

		// verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equal(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)

		assert.Equal(t, test.expectedBody, string(body))
	}
}
