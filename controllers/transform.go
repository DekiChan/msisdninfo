package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dekichan/msisdninfo/services"
	"github.com/dekichan/msisdninfo/types"
)

func TransformHandler(writer http.ResponseWriter, request *http.Request) {
	msisdn := request.URL.Query().Get("msisdn")
	fmt.Println(fmt.Sprintf("Parsing transform request for %s", msisdn))

	if len(msisdn) != 0 {
		msisdnService := services.CreateMsisdnService()

		resp, err := msisdnService.Parse(msisdn)

		if err == nil {
			json.NewEncoder(writer).Encode(resp)
			return
		} else {
			json.NewEncoder(writer).Encode(err)
			return
		}
	}

	json.NewEncoder(writer).Encode(types.ErrorResponseMsg{
		Message: "Error: no msisdn provided",
	})
}
