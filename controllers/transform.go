package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dekichan/msisdninfo/services"
	"github.com/dekichan/msisdninfo/types"
)

func TransformHandler(writer http.ResponseWriter, request *http.Request) {
	msisdn := request.URL.Query().Get("msisdn")

	if len(msisdn) != 0 {
		msisdnService := services.CreateMsisdnService()

		resp := msisdnService.Parse("38640737152")
		json.NewEncoder(writer).Encode(resp)
		return
	}
	// respMsg := types.TransformResponseMsg{
	// 	MnoIdentifier:     "A1",
	// 	CountryCode:       386,
	// 	CountryIdentifier: "SI",
	// 	SubscriberNumber:  "737152",
	// }
	json.NewEncoder(writer).Encode(types.ErrorResponseMsg{
		Message: "Error: no msisdn provided",
	})
}
