package cmd

import (
	"bitbucket.org/beneregistry/common/apiserver"
	"fmt"
	"github.com/lanemets/claim-funnel/interfaces/benerest"
	"github.com/spf13/cobra"
)

func NewServeCommand(s benerest.Server) *cobra.Command {
	serverConfig := apiserver.NewServerConfig("local")
	return &cobra.Command{
		Use:   "serve",
		Short: fmt.Sprintf("Launches the webserver on %s", serverConfig.Addr),
		Run: func(cmd *cobra.Command, args []string) {
			Serve(s, serverConfig)
		}}
}

func Serve(s benerest.Server, serverConfig *apiserver.ServerConfig) {
	s.Start()
}
