package api

import (
	"github.com/hashicorp/go-hclog"
	"github.com/radekg/yugabyte-db-go-client-api/configs"
	"github.com/radekg/yugabyte-db-go-client/client"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// YBRpcAPI is a client implementing the CLI functionality.
type YBRpcAPI interface {
	CheckExists(*configs.OpGetTableSchemaConfig) (*ybApi.GetTableSchemaResponsePB, error)
	DescribeTable(*configs.OpGetTableSchemaConfig) (*ybApi.GetTableSchemaResponsePB, error)

	Execute(payload, response protoreflect.ProtoMessage) error

	GetIsLoadBalancerIdle() (*ybApi.IsLoadBalancerIdleResponsePB, error)
	GetLeaderBlacklistCompletion() (*ybApi.GetLoadMovePercentResponsePB, error)
	GetLoadMoveCompletion() (*ybApi.GetLoadMovePercentResponsePB, error)
	GetMasterRegistration() (*ybApi.GetMasterRegistrationResponsePB, error)
	GetTabletsForTable(*configs.OpGetTableLocationsConfig) (*ybApi.GetTableLocationsResponsePB, error)
	GetUniverseConfig() (*ybApi.GetMasterClusterConfigResponsePB, error)
	IsLoadBalanced(*configs.OpIsLoadBalancedConfig) (*ybApi.IsLoadBalancedResponsePB, error)
	IsTabletServerReady() (*ybApi.IsTabletServerReadyResponsePB, error)
	LeaderStepDown(*configs.OpLeaderStepDownConfig) (*ybApi.LeaderStepDownResponsePB, error)
	ListMasters() (*ybApi.ListMastersResponsePB, error)
	ListTables(*configs.OpListTablesConfig) (*ybApi.ListTablesResponsePB, error)
	ListTabletServers(*configs.OpListTabletServersConfig) (*ybApi.ListTabletServersResponsePB, error)
	MasterLeaderStepDown(*configs.OpMMasterLeaderStepdownConfig) (*ybApi.GetMasterRegistrationResponsePB, error)
	ModifyPlacementInfo(*configs.OpModifyPlacementInfoConfig) (*ybApi.ChangeMasterClusterConfigResponsePB, error)
	Ping() (*ybApi.PingResponsePB, error)
	SetLoadBalancerState(bool) (*ybApi.ChangeLoadBalancerStateResponsePB, error)
	SetPreferredZones(*configs.OpSetPreferredZonesConfig) (*ybApi.SetPreferredZonesResponsePB, error)

	ServerClock() (*ybApi.ServerClockResponsePB, error)

	SnapshotsCreateSchedule(*configs.OpSnapshotCreateScheduleConfig) (*ybApi.CreateSnapshotScheduleResponsePB, error)
	SnapshotsCreate(*configs.OpSnapshotCreateConfig) (*ybApi.CreateSnapshotResponsePB, error)
	SnapshotsDeleteSchedule(*configs.OpSnapshotDeleteScheduleConfig) (*ybApi.DeleteSnapshotScheduleResponsePB, error)
	SnapshotsDelete(*configs.OpSnapshotDeleteConfig) (*ybApi.DeleteSnapshotResponsePB, error)
	SnapshotsExport(*configs.OpSnapshotExportConfig) (*SnapshotExportData, error)
	PreProcessSnapshotsImportFromBytes(*configs.OpSnapshotImportConfig, []byte) (*ybApi.ImportSnapshotMetaRequestPB, error)
	PreProcessSnapshotsImportFromFile(*configs.OpSnapshotImportConfig) (*ybApi.ImportSnapshotMetaRequestPB, error)
	SnapshotsImport(*configs.OpSnapshotImportConfig) (*ybApi.ImportSnapshotMetaResponsePB, error)
	SnapshotsListSchedules(*configs.OpSnapshotListSchedulesConfig) (*ybApi.ListSnapshotSchedulesResponsePB, error)
	SnapshotsListRestorations(*configs.OpSnapshotListRestorationsConfig) (*ybApi.ListSnapshotRestorationsResponsePB, error)
	SnapshotsList(*configs.OpSnapshotListConfig) (*ybApi.ListSnapshotsResponsePB, error)
	SnapshotsRestoreSchedule(*configs.OpSnapshotRestoreScheduleConfig) (*ybApi.RestoreSnapshotResponsePB, error)
	SnapshotsRestore(*configs.OpSnapshotRestoreConfig) (*ybApi.RestoreSnapshotResponsePB, error)

	YsqlCatalogVersion() (*ybApi.GetYsqlCatalogConfigResponsePB, error)
}

type defaultRpcAPI struct {
	connectedClient client.YBClient
	logger          hclog.Logger
}

// NewRpcAPI returns a configured instance of the default CLI client.
func NewRpcAPI(c client.YBClient, logger hclog.Logger) YBRpcAPI {
	return &defaultRpcAPI{
		connectedClient: c,
		logger:          logger,
	}
}

func (c *defaultRpcAPI) Execute(input, output protoreflect.ProtoMessage) error {
	return c.connectedClient.Execute(input, output)
}
