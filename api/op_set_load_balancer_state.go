package api

import (
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Set load balancer state.
func (c *defaultRpcAPI) SetLoadBalancerState(enable bool) (*ybApi.ChangeLoadBalancerStateResponsePB, error) {
	payload := &ybApi.ChangeLoadBalancerStateRequestPB{
		IsEnabled: utils.PBool(enable),
	}
	responsePayload := &ybApi.ChangeLoadBalancerStateResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, clientErrors.NewMasterError(responsePayload.Error)
}
