package api

import (
	"github.com/radekg/yugabyte-db-go-client-api/configs"
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	"github.com/radekg/yugabyte-db-go-client/utils/ybdbid"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Delete snapshot.
func (c *defaultRpcAPI) SnapshotsDelete(opConfig *configs.OpSnapshotDeleteConfig) (*ybApi.DeleteSnapshotResponsePB, error) {

	ybDbID, err := ybdbid.TryParseSnapshotIDFromString(opConfig.SnapshotID)
	if err != nil {
		c.logger.Error("given snapshot id is not valid",
			"original-value", opConfig.SnapshotID,
			"reason", err)
		return nil, err
	}

	payload := &ybApi.DeleteSnapshotRequestPB{
		SnapshotId: ybDbID.Bytes(),
	}
	responsePayload := &ybApi.DeleteSnapshotResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, clientErrors.NewMasterError(responsePayload.Error)
}
