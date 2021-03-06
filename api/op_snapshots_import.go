package api

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/radekg/yugabyte-db-go-client-api/configs"
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

type ybTableName struct {
	YQLDatabaseType string
	KeyspaceName    string
	TableName       string
}

func newEmptyYBTableName() *ybTableName {
	return newYBTableName("", "")
}
func newYBTableName(keyspace, tableName string) *ybTableName {
	if keyspace == "" {
		return &ybTableName{
			TableName: tableName,
		}
	}
	parsedKeyspace := parseKeyspace(keyspace)
	return &ybTableName{
		YQLDatabaseType: parsedKeyspace.YQLDatabaseType,
		KeyspaceName:    parsedKeyspace.Keyspace,
		TableName:       tableName,
	}
}

func (tn *ybTableName) empty() bool {
	return tn.YQLDatabaseType == "" && tn.KeyspaceName == "" && tn.TableName == ""
}
func (tn *ybTableName) hasNamespace() bool {
	return tn.YQLDatabaseType != ""
}
func (tn *ybTableName) hasTable() bool {
	return tn.TableName != ""
}

// Pre-process snapshot import metadata file from input bytes.
func (c *defaultRpcAPI) PreProcessSnapshotsImportFromBytes(opConfig *configs.OpSnapshotImportConfig, rawProtoBytes []byte) (*ybApi.ImportSnapshotMetaRequestPB, error) {

	givenKeyspace := ""
	if opConfig.Keyspace != "" {
		parsedKeyspace := parseKeyspace(opConfig.Keyspace)
		givenKeyspace = parsedKeyspace.Keyspace
	}

	tables := []*ybTableName{}
	for _, tableName := range opConfig.TableName {
		tables = append(tables, newYBTableName(opConfig.Keyspace, tableName))
	}

	snapshotInfo := &ybApi.SnapshotInfoPB{}
	if err := utils.DeserializeProto(rawProtoBytes, snapshotInfo); err != nil {
		return nil, err
	}

	for _, backupEntry := range snapshotInfo.BackupEntries {

		wasTableRenamed := false
		tableIndex := 0
		tableName := newEmptyYBTableName()
		if tableIndex < len(tables) {
			tableName = tables[tableIndex]
		}

		sysRowEntry := backupEntry.Entry
		switch *sysRowEntry.Type {
		case ybApi.SysRowEntry_NAMESPACE:

			meta := &ybApi.SysNamespaceEntryPB{}
			if err := utils.DeserializeProto(sysRowEntry.Data, meta); err != nil {
				return nil, err
			}

			if givenKeyspace != "" && givenKeyspace != string(meta.Name) {
				meta.Name = []byte(givenKeyspace)
				metaBytes, err := utils.SerializeProto(meta)
				if err != nil {
					return nil, err
				}
				sysRowEntry.Data = metaBytes
			}

		case ybApi.SysRowEntry_TABLE:

			if wasTableRenamed && tableName.TableName == "" {
				return nil, fmt.Errorf("there is no name for table (including indexes) number: %d", tableIndex)
			}
			meta := &ybApi.SysTablesEntryPB{}
			if err := utils.DeserializeProto(sysRowEntry.Data, meta); err != nil {
				return nil, err
			}
			updateMeta := false
			if !tableName.empty() && tableName.TableName != string(meta.Name) {
				meta.Name = []byte(tableName.TableName)
				updateMeta = true
				wasTableRenamed = true
			}
			if givenKeyspace != "" && givenKeyspace != string(meta.NamespaceName) {
				meta.NamespaceName = []byte(givenKeyspace)
				updateMeta = true
			}
			if len(meta.Name) == 0 {
				return nil, fmt.Errorf("could not find table name from snapshot metadata")
			}
			if updateMeta {
				metaBytes, err := utils.SerializeProto(meta)
				if err != nil {
					return nil, err
				}
				sysRowEntry.Data = metaBytes
			}

			colocatedPrefix := ""
			if meta.Colocated != nil && *meta.Colocated {
				colocatedPrefix = "colocatated "
			}

			if len(meta.IndexedTableId) == 0 {
				c.logger.Info(fmt.Sprintf("table type: %stable", colocatedPrefix))
			} else {
				c.logger.Info(fmt.Sprintf("table type: %sindex (attaching the old table id)", colocatedPrefix),
					"old-table-id", string(meta.IndexedTableId))
			}

			if !tableName.empty() {
				c.logger.Info(fmt.Sprintf("target imported %s", colocatedPrefix),
					"table-name", tableName.TableName)
			} else if givenKeyspace != "" {
				c.logger.Info(fmt.Sprintf("target imported %s (attaching the old table id)", colocatedPrefix),
					"keyspace-name", givenKeyspace)
			}

			if meta.Colocated != nil && *meta.Colocated {
				c.logger.Info("Colocated table being imported",
					"namespace-name", string(meta.NamespaceName),
					"namespace-id", string(meta.NamespaceId),
					"table-name", tableName.TableName)
			}

			tableIndex = tableIndex + 1
		}
	}

	return &ybApi.ImportSnapshotMetaRequestPB{
		Snapshot: snapshotInfo,
	}, nil

}

// Pre-process snapshot import metadata file.
func (c *defaultRpcAPI) PreProcessSnapshotsImportFromFile(opConfig *configs.OpSnapshotImportConfig) (*ybApi.ImportSnapshotMetaRequestPB, error) {
	statResult, err := os.Stat(opConfig.FilePath)
	if err != nil {
		return nil, err
	}
	if statResult.IsDir() {
		return nil, fmt.Errorf("path %s points at a directory", opConfig.FilePath)
	}

	rawProtoBytes, err := ioutil.ReadFile(opConfig.FilePath)
	if err != nil {
		return nil, err
	}

	return c.PreProcessSnapshotsImportFromBytes(opConfig, rawProtoBytes)
}

// Import snapshot.
func (c *defaultRpcAPI) SnapshotsImport(opConfig *configs.OpSnapshotImportConfig) (*ybApi.ImportSnapshotMetaResponsePB, error) {

	payload, err := c.PreProcessSnapshotsImportFromFile(opConfig)
	if err != nil {
		return nil, err
	}
	responsePayload := &ybApi.ImportSnapshotMetaResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}

	return responsePayload, clientErrors.NewMasterError(responsePayload.Error)
}
