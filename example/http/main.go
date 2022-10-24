// Package main
package main

import (
	"context"
	"net/http"

	"gitlab.privy.id/go_graphql/pkg/httpclient"
	"gitlab.privy.id/go_graphql/pkg/logger"
	"gitlab.privy.id/go_graphql/pkg/util"
)

type (
	Payload struct {
		Name   string `json:"name"`
		Rating int64  `json:"rating"`
	}

	Response struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		Rating int64  `json:"rating"`
	}
)

func main() {

	var result Response
	param := Payload{
		Name:   "test",
		Rating: 8,
	}

	header := httpclient.Headers{}
	header.Add("Content-Type", "application/json")

	endpoint := "https://jsonplaceholder.typicode.com/todos"

	resp, err := httpclient.Request(httpclient.RequestOptions{
		Payload:       param,
		URL:           endpoint, // set to your config
		Method:        http.MethodPost,
		TimeoutSecond: 3,                    // set to your config
		Context:       context.Background(), // this context for http propagation
		Header:        header,
	})

	if err != nil {
		logger.Error(logger.MessageFormat("request http error : %v", err))

		// TODO: do something when got error
	}

	err = resp.DecodeJSON(&result)

	if err != nil {
		logger.Error(logger.MessageFormat("request http error : %v", err))
		// TODO: do something when got error
	}

	// TODO: do something when success

	util.DebugPrint(result)

}
