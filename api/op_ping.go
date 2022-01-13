package api

import (
	"github.com/radekg/yugabyte-db-go-client/client"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Ping pings a certain YB server.
func (c *defaultRpcAPI) Ping() (*ybApi.PingResponsePB, error) {
	payload := &ybApi.PingRequestPB{}
	responsePayload := &ybApi.PingResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}

// Ping pings a certain YB server.
func Ping(c client.YBConnectedClient) (*ybApi.PingResponsePB, error) {
	payload := &ybApi.PingRequestPB{}
	responsePayload := &ybApi.PingResponsePB{}
	if err := c.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
