package getisloadbalanceridle

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/radekg/yugabyte-db-go-client-api/api"
	"github.com/radekg/yugabyte-db-go-client-api/configs"
	"github.com/spf13/cobra"
)

// Command is the command declaration.
var Command = &cobra.Command{
	Use:   "get-is-load-balancer-idle",
	Short: "Finds out if the load balancer is idle",
	Run:   run,
	Long:  ``,
}

var (
	commandConfig = configs.NewCliConfig()
	logConfig     = configs.NewLogginConfig()
)

func initFlags() {
	Command.Flags().AddFlagSet(commandConfig.FlagSet())
	Command.Flags().AddFlagSet(logConfig.FlagSet())
}

func init() {
	initFlags()
}

func run(cobraCommand *cobra.Command, _ []string) {
	os.Exit(processCommand())
}

func processCommand() int {

	logger := logConfig.NewLogger("get-is-load-balancer-idle")

	for _, validatingConfig := range []configs.ValidatingConfig{commandConfig} {
		if err := validatingConfig.Validate(); err != nil {
			logger.Error("configuration is invalid", "reason", err)
			return 1
		}
	}

	cliClient, err := api.MasterLeaderConnectedClient(commandConfig, logger.Named("client"))
	if err != nil {
		logger.Error("could not connect to a leader master", "reason", err)
		return 1
	}
	defer cliClient.Close()

	registration, err := cliClient.GetIsLoadBalancerIdle()
	if err != nil {
		logger.Error("failed reading load balancer idle state", "reason", err)
		return 1
	}

	jsonBytes, err := json.MarshalIndent(registration, "", "  ")
	if err != nil {
		logger.Error("failed marshaling JSON response", "reason", err)
		return 1
	}

	fmt.Println(string(jsonBytes))

	return 0
}
