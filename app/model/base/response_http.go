package base

import (
	"context"
	"dating-apps/helper/message"
	"encoding/json"
	"net/http"
	"reflect"
)

type metaResponse struct {
	// Code is the response code
	// example: 1000
	Code int `json:"code"`
	// Message is the response message
	// example: Success
	Message string `json:"message"`
}

type responseHttp struct {
	// Meta is the API response information
	// in: struct{}
	Meta metaResponse `json:"meta"`

	// Data is our data
	// in: struct{}
	Data data `json:"data"`
	// Errors is the response message
	// in: string
	Errors interface{} `json:"errors,omitempty"`
}

type data struct {
	Records interface{} `json:"records,omitempty"`
	Record  interface{} `json:"record,omitempty"`
}

func SetDefaultResponse(ctx context.Context, msg message.Message) interface{} {
	return responseHttp{
		Meta: metaResponse{
			Code:    msg.Code,
			Message: msg.Message,
		},
	}
}

func GetHttpResponse(resp interface{}) *responseHttp {
	result, ok := resp.(responseHttp)
	if ok {
		return &result
	}
	return nil
}

func SetHttpResponse(ctx context.Context, msg message.Message, result interface{}) interface{} {
	dt := data{}
	isSlice := reflect.ValueOf(result).Kind() == reflect.Slice
	if isSlice {
		dt.Records = result
		dt.Record = nil
	} else {
		dt.Records = nil
		dt.Record = result
	}

	return responseHttp{
		Meta: metaResponse{
			Code:    msg.Code,
			Message: msg.Message,
		},

		Data: dt,
	}
}

func ResponseWriter(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
	return
}
