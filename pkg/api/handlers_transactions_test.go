package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/api"
	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core"
	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core/log"
	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthoriseTransaction(t *testing.T) {

	type CreditCard struct {
		Name        string `json:"name"`
		Number      int64  `json:"number"`
		ExpiryMonth int    `json:"expiry_month"`
		ExpiryYear  int    `json:"expiry_year"`
		CVV         int    `json:"cvv"`
	}

	type RequestBody struct {
		Currency   string     `json:"currency"`
		Amount     float64    `json:"amount"`
		CreditCard CreditCard `json:"credit_card"`
	}

	type ResponseBody struct {
		Code            uint   `json:"code"`
		AuthorisationID string `json:"authorisation_id,omitempty"`
	}

	// Setup
	assert := assert.New(t)
	logger := log.NullLogger{}
	ccfc := createCreditCardFileChecker()
	at := repository.NewAuthoriserInMemoryTracker()
	server := api.NewServer("", 9999, false, logger, ccfc, at)
	router := server.Router

	// Table driven testing
	tests := map[string]struct {
		RequestBody          RequestBody
		expectedStatusCode   int
		expectedResponseBody ResponseBody
	}{
		"empty body": {
			RequestBody:        RequestBody{},
			expectedStatusCode: 400,
		},
		"valid request": {
			RequestBody: RequestBody{
				Currency: "EUR",
				Amount:   10.50,
				CreditCard: CreditCard{
					Name:        "customer1",
					Number:      1111222233334444,
					ExpiryMonth: 10,
					ExpiryYear:  2025,
					CVV:         123},
			},
			expectedStatusCode: 200,
			expectedResponseBody: ResponseBody{
				Code: 1,
			},
		},
		"failed request": {
			RequestBody: RequestBody{
				Currency: "EUR",
				Amount:   10.50,
				CreditCard: CreditCard{
					Name:        "customer1",
					Number:      4000000000000119,
					ExpiryMonth: 10,
					ExpiryYear:  2025,
					CVV:         123},
			},
			expectedStatusCode: 200,
			expectedResponseBody: ResponseBody{
				Code: 2,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			requestBodyBytes, err := json.Marshal(test.RequestBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/v1/authorise", bytes.NewBuffer(requestBodyBytes))
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatusCode, w.Code)

			if test.expectedStatusCode == 200 {
				var response ResponseBody
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(test.expectedResponseBody.Code, response.Code)
			}
		})
	}
}

func TestCaptureTransaction(t *testing.T) {

	type RequestBody struct {
		AuthorisationID string  `json:"authorisation_id"`
		Amount          float64 `json:"amount"`
	}

	type ResponseBody struct {
		Code uint `json:"code"`
	}

	// Setup
	assert := assert.New(t)
	logger := log.NullLogger{}
	ccfc := createCreditCardFileChecker()
	at := repository.NewAuthoriserInMemoryTracker()
	uid1 := "53871001-f41a-4b87-9179-38d531bacece"
	at.Authorisations[uid1] = 4000000000000001
	uid2 := "53871001-f41a-4b87-9179-38d531baaaaa"
	at.Authorisations[uid2] = 4000000000000259

	server := api.NewServer("", 9999, false, logger, ccfc, at)
	router := server.Router

	// Table driven testing
	tests := map[string]struct {
		RequestBody          RequestBody
		expectedStatusCode   int
		expectedResponseBody ResponseBody
	}{
		"empty body": {
			RequestBody:        RequestBody{},
			expectedStatusCode: 400,
		},
		"valid request": {
			RequestBody: RequestBody{
				AuthorisationID: uid1,
				Amount:          10.50,
			},
			expectedStatusCode: 200,
			expectedResponseBody: ResponseBody{
				Code: 1,
			},
		},
		"failed request": {
			RequestBody: RequestBody{
				AuthorisationID: uid2,
				Amount:          10.50,
			},
			expectedStatusCode: 200,
			expectedResponseBody: ResponseBody{
				Code: 2,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			requestBodyBytes, err := json.Marshal(test.RequestBody)
			require.NoError(t, err)
			responseBodyBytes, err := json.Marshal(test.expectedResponseBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/v1/capture", bytes.NewBuffer(requestBodyBytes))
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatusCode, w.Code)

			if test.expectedStatusCode == 200 {
				assert.JSONEq(string(responseBodyBytes), w.Body.String())
			}
		})
	}
}

