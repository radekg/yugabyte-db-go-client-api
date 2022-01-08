package main

import (
	"fmt"
	"os"

	"github.com/radekg/yugabyte-db-go-client-api/cmd/checkexists"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/describetable"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/getisloadbalanceridle"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/getleaderblacklistcompletion"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/getloadmovecompletion"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/getmasterregistration"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/gettabletsfortable"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/getuniverseconfig"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/isloadbalanced"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/isserverready"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/leaderstepdown"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/listmasters"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/listtables"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/listtabletservers"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/masterleaderstepdown"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/modifyplacementinfo"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/ping"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/setloadbalancerstate"
	"github.com/radekg/yugabyte-db-go-client-api/cmd/setpreferredzones"

	snapshotscreate "github.com/radekg/yugabyte-db-go-client-api/cmd/snapshots/create"
	snapshotscreateschedule "github.com/radekg/yugabyte-db-go-client-api/cmd/snapshots/createschedule"
	snapshotsdelete "github.com/radekg/yugabyte-db-go-client-api/cmd/snapshots/delete"
	snapshotsdeleteschedule "github.com/radekg/yugabyte-db-go-client-api/cmd/snapshots/deleteschedule"
	snapshotsexport "github.com/radekg/yugabyte-db-go-client-api/cmd/snapshots/export"
	snapshotsimport "github.com/radekg/yugabyte-db-go-client-api/cmd/snapshots/import"
	snapshotslist "github.com/radekg/yugabyte-db-go-client-api/cmd/snapshots/list"
	snapshotslistrestorations "github.com/radekg/yugabyte-db-go-client-api/cmd/snapshots/listrestorations"
	snapshotslistschedules "github.com/radekg/yugabyte-db-go-client-api/cmd/snapshots/listschedules"
	snapshotsrestore "github.com/radekg/yugabyte-db-go-client-api/cmd/snapshots/restore"
	snapshotsrestoreschedule "github.com/radekg/yugabyte-db-go-client-api/cmd/snapshots/restoreschedule"

	"github.com/radekg/yugabyte-db-go-client-api/cmd/ysqlcatalogversion"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ybcli",
	Short: "ybcli",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(checkexists.Command)
	rootCmd.AddCommand(describetable.Command)
	rootCmd.AddCommand(getisloadbalanceridle.Command)
	rootCmd.AddCommand(getleaderblacklistcompletion.Command)
	rootCmd.AddCommand(getloadmovecompletion.Command)
	rootCmd.AddCommand(getmasterregistration.Command)
	rootCmd.AddCommand(gettabletsfortable.Command)
	rootCmd.AddCommand(getuniverseconfig.Command)
	rootCmd.AddCommand(isloadbalanced.Command)
	rootCmd.AddCommand(isserverready.Command)
	rootCmd.AddCommand(leaderstepdown.Command)
	rootCmd.AddCommand(listmasters.Command)
	rootCmd.AddCommand(listtables.Command)
	rootCmd.AddCommand(listtabletservers.Command)
	rootCmd.AddCommand(masterleaderstepdown.Command)
	rootCmd.AddCommand(modifyplacementinfo.Command)
	rootCmd.AddCommand(ping.Command)
	rootCmd.AddCommand(setloadbalancerstate.Command)
	rootCmd.AddCommand(setpreferredzones.Command)

	rootCmd.AddCommand(snapshotscreateschedule.Command)
	rootCmd.AddCommand(snapshotscreate.Command)
	rootCmd.AddCommand(snapshotsdeleteschedule.Command)
	rootCmd.AddCommand(snapshotsdelete.Command)
	rootCmd.AddCommand(snapshotsexport.Command)
	rootCmd.AddCommand(snapshotsimport.Command)
	rootCmd.AddCommand(snapshotslistrestorations.Command)
	rootCmd.AddCommand(snapshotslistschedules.Command)
	rootCmd.AddCommand(snapshotslist.Command)
	rootCmd.AddCommand(snapshotsrestoreschedule.Command)
	rootCmd.AddCommand(snapshotsrestore.Command)

	rootCmd.AddCommand(ysqlcatalogversion.Command)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
