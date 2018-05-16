package serve

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func TestInvalidMsisdnResponses(t *testing.T) {
	testMsisdns := services.GetTestMsisdns()

	for _, msisdn := range testMsisdns.Invalid {
		_, res := makeRequest(fmt.Sprintf("/transfrom?msisdn=%s", msisdn))

		codeOk, _ := responseHas(res, http.StatusBadRequest, "")

		if !codeOk {
			t.Error("Responses to requests with invalid msisdns should have status code 400")
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

// func TestValidMsisdnResponses() {

// }

// /transform
// {
// 	"message": "Error: no msisdn provided"
// 	}

// // /trandform +38...
// {
// 	"Message": "Oops, invalid msisdn."
// 	}

// 38642737152
// {
// 	"Message": "Oops, invalid msisdn: no matching network operator found"
// 	}