func TestVoidTransaction(t *testing.T) {

	type RequestBody struct {
		AuthorisationID string `json:"authorisation_id"`
	}

	type ResponseBody struct {
		Code uint `json:"code"`
	}

	// Setup
	assert := assert.New(t)
	logger := log.NullLogger{}
	ccfc := createCreditCardFileChecker()
	at := repository.NewAuthoriserInMemoryTracker()
	uid1 := "53871001-f41a-4b87-9179-38d531bacece"
	at.Authorisations[uid1] = 4000000000000001
	uid2 := "53871001-f41a-4b87-9179-38d531baaaaa"
	at.Authorisations[uid2] = 4000000000000500

	server := api.NewServer("", 9999, false, logger, ccfc, at)
	router := server.Router

	// Table driven testing
	tests := map[string]struct {
		RequestBody          RequestBody
		expectedStatusCode   int
		expectedResponseBody ResponseBody
	}{
		"empty body": {
			RequestBody:        RequestBody{},
			expectedStatusCode: 400,
		},
		"valid request": {
			RequestBody: RequestBody{
				AuthorisationID: uid1,
			},
			expectedStatusCode: 200,
			expectedResponseBody: ResponseBody{
				Code: 1,
			},
		},
		"failed request": {
			RequestBody: RequestBody{
				AuthorisationID: uid2,
			},
			expectedStatusCode: 200,
			expectedResponseBody: ResponseBody{
				Code: 2,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			requestBodyBytes, err := json.Marshal(test.RequestBody)
			require.NoError(t, err)
			responseBodyBytes, err := json.Marshal(test.expectedResponseBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/v1/void", bytes.NewBuffer(requestBodyBytes))
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatusCode, w.Code)

			if test.expectedStatusCode == 200 {
				assert.JSONEq(string(responseBodyBytes), w.Body.String())
			}
		})
	}
}

func TestRefundTransaction(t *testing.T) {

	type RequestBody struct {
		AuthorisationID string  `json:"authorisation_id"`
		Amount          float64 `json:"amount"`
	}

	type ResponseBody struct {
		Code uint `json:"code"`
	}

	// Setup
	assert := assert.New(t)
	logger := log.NullLogger{}
	ccfc := createCreditCardFileChecker()
	at := repository.NewAuthoriserInMemoryTracker()
	uid1 := "53871001-f41a-4b87-9179-38d531bacece"
	at.Authorisations[uid1] = 4000000000000001
	uid2 := "53871001-f41a-4b87-9179-38d531baaaaa"
	at.Authorisations[uid2] = 4000000000003238

	server := api.NewServer("", 9999, false, logger, ccfc, at)
	router := server.Router

	// Table driven testing
	tests := map[string]struct {
		RequestBody          RequestBody
		expectedStatusCode   int
		expectedResponseBody ResponseBody
	}{
		"empty body": {
			RequestBody:        RequestBody{},
			expectedStatusCode: 400,
		},
		"valid request": {
			RequestBody: RequestBody{
				AuthorisationID: uid1,
				Amount:          10.50,
			},
			expectedStatusCode: 200,
			expectedResponseBody: ResponseBody{
				Code: 1,
			},
		},
		"failed request": {
			RequestBody: RequestBody{
				AuthorisationID: uid2,
				Amount:          10.50,
			},
			expectedStatusCode: 200,
			expectedResponseBody: ResponseBody{
				Code: 2,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			requestBodyBytes, err := json.Marshal(test.RequestBody)
			require.NoError(t, err)
			responseBodyBytes, err := json.Marshal(test.expectedResponseBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/v1/refund", bytes.NewBuffer(requestBodyBytes))
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatusCode, w.Code)

			if test.expectedStatusCode == 200 {
				assert.JSONEq(string(responseBodyBytes), w.Body.String())
			}
		})
	}
}

func createCreditCardFileChecker() *repository.CreditCardFileChecker {
	ccfc := repository.NewCreditCardFileChecker()

	ccfc.CreditCards[4000000000000119] = core.CCFailReason_Authorise
	ccfc.CreditCards[4000000000000259] = core.CCFailReason_Capture
	ccfc.CreditCards[4000000000000500] = core.CCFailReason_Void
	ccfc.CreditCards[4000000000003238] = core.CCFailReason_Refund

	return ccfc
}
