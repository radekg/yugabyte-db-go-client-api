package api

import (
	"github.com/radekg/yugabyte-db-go-client-api/configs"
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	"github.com/radekg/yugabyte-db-go-client/utils/ybdbid"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// List snapshots.
func (c *defaultRpcAPI) SnapshotsListSchedules(opConfig *configs.OpSnapshotListSchedulesConfig) (*ybApi.ListSnapshotSchedulesResponsePB, error) {
	payload := &ybApi.ListSnapshotSchedulesRequestPB{}
	if len(opConfig.ScheduleID) > 0 {
		ybDbID, err := ybdbid.TryParseSnapshotIDFromString(opConfig.ScheduleID)
		if err != nil {
			c.logger.Error("given schedule id is not valid",
				"original-value", opConfig.ScheduleID,
				"reason", err)
			return nil, err
		}
		payload.SnapshotScheduleId = ybDbID.Bytes()
	}
	responsePayload := &ybApi.ListSnapshotSchedulesResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, clientErrors.NewMasterError(responsePayload.Error)
}
