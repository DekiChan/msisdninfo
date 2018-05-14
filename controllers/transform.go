package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dekichan/msisdninfo/types"
)

func TransformHandler(writer http.ResponseWriter, request *http.Request) {
	respMsg := types.TransformResponseMsg{
		MnoIdentifier:     "A1",
		CountryCode:       386,
		CountryIdentifier: "SI",
		SubscriberNumber:  "737152",
	}

	json.NewEncoder(writer).Encode(respMsg)
}
