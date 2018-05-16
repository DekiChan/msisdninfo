package serve

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dekichan/msisdninfo/services"
	"github.com/gorilla/mux"
)

func createRouter() *mux.Router {
	router := mux.NewRouter()
	setupRoutes(router)
	return router
}

func makeRequest(route string) (req *http.Request, res *http.Response) {
	req, _ = http.NewRequest("GET", route, nil)
	recorder := httptest.NewRecorder()

	createRouter().ServeHTTP(recorder, req)

	res = recorder.Result()

	return
}

func TestUnexistentRoute(t *testing.T) {
	_, res := makeRequest("/unexistent")
	codeOk, bodyOk := responseHas(res, 404, "404 page not found")

	if !codeOk || !bodyOk {
		t.Error("Responses to requests for unexistent routes should have status code 404 Not Found and string 404 page not found in body")
	}
}

func TestWithoutMsisdn(t *testing.T) {
	_, res := makeRequest("/transform")

	codeOk, bodyOk := responseHas(res, http.StatusBadRequest, `{"message":"Error: no msisdn provided"}`)

	if !codeOk {
		t.Error("Request without msisdn should have code 400")
	} else if !bodyOk {
		t.Error("Request without msisdn has wrong body")
	}
}

func TestInvalidMsisdnResponses(t *testing.T) {
	testMsisdns := services.GetTestMsisdns()

	for _, msisdn := range testMsisdns.Invalid {
		_, res := makeRequest(fmt.Sprintf("/transform?msisdn=%s", msisdn))

		codeOk, _ := responseHas(res, http.StatusBadRequest, "")

		if !codeOk {
			t.Error("Responses to requests with invalid msisdns should have status code 400")
		}
	}
}

func TestValidMsisdnResponses(t *testing.T) {
	testMsisdns := services.GetTestMsisdns()

	for _, msisdn := range testMsisdns.ValidSloA1 {
		_, res := makeRequest(fmt.Sprintf("/transform?msisdn=%s", url.QueryEscape(msisdn)))

		codeOk, bodyOk := responseHas(res, http.StatusOK, `{"mno_identifier":"A1","country_code":386,"country_identifier":"SI","subscriber_number":"123456"}`)

		if !codeOk {
			t.Error("Responses to requests with valid msisdns should have status code 200")
		} else if !bodyOk {
			t.Error("Response to request with valid A1 msisdn got wrong result")
		}
	}
}

func responseHas(res *http.Response, code int, body string) (codeOk bool, bodyOk bool) {
	codeOk = res.StatusCode == code

	bodyBytes, _ := ioutil.ReadAll(res.Body)
	bodyString := strings.Trim(string(bodyBytes), "\r\n")
	bodyOk = bodyString == body

	return
}
