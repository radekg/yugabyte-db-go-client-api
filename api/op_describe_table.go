package api

import (
	"github.com/radekg/yugabyte-db-go-client-api/configs"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// DescribeTable returns info on a table in this database.
func (c *defaultRpcAPI) DescribeTable(opConfig *configs.OpGetTableSchemaConfig) (*ybApi.GetTableSchemaResponsePB, error) {
	if opConfig.UUID != "" {
		// we can short circuit everything below:
		return c.getTableSchemaByUUID([]byte(opConfig.UUID))
	}
	return c.lookupTableByName(opConfig.Keyspace, opConfig.Name)
}
