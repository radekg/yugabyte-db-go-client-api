package api

import (
	"github.com/radekg/yugabyte-db-go-client-api/configs"
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// IsLoadBalanced returns a list of masters or an error if call failed.
func (c *defaultRpcAPI) IsLoadBalanced(opConfig *configs.OpIsLoadBalancedConfig) (*ybApi.IsLoadBalancedResponsePB, error) {
	payload := &ybApi.IsLoadBalancedRequestPB{
		ExpectedNumServers: func() *int32 {
			if opConfig.ExpectedNumServers > 0 {
				return utils.PInt32(int32(opConfig.ExpectedNumServers))
			}
			return nil
		}(),
	}
	responsePayload := &ybApi.IsLoadBalancedResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, clientErrors.NewMasterError(responsePayload.Error)
}
