package cmd

import (
	"fmt"
	"github.com/lanemets/claim-funnel/interfaces/benerest"
	"os"

	"bitbucket.org/beneregistry/common/viperman"
	"bitbucket.org/beneregistry/common/zaplogger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var _ *zap.SugaredLogger

func init() {
	_ = zaplogger.GetDefaultLogger()
}

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "claim-funnel",
	Short: "Claim Processor",
	Long: `To get started run the serve subcommand which will start a server
on localhost:8080:

    claim-funnel serve

Hit over HTTP 1.1 with curl:

    curl -X GET -k http://localhost:8080/v1/version
`,
}

func Execute(s benerest.Server) {
	RootCmd.AddCommand(NewServeCommand(s))
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {

	initConfig()

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	viperman.InitViper("claim-funnel", viper.GetViper())
}
