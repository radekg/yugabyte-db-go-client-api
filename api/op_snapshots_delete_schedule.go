package api

import (
	"github.com/radekg/yugabyte-db-go-client-api/configs"
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	"github.com/radekg/yugabyte-db-go-client/utils/ybdbid"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Delete snapshot.
func (c *defaultRpcAPI) SnapshotsDeleteSchedule(opConfig *configs.OpSnapshotDeleteScheduleConfig) (*ybApi.DeleteSnapshotScheduleResponsePB, error) {

	ybDbID, err := ybdbid.TryParseSnapshotIDFromString(opConfig.ScheduleID)
	if err != nil {
		c.logger.Error("given schedule id is not valid",
			"original-value", opConfig.ScheduleID,
			"reason", err)
		return nil, err
	}

	payload := &ybApi.DeleteSnapshotScheduleRequestPB{
		SnapshotScheduleId: ybDbID.Bytes(),
	}
	responsePayload := &ybApi.DeleteSnapshotScheduleResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, clientErrors.NewMasterError(responsePayload.Error)
}
