package api

import (
	"github.com/radekg/yugabyte-db-go-client-api/configs"
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	"github.com/radekg/yugabyte-db-go-client/utils"
	"github.com/radekg/yugabyte-db-go-client/utils/relativetime"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Create a snapshot.
func (c *defaultRpcAPI) SnapshotsCreateSchedule(opConfig *configs.OpSnapshotCreateScheduleConfig) (*ybApi.CreateSnapshotScheduleResponsePB, error) {

	payload := &ybApi.CreateSnapshotScheduleRequestPB{
		Options: &ybApi.SnapshotScheduleOptionsPB{
			Filter: &ybApi.SnapshotScheduleFilterPB{
				Filter: &ybApi.SnapshotScheduleFilterPB_Tables{
					Tables: &ybApi.TableIdentifiersPB{
						Tables: []*ybApi.TableIdentifierPB{
							{
								Namespace: parseKeyspace(opConfig.Keyspace).toProtoKeyspace(),
							},
						},
					},
				},
			},
		},
	}

	if opConfig.IntervalSecs > 0 {
		payload.Options.IntervalSec = func() *uint64 {
			v := uint64(opConfig.IntervalSecs.Seconds())
			return &v
		}()
	}
	if opConfig.RetendionDurationSecs > 0 {
		payload.Options.RetentionDurationSec = func() *uint64 {
			v := uint64(opConfig.RetendionDurationSecs.Seconds())
			return &v
		}()
	}

	deleteFixedTime, deleteDuration, err := relativetime.ParseTimeOrDuration(opConfig.DeleteAfter)
	if err != nil {
		c.logger.Error("invalid delete after expression", "expression", opConfig.DeleteAfter, "reason", err)
		return nil, err
	}

	futureTime, err := relativetime.RelativeOrFixedFuture(deleteFixedTime,
		deleteDuration,
		c.defaultServerClockResolver)
	if err != nil {
		c.logger.Error("failed resolving delete at time", "reason", err)
		return nil, err
	}
	if futureTime > 0 {
		payload.Options.DeleteTime = utils.PUint64(futureTime)
	}

	responsePayload := &ybApi.CreateSnapshotScheduleResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}

	return responsePayload, clientErrors.NewMasterError(responsePayload.Error)
}
