package cmd

import (
	"github.com/JoeReid/slb_websocket_server/server"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the websockets server",
	Long:  `Start the websockets server`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Run()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
