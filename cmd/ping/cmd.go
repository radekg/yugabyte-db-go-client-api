package ping

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/radekg/yugabyte-db-go-client-api/api"
	"github.com/radekg/yugabyte-db-go-client-api/configs"
	"github.com/radekg/yugabyte-db-go-client/client"
	ybClientConfigs "github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/spf13/cobra"
)

// Command is the command declaration.
var Command = &cobra.Command{
	Use:   "ping",
	Short: "Ping a certain YB server",
	Run:   run,
	Long:  ``,
}

var (
	commandConfig = configs.NewCliConfig()
	logConfig     = configs.NewLogginConfig()
	opConfig      = configs.NewOpPingConfig()
)

func initFlags() {
	Command.Flags().AddFlagSet(commandConfig.FlagSet())
	Command.Flags().AddFlagSet(logConfig.FlagSet())
	Command.Flags().AddFlagSet(opConfig.FlagSet())
}

func init() {
	initFlags()
}

func run(cobraCommand *cobra.Command, _ []string) {
	os.Exit(processCommand())
}

func processCommand() int {

	logger := logConfig.NewLogger("ping")

	for _, validatingConfig := range []configs.ValidatingConfig{commandConfig, opConfig} {
		if err := validatingConfig.Validate(); err != nil {
			logger.Error("configuration is invalid", "reason", err)
			return 1
		}
	}

	tlsConfig, tlsConfigErr := commandConfig.TLSConfig()
	if tlsConfigErr != nil {
		logger.Error("TLS configuration failed", "reason", tlsConfigErr)
		return 1
	}

	cliClient, err := client.NewDefaultConnector().Connect(&ybClientConfigs.YBSingleNodeClientConfig{
		MasterHostPort: fmt.Sprintf("%s:%d", opConfig.Host, opConfig.Port),
		TLSConfig:      tlsConfig,
		OpTimeout:      uint32(commandConfig.OpTimeout),
	})
	if err != nil {
		// careful: different than other commands:
		logger.Error("server not reachable", "reason", err)
		return 2
	}
	select {
	case err := <-cliClient.OnConnectError():
		// TODO: LATER: in this case, this may indicate the service unavailability
		logger.Error("failed connecting a client", "reason", err)
		return 1
	case <-cliClient.OnConnected():
		logger.Debug("client connected")
	}
	defer cliClient.Close()

	responsePayload, err := api.Ping(cliClient)
	if err != nil {
		// TODO: LATER: in this case, this may indicate the service unavailability
		logger.Error("failed reading ping response", "reason", err)
		return 1
	}

	jsonBytes, err := json.MarshalIndent(responsePayload, "", "  ")
	if err != nil {
		logger.Error("failed marshaling JSON response", "reason", err)
		return 1
	}

	fmt.Println(string(jsonBytes))

	return 0
}
