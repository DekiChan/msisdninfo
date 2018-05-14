package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dekichan/msisdninfo/types"
)

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	respMsg := types.HomeResponseMsg{
		Status:  200,
		Message: "Kaj zdaj?",
	}

	json.NewEncoder(writer).Encode(respMsg)
}
